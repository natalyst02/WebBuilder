package mediasvc

import (
	"io"
	"mime/multipart"

	"appota/web-builder/media"
)

type (
	UploadRequest struct {
		File        multipart.File
		Filename    string
		Title       string `json:"title"`
		Description string `json:"description"`
		Tags        string `json:"tags"`
		ProjectID   string `json:"projectId"`
	}
	UploadResponse struct {
		Status string       `json:"status,omitempty"`
		Item   *media.Media `json:"item"`
	}
)

type (
	GetFilesInfoRequest struct {
		ProjectID string
		Page      int64
		Limit     int64
		SortField string
		SortOrder int64
	}

	GetFilesInfoResponse struct {
		Total int           `json:"total"`
		Page  int           `json:"page"`
		Items []media.Media `json:"items"`
	}
)

type (
	UpdateRequest struct {
		ID          string
		Title       string `json:"title,omitempty"`
		Description string `json:"description,omitempty"`
		Tags        string `json:"tags,omitempty"`
		ProjectID   string `json:"projectId,omitempty"`
	}
	UpdateResponse struct {
		Status  string            `json:"status,omitempty"`
		Content media.UpdateMedia `json:"content,omitempty"`
	}
)

type (
	DeleteRequest struct {
		ID string
	}
	DeleteResponse struct {
		Status string `json:"status,omitempty"`
	}
)

type (
	GetFileRequest struct {
		Filename string
	}
	GetFileResponse struct {
		File io.Reader
	}
)
