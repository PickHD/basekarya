package payroll

import (
	"fmt"
	"hris-backend/pkg/logger"
	"hris-backend/pkg/response"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

func (h *Handler) Generate(ctx echo.Context) error {
	var req GenerateRequest
	if err := ctx.Bind(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}

	resp, err := h.service.GenerateAll(ctx.Request().Context(), &req)
	if err != nil {
		logger.Errorw("Generate All Payroll Employees failed: %w", err)

		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Generate All Payroll Employees Successfully", resp, nil, nil)
}

func (h *Handler) GetList(ctx echo.Context) error {
	filter := h.parseFilter(ctx)

	data, meta, err := h.service.GetList(ctx.Request().Context(), filter)
	if err != nil {
		logger.Errorw("Failed to fetch payroll list: %w", err)

		return response.NewResponses[any](ctx, http.StatusInternalServerError, "Failed to fetch payroll list", nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Fetch Payroll List Success", data, nil, meta)
}

func (h *Handler) GetDetail(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "invalid id", nil, err, nil)
	}

	data, err := h.service.GetDetail(ctx.Request().Context(), uint(id))
	if err != nil {
		logger.Errorw("Failed to fetch detail payroll: %w", err)

		return response.NewResponses[any](ctx, http.StatusInternalServerError, "Failed to fetch payroll detail", nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Fetch Payroll Detail Success", data, nil, nil)
}

func (h *Handler) DownloadPayslipPDF(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "invalid id", nil, err, nil)
	}

	pdf, data, err := h.service.GeneratePayslipPDF(ctx.Request().Context(), uint(id))
	if err != nil {
		logger.Errorw("Failed to fetch detail payroll: %w", err)

		return response.NewResponses[any](ctx, http.StatusInternalServerError, "Failed to fetch payroll detail", nil, err, nil)
	}

	filename := fmt.Sprintf("Payslip-%s-%s.pdf", data.Employee.NIK, data.PeriodDate.Format("Jan2006"))
	ctx.Response().Header().Set("Content-Type", "application/pdf")
	ctx.Response().Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	err = pdf.Output(ctx.Response().Writer)
	if err != nil {
		return err
	}
	return nil
}

func (h *Handler) MarkAsPaid(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "invalid id", nil, err, nil)
	}

	err = h.service.MarkAsPaid(ctx.Request().Context(), uint(id))
	if err != nil {
		logger.Errorw("Failed to mark as paid payroll: %w", err)

		return response.NewResponses[any](ctx, http.StatusInternalServerError, "Failed to mark as paid payroll", nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Mark As Paid Payroll Success", nil, nil, nil)
}

func (h *Handler) parseFilter(ctx echo.Context) *PayrollFilter {
	month := int(time.Now().Month())
	year := time.Now().Year()
	page := 1
	limit := 10
	search := ""

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
	if s := ctx.QueryParam("search"); s != "" {
		fmt.Sscanf(s, "%s", &search)
	}

	return &PayrollFilter{
		Page:    page,
		Limit:   limit,
		Month:   month,
		Year:    year,
		Keyword: search,
	}
}
