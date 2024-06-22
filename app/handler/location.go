package handler

import (
	"encoding/json"
	"net/http"

	"github.com/jamesparry2/Muzz/app/core"
	"github.com/labstack/echo/v4"
)

type LocationRequest struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}

// @Location Set Location
// @Description Allows for a user to set their current location to allow the dicovery to order results based on distance
// @Accept json
// @Produce json
// @Param id path int  true  "User ID"
// @Param LocationRequest body LocationRequest true "lat long"
// @Success 200
// @Failure 400 {object} APIError
// @Failure 500 {object} APIError
// @Router /user/{id}/location [post]
func (h *Handler) Location(ctx echo.Context) error {
	userId, err := GetUserIDPathParam(ctx)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, NewAPIError(http.StatusBadRequest, "location", err.Error()))
	}

	locationRequest := LocationRequest{}
	if err := json.NewDecoder(ctx.Request().Body).Decode(&locationRequest); err != nil {
		return ctx.JSON(http.StatusBadRequest, NewAPIError(http.StatusBadRequest, "location", "invalid body request sent"))
	}

	if err := h.core.Location(ctx.Request().Context(), &core.LocationRequest{
		UserID: userId,
		Lat:    locationRequest.Lat,
		Long:   locationRequest.Long,
	}); err != nil {
		return ctx.JSON(http.StatusInternalServerError, NewAPIError(http.StatusInternalServerError, "location", err.Error()))
	}

	return ctx.JSON(http.StatusOK, nil)
}
