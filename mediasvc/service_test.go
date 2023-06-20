package mediasvc

import (
	"context"
	"os"
	"testing"

	"appota/web-builder/media"
	"appota/web-builder/testing/mocks/mock_media"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type test struct {
	name      string
	request   interface{}
	expect    interface{}
	expectErr error
}

func TestUpload(t *testing.T) {
	repo := new(mock_media.Repository)
	s := New(repo)

	repo.On("SaveMedia", mock.IsType([]byte{}), mock.MatchedBy(func(arg interface{}) bool {
		actual := arg.(string)
		expect := actual

		return assert.Equal(t, expect, actual)
	})).Return(nil).Once()
	repo.On("InsertMediaData", mock.MatchedBy(func(arg interface{}) bool {
		actual := arg.(*media.Media)
		expect := &media.Media{
			ID:         actual.ID,
			File:       actual.File,
			Filename:   actual.Filename,
			Filepath:   actual.Filepath,
			UploadedAt: actual.UploadedAt,
		}

		return assert.Equal(t, expect, actual)
	})).Return(nil).Once()

	file, err := os.Create("/tmp/webbuilder-test.png")
	if err != nil {
		log.Error(err)
		t.Error(err)
	}
	r := UploadRequest{
		File: file,
	}

	res, err := s.UploadFiles(context.TODO(), r)

	require.Nil(t, err)
	require.Equal(t, &UploadResponse{Status: "ok", Item: res.Item}, res)

	repo.AssertExpectations(t)
}

func TestGetFilesInfo(t *testing.T) {
	repo := new(mock_media.Repository)
	s := New(repo)

	repo.On("FindMediaData", mock.MatchedBy(func(arg interface{}) bool {
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

	r := GetFilesInfoRequest{
		ProjectID: "c694fdcf-c3ce-496a-b498-13f7dd2facb0",
		Page:      1,
		Limit:     10,
		SortField: "uploadedAt",
		SortOrder: -1,
	}

	_, err := s.GetFilesInfo(context.TODO(), r)

	require.Nil(t, err)
}

func TestGetFile(t *testing.T) {
	tests := []test{
		{
			name: "test 1",
			request: GetFileRequest{
				Filename: "636ca667c052f1a183fc657a.png",
			},
			expect:    "636ca667c052f1a183fc657a.png",
			expectErr: nil,
		},
		{
			name: "test 2",
			request: GetFileRequest{
				Filename: "6371ac9c894feea21eefdc26.png",
			},
			expect:    "6371ac9c894feea21eefdc26.png",
			expectErr: nil,
		},
	}

	repo := new(mock_media.Repository)
	s := New(repo)

	for _, test := range tests {
		repo.On("GetMedia", mock.MatchedBy(func(arg interface{}) bool {
			actual := arg.(string)
			expect := test.expect.(string)

			return assert.Equal(t, expect, actual)
		})).Return([]byte{}, "", nil).Once()

		t.Run(test.name, func(t *testing.T) {
			_, err := s.GetFiles(context.TODO(), test.request.(GetFileRequest))
			require.Nil(t, err)
		})
	}
}

func TestUpdate(t *testing.T) {
	tests := []test{
		{
			name: "test 1",
			request: UpdateRequest{
				ID: "636ca667c052f1a183fc657a",
			},
			expect:    "",
			expectErr: nil,
		},
	}

	repo := new(mock_media.Repository)
	s := New(repo)

	for _, test := range tests {
		repo.On("UpdateMediaData", mock.MatchedBy(func(arg interface{}) bool {
			id, _ := primitive.ObjectIDFromHex(test.request.(UpdateRequest).ID)

			actual := arg.(primitive.D)
			expect := bson.D{{Key: "_id", Value: id}}

			return assert.Equal(t, expect, actual)
		}),
			mock.MatchedBy(func(arg interface{}) bool {
				actual := arg.(*media.UpdateMedia)
				expect := &media.UpdateMedia{}

				return assert.Equal(t, expect, actual)
			})).Return(&mongo.UpdateResult{}, nil)

		t.Run(test.name, func(t *testing.T) {
			res, err := s.UpdateFiles(context.TODO(), test.request.(UpdateRequest))

			require.Nil(t, err)
			require.Equal(t, &UpdateResponse{Status: "ok"}, res)
		})
	}
}

func TestDelete(t *testing.T) {
	tests := []test{
		{
			name: "test_1",
			request: DeleteRequest{
				ID: "636ca667c052f1a183fc657a",
			},
			expect:    "",
			expectErr: nil,
		},
	}

	repo := new(mock_media.Repository)
	s := New(repo)

	for _, test := range tests {
		repo.On("DeleteMediaData", mock.MatchedBy(func(arg interface{}) bool {
			id := test.request.(DeleteRequest).ID
			objectID, _ := primitive.ObjectIDFromHex(id)

			actual := arg.(primitive.D)
			expect := bson.D{{Key: "_id", Value: objectID}}

			return assert.Equal(t, expect, actual)
		})).Return(&mongo.DeleteResult{}, nil)

		t.Run(test.name, func(t *testing.T) {
			res, err := s.DeleteFiles(context.TODO(), test.request.(DeleteRequest))

			require.Nil(t, err)
			require.Equal(t, &DeleteResponse{Status: "ok"}, res)
		})
	}
}
