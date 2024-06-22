package handler_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jamesparry2/Muzz/app/core"
	"github.com/jamesparry2/Muzz/app/core/mock"
	"github.com/jamesparry2/Muzz/app/handler"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestPreference(t *testing.T) {
	validJSON := `{"gender": "male", "max_age": 0, "min_age": 0}`

	t.Run("should return an error when making a request without a id in the path", func(t *testing.T) {
		// Setup HTTP Handlers
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/user/:id/preference", nil)
		c := echo.New().NewContext(req, rec)

		client := handler.NewHandler(&handler.HandlerOption{})

		expectedBody := `{"message":"missing id in path param","is_retryable":true,"code":"preference"}
`

		if assert.NoError(t, client.Preference(c), "unexpected error was returned") {
			assert.Equal(t, http.StatusBadRequest, rec.Code, "unexpected status code was returned")
			assert.Equal(t, expectedBody, rec.Body.String(), "unexpected JSON was returned")
		}
	})

	t.Run("should return an error when the provided request body is invalid", func(t *testing.T) {
		// Setup HTTP Handlers
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/user/:id/preference", nil)
		c := echo.New().NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")

		client := handler.NewHandler(&handler.HandlerOption{})

		expectedBody := `{"message":"invalid body request sent","is_retryable":true,"code":"preference"}
`

		if assert.NoError(t, client.Preference(c), "unexpected error was returned") {
			assert.Equal(t, http.StatusBadRequest, rec.Code, "unexpected status code was returned")
			assert.Equal(t, expectedBody, rec.Body.String(), "unexpected JSON was returned")
		}
	})

	t.Run("should return an error when the location fails to update", func(t *testing.T) {
		// Setup HTTP Handlers
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/user/:id/preference", strings.NewReader(validJSON))
		c := echo.New().NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")

		client := handler.NewHandler(&handler.HandlerOption{
			Core: &mock.MockCore{
				MockPreference: func(ctx context.Context, request *core.PreferenceRequest) error {
					return errors.New("failed to update")
				},
			},
		})

		expectedBody := `{"message":"failed to update","is_retryable":false,"code":"preference"}
`

		if assert.NoError(t, client.Preference(c), "unexpected error was returned") {
			assert.Equal(t, http.StatusInternalServerError, rec.Code, "unexpected status code was returned")
			assert.Equal(t, expectedBody, rec.Body.String(), "unexpected JSON was returned")
		}
	})

	t.Run("should return an success", func(t *testing.T) {
		// Setup HTTP Handlers
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/user/:id/preference", strings.NewReader(validJSON))
		c := echo.New().NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")

		client := handler.NewHandler(&handler.HandlerOption{
			Core: &mock.MockCore{
				MockPreference: func(ctx context.Context, request *core.PreferenceRequest) error {
					return nil
				},
			},
		})

		if assert.NoError(t, client.Preference(c), "unexpected error was returned") {
			assert.Equal(t, http.StatusOK, rec.Code, "unexpected status code was returned")
		}
	})
}
