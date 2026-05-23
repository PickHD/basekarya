package payroll

import (
	"context"

	"basekarya-backend/internal/modules/company"
	"basekarya-backend/internal/modules/loan"
	"basekarya-backend/internal/modules/user"
	"basekarya-backend/internal/testutil"
	"basekarya-backend/pkg/constants"
	"basekarya-backend/pkg/response"

	"github.com/signintech/gopdf"
	"github.com/stretchr/testify/mock"
)

type mockRepo struct{ mock.Mock }

func (m *mockRepo) CreateBulk(ctx context.Context, payroll *[]Payroll) error {
	return m.Called(ctx, payroll).Error(0)
}

func (m *mockRepo) FindAll(ctx context.Context, filter *PayrollFilter) ([]Payroll, int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]Payroll), args.Get(1).(int64), args.Error(2)
}

func (m *mockRepo) FindByID(ctx context.Context, id uint) (*Payroll, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Payroll), args.Error(1)
}

func (m *mockRepo) GetExistingEmployeeID(ctx context.Context, month, year int) (map[uint]bool, error) {
	args := m.Called(ctx, month, year)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[uint]bool), args.Error(1)
}

func (m *mockRepo) UpdateStatus(ctx context.Context, id uint, status constants.PayrollStatus) error {
	return m.Called(ctx, id, status).Error(0)
}

type mockUserProvider struct{ mock.Mock }

func (m *mockUserProvider) FindAllEmployeeActive(ctx context.Context) ([]user.Employee, error) {
	args := m.Called(ctx)
	return args.Get(0).([]user.Employee), args.Error(1)
}

type mockAttendanceProvider struct{ mock.Mock }

func (m *mockAttendanceProvider) GetBulkLateDuration(ctx context.Context, month, year int) (map[uint]int, error) {
	args := m.Called(ctx, month, year)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[uint]int), args.Error(1)
}

type mockReimbursementProvider struct{ mock.Mock }

func (m *mockReimbursementProvider) GetBulkApprovedAmount(ctx context.Context, month, year int) (map[uint]float64, error) {
	args := m.Called(ctx, month, year)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[uint]float64), args.Error(1)
}

type mockCompanyProvider struct{ mock.Mock }

func (m *mockCompanyProvider) FindByID(ctx context.Context, id uint) (*company.Company, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*company.Company), args.Error(1)
}

type mockNotificationProvider struct{ mock.Mock }

func (m *mockNotificationProvider) SendNotification(ctx context.Context, userID uint, Type string, Title string, Message string, relatedID uint) error {
	return m.Called(ctx, userID, Type, Title, Message, relatedID).Error(0)
}

type mockEmailProvider struct{ mock.Mock }

func (m *mockEmailProvider) SendWithAttachment(to, subject, htmlBody, fileName string, attachmentBytes []byte) error {
	return m.Called(to, subject, htmlBody, fileName, attachmentBytes).Error(0)
}

type mockLoanProvider struct{ mock.Mock }

func (m *mockLoanProvider) GetBulkActiveLoansByEmployeeIds(ctx context.Context, ids []uint) (map[uint]loan.Loan, error) {
	args := m.Called(ctx, ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[uint]loan.Loan), args.Error(1)
}

func (m *mockLoanProvider) Update(ctx context.Context, l *loan.Loan) error {
	return m.Called(ctx, l).Error(0)
}

type mockOvertimeProvider struct{ mock.Mock }

func (m *mockOvertimeProvider) GetBulkActiveOvertimesByEmployeeIds(ctx context.Context, month, year int, ids []uint) (map[uint]int, error) {
	args := m.Called(ctx, month, year, ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[uint]int), args.Error(1)
}

func (m *mockOvertimeProvider) UpdateBulkStatusByEmployeeId(ctx context.Context, employeeID uint, periodMonth, periodYear int, status constants.OvertimeStatus) error {
	return m.Called(ctx, employeeID, periodMonth, periodYear, status).Error(0)
}

type mockService struct{ mock.Mock }

func (m *mockService) GenerateAll(ctx context.Context, req *GenerateRequest) (*GenerateResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GenerateResponse), args.Error(1)
}

func (m *mockService) GetList(ctx context.Context, filter *PayrollFilter) ([]PayrollListResponse, *response.Meta, error) {
	args := m.Called(ctx, filter)
	var meta *response.Meta
	if args.Get(1) != nil {
		meta = args.Get(1).(*response.Meta)
	}
	return args.Get(0).([]PayrollListResponse), meta, args.Error(2)
}

func (m *mockService) GetDetail(ctx context.Context, id uint) (*PayrollDetailResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*PayrollDetailResponse), args.Error(1)
}

func (m *mockService) GeneratePayslipPDF(ctx context.Context, id uint) (*gopdf.GoPdf, *Payroll, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, nil, args.Error(2)
	}
	return args.Get(0).(*gopdf.GoPdf), args.Get(1).(*Payroll), args.Error(2)
}

func (m *mockService) MarkAsPaid(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}

func (m *mockService) BlastPayslipEmail(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}

func newTestService() (Service, *mockRepo, *mockUserProvider, *mockReimbursementProvider, *mockAttendanceProvider, *mockCompanyProvider, *mockNotificationProvider, *testutil.MockTransactionManager, *mockEmailProvider, *mockLoanProvider, *mockOvertimeProvider) {
	repo := new(mockRepo)
	userP := new(mockUserProvider)
	reimburse := new(mockReimbursementProvider)
	attend := new(mockAttendanceProvider)
	comp := new(mockCompanyProvider)
	notif := new(mockNotificationProvider)
	tm := testutil.NewMockTransactionManager()
	email := new(mockEmailProvider)
	loanP := new(mockLoanProvider)
	overtimeP := new(mockOvertimeProvider)

	svc := NewService(repo, userP, reimburse, attend, comp, notif, tm, nil, email, loanP, overtimeP)
	return svc, repo, userP, reimburse, attend, comp, notif, tm, email, loanP, overtimeP
}
