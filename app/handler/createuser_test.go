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

func TestCreateUser(t *testing.T) {
	t.Run("should return an error when making a request with a poorly formatted body", func(t *testing.T) {
		// Setup HTTP Handlers
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(httptest.NewRequest(http.MethodPost, "/", nil), rec)

		client := handler.NewHandler(&handler.HandlerOption{})

		// Random new line needed because Echo is adding a new line in the body response??
		expectedBody := `{"message":"invalid body request sent","is_retryable":true,"code":"create_user"}
`

		if assert.NoError(t, client.CreateUser(c), "unexpected error was returned") {
			assert.Equal(t, http.StatusBadRequest, rec.Code, "unexpected status code was returned")
			assert.Equal(t, expectedBody, rec.Body.String(), "unexpected JSON was returned")
		}
	})

	t.Run("should return an error when the user fails to create", func(t *testing.T) {
		// Setup HTTP Handlers
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{ "email": "jamesparr2@gmail.com" }`)), rec)

		client := handler.NewHandler(&handler.HandlerOption{
			Core: &mock.MockCore{
				MockCreateUser: func(ctx context.Context, request *core.CreateUserRequest) (*core.CreateUserResponse, error) {
					return nil, errors.New("failed to create user")
				},
			},
		})

		expectedBody := `{"message":"failed to create user","is_retryable":false,"code":"create_user"}
`

		if assert.NoError(t, client.CreateUser(c), "unexpected error was returned") {
			assert.Equal(t, http.StatusInternalServerError, rec.Code, "unexpected status code was returned")
			assert.Equal(t, expectedBody, rec.Body.String(), "unexpected JSON was returned")
		}
	})

	t.Run("should return the newly created user on success", func(t *testing.T) {
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{ "email": "jamesparr2@gmail.com" }`)), rec)

		client := handler.NewHandler(&handler.HandlerOption{
			Core: &mock.MockCore{
				MockCreateUser: func(ctx context.Context, request *core.CreateUserRequest) (*core.CreateUserResponse, error) {
					return &core.CreateUserResponse{Email: "james", Password: "some-password", ID: 1, Name: "James", Gender: "male", Age: 28}, nil
				},
			},
		})

		expectedBody := `{"result":{"id":1,"email":"james","password":"some-password","name":"James","gender":"male","age":28}}
`

		if assert.NoError(t, client.CreateUser(c), "unexpected error was returned") {
			assert.Equal(t, http.StatusCreated, rec.Code, "unexpected status code was returned")
			assert.Equal(t, expectedBody, rec.Body.String(), "unexpected JSON was returned")
		}
	})
}
