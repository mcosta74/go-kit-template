package service

import (
	"context"
	"net/http"
	"strings"

	"github.com/go-kit/log"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

type mwContextKey int

const (
	AcceptLanguageKey mwContextKey = iota
	LoggerKey
	ValidatorKey
	TranslatorKey
)

func acceptLanguageToContext(ctx context.Context, r *http.Request) context.Context {
	return context.WithValue(ctx, AcceptLanguageKey, r.Header.Get("Accept-Language"))
}

func getAcceptLanguage(ctx context.Context) string {
	return ctx.Value(AcceptLanguageKey).(string)
}

func loggerToContext(logger log.Logger) func(ctx context.Context, r *http.Request) context.Context {
	return func(ctx context.Context, r *http.Request) context.Context {
		return context.WithValue(ctx, LoggerKey, logger)
	}
}

func getLogger(ctx context.Context) log.Logger {
	logger, ok := ctx.Value(LoggerKey).(log.Logger)
	if !ok {
		return log.NewNopLogger()
	}
	return logger
}

func validatorToContext(ctx context.Context, validate *validator.Validate, uni *ut.UniversalTranslator) context.Context {
	tmp := context.WithValue(ctx, ValidatorKey, validate)
	return context.WithValue(tmp, TranslatorKey, uni)
}

func getValidator(ctx context.Context) *validator.Validate {
	validate, ok := ctx.Value(ValidatorKey).(*validator.Validate)
	if !ok {
		return nil
	}
	return validate
}

func getTranslator(ctx context.Context) ut.Translator {
	uni, ok := ctx.Value(TranslatorKey).(*ut.UniversalTranslator)
	if !ok {
		return nil
	}

	langs := getAcceptedLanguages(ctx)

	tr, _ := uni.GetTranslator("") // default
	for _, l := range langs {
		if temp, found := uni.GetTranslator(getCountryCode(l)); found {
			tr = temp
			break
		}
	}
	return tr
}

func getAcceptedLanguages(ctx context.Context) []string {
	al := getAcceptLanguage(ctx)
	return strings.Split(al, ",")
}

func getCountryCode(l string) string {
	return strings.TrimSpace(l)[0:2]
}
