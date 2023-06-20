package mediasvc

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func HandleHTTP(s Service, router *httprouter.Router) {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encodeMediaError),
	}

	router.Handler(http.MethodPost, "/api/v1/files", kithttp.NewServer(
		makeUploadEndpoint(s),
		decodeUploadRequest,
		kithttp.EncodeJSONResponse,
		opts...,
	))
	router.Handler(http.MethodGet, "/api/v1/files", kithttp.NewServer(
		makeGetFilesInfoEndpoint(s),
		decodeGetFilesInfoRequest,
		kithttp.EncodeJSONResponse,
		opts...,
	))
	router.Handler(http.MethodPut, "/api/v1/files/:id", kithttp.NewServer(
		makeUpdateFilesEndpoint(s),
		decodeUpdateRequest,
		kithttp.EncodeJSONResponse,
		opts...,
	))
	router.Handler(http.MethodDelete, "/api/v1/files/:id", kithttp.NewServer(
		makeDeleteFilesEndpoint(s),
		decodeDeleteRequest,
		kithttp.EncodeJSONResponse,
		opts...,
	))
	router.Handler(http.MethodGet, "/files/:filename", kithttp.NewServer(
		makeGetFilesEndpoint(s),
		decodeGetFileRequest,
		encodeGetFileResponse,
		opts...,
	))
}

func encodeMediaError(_ context.Context, err error, w http.ResponseWriter) {
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

func encodeGetFileResponse(_ context.Context, w http.ResponseWriter, res interface{}) error {
	resFile := res.(*GetFileResponse).File
	file, err := io.ReadAll(resFile)
	if err != nil {
		log.Error(err)
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	if _, err := w.Write(file); err != nil {
		return err
	}

	return nil
}

func decodeUploadRequest(ctx context.Context, r *http.Request) (any, error) {
	// Maximum upload of 10MB files
	r.ParseMultipartForm(10 << 20)

	file, header, err := r.FormFile("file")
	if err != nil {
		log.Error(err)
		return &UploadRequest{}, err
	}
	defer file.Close()

	return UploadRequest{
		File:        file,
		Filename:    header.Filename,
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
		Tags:        r.FormValue("tags"),
		ProjectID:   r.FormValue("projectId"),
	}, nil
}

func decodeGetFilesInfoRequest(ctx context.Context, r *http.Request) (any, error) {
	// get query and convert
	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	offset, _ := strconv.Atoi(q.Get("offset"))
	sortOrder, _ := strconv.Atoi(q.Get("sortOrder"))

	return GetFilesInfoRequest{
		ProjectID: q.Get("projectId"),
		Page:      int64(page),
		Limit:     int64(offset),
		SortField: q.Get("sortField"),
		SortOrder: int64(sortOrder),
	}, nil
}

func decodeGetFileRequest(ctx context.Context, r *http.Request) (any, error) {
	params := httprouter.ParamsFromContext(r.Context())

	return GetFileRequest{
		Filename: params.ByName("filename"),
	}, nil
}

func decodeUpdateRequest(ctx context.Context, r *http.Request) (any, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	var request UpdateRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	params := httprouter.ParamsFromContext(r.Context())
	request.ID = params.ByName("id")

	return request, nil
}

func decodeDeleteRequest(ctx context.Context, r *http.Request) (any, error) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := primitive.ObjectIDFromHex(params.ByName("id"))
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return DeleteRequest{
		ID: id.Hex(),
	}, nil
}
