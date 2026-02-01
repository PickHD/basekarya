package payroll

import (
	"context"
	"fmt"
	"hris-backend/internal/modules/attendance"
	"hris-backend/internal/modules/reimbursement"
	"hris-backend/internal/modules/user"
	"hris-backend/pkg/constants"
	"hris-backend/pkg/logger"
	"hris-backend/pkg/response"
	"time"
)

type Service interface {
	GenerateAll(ctx context.Context, req *GenerateRequest) (*GenerateResponse, error)
	GetList(ctx context.Context, filter *PayrollFilter) ([]PayrollListResponse, *response.Meta, error)
}

type service struct {
	repo              Repository
	userRepo          user.Repository
	reimbursementRepo reimbursement.Repository
	attendanceRepo    attendance.Repository
}

func NewService(repo Repository,
	userRepo user.Repository,
	reimbursementRepo reimbursement.Repository,
	attendanceRepo attendance.Repository) Service {
	return &service{repo, userRepo, reimbursementRepo, attendanceRepo}
}

func (s *service) GenerateAll(ctx context.Context, req *GenerateRequest) (*GenerateResponse, error) {
	employees, err := s.userRepo.FindAllEmployeeActive()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch all employee active: %w", err)
	}

	existingPayrollMap, err := s.repo.GetExistingEmployeeID(req.Month, req.Year)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch existing employee id: %w", err)
	}

	attendanceMap, err := s.attendanceRepo.GetBulkLateDuration(req.Month, req.Year)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch bulk late duration: %w", err)
	}

	reimburseMap, err := s.reimbursementRepo.GetBulkApprovedAmount(req.Month, req.Year)
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
		reimburseAmount := reimburseMap[emp.ID]

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

	meta := response.NewMeta(filter.Page, filter.Limit, total)
	return responses, meta, nil
}
