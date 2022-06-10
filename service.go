package service

import (
	"context"
	"errors"
)

type Service interface {
	GetHealth(ctx context.Context) string
	Cheers(ctx context.Context, name string) (string, error)
}

func NewService() Service {
	return &service{}
}

type service struct {
}

func (s *service) GetHealth(ctx context.Context) string {
	return "ok"
}

func (s *service) Cheers(ctx context.Context, name string) (string, error) {
	if name == "" {
		return "", errors.New("empty string")
	}
	return "hello " + name, nil
}
