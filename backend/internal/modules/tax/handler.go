package tax

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

func (h *Handler) ListTERBrackets(ctx echo.Context) error {
	category := ctx.QueryParam("category")
	page, _ := strconv.Atoi(ctx.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(ctx.QueryParam("limit"))
	if limit < 1 || limit > 100 {
		limit = 10
	}

	filter := TERBracketFilter{
		Category: category,
		Page:     page,
		Limit:    limit,
	}

	data, total, err := h.service.ListTERBrackets(ctx.Request().Context(), filter)
	if err != nil {
		logger.Errorw("list TER brackets failed: ", err)
		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	meta := response.NewMetaOffset(page, limit, total)
	return response.NewResponses[any](ctx, http.StatusOK, "Get TER Brackets Success", data, nil, meta)
}

func (h *Handler) CreateTERBracket(ctx echo.Context) error {
	var req TERBracketRequest
	if err := ctx.Bind(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}

	if err := ctx.Validate(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}

	err := h.service.CreateTERBracket(ctx.Request().Context(), &req)
	if err != nil {
		logger.Errorw("create TER bracket failed: ", err)
		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusCreated, "TER bracket created successfully", nil, nil, nil)
}

func (h *Handler) GetTERBracketByID(ctx echo.Context) error {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "invalid id", nil, err, nil)
	}

	data, err := h.service.GetTERBracketByID(ctx.Request().Context(), uint(id))
	if err != nil {
		logger.Errorw("get TER bracket failed: ", err)
		return response.NewResponses[any](ctx, http.StatusBadRequest, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Get TER Bracket Success", data, nil, nil)
}

func (h *Handler) UpdateTERBracket(ctx echo.Context) error {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "invalid id", nil, err, nil)
	}

	var req TERBracketRequest
	if err := ctx.Bind(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}

	if err := ctx.Validate(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}

	err = h.service.UpdateTERBracket(ctx.Request().Context(), uint(id), &req)
	if err != nil {
		logger.Errorw("update TER bracket failed: ", err)
		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "TER bracket updated successfully", nil, nil, nil)
}

func (h *Handler) DeleteTERBracket(ctx echo.Context) error {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "invalid id", nil, err, nil)
	}

	err = h.service.DeleteTERBracket(ctx.Request().Context(), uint(id))
	if err != nil {
		logger.Errorw("delete TER bracket failed: ", err)
		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "TER bracket deleted successfully", nil, nil, nil)
}

func (h *Handler) ListPTKPConfigs(ctx echo.Context) error {
	year, _ := strconv.Atoi(ctx.QueryParam("year"))
	if year < 1 {
		year = 2026
	}

	data, total, err := h.service.ListPTKPConfigs(ctx.Request().Context(), year)
	if err != nil {
		logger.Errorw("list PTKP configs failed: ", err)
		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	meta := response.NewMetaOffset(1, int(total), total)
	return response.NewResponses[any](ctx, http.StatusOK, "Get PTKP Configs Success", data, nil, meta)
}

func (h *Handler) CreatePTKPConfig(ctx echo.Context) error {
	var req PTKPConfigRequest
	if err := ctx.Bind(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}

	if err := ctx.Validate(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}

	err := h.service.CreatePTKPConfig(ctx.Request().Context(), &req)
	if err != nil {
		logger.Errorw("create PTKP config failed: ", err)
		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusCreated, "PTKP config created successfully", nil, nil, nil)
}

func (h *Handler) UpdatePTKPConfig(ctx echo.Context) error {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "invalid id", nil, err, nil)
	}

	var req PTKPConfigRequest
	if err := ctx.Bind(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}

	if err := ctx.Validate(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}

	err = h.service.UpdatePTKPConfig(ctx.Request().Context(), uint(id), &req)
	if err != nil {
		logger.Errorw("update PTKP config failed: ", err)
		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "PTKP config updated successfully", nil, nil, nil)
}

func (h *Handler) DeletePTKPConfig(ctx echo.Context) error {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "invalid id", nil, err, nil)
	}

	err = h.service.DeletePTKPConfig(ctx.Request().Context(), uint(id))
	if err != nil {
		logger.Errorw("delete PTKP config failed: ", err)
		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "PTKP config deleted successfully", nil, nil, nil)
}
