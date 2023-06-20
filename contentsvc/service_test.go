package contentsvc

import (
	"context"
	"testing"

	"appota/web-builder/content"
	"appota/web-builder/testing/mocks/mock_content"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestGetContent(t *testing.T) {
	const projectId = "c694fdcf-c3ce-496a-b498-13f7dd2facb0"

	repo := new(mock_content.Repository)
	repo.On("GetJSONFile", projectId, "contents").Return([]byte{}, nil).Once()

	s := New(repo)

	r := GetContentsRequest{
		ProjectID: projectId,
	}
	_, err := s.GetContents(context.TODO(), r)

	require.Nil(t, err)
}

func TestSaveContent(t *testing.T) {
	const filename = "c694fdcf-c3ce-496a-b498-13f7dd2facb0.json"
	const filefolder = "contents"

	repo := new(mock_content.Repository)
	repo.On("SaveJSONFile", filename, []byte(nil), filefolder).Return(nil).Once()

	s := New(repo)

	r := SaveContentsRequest{
		ProjectID: "c694fdcf-c3ce-496a-b498-13f7dd2facb0",
	}
	res, err := s.SaveContents(context.TODO(), r)

	require.Nil(t, err)
	require.Equal(t, &SaveContentsResponse{Status: "ok"}, res)
}

func TestGetTemplates(t *testing.T) {
	const id = "6380ddfa147a5e7185f6bcd3"
	const filefolder = "templates"

	repo := new(mock_content.Repository)
	repo.On("GetJSONFile", id, filefolder).Return([]byte{}, nil).Once()

	s := New(repo)

	r := GetTemplatesRequest{
		ID: id,
	}
	_, err := s.GetTemplates(context.TODO(), r)
	if err != nil {
		t.Error(err)
	}

	require.Nil(t, err)
}

func TestGetTemplatesInfo(t *testing.T) {
	repo := new(mock_content.Repository)

	repo.On("FindTemplatesData", mock.MatchedBy(func(arg interface{}) bool {
		actual := arg.(primitive.D)
		expect := bson.D{{Key: "projectId", Value: "c694fdcf-c3ce-496a-b498-13f7dd2facb0"}}

		return assert.Equal(t, expect, actual)
	}), mock.MatchedBy(func(arg interface{}) bool {
		actual := arg.(options.FindOptions)
		expect := options.FindOptions{
			Skip:  actual.Skip,
			Limit: actual.Limit,
			Sort:  actual.Sort,
		}

		return assert.Equal(t, expect, actual)
	})).Return(&mongo.Cursor{}, nil).Once()

	s := New(repo)

	r := GetTemplatesInfoRequest{
		ProjectID: "c694fdcf-c3ce-496a-b498-13f7dd2facb0",
		Page:      1,
		Limit:     10,
		SortField: "uploadedAt",
		SortOrder: -1,
	}

	_, err := s.GetTemplatesInfo(context.TODO(), r)

	require.Nil(t, err)
}

func TestSaveTemplates(t *testing.T) {
	const filefolder = "templates"

	data := &content.Template{
		Name:      "Navigation Bar",
		Type:      "landing-page",
		Tags:      "home;main",
		ProjectID: "c694fdcf-c3ce-496a-b498-13f7dd2facb0",
	}

	repo := new(mock_content.Repository)
	repo.On("SaveJSONFile", mock.MatchedBy(func(arg interface{}) bool {
		actual := arg.(string)
		expect := actual

		return assert.Equal(t, expect, actual)
	}), mock.MatchedBy(func(arg interface{}) bool {
		actual := arg.([]byte)
		expect := actual

		return assert.Equal(t, expect, actual)
	}), filefolder).Return(nil).Once()
	repo.On("InsertTemplatesData", mock.MatchedBy(func(arg interface{}) bool {
		actual := arg.(*content.Template)
		expect := actual

		return assert.Equal(t, expect, actual)
	})).Return(nil).Once()

	s := New(repo)

	r := SaveTemplatesRequest{
		ProjectID: data.ProjectID,
		Name:      data.Name,
		Type:      data.Type,
		Tags:      data.Tags,
		Content:   nil,
	}

	res, err := s.SaveTemplates(context.TODO(), r)
	if err != nil {
		t.Error(err)
	}

	expect := &SaveTemplatesResponse{
		Status:  "ok",
		Content: *data,
	}
	expect.Content.ID = res.Content.ID

	require.Nil(t, err)
	require.Equal(t, expect, res)
}

func TestDeleteTemplates(t *testing.T) {
	const id = "6380ddfa147a5e7185f6bcd3"
	const filefolder = "templates"

	repo := new(mock_content.Repository)
	repo.On("DeleteJSONFile", id, filefolder).Return(nil).Once()
	repo.On("DeleteTemplatesData", id).Return(&mongo.DeleteResult{}, nil).Once()

	s := New(repo)

	r := DeleteTemplatesRequest{
		ID: id,
	}

	_, err := s.DeleteTemplates(context.TODO(), r)
	if err != nil {
		t.Error(err)
	}

	require.Nil(t, err)
}

func TestUpdate(t *testing.T) {
	const id = "6380ddfa147a5e7185f6bcd3"
	data := &content.Template{
		Name:      "Toolbar",
		Type:      "dashboard",
		Tags:      "utils",
		ProjectID: "c694fdcf-c3ce-496a-b498-13f7dd2facb0",
	}

	repo := new(mock_content.Repository)
	s := New(repo)

	repo.On("UpdateTemplatesData", id, data).Return(&content.Template{}, nil).Once()

	r := UpdateTemplatesRequest{
		ID:        id,
		Type:      data.Type,
		Name:      data.Name,
		Tags:      data.Tags,
		ProjectID: data.ProjectID,
	}

	_, err := s.UpdateTemplates(context.TODO(), r)
	require.Nil(t, err)
}
