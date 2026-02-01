package payroll

import (
	"fmt"
	"hris-backend/pkg/logger"
	"hris-backend/pkg/response"
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

	// data, total, err := h.service.GetList(ctx.Request().Context(), page, limit, month, year, keyword)
	// if err != nil {
	// 	response.Error(c, http.StatusInternalServerError, "Failed to fetch payroll list", err)
	// 	return
	// }
}
