package loan

import (
	"context"

	"basekarya-backend/pkg/constants"
	"basekarya-backend/pkg/response"

	"github.com/stretchr/testify/mock"
	"github.com/xuri/excelize/v2"
)

type mockRepo struct{ mock.Mock }

func (m *mockRepo) Create(ctx context.Context, loan *Loan) error {
	return m.Called(ctx, loan).Error(0)
}

func (m *mockRepo) FindByID(ctx context.Context, id uint) (*Loan, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Loan), args.Error(1)
}

func (m *mockRepo) FindActiveLoanByUserID(ctx context.Context, userID uint) (*Loan, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Loan), args.Error(1)
}

func (m *mockRepo) FindAll(ctx context.Context, filter LoanFilter) ([]Loan, int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]Loan), args.Get(1).(int64), args.Error(2)
}

func (m *mockRepo) Update(ctx context.Context, loan *Loan) error {
	return m.Called(ctx, loan).Error(0)
}

func (m *mockRepo) GetBulkActiveLoansByEmployeeIds(ctx context.Context, ids []uint) (map[uint]Loan, error) {
	args := m.Called(ctx, ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[uint]Loan), args.Error(1)
}

type mockNotification struct{ mock.Mock }

func (m *mockNotification) SendNotification(ctx context.Context, userID uint, notifType string, title string, message string, relatedID uint) error {
	return m.Called(ctx, userID, notifType, title, message, relatedID).Error(0)
}

func (m *mockNotification) BlastNotification(ctx context.Context, userIDs []uint, notifType string, title string, message string, relatedID uint) error {
	return m.Called(ctx, userIDs, notifType, title, message, relatedID).Error(0)
}

type mockUserProvider struct{ mock.Mock }

func (m *mockUserProvider) FindApprovalUsers(ctx context.Context, permissionApprovalName string) ([]uint, error) {
	args := m.Called(ctx, permissionApprovalName)
	return args.Get(0).([]uint), args.Error(1)
}

type mockExcel struct{ mock.Mock }

func (m *mockExcel) GenerateSimpleExcel(sheetName string, headers []string, rows [][]interface{}) ([]byte, error) {
	args := m.Called(sheetName, headers, rows)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *mockExcel) NewFile() *excelize.File {
	args := m.Called()
	if args.Get(0) == nil {
		return excelize.NewFile()
	}
	return args.Get(0).(*excelize.File)
}

func (m *mockExcel) WriteToBuffer(file *excelize.File) ([]byte, error) {
	args := m.Called(file)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

type mockService struct{ mock.Mock }

func (m *mockService) Create(ctx context.Context, req *LoanRequest) error {
	return m.Called(ctx, req).Error(0)
}

func (m *mockService) GetLoanDetail(ctx context.Context, id uint) (*LoanDetailResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*LoanDetailResponse), args.Error(1)
}

func (m *mockService) GetLoans(ctx context.Context, filter LoanFilter) ([]LoanListResponse, *response.Meta, error) {
	args := m.Called(ctx, filter)
	var meta *response.Meta
	if args.Get(1) != nil {
		meta = args.Get(1).(*response.Meta)
	}
	return args.Get(0).([]LoanListResponse), meta, args.Error(2)
}

func (m *mockService) ProcessAction(ctx context.Context, req *ActionRequest) error {
	return m.Called(ctx, req).Error(0)
}

func (m *mockService) Export(ctx context.Context, filter LoanFilter) ([]byte, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func approvalLoanKey() string { return string(constants.APPROVAL_LOAN) }
