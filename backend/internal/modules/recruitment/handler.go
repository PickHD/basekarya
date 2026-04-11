package recruitment

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

// ── Requisition Handlers ──────────────────────────────────────────────────────

func (h *Handler) CreateRequisition(ctx echo.Context) error {
	userCtx, err := utils.GetUserContext(ctx)
	if err != nil {
		return response.NewResponses[any](ctx, http.StatusUnauthorized, err.Error(), nil, err, nil)
	}

	var req CreateRequisitionRequest
	if err := ctx.Bind(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}
	if err := ctx.Validate(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}

	if err := h.service.CreateRequisition(ctx.Request().Context(), userCtx.UserID, &req); err != nil {
		logger.Errorw("CreateRequisition failed:", err)
		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusCreated, "Requisition created", nil, nil, nil)
}

func (h *Handler) SubmitRequisition(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "invalid id", nil, err, nil)
	}

	userCtx, err := utils.GetUserContext(ctx)
	if err != nil {
		return response.NewResponses[any](ctx, http.StatusUnauthorized, err.Error(), nil, err, nil)
	}

	if err := h.service.SubmitRequisition(ctx.Request().Context(), uint(id), userCtx.UserID); err != nil {
		logger.Errorw("SubmitRequisition failed:", err)
		return response.NewResponses[any](ctx, http.StatusBadRequest, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Requisition submitted for approval", nil, nil, nil)
}

func (h *Handler) RequisitionAction(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "invalid id", nil, err, nil)
	}

	userCtx, err := utils.GetUserContext(ctx)
	if err != nil {
		return response.NewResponses[any](ctx, http.StatusUnauthorized, err.Error(), nil, err, nil)
	}

	var req RequisitionActionRequest
	if err := ctx.Bind(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}
	if err := ctx.Validate(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}

	if err := h.service.RequisitionAction(ctx.Request().Context(), uint(id), userCtx.UserID, &req); err != nil {
		logger.Errorw("RequisitionAction failed:", err)
		return response.NewResponses[any](ctx, http.StatusBadRequest, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Requisition action processed", nil, nil, nil)
}

func (h *Handler) GetRequisitions(ctx echo.Context) error {
	page, _ := strconv.Atoi(ctx.QueryParam("page"))
	limit, _ := strconv.Atoi(ctx.QueryParam("limit"))
	departmentID, _ := strconv.Atoi(ctx.QueryParam("department_id"))

	filter := &RequisitionFilter{
		Status:       ctx.QueryParam("status"),
		Priority:     ctx.QueryParam("priority"),
		Search:       ctx.QueryParam("search"),
		DepartmentID: uint(departmentID),
		Page:         page,
		Limit:        limit,
	}

	data, meta, err := h.service.GetRequisitions(ctx.Request().Context(), filter)
	if err != nil {
		logger.Errorw("GetRequisitions failed:", err)
		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Get Requisitions Success", data, nil, meta)
}

func (h *Handler) GetRequisitionDetail(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "invalid id", nil, err, nil)
	}

	data, err := h.service.GetRequisitionDetail(ctx.Request().Context(), uint(id))
	if err != nil {
		logger.Errorw("GetRequisitionDetail failed:", err)
		return response.NewResponses[any](ctx, http.StatusNotFound, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Get Requisition Detail Success", data, nil, nil)
}

func (h *Handler) CloseRequisition(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "invalid id", nil, err, nil)
	}

	if err := h.service.CloseRequisition(ctx.Request().Context(), uint(id)); err != nil {
		logger.Errorw("CloseRequisition failed:", err)
		return response.NewResponses[any](ctx, http.StatusBadRequest, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Requisition closed", nil, nil, nil)
}

func (h *Handler) DeleteRequisition(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "invalid id", nil, err, nil)
	}

	if err := h.service.DeleteRequisition(ctx.Request().Context(), uint(id)); err != nil {
		logger.Errorw("DeleteRequisition failed:", err)
		return response.NewResponses[any](ctx, http.StatusBadRequest, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Requisition deleted", nil, nil, nil)
}

// ── Applicant Handlers ────────────────────────────────────────────────────────

func (h *Handler) AddApplicant(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "invalid requisition id", nil, err, nil)
	}

	var req CreateApplicantRequest
	if err := ctx.Bind(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}
	if err := ctx.Validate(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}

	if err := h.service.AddApplicant(ctx.Request().Context(), uint(id), &req); err != nil {
		logger.Errorw("AddApplicant failed:", err)
		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusCreated, "Applicant added", nil, nil, nil)
}

func (h *Handler) UpdateApplicantStage(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "invalid applicant id", nil, err, nil)
	}

	userCtx, err := utils.GetUserContext(ctx)
	if err != nil {
		return response.NewResponses[any](ctx, http.StatusUnauthorized, err.Error(), nil, err, nil)
	}

	var req UpdateApplicantStageRequest
	if err := ctx.Bind(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}
	if err := ctx.Validate(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}

	if err := h.service.UpdateStage(ctx.Request().Context(), uint(id), userCtx.UserID, &req); err != nil {
		logger.Errorw("UpdateApplicantStage failed:", err)
		return response.NewResponses[any](ctx, http.StatusBadRequest, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Applicant stage updated", nil, nil, nil)
}

func (h *Handler) GetApplicants(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "invalid requisition id", nil, err, nil)
	}

	data, err := h.service.GetApplicantsByRequisition(ctx.Request().Context(), uint(id))
	if err != nil {
		logger.Errorw("GetApplicants failed:", err)
		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Get Applicants Success", data, nil, nil)
}

func (h *Handler) GetApplicantDetail(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "invalid applicant id", nil, err, nil)
	}

	data, err := h.service.GetApplicantDetail(ctx.Request().Context(), uint(id))
	if err != nil {
		logger.Errorw("GetApplicantDetail failed:", err)
		return response.NewResponses[any](ctx, http.StatusNotFound, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Get Applicant Detail Success", data, nil, nil)
}
