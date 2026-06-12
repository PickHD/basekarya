package asset

import (
	"basekarya-backend/pkg/constants"
	"basekarya-backend/pkg/logger"
	"basekarya-backend/pkg/response"
	"basekarya-backend/pkg/utils"
	"net/http"
	"slices"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

func (h *Handler) CreateCategory(ctx echo.Context) error {
	var req CreateAssetCategoryRequest
	if err := ctx.Bind(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}

	if err := ctx.Validate(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}

	err := h.service.CreateCategory(ctx.Request().Context(), &req)
	if err != nil {
		logger.Errorw("create asset category failed: ", err)
		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusCreated, "Asset category created successfully", nil, nil, nil)
}

func (h *Handler) GetAllCategories(ctx echo.Context) error {
	page, _ := strconv.Atoi(ctx.QueryParam("page"))
	limit, _ := strconv.Atoi(ctx.QueryParam("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	filter := AssetCategoryFilter{
		Page:  page,
		Limit: limit,
	}

	data, meta, err := h.service.GetCategories(ctx.Request().Context(), filter)
	if err != nil {
		logger.Errorw("get asset categories failed: ", err)
		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Get Asset Categories Success", data, nil, meta)
}

func (h *Handler) GetCategoryDetail(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "invalid id", nil, err, nil)
	}

	data, err := h.service.GetCategoryDetail(ctx.Request().Context(), uint(id))
	if err != nil {
		logger.Errorw("get asset category detail failed: ", err)
		return response.NewResponses[any](ctx, http.StatusBadRequest, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Get Asset Category Detail Success", data, nil, nil)
}

func (h *Handler) UpdateCategory(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "invalid id", nil, err, nil)
	}

	var req UpdateAssetCategoryRequest
	if err := ctx.Bind(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}

	req.ID = uint(id)

	if err := ctx.Validate(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}

	err = h.service.UpdateCategory(ctx.Request().Context(), &req)
	if err != nil {
		logger.Errorw("update asset category failed: ", err)
		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Asset category updated successfully", nil, nil, nil)
}

func (h *Handler) DeleteCategory(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "invalid id", nil, err, nil)
	}

	err = h.service.DeleteCategory(ctx.Request().Context(), uint(id))
	if err != nil {
		logger.Errorw("delete asset category failed: ", err)
		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Asset category deleted successfully", nil, nil, nil)
}

func (h *Handler) CreateAsset(ctx echo.Context) error {
	var req CreateAssetRequest
	if err := ctx.Bind(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}

	if err := ctx.Validate(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}

	err := h.service.CreateAsset(ctx.Request().Context(), &req)
	if err != nil {
		logger.Errorw("create asset failed: ", err)
		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusCreated, "Asset created successfully", nil, nil, nil)
}

func (h *Handler) GetAllAssets(ctx echo.Context) error {
	status := ctx.QueryParam("status")
	condition := ctx.QueryParam("condition")
	categoryID, _ := strconv.Atoi(ctx.QueryParam("category_id"))
	page, _ := strconv.Atoi(ctx.QueryParam("page"))
	limit, _ := strconv.Atoi(ctx.QueryParam("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	filter := AssetFilter{
		Status:     status,
		Condition:  condition,
		CategoryID: uint(categoryID),
		Page:       page,
		Limit:      limit,
	}

	data, meta, err := h.service.GetAssets(ctx.Request().Context(), filter)
	if err != nil {
		logger.Errorw("get assets failed: ", err)
		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Get Assets Success", data, nil, meta)
}

func (h *Handler) GetAssetDetail(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "invalid id", nil, err, nil)
	}

	data, err := h.service.GetAssetDetail(ctx.Request().Context(), uint(id))
	if err != nil {
		logger.Errorw("get asset detail failed: ", err)
		return response.NewResponses[any](ctx, http.StatusBadRequest, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Get Asset Detail Success", data, nil, nil)
}

func (h *Handler) UpdateAsset(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "invalid id", nil, err, nil)
	}

	var req UpdateAssetRequest
	if err := ctx.Bind(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}

	req.ID = uint(id)

	err = h.service.UpdateAsset(ctx.Request().Context(), &req)
	if err != nil {
		logger.Errorw("update asset failed: ", err)
		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Asset updated successfully", nil, nil, nil)
}

func (h *Handler) DeleteAsset(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "invalid id", nil, err, nil)
	}

	err = h.service.DeleteAsset(ctx.Request().Context(), uint(id))
	if err != nil {
		logger.Errorw("delete asset failed: ", err)
		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Asset deleted successfully", nil, nil, nil)
}

func (h *Handler) CreateAssignment(ctx echo.Context) error {
	userContext, err := utils.GetUserContext(ctx)
	if err != nil {
		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	var req CreateAssetAssignmentRequest
	if err := ctx.Bind(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}

	req.UserID = userContext.UserID
	req.EmployeeID = *userContext.EmployeeID

	if err := ctx.Validate(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}

	err = h.service.CreateAssignment(ctx.Request().Context(), &req)
	if err != nil {
		logger.Errorw("create asset assignment failed: ", err)
		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusCreated, "Asset assignment created successfully", nil, nil, nil)
}

func (h *Handler) GetAllAssignments(ctx echo.Context) error {
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

	filter := AssetAssignmentFilter{
		Status: status,
		Page:   page,
		Limit:  limit,
	}

	if !slices.Contains(userContext.Permissions, constants.VIEW_ASSET) && slices.Contains(userContext.Permissions, constants.VIEW_SELF_ASSET) {
		filter.UserID = userContext.UserID
	}

	data, meta, err := h.service.GetAssignments(ctx.Request().Context(), filter)
	if err != nil {
		logger.Errorw("get asset assignments failed: ", err)
		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Get Asset Assignments Success", data, nil, meta)
}

func (h *Handler) GetAssignmentDetail(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "invalid id", nil, err, nil)
	}

	data, err := h.service.GetAssignmentDetail(ctx.Request().Context(), uint(id))
	if err != nil {
		logger.Errorw("get asset assignment detail failed: ", err)
		return response.NewResponses[any](ctx, http.StatusBadRequest, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Get Asset Assignment Detail Success", data, nil, nil)
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
		logger.Errorw("Process approval action asset assignment failed: %w", err)
		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Process Approval Action Asset Assignment Success", nil, nil, nil)
}

func (h *Handler) ProcessReturn(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "invalid id", nil, err, nil)
	}

	userContext, err := utils.GetUserContext(ctx)
	if err != nil {
		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	req := ReturnRequest{
		ID:      uint(id),
		UserID:   userContext.UserID,
	}

	err = h.service.ProcessReturn(ctx.Request().Context(), &req)
	if err != nil {
		logger.Errorw("Process return asset assignment failed: %w", err)
		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Asset returned successfully", nil, nil, nil)
}

func (h *Handler) ExportAssets(ctx echo.Context) error {
	status := ctx.QueryParam("status")
	condition := ctx.QueryParam("condition")
	categoryID, _ := strconv.Atoi(ctx.QueryParam("category_id"))

	filter := AssetFilter{
		Status:     status,
		Condition:  condition,
		CategoryID: uint(categoryID),
	}

	excelFile, err := h.service.Export(ctx.Request().Context(), filter)
	if err != nil {
		logger.Errorw("export assets failed: ", err)
		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	ctx.Response().Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Response().Header().Set("Content-Disposition", "attachment; filename=assets.xlsx")
	return ctx.Blob(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", excelFile)
}
