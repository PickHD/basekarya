package attendance

import (
	"context"
	"io"

	"basekarya-backend/internal/modules/user"
	"basekarya-backend/pkg/constants"
	"basekarya-backend/pkg/response"

	"github.com/stretchr/testify/mock"
	"github.com/xuri/excelize/v2"
)

type mockRepo struct{ mock.Mock }

func (m *mockRepo) GetTodayAttendance(ctx context.Context, employeeID uint) (*Attendance, error) {
	args := m.Called(ctx, employeeID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Attendance), args.Error(1)
}

func (m *mockRepo) Create(ctx context.Context, attendance *Attendance) error {
	return m.Called(ctx, attendance).Error(0)
}

func (m *mockRepo) Update(ctx context.Context, attendance *Attendance) error {
	return m.Called(ctx, attendance).Error(0)
}

func (m *mockRepo) GetHistory(ctx context.Context, employeeID uint, month, year, limit int, cursor string) ([]Attendance, *response.Cursor, error) {
	args := m.Called(ctx, employeeID, month, year, limit, cursor)
	var cur *response.Cursor
	if args.Get(1) != nil {
		cur = args.Get(1).(*response.Cursor)
	}
	if args.Get(0) == nil {
		return nil, cur, args.Error(2)
	}
	return args.Get(0).([]Attendance), cur, args.Error(2)
}

func (m *mockRepo) FindAll(ctx context.Context, filter *FilterParams) ([]Attendance, *response.Cursor, error) {
	args := m.Called(ctx, filter)
	var cur *response.Cursor
	if args.Get(1) != nil {
		cur = args.Get(1).(*response.Cursor)
	}
	if args.Get(0) == nil {
		return nil, cur, args.Error(2)
	}
	return args.Get(0).([]Attendance), cur, args.Error(2)
}

func (m *mockRepo) CountByStatus(ctx context.Context, status constants.AttendanceStatus, todayDate string) (int64, error) {
	args := m.Called(ctx, status, todayDate)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockRepo) CountAttendanceToday(ctx context.Context, todayDate string) (int64, error) {
	args := m.Called(ctx, todayDate)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockRepo) GetBulkLateDuration(ctx context.Context, month, year int) (map[uint]int, error) {
	args := m.Called(ctx, month, year)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[uint]int), args.Error(1)
}

type mockStorage struct{ mock.Mock }

func (m *mockStorage) UploadFileByte(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) (string, error) {
	args := m.Called(ctx, objectName, reader, size, contentType)
	return args.String(0), args.Error(1)
}

type mockLocationFetcher struct{ mock.Mock }

func (m *mockLocationFetcher) GetAddressFromCoords(lat, long float64) string {
	args := m.Called(lat, long)
	return args.String(0)
}

type mockUserProvider struct{ mock.Mock }

func (m *mockUserProvider) FindByID(ctx context.Context, id uint) (*user.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *mockUserProvider) CountActiveEmployee(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

type mockGeocodeWorker struct{ mock.Mock }

func (m *mockGeocodeWorker) Start(workerCount int) {
	m.Called(workerCount)
}

func (m *mockGeocodeWorker) Enqueue(job GeocodeJob) {
	m.Called(job)
}

func (m *mockGeocodeWorker) Stop() {
	m.Called()
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

func (m *mockService) Clock(ctx context.Context, userID uint, req *ClockRequest) (*AttendanceResponse, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*AttendanceResponse), args.Error(1)
}

func (m *mockService) GetTodayStatus(ctx context.Context, userID uint) (*TodayStatusResponse, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*TodayStatusResponse), args.Error(1)
}

func (m *mockService) GetMyHistory(ctx context.Context, userID uint, month, year, limit int, cursor string) ([]Attendance, *response.Meta, error) {
	args := m.Called(ctx, userID, month, year, limit, cursor)
	var meta *response.Meta
	if args.Get(1) != nil {
		meta = args.Get(1).(*response.Meta)
	}
	if args.Get(0) == nil {
		return nil, meta, args.Error(2)
	}
	return args.Get(0).([]Attendance), meta, args.Error(2)
}

func (m *mockService) GetAllRecap(ctx context.Context, filter *FilterParams) ([]RecapResponse, *response.Meta, error) {
	args := m.Called(ctx, filter)
	var meta *response.Meta
	if args.Get(1) != nil {
		meta = args.Get(1).(*response.Meta)
	}
	if args.Get(0) == nil {
		return nil, meta, args.Error(2)
	}
	return args.Get(0).([]RecapResponse), meta, args.Error(2)
}

func (m *mockService) GenerateExcel(ctx context.Context, filter *FilterParams) ([]byte, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *mockService) GetDashboardStats(ctx context.Context) (*DashboardStatResponse, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*DashboardStatResponse), args.Error(1)
}
