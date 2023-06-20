package contentsvc

import "appota/web-builder/content"

type (
	SaveContentsRequest struct {
		ProjectID string `validate:"required,uuid"`
		JSONFile  []byte
	}
	SaveContentsResponse struct {
		Status string `json:"status,omitempty"`
		Error  string `json:"error,omitempty"`
	}
)

type (
	GetContentsRequest struct {
		ProjectID string `validate:"required,uuid"`
	}
	GetContentsResponse struct {
		File []byte `json:"-"`
	}
)

type (
	GetTemplatesInfoRequest struct {
		ProjectID string
		Page      int64
		Limit     int64
		SortField string
		SortOrder int64
	}
	GetTemplatesInfoResponse struct {
		Total int                `json:"total"`
		Page  int                `json:"page"`
		Items []content.Template `json:"items"`
	}
)

type (
	GetTemplatesRequest struct {
		ID string
	}
	GetTemplatesResponse struct {
		File []byte
	}
)

type (
	SaveTemplatesRequest struct {
		Name      string      `json:"name"`
		Type      string      `json:"type"`
		Tags      string      `json:"tags"`
		ProjectID string      `json:"projectId"`
		Content   interface{} `json:"content"`
	}
	SaveTemplatesResponse struct {
		Status  string           `json:"status"`
		Content content.Template `json:"content"`
	}
)

type (
	DeleteTemplatesRequest struct {
		ID string
	}
	DeleteTemplatesResponse struct {
		Status string `json:"status,omitempty"`
	}
)

type (
	UpdateTemplatesRequest struct {
		ID        string
		Type      string `json:"type,omitempty"`
		Name      string `json:"name,omitempty"`
		Tags      string `json:"tags,omitempty"`
		ProjectID string `json:"projectId,omitempty"`
	}
	UpdateTemplatesResponse struct {
		Status  string            `json:"status,omitempty"`
		Content *content.Template `json:"content"`
	}
)
