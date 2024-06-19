package handler

import (
	"net/http"

	"github.com/jamesparry2/Muzz/app/core"
	"github.com/jamesparry2/Muzz/app/store"
	"github.com/labstack/echo/v4"
)

func (h *Handler) Discovery(ctx echo.Context) error {
	userId, err := GetUserIDPathParam(ctx)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, NewAPIError(http.StatusBadRequest, "discovery", err.Error()))
	}

	response, err := h.core.Discovery(ctx.Request().Context(), &core.DiscoveryRequest{
		UserID: userId,
	})

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, NewAPIError(http.StatusInternalServerError, "discovery", err.Error()))
	}

	return ctx.JSON(http.StatusOK, CollectionResponse[store.User]{
		Results: response.Users,
	})
}
