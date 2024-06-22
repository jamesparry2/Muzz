package handler

import (
	"encoding/json"
	"net/http"

	"github.com/jamesparry2/Muzz/app/core"
	"github.com/labstack/echo/v4"
)

type CreateUserRequest struct {
	Email string `json:"email"`
}

// @CreateUser Create User
// @Description for Testing Purposes for Testing Purposes with a provided Email
// @Accept json
// @Produce json
// @Param CreateUserRequest body CreateUserRequest true "email"
// @Success 200 {object} SingleResponse
// @Failure 400 {object} APIError
// @Failure 500 {object} APIError
// @Router /user/create [post]
func (h *Handler) CreateUser(ctx echo.Context) error {
	// Make this a randomizer
	createUserRequest := CreateUserRequest{}
	if err := json.NewDecoder(ctx.Request().Body).Decode(&createUserRequest); err != nil {
		return ctx.JSON(http.StatusBadRequest, NewAPIError(http.StatusBadRequest, "create_user", "invalid body request sent"))
	}

	response, err := h.core.CreateUser(
		ctx.Request().Context(),
		&core.CreateUserRequest{
			Email:    createUserRequest.Email,
			Password: "superstring",
			Name:     "James",
			Gender:   "Male",
			Age:      28,
		})

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, NewAPIError(http.StatusInternalServerError, "create_user", err.Error()))
	}

	// Need to create a mapper to confirm to the seperation
	return ctx.JSON(http.StatusCreated, SingleResponse{Result: response})
}
