package master

import (
	"basekarya-backend/pkg/logger"
	"basekarya-backend/pkg/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

func (h *Handler) GetShifts(ctx echo.Context) error {
	resp, err := h.service.GetAllShifts(ctx.Request().Context())
	if err != nil {
		logger.Errorw("get shifts failed: ", err)

		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Get Shifts Successfully", resp, nil, nil)
}

func (h *Handler) GetLeaveTypes(ctx echo.Context) error {
	resp, err := h.service.GetAllLeaveTypes(ctx.Request().Context())

	if err != nil {
		logger.Errorw("get leave types failed: ", err)

		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Get Leave Types Successfully", resp, nil, nil)
}
