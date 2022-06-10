package service_test

import (
	"context"
	"testing"

	"github.com/mcosta74/service"
)

func Test_service_GetHealth(t *testing.T) {
	s := service.NewService()

	t.Run("default", func(t *testing.T) {
		ctx := context.Background()
		want := "ok"
		if got := s.GetHealth(ctx); got != want {
			t.Errorf("service.GetHealth() = %v, want %v", got, want)
		}
	})
}

func Test_service_Cheers(t *testing.T) {
	s := service.NewService()
	ctx := context.Background()

	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{"massimo", "massimo", "hello massimo", false},
		{"thierry", "thierry", "hello thierry", false},
		{"empty", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.Cheers(ctx, tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Error mismatch got=(%v != nil) want=%v", err, tt.wantErr)
			}

			if got != tt.want {
				t.Errorf("Cheers: got=%q want=%q", got, tt.want)
			}
		})
	}
}
