package payroll

import (
	"context"
	"errors"
	"fmt"
	"hris-backend/pkg/constants"
	"hris-backend/pkg/logger"
	"hris-backend/pkg/response"
	"hris-backend/pkg/utils"
	"time"

	"github.com/go-pdf/fpdf"
)

type Service interface {
	GenerateAll(ctx context.Context, req *GenerateRequest) (*GenerateResponse, error)
	GetList(ctx context.Context, filter *PayrollFilter) ([]PayrollListResponse, *response.Meta, error)
	GetDetail(ctx context.Context, id uint) (*PayrollDetailResponse, error)
	GeneratePayslipPDF(ctx context.Context, id uint) (*fpdf.Fpdf, *Payroll, error)
	MarkAsPaid(ctx context.Context, id uint) error
}

type service struct {
	repo          Repository
	user          UserProvider
	reimbursement ReimbursementProvider
	attendance    AttendanceProvider
}

func NewService(repo Repository,
	user UserProvider,
	reimbursement ReimbursementProvider,
	attendance AttendanceProvider) Service {
	return &service{repo, user, reimbursement, attendance}
}

func (s *service) GenerateAll(ctx context.Context, req *GenerateRequest) (*GenerateResponse, error) {
	employees, err := s.user.FindAllEmployeeActive()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch all employee active: %w", err)
	}

	existingPayrollMap, err := s.repo.GetExistingEmployeeID(req.Month, req.Year)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch existing employee id: %w", err)
	}

	attendanceMap, err := s.attendance.GetBulkLateDuration(req.Month, req.Year)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch bulk late duration: %w", err)
	}

	reimburseMap, err := s.reimbursement.GetBulkApprovedAmount(req.Month, req.Year)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch bulk approved amount: %w", err)
	}

	successCount := 0
	periodDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.Local)

	var payrollsToInsert []Payroll

	for _, emp := range employees {
		// if already exist on this year & month, skip
		if existingPayrollMap[emp.ID] {
			continue
		}

		// take data with O(1) lookup
		baseSalary := emp.BaseSalary
		totalLateMinutes := attendanceMap[emp.ID]
		reimburseAmount := reimburseMap[emp.UserID]

		// calculate net salary
		latePenaltyAmount := float64(totalLateMinutes * constants.PenaltyPerMinuteLate)
		totalAllowance := reimburseAmount
		totalDeduction := latePenaltyAmount
		netSalary := baseSalary + totalAllowance - totalDeduction

		// construct object
		payroll := Payroll{
			EmployeeID:     emp.ID,
			PeriodDate:     periodDate,
			BaseSalary:     baseSalary,
			TotalAllowance: totalAllowance,
			TotalDeduction: totalDeduction,
			NetSalary:      netSalary,
			Status:         constants.PayrollStatusDraft,
			Details:        []PayrollDetail{},
		}

		payroll.Details = append(payroll.Details, PayrollDetail{
			Title:  "Base Salary",
			Type:   constants.DetailTypeAllowance,
			Amount: baseSalary,
		})

		// check if reimburse amount not zero
		if reimburseAmount > 0 {
			payroll.Details = append(payroll.Details, PayrollDetail{
				Title:  "Reimbursement Approved",
				Type:   constants.DetailTypeAllowance,
				Amount: reimburseAmount,
			})
		}

		// check if late penalty amount not zero
		if latePenaltyAmount > 0 {
			payroll.Details = append(payroll.Details, PayrollDetail{
				Title:  fmt.Sprintf("Potongan Terlambat (%d menit)", totalLateMinutes),
				Type:   constants.DetailTypeDeduction,
				Amount: latePenaltyAmount,
			})
		}

		// insert to slice & update success count
		payrollsToInsert = append(payrollsToInsert, payroll)
		successCount++
	}

	// check if payrollsToInsert empty, return 0
	if len(payrollsToInsert) == 0 {
		return nil, nil
	}

	// bulk insert payrolls
	if err := s.repo.CreateBulk(&payrollsToInsert); err != nil {
		logger.Errorf("Failed create bulk payrolls %w", err)

		successCount = 0
		return nil, err
	}

	return &GenerateResponse{
		SuccessCount: successCount,
		Year:         req.Year,
		Month:        req.Month,
	}, nil
}

func (s *service) GetList(ctx context.Context, filter *PayrollFilter) ([]PayrollListResponse, *response.Meta, error) {
	data, total, err := s.repo.FindAll(filter)
	if err != nil {
		return nil, nil, err
	}

	if len(data) == 0 {
		return []PayrollListResponse{}, nil, nil
	}

	var responses []PayrollListResponse
	for _, p := range data {
		empName := "Unknown"
		empNIK := "-"
		if p.Employee != nil {
			empName = p.Employee.FullName
			empNIK = p.Employee.NIK
		}

		responses = append(responses, PayrollListResponse{
			ID:           p.ID,
			EmployeeName: empName,
			EmployeeNIK:  empNIK,
			PeriodDate:   p.PeriodDate.Format("2006-01-02"),
			NetSalary:    p.NetSalary,
			Status:       string(p.Status),
			CreatedAt:    p.CreatedAt,
		})
	}

	meta := response.NewMetaOffset(filter.Page, filter.Limit, total)
	return responses, meta, nil
}

