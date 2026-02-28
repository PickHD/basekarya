package leave

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

func (h *Handler) Apply(ctx echo.Context) error {
	userContext, err := utils.GetUserContext(ctx)
	if err != nil {
		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	var req ApplyRequest
	if err := ctx.Bind(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}

	req.UserID = userContext.UserID
	req.EmployeeID = *userContext.EmployeeID

	if err := ctx.Validate(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}

	err = h.service.Apply(ctx.Request().Context(), &req)
	if err != nil {
		logger.Errorw("Apply Leave Request failed: %w", err)

		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusCreated, "Apply Leave Request Success", nil, nil, nil)
}

func (h *Handler) RequestAction(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "invalid id", nil, err, nil)
	}

	userContext, err := utils.GetUserContext(ctx)
	if err != nil {
		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	var req LeaveActionRequest
	if err := ctx.Bind(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}

	req.RequestID = uint(id)
	req.ApproverID = userContext.UserID

	if err := ctx.Validate(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}

	err = h.service.RequestAction(ctx.Request().Context(), &req)
	if err != nil {
		logger.Errorw("Process approval action leave request failed: %w", err)

		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Process Approval Action Leave Request Success", nil, nil, nil)
}

func (h *Handler) GetAll(ctx echo.Context) error {
	userContext, err := utils.GetUserContext(ctx)
	if err != nil {
		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	status := ctx.QueryParam("status")
	page, _ := strconv.Atoi(ctx.QueryParam("page"))
	limit, _ := strconv.Atoi(ctx.QueryParam("limit"))
	search := ctx.QueryParam("search")

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	filter := LeaveFilter{
		Status: status,
		Page:   page,
		Limit:  limit,
		Search: search,
	}

	if userContext.Role != string(constants.UserRoleSuperadmin) {
		filter.UserID = userContext.UserID
	}

	data, meta, err := h.service.GetList(ctx.Request().Context(), &filter)
	if err != nil {
		logger.Errorw("get leave requests failed: ", err)

		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Get Leave Request List Success", data, nil, meta)
}

func (h *Handler) GetDetail(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "invalid id", nil, err, nil)
	}

	data, err := h.service.GetDetail(ctx.Request().Context(), uint(id))
	if err != nil {
		logger.Errorw("get leave request detail failed: ", err)

		return response.NewResponses[any](ctx, http.StatusBadRequest, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Get Leave Request Detail Success", data, nil, nil)
}

func (h *Handler) Export(ctx echo.Context) error {
	userContext, err := utils.GetUserContext(ctx)
	if err != nil {
		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	status := ctx.QueryParam("status")
	search := ctx.QueryParam("search")

	filter := LeaveFilter{
		Status: status,
		Search: search,
	}

	if userContext.Role != string(constants.UserRoleSuperadmin) {
		filter.UserID = userContext.UserID
	}

	excelFile, err := h.service.Export(ctx.Request().Context(), &filter)
	if err != nil {
		logger.Errorw("export leave requests failed: ", err)
		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	ctx.Response().Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Response().Header().Set("Content-Disposition", "attachment; filename=leaves.xlsx")
	return ctx.Blob(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", excelFile)
}
