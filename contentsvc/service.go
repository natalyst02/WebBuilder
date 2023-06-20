package contentsvc

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"appota/web-builder/content"

	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	contentStoreSubdirectoryName  = "contents"
	templateStoreSubdirectoryName = "templates"
)

type Service interface {
	GetContents(context.Context, GetContentsRequest) (*GetContentsResponse, error)
	SaveContents(context.Context, SaveContentsRequest) (*SaveContentsResponse, error)
	GetTemplatesInfo(context.Context, GetTemplatesInfoRequest) (any, error)
	GetTemplates(context.Context, GetTemplatesRequest) (*GetTemplatesResponse, error)
	SaveTemplates(context.Context, SaveTemplatesRequest) (*SaveTemplatesResponse, error)
	DeleteTemplates(context.Context, DeleteTemplatesRequest) (*DeleteTemplatesResponse, error)
	UpdateTemplates(context.Context, UpdateTemplatesRequest) (*UpdateTemplatesResponse, error)
}

type service struct {
	contentRepo content.Repository
}

func (s *service) GetContents(ctx context.Context, r GetContentsRequest) (*GetContentsResponse, error) {
	err := validator.New().Var(r.ProjectID, "required,uuid")
	if err != nil {
		e := errors.New("parameter have incorrect format")
		log.Error(e)

		return nil, e
	}

	file, err := s.contentRepo.GetJSONFile(r.ProjectID, contentStoreSubdirectoryName)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	log.Infof("got a get %s.json request, returning response...", r.ProjectID)
	return &GetContentsResponse{
		File: file,
	}, nil
}

func (s *service) SaveContents(_ context.Context, r SaveContentsRequest) (*SaveContentsResponse, error) {
	err := validator.New().Var(r.ProjectID, "required,uuid")
	if err != nil {
		e := errors.New("parameter have incorrect format")
		log.Error(e)

		return nil, e
	}

	filename := strings.Join([]string{r.ProjectID, "json"}, ".")
	if err := s.contentRepo.SaveJSONFile(filename, r.JSONFile, contentStoreSubdirectoryName); err != nil {
		log.Error(err)
		return nil, err
	}

	log.Infof("%s uploaded successfully", filename)
	return &SaveContentsResponse{
		Status: "ok",
	}, nil
}

func (s *service) GetTemplatesInfo(ctx context.Context, r GetTemplatesInfoRequest) (any, error) {
	skip := (r.Page - 1) * r.Limit
	opts := options.FindOptions{
		Skip:  &skip,
		Limit: &r.Limit,
	}
	if r.SortField != "" {
		opts.Sort = bson.D{{Key: r.SortField, Value: r.SortOrder}}
	}

	result, err := s.contentRepo.FindTemplatesData(bson.D{{Key: "projectId", Value: r.ProjectID}}, opts)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	var items []content.Template
	if result.RemainingBatchLength() == 0 {
		return &GetTemplatesInfoResponse{
			Total: 0,
			Page:  int(r.Page),
			Items: nil,
		}, nil
	}

	if err := result.All(ctx, &items); err != nil {
		return nil, err
	}

	log.Info("got a get templates request, returning data...")
	return &GetTemplatesInfoResponse{
		Total: len(items),
		Page:  int(r.Page),
		Items: items,
	}, nil
}

func (s *service) GetTemplates(ctx context.Context, content GetTemplatesRequest) (*GetTemplatesResponse, error) {
	file, err := s.contentRepo.GetJSONFile(content.ID, templateStoreSubdirectoryName)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	log.Infof("got a %s.json request, returning response...", content.ID)
	return &GetTemplatesResponse{file}, nil
}

func (s *service) SaveTemplates(ctx context.Context, r SaveTemplatesRequest) (*SaveTemplatesResponse, error) {
	templateID := primitive.NewObjectID()
	templateName := strings.Join([]string{templateID.Hex(), "json"}, ".")

	b, err := json.Marshal(r.Content)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	s.contentRepo.SaveJSONFile(templateName, b, templateStoreSubdirectoryName)

	metadata := &content.Template{
		ID:        templateID,
		Type:      r.Type,
		Name:      r.Name,
		Tags:      r.Tags,
		ProjectID: r.ProjectID,
	}
	err = s.contentRepo.InsertTemplatesData(metadata)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	log.Infof("%s saved successfully", metadata.Name)
	return &SaveTemplatesResponse{
		Status:  "ok",
		Content: *metadata,
	}, nil
}

func (s *service) DeleteTemplates(ctx context.Context, data DeleteTemplatesRequest) (*DeleteTemplatesResponse, error) {
	err := s.contentRepo.DeleteJSONFile(data.ID, templateStoreSubdirectoryName)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	_, err = s.contentRepo.DeleteTemplatesData(data.ID)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	log.Infof("%s deleted successfully", data.ID)
	return &DeleteTemplatesResponse{
		Status: "ok",
	}, nil
}

func (s *service) UpdateTemplates(ctx context.Context, data UpdateTemplatesRequest) (*UpdateTemplatesResponse, error) {
	updateData := &content.Template{
		Name:      data.Name,
		Type:      data.Type,
		Tags:      data.Tags,
		ProjectID: data.ProjectID,
	}

	res, err := s.contentRepo.UpdateTemplatesData(data.ID, updateData)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	log.Info("update request has succeeded")
	return &UpdateTemplatesResponse{
		Status:  "ok",
		Content: res,
	}, nil
}

func New(repo content.Repository) Service {
	s := &service{
		contentRepo: repo,
	}
	return s
}
