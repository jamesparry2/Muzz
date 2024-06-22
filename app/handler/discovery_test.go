package handler_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jamesparry2/Muzz/app/core"
	"github.com/jamesparry2/Muzz/app/core/mock"
	"github.com/jamesparry2/Muzz/app/handler"
	"github.com/jamesparry2/Muzz/app/store"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestDiscovery(t *testing.T) {
	t.Run("should return an error when making a request without a id in the path", func(t *testing.T) {
		// Setup HTTP Handlers
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/user/:id/discovery", nil)
		c := echo.New().NewContext(req, rec)

		client := handler.NewHandler(&handler.HandlerOption{})

		expectedBody := `{"message":"missing id in path param","is_retryable":true,"code":"discovery"}
`

		if assert.NoError(t, client.Discovery(c), "unexpected error was returned") {
			assert.Equal(t, http.StatusBadRequest, rec.Code, "unexpected status code was returned")
			assert.Equal(t, expectedBody, rec.Body.String(), "unexpected JSON was returned")
		}
	})

	t.Run("should return an error when attempting to find new users to match with", func(t *testing.T) {
		// Setup HTTP Handlers
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/user/:id/discovery", nil)
		c := echo.New().NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")

		client := handler.NewHandler(&handler.HandlerOption{Core: &mock.MockCore{
			MockDiscovery: func(ctx context.Context, request *core.DiscoveryRequest) (*core.DiscoveryResponse, error) {
				return nil, errors.New("unable to do the searching")
			},
		}})

		expectedBody := `{"message":"unable to do the searching","is_retryable":false,"code":"discovery"}
`

		if assert.NoError(t, client.Discovery(c), "unexpected error was returned") {
			assert.Equal(t, http.StatusInternalServerError, rec.Code, "unexpected status code was returned")
			assert.Equal(t, expectedBody, rec.Body.String(), "unexpected JSON was returned")
		}
	})

	t.Run("should return the found users on success", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/user/:id/discovery", nil)
		c := echo.New().NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")

		client := handler.NewHandler(&handler.HandlerOption{Core: &mock.MockCore{
			MockDiscovery: func(ctx context.Context, request *core.DiscoveryRequest) (*core.DiscoveryResponse, error) {
				return &core.DiscoveryResponse{Users: []store.User{{
					Email: "james", Password: "some-password", Name: "James", Gender: "male", Age: 28,
				}}}, nil
			},
		}})
		expectedBody := `{"results":[{"id":0,"name":"James","gender":"male","age":28,"distance_from_me":0}]}
`

		if assert.NoError(t, client.Discovery(c), "unexpected error was returned") {
			assert.Equal(t, http.StatusOK, rec.Code, "unexpected status code was returned")
			assert.Equal(t, expectedBody, rec.Body.String(), "unexpected JSON was returned")
		}
	})
}
