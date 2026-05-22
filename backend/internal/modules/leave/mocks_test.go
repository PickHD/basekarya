package leave

import (
	"context"
	"io"

	"basekarya-backend/internal/modules/attendance"
	"basekarya-backend/internal/modules/master"
	"basekarya-backend/pkg/constants"
	"basekarya-backend/pkg/response"

	"github.com/stretchr/testify/mock"
	"github.com/xuri/excelize/v2"
)

// --- Repository Mock ---

type mockRepo struct{ mock.Mock }

func (m *mockRepo) CreateRequest(ctx context.Context, req *LeaveRequest) error {
	return m.Called(ctx, req).Error(0)
}

func (m *mockRepo) FindRequestByID(ctx context.Context, id uint) (*LeaveRequest, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*LeaveRequest), args.Error(1)
}

func (m *mockRepo) FindAllRequests(ctx context.Context, filter *LeaveFilter) ([]LeaveRequest, int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]LeaveRequest), args.Get(1).(int64), args.Error(2)
}

func (m *mockRepo) GetBalance(ctx context.Context, employeeID, leaveTypeID uint, year int) (*LeaveBalance, error) {
	args := m.Called(ctx, employeeID, leaveTypeID, year)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*LeaveBalance), args.Error(1)
}

func (m *mockRepo) ApproveRequest(ctx context.Context, requestID, approverID uint, attendanceRecords []attendance.Attendance, shouldDeduct bool, days int) error {
	return m.Called(ctx, requestID, approverID, attendanceRecords, shouldDeduct, days).Error(0)
}

func (m *mockRepo) RejectRequest(ctx context.Context, requestID, approverID uint, reason string) error {
	return m.Called(ctx, requestID, approverID, reason).Error(0)
}

func (m *mockRepo) FindAllLeaveTypes(ctx context.Context) ([]master.LeaveType, error) {
	args := m.Called(ctx)
	return args.Get(0).([]master.LeaveType), args.Error(1)
}

func (m *mockRepo) CreateLeaveBalances(ctx context.Context, balances []LeaveBalance) error {
	return m.Called(ctx, balances).Error(0)
}

// --- StorageProvider Mock ---

type mockStorage struct{ mock.Mock }

func (m *mockStorage) UploadFileByte(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) (string, error) {
	args := m.Called(ctx, objectName, reader, size, contentType)
	return args.String(0), args.Error(1)
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

// --- Service Mock (for handler tests) ---

type mockService struct{ mock.Mock }

func (m *mockService) Apply(ctx context.Context, req *ApplyRequest) error {
	return m.Called(ctx, req).Error(0)
}

func (m *mockService) RequestAction(ctx context.Context, req *LeaveActionRequest) error {
	return m.Called(ctx, req).Error(0)
}

func (m *mockService) GetList(ctx context.Context, filter *LeaveFilter) ([]LeaveRequestListResponse, *response.Meta, error) {
	args := m.Called(ctx, filter)
	var meta *response.Meta
	if args.Get(1) != nil {
		meta = args.Get(1).(*response.Meta)
	}
	return args.Get(0).([]LeaveRequestListResponse), meta, args.Error(2)
}

func (m *mockService) GetDetail(ctx context.Context, id uint) (*LeaveRequestDetailResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*LeaveRequestDetailResponse), args.Error(1)
}

func (m *mockService) GenerateInitialBalance(ctx context.Context, employeeID uint) error {
	return m.Called(ctx, employeeID).Error(0)
}

func (m *mockService) GenerateAnnualBalance(ctx context.Context) error {
	return m.Called(ctx).Error(0)
}

func (m *mockService) Export(ctx context.Context, filter *LeaveFilter) ([]byte, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

// helper to get constant values in tests
func approvalLeaveKey() string { return string(constants.APPROVAL_LEAVE) }
