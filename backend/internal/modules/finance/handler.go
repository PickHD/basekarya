package finance

import (
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

func (h *Handler) CreateTransaction(ctx echo.Context) error {
	userContext, err := utils.GetUserContext(ctx)
	if err != nil {
		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	var req CreateTransactionRequest
	if err := ctx.Bind(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}

	req.CreatedBy = userContext.UserID

	if err := ctx.Validate(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}

	err = h.service.CreateTransaction(ctx.Request().Context(), &req)
	if err != nil {
		logger.Errorw("finance transaction create failed: ", err)

		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusCreated, "Finance transaction created successfully", nil, nil, nil)
}

func (h *Handler) GetAllTransactions(ctx echo.Context) error {
	txType := ctx.QueryParam("type")
	status := ctx.QueryParam("status")
	startDate := ctx.QueryParam("start_date")
	endDate := ctx.QueryParam("end_date")
	page, _ := strconv.Atoi(ctx.QueryParam("page"))
	limit, _ := strconv.Atoi(ctx.QueryParam("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	filter := TransactionFilter{
		Type:      txType,
		Status:    status,
		StartDate: startDate,
		EndDate:   endDate,
		Page:      page,
		Limit:     limit,
	}

	data, meta, err := h.service.GetTransactions(ctx.Request().Context(), filter)
	if err != nil {
		logger.Errorw("get finance transactions failed: ", err)

		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Get Finance Transactions Success", data, nil, meta)
}

func (h *Handler) GetTransactionDetail(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "invalid id", nil, err, nil)
	}

	data, err := h.service.GetTransactionDetail(ctx.Request().Context(), uint(id))
	if err != nil {
		logger.Errorw("get finance transaction detail failed: ", err)

		return response.NewResponses[any](ctx, http.StatusBadRequest, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Get Finance Transaction Detail Success", data, nil, nil)
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
		logger.Errorw("process approval action finance failed: %w", err)

		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Process Approval Action Finance Success", nil, nil, nil)
}

func (h *Handler) ExportTransactions(ctx echo.Context) error {
	txType := ctx.QueryParam("type")
	status := ctx.QueryParam("status")
	startDate := ctx.QueryParam("start_date")
	endDate := ctx.QueryParam("end_date")

	filter := TransactionFilter{
		Type:      txType,
		Status:    status,
		StartDate: startDate,
		EndDate:   endDate,
	}

	excelFile, err := h.service.ExportTransactions(ctx.Request().Context(), filter)
	if err != nil {
		logger.Errorw("export finance transactions failed: ", err)
		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	ctx.Response().Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Response().Header().Set("Content-Disposition", "attachment; filename=finance_transactions.xlsx")
	return ctx.Blob(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", excelFile)
}

func (h *Handler) CreateCategory(ctx echo.Context) error {
	var req CategoryRequest
	if err := ctx.Bind(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}

	if err := ctx.Validate(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}

	err := h.service.CreateCategory(ctx.Request().Context(), &req)
	if err != nil {
		logger.Errorw("finance category create failed: ", err)

		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusCreated, "Finance category created successfully", nil, nil, nil)
}

func (h *Handler) GetCategories(ctx echo.Context) error {
	catType := ctx.QueryParam("type")

	data, err := h.service.GetCategories(ctx.Request().Context(), catType)
	if err != nil {
		logger.Errorw("get finance categories failed: ", err)

		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Get Finance Categories Success", data, nil, nil)
}

func (h *Handler) UpdateCategory(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "invalid id", nil, err, nil)
	}

	var req CategoryRequest
	if err := ctx.Bind(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}

	if err := ctx.Validate(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}

	err = h.service.UpdateCategory(ctx.Request().Context(), uint(id), &req)
	if err != nil {
		logger.Errorw("finance category update failed: ", err)

		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Finance category updated successfully", nil, nil, nil)
}

func (h *Handler) DeleteCategory(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "invalid id", nil, err, nil)
	}

	err = h.service.DeleteCategory(ctx.Request().Context(), uint(id))
	if err != nil {
		logger.Errorw("finance category delete failed: ", err)

		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Finance category deleted successfully", nil, nil, nil)
}

func (h *Handler) GetDashboard(ctx echo.Context) error {
	startDate := ctx.QueryParam("start_date")
	endDate := ctx.QueryParam("end_date")

	data, err := h.service.GetDashboard(ctx.Request().Context(), startDate, endDate)
	if err != nil {
		logger.Errorw("get finance dashboard failed: ", err)

		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Get Finance Dashboard Success", data, nil, nil)
}
