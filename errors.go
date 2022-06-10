package service

import (
	"encoding/json"
	"net/http"
	"strings"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

type APIError struct {
	Code       int
	Err        error
	translator ut.Translator
}

func (e *APIError) StatusCode() int {
	if e.Code == 0 {
		return http.StatusInternalServerError
	}
	return e.Code
}

func (e *APIError) Error() string {
	if errs, ok := e.Err.(validator.ValidationErrors); ok {
		if e.translator != nil {
			errorStrings := e.getStrings(errs)
			return strings.Join(errorStrings, ",")
		}
	}
	return e.Err.Error()
}

func (e *APIError) MarshalJSON() ([]byte, error) {
	if errs, ok := e.Err.(validator.ValidationErrors); ok {
		if e.translator != nil {
			errorStrings := e.getStrings(errs)

			tmp := struct {
				Errors []string `json:"errors"`
			}{Errors: errorStrings}
			return json.Marshal(tmp)
		}
	}

	tmp := struct {
		Err string `json:"err"`
	}{e.Error()}
	return json.Marshal(tmp)
}

func (e *APIError) getStrings(errs validator.ValidationErrors) []string {
	errors := errs.Translate(e.translator)

	errorStrings := make([]string, 0, len(errors))
	for _, v := range errors {
		errorStrings = append(errorStrings, v)
	}
	return errorStrings
}
