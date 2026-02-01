package reimbursement

import (
	"fmt"
	"hris-backend/pkg/constants"
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

func (h *Handler) Create(ctx echo.Context) error {
	userContext, err := utils.GetUserContext(ctx)
	if err != nil {
		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	req, err := h.parseAndValidateFormData(ctx, userContext.UserID)
	if err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, err.Error(), nil, err, nil)
	}

	err = h.service.Create(ctx.Request().Context(), req)
	if err != nil {
		logger.Errorw("reimburstment create failed: ", err)

		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusCreated, "Reimburstment created successfully", nil, nil, nil)
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

	filter := ReimbursementFilter{
		Status: status,
		Page:   page,
		Limit:  limit,
	}

	if userContext.Role != string(constants.UserRoleSuperadmin) {
		filter.UserID = userContext.UserID
	}

	data, meta, err := h.service.GetReimbursements(ctx.Request().Context(), filter)
	if err != nil {
		logger.Errorw("get reimbursements failed: ", err)

		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Get Reimbursements List Success", data, nil, meta)
}

func (h *Handler) GetDetail(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "invalid id", nil, err, nil)
	}

	data, err := h.service.GetReimburseDetail(ctx.Request().Context(), uint(id))
	if err != nil {
		logger.Errorw("get reimburse detail failed: ", err)

		return response.NewResponses[any](ctx, http.StatusBadRequest, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Get Reimbursement Detail Success", data, nil, nil)
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
		logger.Errorw("Process approval action reimbursement failed: %w", err)

		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Process Approval Action Reimbursement Success", nil, nil, nil)
}

func (h *Handler) parseAndValidateFormData(ctx echo.Context, userID uint) (*ReimbursementRequest, error) {
	amount, err := strconv.ParseFloat(ctx.FormValue("amount"), 64)
	if err != nil {
		return nil, fmt.Errorf("invalid amount")
	}

	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		return nil, fmt.Errorf("file required")
	}

	if fileHeader.Size > 5*1024*1024 {
		return nil, fmt.Errorf("File size exceeds 5MB limit")
	}

	return &ReimbursementRequest{
		UserID:      userID,
		Title:       ctx.FormValue("title"),
		Description: ctx.FormValue("description"),
		Date:        ctx.FormValue("date"),
		Amount:      amount,
		File:        fileHeader,
	}, nil
}
