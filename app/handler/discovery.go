package handler

import (
	"net/http"

	"github.com/jamesparry2/Muzz/app/core"
	"github.com/jamesparry2/Muzz/app/store"
	"github.com/labstack/echo/v4"
)

type DiscoveryUser struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Gender string `json:"gender"`
	Age    int    `json:"age"`
}

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

	apiResponse := []DiscoveryUser{}
	for _, user := range response.Users {
		apiResponse = append(apiResponse, mapUserToDiscoveryUser(&user))
	}

	return ctx.JSON(http.StatusOK, CollectionResponse[DiscoveryUser]{
		Results: apiResponse,
	})
}

func mapUserToDiscoveryUser(user *store.User) DiscoveryUser {
	return DiscoveryUser{
		ID:     user.ID,
		Name:   user.Name,
		Gender: user.Gender,
		Age:    user.Age,
	}
}
