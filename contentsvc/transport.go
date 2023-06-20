package contentsvc

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

func HandleHTTP(s Service, router *httprouter.Router) {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encodeContentError),
	}

	router.Handler(http.MethodGet, "/api/v1/contents/:projectId", kithttp.NewServer(
		makeGetContentEnpoint(s),
		decodeGetContentRequest,
		encodeGetContentsRequest,
		opts...,
	))
	router.Handler(http.MethodPost, "/api/v1/contents/:projectId", kithttp.NewServer(
		makeSaveContentEnpoint(s),
		decodeSaveContentRequest,
		kithttp.EncodeJSONResponse,
		opts...,
	))
	router.Handler(http.MethodDelete, "/api/v1/templates/:id", kithttp.NewServer(
		makeDeleteTemplatesEnpoint(s),
		decodeDeleteTemplatesRequest,
		kithttp.EncodeJSONResponse,
		opts...,
	))
	router.Handler(http.MethodPut, "/api/v1/templates/:id", kithttp.NewServer(
		makeUpdateTemplatesEndpoint(s),
		decodeUpdateTemplatesRequest,
		kithttp.EncodeJSONResponse,
		opts...,
	))
	router.Handler(http.MethodPost, "/api/v1/templates", kithttp.NewServer(
		makeSaveTemplatesEnpoint(s),
		decodeSaveTemplatesRequest,
		kithttp.EncodeJSONResponse,
		opts...,
	))
	router.Handler(http.MethodGet, "/api/v1/templates", kithttp.NewServer(
		makeGetTemplatesInfoEndpoint(s),
		decodeGetTemplatesInfoRequest,
		kithttp.EncodeJSONResponse,
		opts...,
	))
	router.Handler(http.MethodGet, "/api/v1/templates/:id", kithttp.NewServer(
		makeGetTemplatesEnpoint(s),
		decodeGetTemplatesRequest,
		encodeGetTemplatesRequest,
		opts...,
	))
}

func encodeContentError(_ context.Context, err error, w http.ResponseWriter) {
	type errorResponse struct {
		Error string `json:"error"`
	}
	res := errorResponse{
		Error: err.Error(),
	}

	b, err := json.Marshal(res)
	if err != nil {
		log.Error(err)
		return
	}

	contentType, body := "application/json; charset=utf-8", []byte(b)

	w.Header().Set("Content-Type", contentType)
	if headerer, ok := err.(kithttp.Headerer); ok {
		for k, values := range headerer.Headers() {
			for _, v := range values {
				w.Header().Add(k, v)
			}
		}
	}

	code := http.StatusInternalServerError
	if sc, ok := err.(kithttp.StatusCoder); ok {
		code = sc.StatusCode()
	}

	w.WriteHeader(code)
	w.Write(body)
}

func encodeGetContentsRequest(_ context.Context, w http.ResponseWriter, res interface{}) error {
	contentType, body := "application/json; charset=utf-8", res.(*GetContentsResponse).File

	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(http.StatusOK)
	w.Write(body)

	return nil
}

func encodeGetTemplatesRequest(_ context.Context, w http.ResponseWriter, res interface{}) error {
	contentType, body := "application/json; charset=utf-8", res.(*GetTemplatesResponse).File

	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(http.StatusOK)
	w.Write(body)

	return nil
}

func decodeGetContentRequest(ctx context.Context, r *http.Request) (any, error) {
	params := httprouter.ParamsFromContext(r.Context())

	id := params.ByName("projectId")
	return GetContentsRequest{
		ProjectID: id,
	}, nil
}

func decodeSaveContentRequest(ctx context.Context, r *http.Request) (any, error) {
	param := httprouter.ParamsFromContext(r.Context())

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return SaveContentsRequest{param.ByName("projectId"), body}, nil
}

func decodeSaveTemplatesRequest(ctx context.Context, r *http.Request) (any, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	var request SaveTemplatesRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return request, nil
}

func decodeGetTemplatesInfoRequest(ctx context.Context, r *http.Request) (any, error) {
	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	limit, _ := strconv.Atoi(q.Get("limit"))
	sortOrder, _ := strconv.Atoi(q.Get("sortOrder"))

	return GetTemplatesInfoRequest{
		ProjectID: q.Get("projectId"),
		Page:      int64(page),
		Limit:     int64(limit),
		SortField: q.Get("sortField"),
		SortOrder: int64(sortOrder),
	}, nil
}

func decodeGetTemplatesRequest(ctx context.Context, r *http.Request) (any, error) {
	params := httprouter.ParamsFromContext(r.Context())

	id := params.ByName("id")
	return GetTemplatesRequest{
		ID: id,
	}, nil
}

func decodeDeleteTemplatesRequest(ctx context.Context, r *http.Request) (any, error) {
	params := httprouter.ParamsFromContext(r.Context())

	return DeleteTemplatesRequest{
		ID: params.ByName("id"),
	}, nil
}

func decodeUpdateTemplatesRequest(ctx context.Context, r *http.Request) (any, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	var request UpdateTemplatesRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	params := httprouter.ParamsFromContext(r.Context())
	request.ID = params.ByName("id")

	return request, nil
}
