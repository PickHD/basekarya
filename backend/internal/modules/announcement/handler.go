package announcement

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

func (h *Handler) PublishAnnouncement(ctx echo.Context) error {
	var req CreateAnnouncementRequest
	if err := ctx.Bind(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}

	if err := ctx.Validate(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}

	err := h.service.Publish(ctx.Request().Context(), &req)
	if err != nil {
		logger.Errorw("Publish Announcement failed : ", err)

		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Publish Announcement Success", nil, nil, nil)
}