func (s *service) GetDetail(ctx context.Context, id uint) (*PayrollDetailResponse, error) {
	payroll, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if payroll.Employee == nil {
		return nil, errors.New("employee not found")
	}

	emp := payroll.Employee

	details := make([]Detail, len(payroll.Details))
	for _, detail := range payroll.Details {
		details = append(details, Detail{
			ID:        detail.ID,
			PayrollID: detail.PayrollID,
			Title:     detail.Title,
			Type:      detail.Type,
			Amount:    detail.Amount,
		})
	}

	payrollDetail := PayrollDetailResponse{
		ID:             payroll.ID,
		EmployeeID:     emp.ID,
		EmployeeName:   emp.FullName,
		EmployeeNIK:    emp.NIK,
		PeriodDate:     payroll.PeriodDate.Format("2006-01-02"),
		BaseSalary:     payroll.BaseSalary,
		TotalAllowance: payroll.TotalAllowance,
		TotalDeduction: payroll.TotalDeduction,
		NetSalary:      payroll.NetSalary,
		Status:         string(payroll.Status),
		CreatedAt:      payroll.CreatedAt,
		Details:        details,
	}

	return &payrollDetail, nil
}

func (s *service) GeneratePayslipPDF(ctx context.Context, id uint) (*fpdf.Fpdf, *Payroll, error) {
	payroll, err := s.repo.FindByID(id)
	if err != nil {
		return nil, nil, err
	}

	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "", 12)

	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(0, 10, "PT. YOUR COMPANY", "", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(0, 5, "Jalan Sudirman No. 123, Jakarta Selatan", "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 5, "Telp: (021) 555-1234 | Email: hr@company.com", "", 1, "C", false, 0, "")

	pdf.SetLineWidth(0.5)
	pdf.Line(10, 30, 200, 30)

	pdf.Ln(10)

	printInfo := func(label, value string, label2, value2 string) {
		pdf.SetFont("Arial", "B", 10)
		pdf.CellFormat(30, 6, label, "", 0, "", false, 0, "")
		pdf.SetFont("Arial", "", 10)
		pdf.CellFormat(60, 6, ": "+value, "", 0, "", false, 0, "")

		if label2 != "" {
			pdf.SetFont("Arial", "B", 10)
			pdf.CellFormat(30, 6, label2, "", 0, "", false, 0, "")
			pdf.SetFont("Arial", "", 10)
			pdf.CellFormat(60, 6, ": "+value2, "", 1, "", false, 0, "")
		} else {
			pdf.Ln(-1)
		}
	}

	periodStr := payroll.PeriodDate.Format("January 2006")

	printInfo("Name", payroll.Employee.FullName, "Period", periodStr)
	printInfo("NIK", payroll.Employee.NIK, "Status", string(payroll.Status))

	pdf.Ln(10)

	pdf.SetFillColor(240, 240, 240)
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(95, 8, " EARNINGS (PENDAPATAN)", "1", 0, "L", true, 0, "")
	pdf.CellFormat(95, 8, " DEDUCTIONS (POTONGAN)", "1", 1, "L", true, 0, "")

	var earnings, deductions []PayrollDetail
	for _, d := range payroll.Details {
		if d.Type == constants.DetailTypeAllowance {
			earnings = append(earnings, d)
		} else {
			deductions = append(deductions, d)
		}
	}

	maxRows := len(earnings)
	if len(deductions) > maxRows {
		maxRows = len(deductions)
	}

	pdf.SetFont("Arial", "", 9)
	formatCurrency := func(amount float64) string {
		return fmt.Sprintf("Rp %s", utils.FormatNumber(amount))
	}

	for i := 0; i < maxRows; i++ {
		if i < len(earnings) {
			pdf.CellFormat(60, 7, " "+earnings[i].Title, "L", 0, "L", false, 0, "")
			pdf.CellFormat(35, 7, formatCurrency(earnings[i].Amount)+" ", "R", 0, "R", false, 0, "")
		} else {
			pdf.CellFormat(60, 7, "", "L", 0, "L", false, 0, "")
			pdf.CellFormat(35, 7, "", "R", 0, "R", false, 0, "")
		}

		if i < len(deductions) {
			pdf.CellFormat(60, 7, " "+deductions[i].Title, "L", 0, "L", false, 0, "")
			pdf.CellFormat(35, 7, formatCurrency(deductions[i].Amount)+" ", "R", 1, "R", false, 0, "")
		} else {
			pdf.CellFormat(60, 7, "", "L", 0, "L", false, 0, "")
			pdf.CellFormat(35, 7, "", "R", 1, "R", false, 0, "")
		}
	}

	pdf.CellFormat(95, 0, "", "T", 0, "", false, 0, "")
	pdf.CellFormat(95, 0, "", "T", 1, "", false, 0, "")

	pdf.Ln(2)
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(60, 8, " Total Earnings", "", 0, "L", false, 0, "")
	pdf.CellFormat(35, 8, formatCurrency(payroll.TotalAllowance)+" ", "", 0, "R", false, 0, "")

	pdf.CellFormat(60, 8, " Total Deductions", "", 0, "L", false, 0, "")
	pdf.CellFormat(35, 8, formatCurrency(payroll.TotalDeduction)+" ", "", 1, "R", false, 0, "")

	pdf.Ln(5)
	pdf.SetFillColor(220, 230, 241)
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(130, 12, "  TAKE HOME PAY ", "1", 0, "L", true, 0, "")
	pdf.CellFormat(60, 12, formatCurrency(payroll.NetSalary)+"  ", "1", 1, "R", true, 0, "")

	pdf.Ln(20)
	pdf.SetFont("Arial", "", 10)

	xAuth := 140.0
	pdf.SetX(xAuth)
	pdf.CellFormat(50, 5, "Authorized Signature,", "", 1, "C", false, 0, "")
	pdf.Ln(20)
	pdf.SetX(xAuth)
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(50, 5, "( HR Manager )", "T", 1, "C", false, 0, "")

	return pdf, payroll, nil
}

func (s *service) MarkAsPaid(ctx context.Context, id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	return s.repo.UpdateStatus(id, constants.PayrollStatusPaid)
}
