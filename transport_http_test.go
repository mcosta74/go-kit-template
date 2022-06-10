package service_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/go-kit/log"
	"github.com/mcosta74/service"
)

func makeServer(s service.Service) http.Handler {
	if s == nil {
		s = service.NewService()
	}
	return service.MakeHTTPHandler(service.MakeEndpoints(s, nil, nil), log.NewNopLogger())
}

func TestGetHealth(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/health", nil)
		response := httptest.NewRecorder()

		server := makeServer(nil)

		server.ServeHTTP(response, request)

		want := map[string]any{"status": "ok"}
		var got map[string]any

		if err := json.NewDecoder(response.Body).Decode(&got); err != nil {
			t.Errorf("/health: error unmarshaling response=%v", err)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("/health: got=%+v want=%+v", got, want)
		}
	})
}
