package reimbursement

import (
	"context"
	"mime/multipart"

	"basekarya-backend/pkg/response"

	"github.com/stretchr/testify/mock"
	"github.com/xuri/excelize/v2"
)

type mockRepo struct{ mock.Mock }

func (m *mockRepo) Create(ctx context.Context, reimbursement *Reimbursement) error {
	return m.Called(ctx, reimbursement).Error(0)
}

func (m *mockRepo) FindByID(ctx context.Context, id uint) (*Reimbursement, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Reimbursement), args.Error(1)
}

func (m *mockRepo) FindAll(ctx context.Context, filter ReimbursementFilter) ([]Reimbursement, int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]Reimbursement), args.Get(1).(int64), args.Error(2)
}

func (m *mockRepo) Update(ctx context.Context, reimbursement *Reimbursement) error {
	return m.Called(ctx, reimbursement).Error(0)
}

func (m *mockRepo) GetBulkApprovedAmount(ctx context.Context, month, year int) (map[uint]float64, error) {
	args := m.Called(ctx, month, year)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[uint]float64), args.Error(1)
}

type mockStorage struct{ mock.Mock }

func (m *mockStorage) UploadFileMultipart(ctx context.Context, file *multipart.FileHeader, objectName string) (string, error) {
	args := m.Called(ctx, file, objectName)
	return args.String(0), args.Error(1)
}

type mockNotification struct{ mock.Mock }

func (m *mockNotification) SendNotification(ctx context.Context, userID uint, Type string, Title string, Message string, relatedID uint) error {
	return m.Called(ctx, userID, Type, Title, Message, relatedID).Error(0)
}

func (m *mockNotification) BlastNotification(ctx context.Context, userIDs []uint, Type string, Title string, Message string, relatedID uint) error {
	return m.Called(ctx, userIDs, Type, Title, Message, relatedID).Error(0)
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

func (m *mockService) Create(ctx context.Context, req *ReimbursementRequest) error {
	return m.Called(ctx, req).Error(0)
}

func (m *mockService) GetReimburseDetail(ctx context.Context, id uint) (*ReimbursementDetailResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ReimbursementDetailResponse), args.Error(1)
}

func (m *mockService) GetReimbursements(ctx context.Context, filter ReimbursementFilter) ([]ReimbursementListResponse, *response.Meta, error) {
	args := m.Called(ctx, filter)
	var meta *response.Meta
	if args.Get(1) != nil {
		meta = args.Get(1).(*response.Meta)
	}
	return args.Get(0).([]ReimbursementListResponse), meta, args.Error(2)
}

func (m *mockService) ProcessAction(ctx context.Context, req *ActionRequest) error {
	return m.Called(ctx, req).Error(0)
}

func (m *mockService) Export(ctx context.Context, filter ReimbursementFilter) ([]byte, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}
