package bpjs

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

func (h *Handler) List(ctx echo.Context) error {
	page, _ := strconv.Atoi(ctx.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(ctx.QueryParam("limit"))
	if limit < 1 || limit > 100 {
		limit = 100
	}

	filter := BPJSRateConfigFilter{
		Page:  page,
		Limit: limit,
	}

	data, total, err := h.service.List(ctx.Request().Context(), filter)
	if err != nil {
		logger.Errorw("list BPJS configs failed: ", err)
		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	meta := response.NewMetaOffset(page, limit, total)
	return response.NewResponses[any](ctx, http.StatusOK, "Get BPJS Configs Success", data, nil, meta)
}

func (h *Handler) Create(ctx echo.Context) error {
	var req BPJSRateConfigRequest
	if err := ctx.Bind(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}

	if err := ctx.Validate(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}

	err := h.service.Create(ctx.Request().Context(), &req)
	if err != nil {
		logger.Errorw("create BPJS config failed: ", err)
		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusCreated, "BPJS config created successfully", nil, nil, nil)
}

func (h *Handler) GetByID(ctx echo.Context) error {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "invalid id", nil, err, nil)
	}

	data, err := h.service.GetByID(ctx.Request().Context(), uint(id))
	if err != nil {
		logger.Errorw("get BPJS config failed: ", err)
		return response.NewResponses[any](ctx, http.StatusBadRequest, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Get BPJS Config Success", data, nil, nil)
}

func (h *Handler) Update(ctx echo.Context) error {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "invalid id", nil, err, nil)
	}

	var req BPJSRateConfigRequest
	if err := ctx.Bind(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}

	if err := ctx.Validate(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}

	err = h.service.Update(ctx.Request().Context(), uint(id), &req)
	if err != nil {
		logger.Errorw("update BPJS config failed: ", err)
		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "BPJS config updated successfully", nil, nil, nil)
}

func (h *Handler) Delete(ctx echo.Context) error {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "invalid id", nil, err, nil)
	}

	err = h.service.Delete(ctx.Request().Context(), uint(id))
	if err != nil {
		logger.Errorw("delete BPJS config failed: ", err)
		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "BPJS config deleted successfully", nil, nil, nil)
}
