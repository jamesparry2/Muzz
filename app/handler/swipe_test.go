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

func TestSwipe(t *testing.T) {
	validJSON := `{ "matched_id": 1, "is_desired": "YES" }`

	t.Run("should return an error when an invalid request body is sent", func(t *testing.T) {
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(httptest.NewRequest(http.MethodPost, "/", nil), rec)

		client := handler.NewHandler(&handler.HandlerOption{})

		// Random new line needed because Echo is adding a new line in the body response??
		expectedBody := `{"message":"invalid body request sent","is_retryable":true,"code":"swipe"}
`

		if assert.NoError(t, client.Swipe(c), "unexpected error was returned") {
			assert.Equal(t, http.StatusBadRequest, rec.Code, "unexpected status code was returned")
			assert.Equal(t, expectedBody, rec.Body.String(), "unexpected JSON was returned")
		}
	})

	t.Run("should return an error when making a request without a id in the path", func(t *testing.T) {
		// Setup HTTP Handlers
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(validJSON))
		c := echo.New().NewContext(req, rec)

		client := handler.NewHandler(&handler.HandlerOption{})

		expectedBody := `{"message":"missing id in path param","is_retryable":true,"code":"swipe"}
`

		if assert.NoError(t, client.Swipe(c), "unexpected error was returned") {
			assert.Equal(t, http.StatusBadRequest, rec.Code, "unexpected status code was returned")
			assert.Equal(t, expectedBody, rec.Body.String(), "unexpected JSON was returned")
		}
	})

	t.Run("should return an error when the swipe fails to process the request", func(t *testing.T) {
		// Setup HTTP Handlers
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(validJSON))
		c := echo.New().NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")

		client := handler.NewHandler(&handler.HandlerOption{
			Core: &mock.MockCore{
				MockSwipe: func(ctx context.Context, request *core.SwipeRequest) (*core.SwipeResponse, error) {
					return nil, errors.New("even batman couldn't fix this")
				},
			},
		})

		expectedBody := `{"message":"even batman couldn't fix this","is_retryable":false,"code":"swipe"}
`

		if assert.NoError(t, client.Swipe(c), "unexpected error was returned") {
			assert.Equal(t, http.StatusInternalServerError, rec.Code, "unexpected status code was returned")
			assert.Equal(t, expectedBody, rec.Body.String(), "unexpected JSON was returned")
		}
	})

	t.Run("should return a success response when the swipe has been processed", func(t *testing.T) {
		// Setup HTTP Handlers
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(validJSON))
		c := echo.New().NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")

		client := handler.NewHandler(&handler.HandlerOption{
			Core: &mock.MockCore{
				MockSwipe: func(ctx context.Context, request *core.SwipeRequest) (*core.SwipeResponse, error) {
					return &core.SwipeResponse{HasMatched: true, MatchedID: 22}, nil
				},
			},
		})

		expectedBody := `{"result":{"matched":true,"matched_id":22}}
`

		if assert.NoError(t, client.Swipe(c), "unexpected error was returned") {
			assert.Equal(t, http.StatusOK, rec.Code, "unexpected status code was returned")
			assert.Equal(t, expectedBody, rec.Body.String(), "unexpected JSON was returned")
		}
	})
}
