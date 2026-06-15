package subscription

import (
	"basekarya-backend/pkg/logger"
	"basekarya-backend/pkg/response"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{s}
}

func (h *Handler) ListPlans(ctx echo.Context) error {
	plans, err := h.service.ListPlans(ctx.Request().Context())
	if err != nil {
		logger.Errorw("Failed to list plans: ", err)
		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}
	return response.NewResponses[any](ctx, http.StatusOK, "Success", plans, nil, nil)
}

func (h *Handler) RequestUpgrade(ctx echo.Context) error {
	var req UpgradeRequest
	if err := ctx.Bind(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}

	if err := ctx.Validate(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}

	result, err := h.service.RequestUpgrade(ctx.Request().Context(), &req)
	if err != nil {
		logger.Errorw("Failed to request upgrade: ", err)
		return response.NewResponses[any](ctx, http.StatusBadRequest, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusCreated, "Upgrade request submitted", result, nil, nil)
}

func (h *Handler) ListPendingRequests(ctx echo.Context) error {
	requests, err := h.service.ListPendingRequests(ctx.Request().Context())
	if err != nil {
		logger.Errorw("Failed to list pending requests: ", err)
		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}
	return response.NewResponses[any](ctx, http.StatusOK, "Success", requests, nil, nil)
}

func (h *Handler) ReviewRequest(ctx echo.Context) error {
	id := ctx.Param("id")
	var requestID uint
	if _, err := fmt.Sscanf(id, "%d", &requestID); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid request ID", nil, err, nil)
	}

	var req ReviewRequest
	if err := ctx.Bind(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}

	if err := ctx.Validate(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}

	if err := h.service.ReviewRequest(ctx.Request().Context(), requestID, &req); err != nil {
		logger.Errorw("Failed to review request: ", err)
		return response.NewResponses[any](ctx, http.StatusBadRequest, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Request reviewed successfully", nil, nil, nil)
}

func (h *Handler) ListAllRequests(ctx echo.Context) error {
	requests, err := h.service.ListAllRequests(ctx.Request().Context())
	if err != nil {
		logger.Errorw("Failed to list all requests: ", err)
		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}
	return response.NewResponses[any](ctx, http.StatusOK, "Success", requests, nil, nil)
}

func (h *Handler) ListCompanies(ctx echo.Context) error {
	search := ctx.QueryParam("search")
	companies, err := h.service.ListCompanies(ctx.Request().Context(), search)
	if err != nil {
		logger.Errorw("Failed to list companies: ", err)
		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}
	return response.NewResponses[any](ctx, http.StatusOK, "Success", companies, nil, nil)
}

func (h *Handler) GetCompanyDetail(ctx echo.Context) error {
	id := ctx.Param("id")
	var companyID uint
	if _, err := fmt.Sscanf(id, "%d", &companyID); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid company ID", nil, err, nil)
	}

	detail, err := h.service.GetCompanyDetail(ctx.Request().Context(), companyID)
	if err != nil {
		logger.Errorw("Failed to get company detail: ", err)
		return response.NewResponses[any](ctx, http.StatusNotFound, "Company not found", nil, err, nil)
	}
	return response.NewResponses[any](ctx, http.StatusOK, "Success", detail, nil, nil)
}

func (h *Handler) UpdateCompanyStatus(ctx echo.Context) error {
	id := ctx.Param("id")
	var companyID uint
	if _, err := fmt.Sscanf(id, "%d", &companyID); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid company ID", nil, err, nil)
	}

	var req UpdateCompanyStatusRequest
	if err := ctx.Bind(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}

	if err := ctx.Validate(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}

	if err := h.service.UpdateCompanyStatus(ctx.Request().Context(), companyID, &req); err != nil {
		logger.Errorw("Failed to update company status: ", err)
		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}
	return response.NewResponses[any](ctx, http.StatusOK, "Company status updated", nil, nil, nil)
}

func (h *Handler) GetDashboardStats(ctx echo.Context) error {
	stats, err := h.service.GetDashboardStats(ctx.Request().Context())
	if err != nil {
		logger.Errorw("Failed to get dashboard stats: ", err)
		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}
	return response.NewResponses[any](ctx, http.StatusOK, "Success", stats, nil, nil)
}

func (h *Handler) RefreshCompanyCache(ctx echo.Context) error {
	id := ctx.Param("id")
	var companyID uint
	if _, err := fmt.Sscanf(id, "%d", &companyID); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid company ID", nil, err, nil)
	}

	_ = h.service.RefreshCompanyCache(ctx.Request().Context(), companyID)

	return response.NewResponses[any](ctx, http.StatusOK, "Cache refreshed for company", nil, nil, nil)
}
