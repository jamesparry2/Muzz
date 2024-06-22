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

func TestLocation(t *testing.T) {
	validJSON := `{"lat": 0.0, "long": 0.0}`

	t.Run("should return an error when making a request without a id in the path", func(t *testing.T) {
		// Setup HTTP Handlers
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/user/:id/location", nil)
		c := echo.New().NewContext(req, rec)

		client := handler.NewHandler(&handler.HandlerOption{})

		expectedBody := `{"message":"missing id in path param","is_retryable":true,"code":"location"}
`

		if assert.NoError(t, client.Location(c), "unexpected error was returned") {
			assert.Equal(t, http.StatusBadRequest, rec.Code, "unexpected status code was returned")
			assert.Equal(t, expectedBody, rec.Body.String(), "unexpected JSON was returned")
		}
	})

	t.Run("should return an error when the provided request body is invalid", func(t *testing.T) {
		// Setup HTTP Handlers
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/user/:id/location", nil)
		c := echo.New().NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")

		client := handler.NewHandler(&handler.HandlerOption{})

		expectedBody := `{"message":"invalid body request sent","is_retryable":true,"code":"location"}
`

		if assert.NoError(t, client.Location(c), "unexpected error was returned") {
			assert.Equal(t, http.StatusBadRequest, rec.Code, "unexpected status code was returned")
			assert.Equal(t, expectedBody, rec.Body.String(), "unexpected JSON was returned")
		}
	})

	t.Run("should return an error when the location fails to update", func(t *testing.T) {
		// Setup HTTP Handlers
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/user/:id/location", strings.NewReader(validJSON))
		c := echo.New().NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")

		client := handler.NewHandler(&handler.HandlerOption{
			Core: &mock.MockCore{
				MockLocation: func(ctx context.Context, request *core.LocationRequest) error {
					return errors.New("failed to update")
				},
			},
		})

		expectedBody := `{"message":"failed to update","is_retryable":false,"code":"location"}
`

		if assert.NoError(t, client.Location(c), "unexpected error was returned") {
			assert.Equal(t, http.StatusInternalServerError, rec.Code, "unexpected status code was returned")
			assert.Equal(t, expectedBody, rec.Body.String(), "unexpected JSON was returned")
		}
	})

	t.Run("should return an success", func(t *testing.T) {
		// Setup HTTP Handlers
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/user/:id/location", strings.NewReader(validJSON))
		c := echo.New().NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")

		client := handler.NewHandler(&handler.HandlerOption{
			Core: &mock.MockCore{
				MockLocation: func(ctx context.Context, request *core.LocationRequest) error {
					return nil
				},
			},
		})

		if assert.NoError(t, client.Location(c), "unexpected error was returned") {
			assert.Equal(t, http.StatusOK, rec.Code, "unexpected status code was returned")
		}
	})
}
