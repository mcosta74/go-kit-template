package service

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

func MakeHTTPHandler(endpoints Endpoints, logger log.Logger) http.Handler {
	mux := chi.NewMux()

	options := []kithttp.ServerOption{
		kithttp.ServerBefore(kithttp.PopulateRequestContext),
		kithttp.ServerBefore(acceptLanguageToContext),
		kithttp.ServerBefore(loggerToContext(logger)),
		kithttp.ServerErrorLogger(level.Error(logger)),
	}

	getHealthHandler := kithttp.NewServer(
		endpoints.GetHealth,
		kithttp.NopRequestDecoder,
		kithttp.EncodeJSONResponse,
		options...,
	)

	cheersHandler := kithttp.NewServer(
		endpoints.Cheers,
		decodeCheersRequest,
		kithttp.EncodeJSONResponse,
		options...,
	)

	mux.Get("/health", getHealthHandler.ServeHTTP)
	mux.Post("/cheers", cheersHandler.ServeHTTP)

	return mux
}

func decodeCheersRequest(ctx context.Context, r *http.Request) (any, error) {
	var req CheersRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, &APIError{Code: http.StatusBadRequest, Err: err}
	}
	return req, nil
}
