package payroll

import (
	"basekarya-backend/internal/modules/bpjs"
	"basekarya-backend/internal/modules/company"
	"basekarya-backend/internal/modules/loan"
	"basekarya-backend/internal/modules/tax"
	"basekarya-backend/internal/modules/user"
	"basekarya-backend/pkg/constants"
	"context"
)

type UserProvider interface {
	FindAllEmployeeActive(ctx context.Context) ([]user.Employee, error)
}

type AttendanceProvider interface {
	GetBulkLateDuration(ctx context.Context, month, year int) (map[uint]int, error)
}

type ReimbursementProvider interface {
	GetBulkApprovedAmount(ctx context.Context, month, year int) (map[uint]float64, error)
}

type CompanyProvider interface {
	FindByID(ctx context.Context, id uint) (*company.Company, error)
}

type NotificationProvider interface {
	SendNotification(ctx context.Context, userID uint,
		Type string,
		Title string,
		Message string, relatedID uint) error
}

type EmailProvider interface {
	SendWithAttachment(to, subject, htmlBody, fileName string, attachmentBytes []byte) error
}

type LoanProvider interface {
	GetBulkActiveLoansByEmployeeIds(ctx context.Context, ids []uint) (map[uint]loan.Loan, error)
	Update(ctx context.Context, loan *loan.Loan) error
}

type OvertimeProvider interface {
	GetBulkActiveOvertimesByEmployeeIds(ctx context.Context, month, year int, ids []uint) (map[uint]int, error)
	UpdateBulkStatusByEmployeeId(ctx context.Context, employeeID uint, periodMonth, periodYear int, status constants.OvertimeStatus) error
}

type TaxProvider interface {
	CalculateTER(ctx context.Context, grossMonthlyIncome float64, maritalStatus constants.MaritalStatus, dependentsCount int) (*tax.PPh21Result, error)
}

type BPJSProvider interface {
	CalculateAll(ctx context.Context, grossMonthlyIncome float64) ([]bpjs.BPJSComponent, error)
}
