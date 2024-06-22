package handler

import (
	"errors"
	"strconv"

	"github.com/jamesparry2/Muzz/app/core"
	"github.com/labstack/echo/v4"
)

type HandlerOption struct {
	Core core.CoreIface
}

type Handler struct {
	core core.CoreIface
}

type APIError struct {
	Message     string `json:"message"`
	IsRetryable bool   `json:"is_retryable"`
	Code        string `json:"code"`
}

func NewAPIError(statusCode int, code, message string) APIError {
	isRetryable := false
	switch {
	case statusCode >= 400 && statusCode < 500:
		isRetryable = true
	}

	return APIError{
		Message:     message,
		Code:        code,
		IsRetryable: isRetryable,
	}
}

func NewHandler(opts *HandlerOption) *Handler {
	return &Handler{
		core: opts.Core,
	}
}

func GetUserIDPathParam(c echo.Context) (uint, error) {
	rawID := c.Param("id")
	if rawID == "" {
		return 0, errors.New("missing id in path param")
	}

	id, err := strconv.ParseUint(rawID, 10, 32)
	if err != nil {
		return 0, err
	}

	return uint(id), nil
}
