package overtime

import (
	"basekarya-backend/pkg/constants"
	"basekarya-backend/pkg/logger"
	"basekarya-backend/pkg/response"
	"basekarya-backend/pkg/utils"
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

func (h *Handler) Create(ctx echo.Context) error {
	userContext, err := utils.GetUserContext(ctx)
	if err != nil {
		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	var req OvertimeRequest
	if err := ctx.Bind(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}

	req.UserID = userContext.UserID
	req.EmployeeID = *userContext.EmployeeID

	if err := ctx.Validate(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}

	err = h.service.Create(ctx.Request().Context(), &req)
	if err != nil {
		logger.Errorw("overtime create failed: ", err)
		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusCreated, "Overtime created successfully", nil, nil, nil)
}

func (h *Handler) GetAll(ctx echo.Context) error {
	userContext, err := utils.GetUserContext(ctx)
	if err != nil {
		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	status := ctx.QueryParam("status")
	page, _ := strconv.Atoi(ctx.QueryParam("page"))
	limit, _ := strconv.Atoi(ctx.QueryParam("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	filter := OvertimeFilter{
		Status: status,
		Page:   page,
		Limit:  limit,
	}

	if userContext.Role != string(constants.UserRoleSuperadmin) {
		filter.UserID = userContext.UserID
	}

	data, meta, err := h.service.GetList(ctx.Request().Context(), filter)
	if err != nil {
		logger.Errorw("get overtimes failed: ", err)
		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Get Overtimes Success", data, nil, meta)
}

func (h *Handler) GetDetail(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "invalid id", nil, err, nil)
	}

	data, err := h.service.GetDetail(ctx.Request().Context(), uint(id))
	if err != nil {
		logger.Errorw("get overtime detail failed: ", err)
		return response.NewResponses[any](ctx, http.StatusBadRequest, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Get Overtime Detail Success", data, nil, nil)
}

func (h *Handler) ProcessAction(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "invalid id", nil, err, nil)
	}

	userContext, err := utils.GetUserContext(ctx)
	if err != nil {
		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	var req ActionRequest
	if err := ctx.Bind(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}

	req.ID = uint(id)
	req.SuperAdminID = userContext.UserID

	if err := ctx.Validate(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}

	err = h.service.ProcessAction(ctx.Request().Context(), &req)
	if err != nil {
		logger.Errorw("Process approval action overtime failed: %w", err)
		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Process Approval Action Overtime Success", nil, nil, nil)
}

func (h *Handler) Export(ctx echo.Context) error {
	userContext, err := utils.GetUserContext(ctx)
	if err != nil {
		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	status := ctx.QueryParam("status")

	filter := OvertimeFilter{
		Status: status,
	}

	if userContext.Role != string(constants.UserRoleSuperadmin) {
		filter.UserID = userContext.UserID
	}

	excelFile, err := h.service.Export(ctx.Request().Context(), filter)
	if err != nil {
		logger.Errorw("export overtimes failed: ", err)
		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	ctx.Response().Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Response().Header().Set("Content-Disposition", "attachment; filename=overtimes.xlsx")
	return ctx.Blob(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", excelFile)
}
