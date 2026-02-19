package company

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

func (h *Handler) GetProfile(ctx echo.Context) error {
	company, err := h.service.GetProfile(ctx.Request().Context())
	if err != nil {
		logger.Errorw("Get company profile failed: ", err)

		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Success get company profile", company, nil, nil)
}

func (h *Handler) UpdateProfile(ctx echo.Context) error {
	var req UpdateCompanyProfileRequest
	if err := ctx.Bind(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}

	if err := ctx.Validate(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}

	file, _ := ctx.FormFile("logo_url")

	err := h.service.UpdateProfile(ctx.Request().Context(), &req, file)
	if err != nil {
		logger.Errorw("Update company profile failed : ", err)

		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Update company profile successfully", nil, nil, nil)
}
