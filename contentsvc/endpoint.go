package contentsvc

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

func makeSaveContentEnpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, r interface{}) (interface{}, error) {
		return s.SaveContents(ctx, r.(SaveContentsRequest))
	}
}

func makeGetContentEnpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, r interface{}) (interface{}, error) {
		return s.GetContents(ctx, r.(GetContentsRequest))
	}
}

func makeGetTemplatesInfoEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, r interface{}) (interface{}, error) {
		return s.GetTemplatesInfo(ctx, r.(GetTemplatesInfoRequest))
	}
}

func makeGetTemplatesEnpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, r interface{}) (interface{}, error) {
		return s.GetTemplates(ctx, r.(GetTemplatesRequest))
	}
}

func makeSaveTemplatesEnpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, r interface{}) (interface{}, error) {
		return s.SaveTemplates(ctx, r.(SaveTemplatesRequest))
	}
}

func makeDeleteTemplatesEnpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, r interface{}) (interface{}, error) {
		return s.DeleteTemplates(ctx, r.(DeleteTemplatesRequest))
	}
}

func makeUpdateTemplatesEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, r interface{}) (interface{}, error) {
		return s.UpdateTemplates(ctx, r.(UpdateTemplatesRequest))
	}
}
