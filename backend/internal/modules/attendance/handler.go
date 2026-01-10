package attendance

import (
	"fmt"
	"hris-backend/pkg/logger"
	"hris-backend/pkg/response"
	"hris-backend/pkg/utils"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

func (h *Handler) Clock(ctx echo.Context) error {
	userContext, err := utils.GetUserContext(ctx)
	if err != nil {
		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	var req ClockRequest
	if err := ctx.Bind(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}

	resp, err := h.service.Clock(ctx.Request().Context(), userContext.UserID, &req)
	if err != nil {
		logger.Errorw("Clock Request failed : ", err)

		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, resp.Message, resp, nil, nil)
}

func (h *Handler) GetTodayStatus(ctx echo.Context) error {
	userContext, err := utils.GetUserContext(ctx)
	if err != nil {
		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	resp, err := h.service.GetTodayStatus(ctx.Request().Context(), userContext.UserID)
	if err != nil {
		logger.Errorw("Get Today Status failed: ", err)

		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Get Today Status Success", resp, nil, nil)
}

func (h *Handler) GetHistory(ctx echo.Context) error {
	userContext, err := utils.GetUserContext(ctx)
	if err != nil {
		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	month := int(time.Now().Month())
	year := time.Now().Year()
	page := 1
	limit := 10

	if m := ctx.QueryParam("month"); m != "" {
		fmt.Sscanf(m, "%d", &month)
	}
	if y := ctx.QueryParam("year"); y != "" {
		fmt.Sscanf(y, "%d", &year)
	}
	if p := ctx.QueryParam("page"); p != "" {
		fmt.Sscanf(p, "%d", &page)
	}
	if l := ctx.QueryParam("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}

	resp, meta, err := h.service.GetMyHistory(ctx.Request().Context(), userContext.UserID, month, year, page, limit)
	if err != nil {
		logger.Errorw("Get My History failed: ", err)

		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Get My History Success", resp, nil, meta)
}

func (h *Handler) GetAllAttendanceRecap(ctx echo.Context) error {
	filter := h.parseFilter(ctx)

	resp, meta, err := h.service.GetAllRecap(ctx.Request().Context(), filter)
	if err != nil {
		logger.Errorw("Get All Attendance Recap Failed: ", err)

		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Get All Attendance Recap Success", resp, nil, meta)
}

func (h *Handler) ExportAttendance(ctx echo.Context) error {
	filter := h.parseFilter(ctx)

	f, err := h.service.GenerateExcel(ctx.Request().Context(), filter)
	if err != nil {
		logger.Errorw("Generate Excel Attendance Failed: ", err)

		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	filename := fmt.Sprintf("Attendance_Recap_%s.xlsx", time.Now().Format("20060102_150405"))
	ctx.Response().Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Response().Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	if err := f.Write(ctx.Response().Writer); err != nil {
		return err
	}
	return nil
}

func (h *Handler) GetDashboardStats(ctx echo.Context) error {
	tz := ctx.QueryParam("timezone")
	if tz == "" {
		tz = "Asia/Jakarta"
	}

	resp, err := h.service.GetDashboardStats(ctx.Request().Context(), tz)
	if err != nil {
		logger.Errorw("Get Dashboard Stats failed: ", err)

		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Get Dashboard Stats Success", resp, nil, nil)
}

func (h *Handler) parseFilter(ctx echo.Context) *FilterParams {
	page := 1
	limit := 10
	if p := ctx.QueryParam("page"); p != "" {
		fmt.Sscanf(p, "%d", &page)
	}
	if l := ctx.QueryParam("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}

	deptID := 0
	if d := ctx.QueryParam("department_id"); d != "" {
		fmt.Sscanf(d, "%d", &deptID)
	}

	tz := ctx.QueryParam("timezone")
	if tz == "" {
		tz = "Asia/Jakarta"
	}

	return &FilterParams{
		Page:         page,
		Limit:        limit,
		StartDate:    ctx.QueryParam("start_date"),
		EndDate:      ctx.QueryParam("end_date"),
		Search:       ctx.QueryParam("search"),
		Timezone:     tz,
		DepartmentID: uint(deptID),
	}
}
