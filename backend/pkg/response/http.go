package response

import (
	"github.com/labstack/echo/v4"
)

type (
	baseResponse struct {
		Message string `json:"message"`
		Data    any    `json:"data"`
		Error   error  `json:"error"`
		Meta    *Meta  `json:"meta,omitempty"`
	}

	Meta struct {
		Page      int64 `json:"page,omitempty"`
		TotalPage int64 `json:"total_page,omitempty"`
		TotalData int64 `json:"total_data,omitempty"`
	}
)

// NewResponses return dynamic JSON responses
func NewResponses[T any](ctx echo.Context, statusCode int, message string, data T, err error, meta *Meta) error {
	if statusCode < 400 {
		return ctx.JSON(statusCode, &baseResponse{
			Message: message,
			Data:    data,
			Error:   nil,
			Meta:    meta,
		})

	}

	return ctx.JSON(statusCode, &baseResponse{
		Message: message,
		Data:    data,
		Error:   err,
		Meta:    nil,
	})
}
