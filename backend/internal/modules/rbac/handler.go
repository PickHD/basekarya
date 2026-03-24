package rbac

import (
	"net/http"
	"strconv"

	"basekarya-backend/pkg/logger"
	"basekarya-backend/pkg/response"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

func (h *Handler) CreateRole(ctx echo.Context) error {
	var req CreateRoleRequest
	if err := ctx.Bind(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}

	if err := ctx.Validate(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Validation Error", nil, err, nil)
	}

	err := h.service.CreateRole(ctx.Request().Context(), &req)
	if err != nil {
		return response.NewResponses[any](ctx, http.StatusInternalServerError, "Failed to create role", nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Role created successfully", nil, nil, nil)
}

func (h *Handler) GetRolePermissions(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Role ID", nil, err, nil)
	}

	data, err := h.service.GetRolePermissions(ctx.Request().Context(), uint(id))
	if err != nil {
		logger.Errorw("Get role permissions failed: ", err)
		return response.NewResponses[any](ctx, http.StatusInternalServerError, "Failed to get role permissions", nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Success getting role permissions", data, nil, nil)
}

func (h *Handler) AssignPermissions(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Role ID", nil, err, nil)
	}

	var req AssignPermissionsRequest
	if err := ctx.Bind(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}

	if err := ctx.Validate(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Validation Error", nil, err, nil)
	}

	err = h.service.AssignPermissions(ctx.Request().Context(), uint(id), &req)
	if err != nil {
		logger.Errorw("failed to assign permissions: ", err)
		return response.NewResponses[any](ctx, http.StatusInternalServerError, "Failed to assign permissions", nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Permissions assigned successfully", nil, nil, nil)
}
