package health

import (
	"basekarya-backend/pkg/logger"
	"basekarya-backend/pkg/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{s}
}

func (h *Handler) HealthCheck(ctx echo.Context) error {
	if err := h.service.Check(); err != nil {
		logger.Errorw("health check failed :", err)

		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "OK", true, nil, nil)
}
