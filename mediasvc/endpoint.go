package mediasvc

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

func makeUploadEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, r interface{}) (interface{}, error) {
		return s.UploadFiles(ctx, r.(UploadRequest))
	}
}

func makeGetFilesInfoEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, r interface{}) (interface{}, error) {
		return s.GetFilesInfo(ctx, r.(GetFilesInfoRequest))
	}
}

func makeUpdateFilesEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, r interface{}) (interface{}, error) {
		return s.UpdateFiles(ctx, r.(UpdateRequest))
	}
}

func makeDeleteFilesEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, r interface{}) (interface{}, error) {
		return s.DeleteFiles(ctx, r.(DeleteRequest))
	}
}

func makeGetFilesEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, r interface{}) (interface{}, error) {
		return s.GetFiles(ctx, r.(GetFileRequest))
	}
}
