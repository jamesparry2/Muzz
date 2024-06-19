package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jamesparry2/Muzz/app/core"
	"github.com/labstack/echo/v4"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func (h *Handler) Login(ctx echo.Context) error {
	body := LoginRequest{}
	if err := json.NewDecoder(ctx.Request().Body).Decode(&body); err != nil {
		return ctx.JSON(http.StatusBadRequest, NewAPIError(http.StatusBadRequest, "login", "invalid request body sent"))
	}

	// Add some base validation, such as min lengths
	loginResponse, err := h.core.Login(ctx.Request().Context(), &core.LoginRequest{
		Email:    body.Email,
		Password: body.Password,
	})

	if err != nil {
		switch {
		case errors.Is(err, core.ErrLoginInvalidCreds):
			return ctx.JSON(http.StatusUnauthorized, NewAPIError(http.StatusUnauthorized, "login", err.Error()))
		default:
			return ctx.JSON(http.StatusInternalServerError, NewAPIError(http.StatusInternalServerError, "login", err.Error()))
		}
	}

	return ctx.JSON(http.StatusOK, LoginResponse{Token: loginResponse.Token})
}
