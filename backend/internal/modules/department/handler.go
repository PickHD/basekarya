package department

import (
	"basekarya-backend/pkg/logger"
	"basekarya-backend/pkg/response"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

func (h *Handler) GetAll(ctx echo.Context) error {
	resp, err := h.service.GetAll(ctx.Request().Context())
	if err != nil {
		logger.Errorw("get departments failed: ", err)
		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}
	return response.NewResponses[any](ctx, http.StatusOK, "Get Departments Successfully", resp, nil, nil)
}

func (h *Handler) GetByID(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "invalid id", nil, err, nil)
	}

	resp, err := h.service.GetByID(ctx.Request().Context(), uint(id))
	if err != nil {
		logger.Errorw("get department failed: ", err)
		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Get Department Successfully", resp, nil, nil)
}

func (h *Handler) Create(ctx echo.Context) error {
	var req CreateDepartmentRequest
	if err := ctx.Bind(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "invalid request", nil, err, nil)
	}

	if err := ctx.Validate(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "validation error", nil, err, nil)
	}

	resp, err := h.service.Create(ctx.Request().Context(), &req)
	if err != nil {
		logger.Errorw("create department failed: ", err)
		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusCreated, "Create Department Successfully", resp, nil, nil)
}

func (h *Handler) Update(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "invalid id", nil, err, nil)
	}

	var req UpdateDepartmentRequest
	if err := ctx.Bind(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "invalid request", nil, err, nil)
	}

	if err := ctx.Validate(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "validation error", nil, err, nil)
	}

	resp, err := h.service.Update(ctx.Request().Context(), uint(id), &req)
	if err != nil {
		logger.Errorw("update department failed: ", err)
		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Update Department Successfully", resp, nil, nil)
}

func (h *Handler) Delete(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "invalid id", nil, err, nil)
	}

	if err := h.service.Delete(ctx.Request().Context(), uint(id)); err != nil {
		logger.Errorw("delete department failed: ", err)
		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Delete Department Successfully", nil, nil, nil)
}
