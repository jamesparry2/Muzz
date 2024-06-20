package handler

import (
	"encoding/json"
	"net/http"

	"github.com/jamesparry2/Muzz/app/core"
	"github.com/labstack/echo/v4"
)

type PreferenceRequest struct {
	Gender string `json:"gender"`
	MaxAge int    `json:"max_age"`
	MinAge int    `json:"min_age"`
}

func (h *Handler) Preference(ctx echo.Context) error {
	userId, err := GetUserIDPathParam(ctx)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, NewAPIError(http.StatusBadRequest, "discovery", err.Error()))
	}

	preferenceRequest := PreferenceRequest{}
	if err := json.NewDecoder(ctx.Request().Body).Decode(&preferenceRequest); err != nil {
		return ctx.JSON(http.StatusBadRequest, NewAPIError(http.StatusBadRequest, "preference", "invalid body request sent"))
	}

	if err := h.core.Preference(ctx.Request().Context(), &core.PreferenceRequest{
		UserID: userId,
		MaxAge: preferenceRequest.MaxAge,
		MinAge: preferenceRequest.MinAge,
		Gender: preferenceRequest.Gender,
	}); err != nil {
		return ctx.JSON(http.StatusInternalServerError, NewAPIError(http.StatusInternalServerError, "preference", "invalid body request sent"))
	}

	return ctx.JSON(http.StatusOK, nil)
}
