package service

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

type Endpoints struct {
	GetHealth endpoint.Endpoint
	Cheers    endpoint.Endpoint
}

type GetHealthResponse struct {
	Status string `json:"status"`
}

type CheersRequest struct {
	Name string `json:"name,omitempty" validate:"required"`
}

type CheersResponse struct {
	Message string `json:"message,omitempty"`
	Err     string `json:"err,omitempty"`
}

func MakeEndpoints(svc Service, validate *validator.Validate, uni *ut.UniversalTranslator) Endpoints {
	return Endpoints{
		GetHealth: makeGetHealth(svc),
		Cheers:    validated(validate, uni)(makeCheers(svc)),
	}
}

func makeGetHealth(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request any) (response any, err error) {
		return GetHealthResponse{
			Status: svc.GetHealth(ctx),
		}, nil
	}
}

func makeCheers(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request any) (response any, err error) {
		req := request.(CheersRequest)

		validate := getValidator(ctx)
		if validate != nil {
			if err := validate.Struct(&req); err != nil {
				return &CheersResponse{}, &APIError{Code: http.StatusBadRequest, Err: err, translator: getTranslator(ctx)}
			}
		}

		r, err := svc.Cheers(ctx, req.Name)
		if err != nil {
			return &CheersResponse{}, &APIError{Err: err}
		}
		return &CheersResponse{Message: r}, nil
	}
}

func validated(validate *validator.Validate, uni *ut.UniversalTranslator) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			return next(validatorToContext(ctx, validate, uni), request)
		}
	}
}
