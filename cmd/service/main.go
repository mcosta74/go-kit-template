package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"reflect"
	"strings"
	"syscall"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/it"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTr "github.com/go-playground/validator/v10/translations/en"
	itTr "github.com/go-playground/validator/v10/translations/it"
	"github.com/mcosta74/service"
	"github.com/oklog/run"
)

func main() {
	logger := setupLogger()
	_ = level.Info(logger).Log("msg", "started")
	defer func() {
		_ = level.Info(logger).Log("msg", "stopped")
	}()

	validate := setupValidator()
	uni := setupTranslator(validate)

	var (
		svc         = service.NewService()
		endpoints   = service.MakeEndpoints(svc, validate, uni)
		httpHandler = service.MakeHTTPHandler(endpoints, log.With(logger, "component", "HTTP"))
	)

	var g run.Group
	{
		// HTTP Handler
		listener, err := net.Listen("tcp", ":8080")
		if err != nil {
			_ = level.Error(logger).Log("component", "HTTP", "during", "Listen", "err", err)
			os.Exit(1)
		}

		g.Add(func() error {
			server := &http.Server{
				Handler:      httpHandler,
				ReadTimeout:  2 * time.Second,
				WriteTimeout: 2 * time.Second,
			}
			return server.Serve(listener)
		}, func(err error) {
			listener.Close()
		})
	}

	{
		// Signal Handler
		g.Add(run.SignalHandler(context.Background(), syscall.SIGTERM, syscall.SIGINT))
	}
	_ = level.Info(logger).Log("msg", g.Run())
}

func setupLogger() log.Logger {
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	logger = level.NewFilter(logger, level.AllowInfo())
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "where", log.DefaultCaller)
	return logger
}

func setupValidator() *validator.Validate {
	validate := validator.New()
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		tag := field.Tag.Get("json")
		return strings.TrimSpace(strings.Split(tag, ",")[0])
	})
	return validate
}

func setupTranslator(validate *validator.Validate) *ut.UniversalTranslator {
	en := en.New()
	uni := ut.New(en, en, it.New())

	tr, _ := uni.GetTranslator("en")
	_ = enTr.RegisterDefaultTranslations(validate, tr)

	tr, _ = uni.GetTranslator("it")
	_ = itTr.RegisterDefaultTranslations(validate, tr)

	return uni
}
