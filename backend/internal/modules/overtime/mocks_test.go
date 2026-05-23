package overtime

import (
	"context"
	"io"

	"basekarya-backend/internal/modules/user"
	"basekarya-backend/pkg/constants"
	"basekarya-backend/pkg/response"

	"github.com/stretchr/testify/mock"
	"github.com/xuri/excelize/v2"
)

// --- Repository Mock ---

type mockRepo struct{ mock.Mock }

func (m *mockRepo) Create(ctx context.Context, overtime *Overtime) error {
	return m.Called(ctx, overtime).Error(0)
}

func (m *mockRepo) FindByID(ctx context.Context, id uint) (*Overtime, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Overtime), args.Error(1)
}

func (m *mockRepo) FindAll(ctx context.Context, filter OvertimeFilter) ([]Overtime, int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]Overtime), args.Get(1).(int64), args.Error(2)
}

func (m *mockRepo) GetBulkActiveOvertimesByEmployeeIds(ctx context.Context, month, year int, ids []uint) (map[uint]int, error) {
	args := m.Called(ctx, month, year, ids)
	if args.Get(0) == nil {
		return map[uint]int{}, args.Error(1)
	}
	return args.Get(0).(map[uint]int), args.Error(1)
}

func (m *mockRepo) UpdateBulkStatusByEmployeeId(ctx context.Context, employeeID uint, periodMonth, periodYear int, status constants.OvertimeStatus) error {
	return m.Called(ctx, employeeID, periodMonth, periodYear, status).Error(0)
}

func (m *mockRepo) Update(ctx context.Context, overtime *Overtime) error {
	return m.Called(ctx, overtime).Error(0)
}

// --- NotificationProvider Mock ---

type mockNotification struct{ mock.Mock }

func (m *mockNotification) SendNotification(ctx context.Context, userID uint, notifType string, title string, message string, relatedID uint) error {
	return m.Called(ctx, userID, notifType, title, message, relatedID).Error(0)
}

func (m *mockNotification) BlastNotification(ctx context.Context, userIDs []uint, notifType string, title string, message string, relatedID uint) error {
	return m.Called(ctx, userIDs, notifType, title, message, relatedID).Error(0)
}

// --- UserProvider Mock ---

type mockUserProvider struct{ mock.Mock }

func (m *mockUserProvider) FindApprovalUsers(ctx context.Context, permissionApprovalName string) ([]uint, error) {
	args := m.Called(ctx, permissionApprovalName)
	return args.Get(0).([]uint), args.Error(1)
}

// --- ExcelProvider Mock ---

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

// --- StorageProvider Mock (for attendance) ---

type mockStorage struct{ mock.Mock }

func (m *mockStorage) UploadFileByte(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) (string, error) {
	args := m.Called(ctx, objectName, reader, size, contentType)
	return args.String(0), args.Error(1)
}

// --- Service Mock (for handler tests) ---

type mockService struct{ mock.Mock }

func (m *mockService) Create(ctx context.Context, req *OvertimeRequest) error {
	return m.Called(ctx, req).Error(0)
}

func (m *mockService) GetDetail(ctx context.Context, id uint) (*OvertimeDetailResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*OvertimeDetailResponse), args.Error(1)
}

func (m *mockService) GetList(ctx context.Context, filter OvertimeFilter) ([]OvertimeListResponse, *response.Meta, error) {
	args := m.Called(ctx, filter)
	var meta *response.Meta
	if args.Get(1) != nil {
		meta = args.Get(1).(*response.Meta)
	}
	return args.Get(0).([]OvertimeListResponse), meta, args.Error(2)
}

func (m *mockService) ProcessAction(ctx context.Context, req *ActionRequest) error {
	return m.Called(ctx, req).Error(0)
}

func (m *mockService) Export(ctx context.Context, filter OvertimeFilter) ([]byte, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

// helper for tests needing user.Employee
func makeEmployee(id, shiftID uint) *user.Employee {
	return &user.Employee{ID: id, ShiftID: shiftID, FullName: "John Doe", NIK: "EMP001"}
}
