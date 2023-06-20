package mediasvc

import (
	"bytes"
	"context"
	"io"
	"path/filepath"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"appota/web-builder/media"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service interface {
	UploadFiles(context.Context, UploadRequest) (*UploadResponse, error)
	GetFiles(context.Context, GetFileRequest) (*GetFileResponse, error)
	GetFilesInfo(context.Context, GetFilesInfoRequest) (*GetFilesInfoResponse, error)
	DeleteFiles(context.Context, DeleteRequest) (*DeleteResponse, error)
	UpdateFiles(context.Context, UpdateRequest) (*UpdateResponse, error)
}

func (s *service) UploadFiles(ctx context.Context, content UploadRequest) (*UploadResponse, error) {
	mediaID := primitive.NewObjectID()
	mediaName := strings.Join([]string{mediaID.Hex(), filepath.Ext(content.Filename)}, "")
	mediaPath := filepath.Join("/files", mediaName)
	mediaFile, err := io.ReadAll(content.File)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	err = s.mediaRepo.SaveMedia(mediaFile, mediaName)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	metadata := &media.Media{
		ID:          mediaID,
		Filename:    mediaName,
		Filepath:    mediaPath,
		Title:       content.Title,
		Description: content.Description,
		Tags:        content.Tags,
		ProjectID:   content.ProjectID,
		UploadedAt:  time.Now(),
	}
	err = s.mediaRepo.InsertMediaData(metadata)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	log.Infof("%s uploaded successfully", mediaName)
	return &UploadResponse{
		Status: "ok",
		Item:   metadata,
	}, nil
}

func (s *service) GetFiles(ctx context.Context, content GetFileRequest) (*GetFileResponse, error) {
	file, _, err := s.mediaRepo.GetMedia(content.Filename)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	f := bytes.NewReader(file)

	log.Info("received get file request")
	return &GetFileResponse{
		File: f,
	}, nil
}

func (s *service) GetFilesInfo(ctx context.Context, r GetFilesInfoRequest) (*GetFilesInfoResponse, error) {
	filter := bson.D{{Key: "projectId", Value: r.ProjectID}}
	skip := (r.Page - 1) * r.Limit
	opts := options.FindOptions{
		Skip:  &skip,
		Limit: &r.Limit,
	}
	if r.SortField != "" {
		opts.Sort = bson.D{{Key: r.SortField, Value: r.SortOrder}}
	}

	result, err := s.mediaRepo.FindMediaData(filter, opts)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	var items []media.Media
	if result.RemainingBatchLength() == 0 {
		log.Info("received get file data request, but there's no item(s) in database")
		return &GetFilesInfoResponse{
			Total: 0,
			Page:  int(r.Page),
			Items: nil,
		}, nil
	}
	if err := result.All(ctx, &items); err != nil {
		return nil, err
	}

	log.Info("received get file data request, returning data...")
	return &GetFilesInfoResponse{
		Total: len(items),
		Page:  int(r.Page),
		Items: items,
	}, nil
}

func (s *service) DeleteFiles(ctx context.Context, content DeleteRequest) (*DeleteResponse, error) {
	id, err := primitive.ObjectIDFromHex(content.ID)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	filter := bson.D{{Key: "_id", Value: id}}
	_, err = s.mediaRepo.DeleteMediaData(filter)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	log.Info("delete request has succeeded")
	return &DeleteResponse{
		Status: "ok",
	}, nil
}

func (s *service) UpdateFiles(ctx context.Context, content UpdateRequest) (*UpdateResponse, error) {
	id, err := primitive.ObjectIDFromHex(content.ID)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	filter := bson.D{{Key: "_id", Value: id}}
	updateData := &media.UpdateMedia{
		Title:       content.Title,
		Description: content.Description,
		Tags:        content.Tags,
		ProjectID:   content.ProjectID,
	}

	_, err = s.mediaRepo.UpdateMediaData(filter, updateData)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	log.Info("update request has succeeded")
	return &UpdateResponse{
		Status:  "ok",
		Content: *updateData,
	}, nil
}

type service struct {
	mediaRepo media.Repository
}

func New(repo media.Repository) Service {
	return &service{
		mediaRepo: repo,
	}
}
