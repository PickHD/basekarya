package onboarding

import (
	"net/http"
	"strconv"

	"basekarya-backend/pkg/logger"
	"basekarya-backend/pkg/response"
	"basekarya-backend/pkg/utils"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

// ── Template Handlers ─────────────────────────────────────────────────────────

func (h *Handler) CreateTemplate(ctx echo.Context) error {
	var req CreateTemplateRequest
	if err := ctx.Bind(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}
	if err := ctx.Validate(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}

	if err := h.service.CreateTemplate(ctx.Request().Context(), &req); err != nil {
		logger.Errorw("CreateTemplate failed:", err)
		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusCreated, "Template created", nil, nil, nil)
}

func (h *Handler) GetTemplates(ctx echo.Context) error {
	data, err := h.service.GetTemplates(ctx.Request().Context())
	if err != nil {
		logger.Errorw("GetTemplates failed:", err)
		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Get Templates Success", data, nil, nil)
}

func (h *Handler) UpdateTemplate(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "invalid id", nil, err, nil)
	}

	var req UpdateTemplateRequest
	if err := ctx.Bind(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}
	if err := ctx.Validate(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}

	if err := h.service.UpdateTemplate(ctx.Request().Context(), uint(id), &req); err != nil {
		logger.Errorw("UpdateTemplate failed:", err)
		return response.NewResponses[any](ctx, http.StatusBadRequest, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Template updated", nil, nil, nil)
}

func (h *Handler) DeleteTemplate(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "invalid id", nil, err, nil)
	}

	if err := h.service.DeleteTemplate(ctx.Request().Context(), uint(id)); err != nil {
		logger.Errorw("DeleteTemplate failed:", err)
		return response.NewResponses[any](ctx, http.StatusBadRequest, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Template deleted", nil, nil, nil)
}

// ── Workflow Handlers ─────────────────────────────────────────────────────────

func (h *Handler) CreateWorkflow(ctx echo.Context) error {
	var req CreateWorkflowRequest
	if err := ctx.Bind(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}
	if err := ctx.Validate(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}

	err := h.service.CreateWorkflow(ctx.Request().Context(), &req)
	if err != nil {
		logger.Errorw("CreateWorkflow failed:", err)
		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusCreated, "Onboarding workflow created", nil, nil, nil)
}

func (h *Handler) GetWorkflows(ctx echo.Context) error {
	page, _ := strconv.Atoi(ctx.QueryParam("page"))
	limit, _ := strconv.Atoi(ctx.QueryParam("limit"))

	filter := &WorkflowFilter{
		Status: ctx.QueryParam("status"),
		Search: ctx.QueryParam("search"),
		Page:   page,
		Limit:  limit,
	}

	data, meta, err := h.service.GetWorkflows(ctx.Request().Context(), filter)
	if err != nil {
		logger.Errorw("GetWorkflows failed:", err)
		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Get Workflows Success", data, nil, meta)
}

func (h *Handler) GetWorkflowDetail(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "invalid id", nil, err, nil)
	}

	data, err := h.service.GetWorkflowDetail(ctx.Request().Context(), uint(id))
	if err != nil {
		logger.Errorw("GetWorkflowDetail failed:", err)
		return response.NewResponses[any](ctx, http.StatusNotFound, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Get Workflow Detail Success", data, nil, nil)
}

// ── Task Handlers ─────────────────────────────────────────────────────────────

func (h *Handler) CompleteTask(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "invalid task id", nil, err, nil)
	}

	userCtx, err := utils.GetUserContext(ctx)
	if err != nil {
		return response.NewResponses[any](ctx, http.StatusUnauthorized, err.Error(), nil, err, nil)
	}

	var req CompleteTaskRequest
	if err := ctx.Bind(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}

	if err := h.service.CompleteTask(ctx.Request().Context(), uint(id), userCtx.UserID, &req); err != nil {
		logger.Errorw("CompleteTask failed:", err)
		return response.NewResponses[any](ctx, http.StatusBadRequest, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Task completed", nil, nil, nil)
}
