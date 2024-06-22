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

// @Preference Set Preference
// @Description Allows for a user to filter discoveries by applying specfic preferences
// @Accept json
// @Produce json
// @Param id path int  true  "User ID"
// @Param PreferenceRequest body PreferenceRequest true "gender max_age min_age"
// @Success 200
// @Failure 400 {object} APIError
// @Failure 500 {object} APIError
// @Router /user/{id}/preference [post]
func (h *Handler) Preference(ctx echo.Context) error {
	userId, err := GetUserIDPathParam(ctx)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, NewAPIError(http.StatusBadRequest, "preference", err.Error()))
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
		return ctx.JSON(http.StatusInternalServerError, NewAPIError(http.StatusInternalServerError, "preference", err.Error()))
	}

	return ctx.JSON(http.StatusOK, nil)
}
