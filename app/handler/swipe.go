package handler

import (
	"encoding/json"
	"net/http"

	"github.com/jamesparry2/Muzz/app/core"
	"github.com/labstack/echo/v4"
)

type SwipeRequest struct {
	MatchedID uint   `json:"matched_id"`
	IsDesired string `json:"is_desired"`
}

type SwipeResponse struct {
	Matched   bool `json:"matched"`
	MatchedID uint `json:"matched_id,omitempty"`
}

func (h *Handler) Swipe(ctx echo.Context) error {
	swipeRequest := SwipeRequest{}
	if err := json.NewDecoder(ctx.Request().Body).Decode(&swipeRequest); err != nil {
		return ctx.JSON(http.StatusBadRequest, NewAPIError(http.StatusBadRequest, "swipe", "invalid body request sent"))
	}

	// Again, add validation for minimum length and enum type on is_desired
	userId, err := GetUserIDPathParam(ctx)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, NewAPIError(http.StatusBadRequest, "swipe", err.Error()))
	}

	resp, err := h.core.Swipe(ctx.Request().Context(), &core.SwipeRequest{
		MatchedID: swipeRequest.MatchedID,
		IsDesired: swipeRequest.IsDesired,
		UserID:    userId,
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, NewAPIError(http.StatusInternalServerError, "swipe", err.Error()))
	}

	return ctx.JSON(http.StatusOK,
		SingleResponse{
			Result: SwipeResponse{
				Matched:   resp.HasMatched,
				MatchedID: resp.MatchedID,
			},
		})
}
