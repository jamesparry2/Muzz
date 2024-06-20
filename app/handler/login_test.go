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

func TestLogin(t *testing.T) {
	validJSON := `{"email": "jamesparry2@gmail.com", "password": "something"}`

	t.Run("should return an error when making a request with a poorly formatted body", func(t *testing.T) {
		// Setup HTTP Handlers
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(httptest.NewRequest(http.MethodPost, "/login", nil), rec)

		client := handler.NewHandler(&handler.HandlerOption{})

		expectedBody := `{"message":"invalid body request sent","is_retryable":true,"code":"login"}
`

		if assert.NoError(t, client.Login(c), "unexpected error was returned") {
			assert.Equal(t, http.StatusBadRequest, rec.Code, "unexpected status code was returned")
			assert.Equal(t, expectedBody, rec.Body.String(), "unexpected JSON was returned")
		}
	})

	t.Run("should return an unauthorized when invalid creds are provided", func(t *testing.T) {
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(validJSON)), rec)

		client := handler.NewHandler(&handler.HandlerOption{
			Core: &mock.MockCore{
				MockLogin: func(ctx context.Context, request *core.LoginRequest) (*core.LoginResponse, error) {
					return nil, core.ErrLoginInvalidCreds
				},
			},
		})

		expectedBody := `{"message":"provided username/password was incorrect","is_retryable":true,"code":"login"}
`

		if assert.NoError(t, client.Login(c), "unexpected error was returned") {
			assert.Equal(t, http.StatusUnauthorized, rec.Code, "unexpected status code was returned")
			assert.Equal(t, expectedBody, rec.Body.String(), "unexpected JSON was returned")
		}
	})

	t.Run("should return an internal error when a terminal error is returned", func(t *testing.T) {
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(validJSON)), rec)

		client := handler.NewHandler(&handler.HandlerOption{
			Core: &mock.MockCore{
				MockLogin: func(ctx context.Context, request *core.LoginRequest) (*core.LoginResponse, error) {
					return nil, errors.New("failure terminally")
				},
			},
		})

		expectedBody := `{"message":"failure terminally","is_retryable":false,"code":"login"}
`

		if assert.NoError(t, client.Login(c), "unexpected error was returned") {
			assert.Equal(t, http.StatusInternalServerError, rec.Code, "unexpected status code was returned")
			assert.Equal(t, expectedBody, rec.Body.String(), "unexpected JSON was returned")
		}
	})

	t.Run("should return a token when a login is successful", func(t *testing.T) {
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(validJSON)), rec)

		client := handler.NewHandler(&handler.HandlerOption{
			Core: &mock.MockCore{
				MockLogin: func(ctx context.Context, request *core.LoginRequest) (*core.LoginResponse, error) {
					return &core.LoginResponse{Token: "supertoken"}, nil
				},
			},
		})

		expectedBody := `{"token":"supertoken"}
`

		if assert.NoError(t, client.Login(c), "unexpected error was returned") {
			assert.Equal(t, http.StatusOK, rec.Code, "unexpected status code was returned")
			assert.Equal(t, expectedBody, rec.Body.String(), "unexpected JSON was returned")
		}
	})
}
