package service_test

import (
	"context"
	"testing"

	"github.com/mcosta74/service"
)

func TestEndpoints(t *testing.T) {
	t.Run("GetHealth", func(t *testing.T) {
		ep := service.MakeEndpoints(service.NewService(), nil, nil)

		got, err := ep.GetHealth(context.Background(), nil)
		if err != nil {
			t.Errorf("endpoints.GetHealth: got unexpected error=%v", err)
		}

		want := service.GetHealthResponse{Status: "ok"}
		if got != want {
			t.Errorf("endpoints.GetHealth: got=%+v, want=%+v", got, want)
		}
	})
}
