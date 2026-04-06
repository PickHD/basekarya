package contract

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

func (h *Handler) Upsert(ctx echo.Context) error {
	var req UpsertContractRequest
	if err := ctx.Bind(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}

	if err := ctx.Validate(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}

	err := h.service.Upsert(ctx.Request().Context(), &req)
	if err != nil {
		logger.Errorw("Upsert Contract failed: %w", err)
		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Upsert Contract Success", nil, nil, nil)
}

func (h *Handler) GetAll(ctx echo.Context) error {
	contractType := ctx.QueryParam("contract_type")
	page, _ := strconv.Atoi(ctx.QueryParam("page"))
	limit, _ := strconv.Atoi(ctx.QueryParam("limit"))
	search := ctx.QueryParam("search")
	expiringWithinDays, _ := strconv.Atoi(ctx.QueryParam("expiring_within_days"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	filter := ContractFilter{
		ContractType:       contractType,
		ExpiringWithinDays: expiringWithinDays,
		Page:               page,
		Limit:              limit,
		Search:             search,
	}

	data, meta, err := h.service.GetList(ctx.Request().Context(), &filter)
	if err != nil {
		logger.Errorw("get contracts failed: ", err)
		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Get Contract List Success", data, nil, meta)
}

func (h *Handler) GetDetail(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "invalid id", nil, err, nil)
	}

	data, err := h.service.GetDetail(ctx.Request().Context(), uint(id))
	if err != nil {
		logger.Errorw("get contract detail failed: ", err)
		return response.NewResponses[any](ctx, http.StatusBadRequest, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Get Contract Detail Success", data, nil, nil)
}

func (h *Handler) GetByEmployee(ctx echo.Context) error {
	empID, err := strconv.Atoi(ctx.Param("employeeId"))
	if err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "invalid employee id", nil, err, nil)
	}

	data, err := h.service.GetByEmployeeID(ctx.Request().Context(), uint(empID))
	if err != nil {
		logger.Errorw("get contract by employee id failed: ", err)
		return response.NewResponses[any](ctx, http.StatusBadRequest, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Get Contract By Employee ID Success", data, nil, nil)
}

func (h *Handler) Delete(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "invalid id", nil, err, nil)
	}

	err = h.service.Delete(ctx.Request().Context(), uint(id))
	if err != nil {
		logger.Errorw("delete contract failed: ", err)
		return response.NewResponses[any](ctx, http.StatusBadRequest, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Delete Contract Success", nil, nil, nil)
}

func (h *Handler) Export(ctx echo.Context) error {
	contractType := ctx.QueryParam("contract_type")
	search := ctx.QueryParam("search")
	expiringWithinDays, _ := strconv.Atoi(ctx.QueryParam("expiring_within_days"))

	filter := ContractFilter{
		ContractType:       contractType,
		ExpiringWithinDays: expiringWithinDays,
		Search:             search,
	}

	excelFile, err := h.service.Export(ctx.Request().Context(), &filter)
	if err != nil {
		logger.Errorw("export contracts failed: ", err)
		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	ctx.Response().Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Response().Header().Set("Content-Disposition", "attachment; filename=contracts.xlsx")
	return ctx.Blob(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", excelFile)
}
