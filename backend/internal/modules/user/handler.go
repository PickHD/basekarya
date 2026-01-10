package user

import (
	"fmt"
	"hris-backend/pkg/logger"
	"hris-backend/pkg/response"
	"hris-backend/pkg/utils"
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

func (h *Handler) GetProfile(ctx echo.Context) error {
	userContext, err := utils.GetUserContext(ctx)
	if err != nil {
		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	user, err := h.service.GetProfile(userContext.UserID)
	if err != nil {
		logger.Errorw("Get profile failed: ", err)

		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Success get profile", user, nil, nil)
}

func (h *Handler) UpdateProfile(ctx echo.Context) error {
	userContext, err := utils.GetUserContext(ctx)
	if err != nil {
		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	var req UpdateProfileRequest
	if err := ctx.Bind(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}

	file, _ := ctx.FormFile("photo")

	err = h.service.UpdateProfile(ctx.Request().Context(), userContext.UserID, &req, file)
	if err != nil {
		logger.Errorw("Update profile failed : ", err)

		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Update profile successfully", nil, nil, nil)
}

func (h *Handler) ChangePassword(ctx echo.Context) error {
	userContext, err := utils.GetUserContext(ctx)
	if err != nil {
		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	var req ChangePasswordRequest
	if err := ctx.Bind(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}

	err = h.service.ChangePassword(userContext.UserID, &req)
	if err != nil {
		logger.Errorw("Change password failed : ", err)

		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Password changed successfully", nil, nil, nil)
}

func (h *Handler) GetAllEmployees(ctx echo.Context) error {
	page := 1
	limit := 10
	search := ctx.QueryParam("search")

	if p := ctx.QueryParam("page"); p != "" {
		fmt.Sscanf(p, "%d", &page)
	}
	if l := ctx.QueryParam("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}

	data, meta, err := h.service.GetAllEmployees(ctx.Request().Context(), page, limit, search)
	if err != nil {
		logger.Errorw("Get All Employees failed: ", err)

		return response.NewResponses[any](ctx, http.StatusInternalServerError, "Failed to get all employees", nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Success get all employees", data, nil, meta)
}

func (h *Handler) CreateEmployee(ctx echo.Context) error {
	var req CreateEmployeeRequest
	if err := ctx.Bind(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}

	err := h.service.CreateEmployee(ctx.Request().Context(), &req)
	if err != nil {
		logger.Errorw("failed to create employee: ", err)

		return response.NewResponses[any](ctx, http.StatusInternalServerError, "Failed to create employee", err.Error(), err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusCreated, "Employee created successfully", nil, nil, nil)
}

func (h *Handler) UpdateEmployee(ctx echo.Context) error {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var req UpdateEmployeeRequest
	if err := ctx.Bind(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}

	err := h.service.UpdateEmployee(ctx.Request().Context(), uint(id), &req)
	if err != nil {
		logger.Errorw("failed to update employee: ", err)

		return response.NewResponses[any](ctx, http.StatusInternalServerError, "Failed to update", err.Error(), err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Employee updated successfully", nil, nil, nil)
}

func (h *Handler) DeleteEmployee(ctx echo.Context) error {
	id, _ := strconv.Atoi(ctx.Param("id"))
	err := h.service.DeleteEmployee(ctx.Request().Context(), uint(id))
	if err != nil {
		logger.Errorw("failed to upddeleteate employee: ", err)

		return response.NewResponses[any](ctx, http.StatusInternalServerError, "Failed to delete", err.Error(), err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Employee deleted successfully", nil, nil, nil)
}
