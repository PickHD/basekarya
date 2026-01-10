package master

import (
	"hris-backend/pkg/logger"
	"hris-backend/pkg/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

func (h *Handler) GetDepartments(ctx echo.Context) error {
	resp, err := h.service.GetAllDepartments()
	if err != nil {
		logger.Errorw("get departments failed: ", err)

		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Get Departments Successfully", resp, nil, nil)
}

func (h *Handler) GetShifts(ctx echo.Context) error {
	resp, err := h.service.GetAllShifts()
	if err != nil {
		logger.Errorw("get shifts failed: ", err)

		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Get Shifts Successfully", resp, nil, nil)
}
