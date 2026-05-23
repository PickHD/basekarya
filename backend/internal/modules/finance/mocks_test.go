package finance

import (
	"context"

	"basekarya-backend/pkg/response"

	"github.com/stretchr/testify/mock"
	"github.com/xuri/excelize/v2"
)

type mockRepo struct{ mock.Mock }

func (m *mockRepo) CreateTransaction(ctx context.Context, tx *FinanceTransaction) error {
	return m.Called(ctx, tx).Error(0)
}

func (m *mockRepo) FindTransactionByID(ctx context.Context, id uint) (*FinanceTransaction, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*FinanceTransaction), args.Error(1)
}

func (m *mockRepo) FindAllTransactions(ctx context.Context, filter TransactionFilter) ([]FinanceTransaction, *response.Cursor, error) {
	args := m.Called(ctx, filter)
	var cursor *response.Cursor
	if args.Get(1) != nil {
		cursor = args.Get(1).(*response.Cursor)
	}
	return args.Get(0).([]FinanceTransaction), cursor, args.Error(2)
}

func (m *mockRepo) UpdateTransaction(ctx context.Context, tx *FinanceTransaction) error {
	return m.Called(ctx, tx).Error(0)
}

func (m *mockRepo) CreateCategory(ctx context.Context, cat *FinanceCategory) error {
	return m.Called(ctx, cat).Error(0)
}

func (m *mockRepo) FindCategoryByID(ctx context.Context, id uint) (*FinanceCategory, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*FinanceCategory), args.Error(1)
}

func (m *mockRepo) FindAllCategories(ctx context.Context, catType string) ([]FinanceCategory, error) {
	args := m.Called(ctx, catType)
	return args.Get(0).([]FinanceCategory), args.Error(1)
}

func (m *mockRepo) UpdateCategory(ctx context.Context, cat *FinanceCategory) error {
	return m.Called(ctx, cat).Error(0)
}

func (m *mockRepo) DeleteCategory(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}

func (m *mockRepo) GetDashboardSummary(ctx context.Context, startDate, endDate string) (*DashboardResponse, error) {
	args := m.Called(ctx, startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*DashboardResponse), args.Error(1)
}

type mockNotificationProvider struct{ mock.Mock }

func (m *mockNotificationProvider) SendNotification(ctx context.Context, userID uint, notifType string, title string, message string, relatedID uint) error {
	return m.Called(ctx, userID, notifType, title, message, relatedID).Error(0)
}

func (m *mockNotificationProvider) BlastNotification(ctx context.Context, userIDs []uint, notifType string, title string, message string, relatedID uint) error {
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

func (m *mockService) CreateTransaction(ctx context.Context, req *CreateTransactionRequest) error {
	return m.Called(ctx, req).Error(0)
}

func (m *mockService) GetTransactionDetail(ctx context.Context, id uint) (*TransactionDetailResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*TransactionDetailResponse), args.Error(1)
}

func (m *mockService) GetTransactions(ctx context.Context, filter TransactionFilter) ([]TransactionListResponse, *response.Meta, error) {
	args := m.Called(ctx, filter)
	var meta *response.Meta
	if args.Get(1) != nil {
		meta = args.Get(1).(*response.Meta)
	}
	return args.Get(0).([]TransactionListResponse), meta, args.Error(2)
}

func (m *mockService) ProcessAction(ctx context.Context, req *ActionRequest) error {
	return m.Called(ctx, req).Error(0)
}

func (m *mockService) ExportTransactions(ctx context.Context, filter TransactionFilter) ([]byte, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *mockService) CreateCategory(ctx context.Context, req *CategoryRequest) error {
	return m.Called(ctx, req).Error(0)
}

func (m *mockService) GetCategories(ctx context.Context, catType string) ([]CategoryResponse, error) {
	args := m.Called(ctx, catType)
	return args.Get(0).([]CategoryResponse), args.Error(1)
}

func (m *mockService) UpdateCategory(ctx context.Context, id uint, req *CategoryRequest) error {
	return m.Called(ctx, id, req).Error(0)
}

func (m *mockService) DeleteCategory(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}

func (m *mockService) GetDashboard(ctx context.Context, startDate, endDate string) (*DashboardResponse, error) {
	args := m.Called(ctx, startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*DashboardResponse), args.Error(1)
}
