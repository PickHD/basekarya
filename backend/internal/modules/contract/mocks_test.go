package contract

import (
	"context"
	"io"

	"basekarya-backend/pkg/response"

	"github.com/stretchr/testify/mock"
	"github.com/xuri/excelize/v2"
)

type mockRepo struct{ mock.Mock }

func (m *mockRepo) Upsert(ctx context.Context, contract *Contract) error {
	return m.Called(ctx, contract).Error(0)
}

func (m *mockRepo) FindByID(ctx context.Context, id uint) (*Contract, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Contract), args.Error(1)
}

func (m *mockRepo) FindByEmployeeID(ctx context.Context, employeeID uint) (*Contract, error) {
	args := m.Called(ctx, employeeID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Contract), args.Error(1)
}

func (m *mockRepo) FindAll(ctx context.Context, filter *ContractFilter) ([]Contract, int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]Contract), args.Get(1).(int64), args.Error(2)
}

func (m *mockRepo) FindExpiringContracts(ctx context.Context, withinDays int) ([]Contract, error) {
	args := m.Called(ctx, withinDays)
	return args.Get(0).([]Contract), args.Error(1)
}

func (m *mockRepo) MarkAlerted(ctx context.Context, ids []uint) error {
	return m.Called(ctx, ids).Error(0)
}

func (m *mockRepo) SoftDelete(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
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

type mockStorageProvider struct{ mock.Mock }

func (m *mockStorageProvider) UploadFileByte(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) (string, error) {
	args := m.Called(ctx, objectName, reader, size, contentType)
	return args.String(0), args.Error(1)
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

func (m *mockService) Upsert(ctx context.Context, req *UpsertContractRequest) error {
	return m.Called(ctx, req).Error(0)
}

func (m *mockService) GetList(ctx context.Context, filter *ContractFilter) ([]ContractListResponse, *response.Meta, error) {
	args := m.Called(ctx, filter)
	var meta *response.Meta
	if args.Get(1) != nil {
		meta = args.Get(1).(*response.Meta)
	}
	return args.Get(0).([]ContractListResponse), meta, args.Error(2)
}

func (m *mockService) GetDetail(ctx context.Context, id uint) (*ContractDetailResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ContractDetailResponse), args.Error(1)
}

func (m *mockService) GetByEmployeeID(ctx context.Context, employeeID uint) (*ContractDetailResponse, error) {
	args := m.Called(ctx, employeeID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ContractDetailResponse), args.Error(1)
}

func (m *mockService) Delete(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}

func (m *mockService) Export(ctx context.Context, filter *ContractFilter) ([]byte, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *mockService) CheckExpiringContracts(ctx context.Context) error {
	return m.Called(ctx).Error(0)
}
