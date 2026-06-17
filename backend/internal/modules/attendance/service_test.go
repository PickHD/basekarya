package attendance

import (
	"errors"
	"testing"
	"time"

	"basekarya-backend/internal/modules/department"
	"basekarya-backend/internal/modules/master"
	"basekarya-backend/internal/modules/user"
	"basekarya-backend/internal/testutil"
	"basekarya-backend/pkg/constants"
	"basekarya-backend/pkg/response"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func newTestAttendanceService() (Service, *mockRepo, *mockUserProvider, *mockStorage, *mockGeocodeWorker, *testutil.MockTransactionManager, *mockExcel) {
	repo := new(mockRepo)
	userProv := new(mockUserProvider)
	storage := new(mockStorage)
	geo := new(mockGeocodeWorker)
	tm := testutil.NewMockTransactionManager()
	excel := new(mockExcel)

	svc := NewService(repo, userProv, storage, geo, tm, excel)
	return svc, repo, userProv, storage, geo, tm, excel
}

func shiftTimeForPresent() string {
	t := time.Now().Add(1 * time.Hour)
	return t.Format("15:04:05")
}

func shiftTimeForLate() string {
	now := time.Now()
	if now.Hour() >= 3 {
		return now.Add(-2 * time.Hour).Format("15:04:05")
	}
	return "02:00:00"
}

func shiftTimeForTooEarly() string {
	t := time.Now().Add(5 * time.Hour)
	return t.Format("15:04:05")
}

func TestService_Clock(t *testing.T) {
	ctx := testutil.CtxWithTenant(1, 1, false)

	tests := []struct {
		name       string
		userID     uint
		req        *ClockRequest
		setupMocks func(*mockRepo, *mockUserProvider, *mockStorage, *mockGeocodeWorker)
		wantErr    bool
		errMsg     string
	}{
		{
			name:   "check-in success present",
			userID: 1,
			req: &ClockRequest{
				Latitude:    -6.2,
				Longitude:   106.8,
				ImageBase64: "aGVsbG8=",
			},
			setupMocks: func(r *mockRepo, u *mockUserProvider, s *mockStorage, g *mockGeocodeWorker) {
				u.On("FindByID", mock.Anything, uint(1)).Return(&user.User{
					Employee: &user.Employee{
						ID:     1,
						ShiftID: 1,
						Shift: &master.Shift{
							ID:        1,
							StartTime: shiftTimeForPresent(),
						},
					},
				}, nil)
				r.On("GetTodayAttendance", mock.Anything, uint(1)).Return(nil, gorm.ErrRecordNotFound)
				s.On("UploadFileByte", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return("http://img.url/in.jpg", nil)
				r.On("Create", mock.Anything, mock.AnythingOfType("*attendance.Attendance")).Return(nil)
				g.On("Enqueue", mock.Anything)
			},
			wantErr: false,
		},
		{
			name:   "check-in success late",
			userID: 1,
			req: &ClockRequest{
				Latitude:    -6.2,
				Longitude:   106.8,
				ImageBase64: "aGVsbG8=",
			},
			setupMocks: func(r *mockRepo, u *mockUserProvider, s *mockStorage, g *mockGeocodeWorker) {
				u.On("FindByID", mock.Anything, uint(1)).Return(&user.User{
					Employee: &user.Employee{
						ID:     1,
						ShiftID: 1,
						Shift: &master.Shift{
							ID:        1,
							StartTime: shiftTimeForLate(),
						},
					},
				}, nil)
				r.On("GetTodayAttendance", mock.Anything, uint(1)).Return(nil, gorm.ErrRecordNotFound)
				s.On("UploadFileByte", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return("http://img.url/in.jpg", nil)
				r.On("Create", mock.Anything, mock.AnythingOfType("*attendance.Attendance")).Return(nil)
				g.On("Enqueue", mock.Anything)
			},
			wantErr: false,
		},
		{
			name:   "check-in too early",
			userID: 1,
			req: &ClockRequest{
				Latitude:    -6.2,
				Longitude:   106.8,
				ImageBase64: "aGVsbG8=",
			},
			setupMocks: func(r *mockRepo, u *mockUserProvider, s *mockStorage, g *mockGeocodeWorker) {
				u.On("FindByID", mock.Anything, uint(1)).Return(&user.User{
					Employee: &user.Employee{
						ID:     1,
						ShiftID: 1,
						Shift: &master.Shift{
							ID:        1,
							StartTime: shiftTimeForTooEarly(),
						},
					},
				}, nil)
				r.On("GetTodayAttendance", mock.Anything, uint(1)).Return(nil, gorm.ErrRecordNotFound)
			},
			wantErr: true,
			errMsg:  "cannot check-in, too early",
		},
		{
			name:   "check-out success",
			userID: 1,
			req: &ClockRequest{
				Latitude:    -6.2,
				Longitude:   106.8,
				ImageBase64: "aGVsbG8=",
			},
			setupMocks: func(r *mockRepo, u *mockUserProvider, s *mockStorage, g *mockGeocodeWorker) {
				u.On("FindByID", mock.Anything, uint(1)).Return(&user.User{
					Employee: &user.Employee{
						ID:     1,
						ShiftID: 1,
						Shift: &master.Shift{
							ID:        1,
							StartTime: "08:00:00",
						},
					},
				}, nil)
				checkInTime := time.Now().Add(-1 * time.Hour)
				r.On("GetTodayAttendance", mock.Anything, uint(1)).Return(&Attendance{
					ID:          1,
					CheckInTime: checkInTime,
					CheckInLat:  -6.2,
					CheckInLong: 106.8,
					Status:      string(constants.AttendanceStatusPresent),
				}, nil)
				s.On("UploadFileByte", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return("http://img.url/out.jpg", nil)
				r.On("Update", mock.Anything, mock.AnythingOfType("*attendance.Attendance")).Return(nil)
				g.On("Enqueue", mock.Anything)
			},
			wantErr: false,
		},
		{
			name:   "already completed attendance today",
			userID: 1,
			req: &ClockRequest{
				Latitude:    -6.2,
				Longitude:   106.8,
				ImageBase64: "aGVsbG8=",
			},
			setupMocks: func(r *mockRepo, u *mockUserProvider, s *mockStorage, g *mockGeocodeWorker) {
				u.On("FindByID", mock.Anything, uint(1)).Return(&user.User{
					Employee: &user.Employee{
						ID:     1,
						ShiftID: 1,
						Shift: &master.Shift{
							ID:        1,
							StartTime: "08:00:00",
						},
					},
				}, nil)
				checkOutTime := time.Now()
				r.On("GetTodayAttendance", mock.Anything, uint(1)).Return(&Attendance{
					ID:           1,
					CheckInTime:  time.Now().Add(-2 * time.Hour),
					CheckOutTime: &checkOutTime,
				}, nil)
			},
			wantErr: true,
			errMsg:  "you have already completed attendance for today",
		},
		{
			name:   "employee data not found",
			userID: 99,
			req: &ClockRequest{
				Latitude:    -6.2,
				Longitude:   106.8,
				ImageBase64: "aGVsbG8=",
			},
			setupMocks: func(r *mockRepo, u *mockUserProvider, s *mockStorage, g *mockGeocodeWorker) {
				u.On("FindByID", mock.Anything, uint(99)).Return(&user.User{}, nil)
			},
			wantErr: true,
			errMsg:  "employee data not found",
		},
		{
			name:   "employee shift not assigned",
			userID: 1,
			req: &ClockRequest{
				Latitude:    -6.2,
				Longitude:   106.8,
				ImageBase64: "aGVsbG8=",
			},
			setupMocks: func(r *mockRepo, u *mockUserProvider, s *mockStorage, g *mockGeocodeWorker) {
				u.On("FindByID", mock.Anything, uint(1)).Return(&user.User{
					Employee: &user.Employee{
						ID:     1,
						Shift:  nil,
					},
				}, nil)
			},
			wantErr: true,
			errMsg:  "employee shift not assigned",
		},
		{
			name:   "user find error",
			userID: 1,
			req: &ClockRequest{
				Latitude:    -6.2,
				Longitude:   106.8,
				ImageBase64: "aGVsbG8=",
			},
			setupMocks: func(r *mockRepo, u *mockUserProvider, s *mockStorage, g *mockGeocodeWorker) {
				u.On("FindByID", mock.Anything, uint(1)).Return(nil, errors.New("db error"))
			},
			wantErr: true,
			errMsg:  "employee data not found",
		},
		{
			name:   "storage upload error on check-in",
			userID: 1,
			req: &ClockRequest{
				Latitude:    -6.2,
				Longitude:   106.8,
				ImageBase64: "aGVsbG8=",
			},
			setupMocks: func(r *mockRepo, u *mockUserProvider, s *mockStorage, g *mockGeocodeWorker) {
				u.On("FindByID", mock.Anything, uint(1)).Return(&user.User{
					Employee: &user.Employee{
						ID:     1,
						ShiftID: 1,
						Shift: &master.Shift{
							ID:        1,
							StartTime: shiftTimeForPresent(),
						},
					},
				}, nil)
				r.On("GetTodayAttendance", mock.Anything, uint(1)).Return(nil, gorm.ErrRecordNotFound)
				s.On("UploadFileByte", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return("", errors.New("upload failed"))
			},
			wantErr: true,
			errMsg:  "upload failed",
		},
		{
			name:   "repo create error on check-in",
			userID: 1,
			req: &ClockRequest{
				Latitude:    -6.2,
				Longitude:   106.8,
				ImageBase64: "aGVsbG8=",
			},
			setupMocks: func(r *mockRepo, u *mockUserProvider, s *mockStorage, g *mockGeocodeWorker) {
				u.On("FindByID", mock.Anything, uint(1)).Return(&user.User{
					Employee: &user.Employee{
						ID:     1,
						ShiftID: 1,
						Shift: &master.Shift{
							ID:        1,
							StartTime: shiftTimeForPresent(),
						},
					},
				}, nil)
				r.On("GetTodayAttendance", mock.Anything, uint(1)).Return(nil, gorm.ErrRecordNotFound)
				s.On("UploadFileByte", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return("http://img.url/in.jpg", nil)
				r.On("Create", mock.Anything, mock.AnythingOfType("*attendance.Attendance")).Return(errors.New("db error"))
			},
			wantErr: true,
			errMsg:  "db error",
		},
		{
			name:   "storage upload error on check-out",
			userID: 1,
			req: &ClockRequest{
				Latitude:    -6.2,
				Longitude:   106.8,
				ImageBase64: "aGVsbG8=",
			},
			setupMocks: func(r *mockRepo, u *mockUserProvider, s *mockStorage, g *mockGeocodeWorker) {
				u.On("FindByID", mock.Anything, uint(1)).Return(&user.User{
					Employee: &user.Employee{
						ID:     1,
						ShiftID: 1,
						Shift: &master.Shift{
							ID:        1,
							StartTime: "08:00:00",
						},
					},
				}, nil)
				checkInTime := time.Now().Add(-1 * time.Hour)
				r.On("GetTodayAttendance", mock.Anything, uint(1)).Return(&Attendance{
					ID:          1,
					CheckInTime: checkInTime,
					CheckInLat:  -6.2,
					CheckInLong: 106.8,
				}, nil)
				s.On("UploadFileByte", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return("", errors.New("upload failed"))
			},
			wantErr: true,
			errMsg:  "upload failed",
		},
		{
			name:   "repo update error on check-out",
			userID: 1,
			req: &ClockRequest{
				Latitude:    -6.2,
				Longitude:   106.8,
				ImageBase64: "aGVsbG8=",
			},
			setupMocks: func(r *mockRepo, u *mockUserProvider, s *mockStorage, g *mockGeocodeWorker) {
				u.On("FindByID", mock.Anything, uint(1)).Return(&user.User{
					Employee: &user.Employee{
						ID:     1,
						ShiftID: 1,
						Shift: &master.Shift{
							ID:        1,
							StartTime: "08:00:00",
						},
					},
				}, nil)
				checkInTime := time.Now().Add(-1 * time.Hour)
				r.On("GetTodayAttendance", mock.Anything, uint(1)).Return(&Attendance{
					ID:          1,
					CheckInTime: checkInTime,
					CheckInLat:  -6.2,
					CheckInLong: 106.8,
				}, nil)
				s.On("UploadFileByte", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return("http://img.url/out.jpg", nil)
				r.On("Update", mock.Anything, mock.AnythingOfType("*attendance.Attendance")).Return(errors.New("db error"))
			},
			wantErr: true,
			errMsg:  "db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, userProv, storage, geo, _, _ := newTestAttendanceService()
			tt.setupMocks(repo, userProv, storage, geo)

			resp, err := svc.Clock(ctx, tt.userID, tt.req)

			if tt.wantErr {
				require.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
				assert.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				assert.NotEmpty(t, resp.Type)
				assert.NotEmpty(t, resp.Status)
			}
		})
	}
}

func TestService_GetTodayStatus(t *testing.T) {
	ctx := testutil.CtxWithTenant(1, 1, false)

	tests := []struct {
		name       string
		userID     uint
		setupMocks func(*mockRepo, *mockUserProvider)
		wantErr    bool
		errMsg     string
		assertFn   func(*testing.T, *TodayStatusResponse)
	}{
		{
			name:   "absent no record",
			userID: 1,
			setupMocks: func(r *mockRepo, u *mockUserProvider) {
				u.On("FindByID", mock.Anything, uint(1)).Return(&user.User{
					Employee: &user.Employee{ID: 1},
				}, nil)
				r.On("GetTodayAttendance", mock.Anything, uint(1)).Return(nil, gorm.ErrRecordNotFound)
			},
			wantErr: false,
			assertFn: func(t *testing.T, resp *TodayStatusResponse) {
				assert.Equal(t, string(constants.AttendanceStatusAbsent), resp.Status)
				assert.Equal(t, string(constants.AttendanceTypeNone), resp.Type)
			},
		},
		{
			name:   "checked in no checkout",
			userID: 1,
			setupMocks: func(r *mockRepo, u *mockUserProvider) {
				u.On("FindByID", mock.Anything, uint(1)).Return(&user.User{
					Employee: &user.Employee{ID: 1},
				}, nil)
				r.On("GetTodayAttendance", mock.Anything, uint(1)).Return(&Attendance{
					ID:          1,
					CheckInTime: time.Now(),
					Status:      string(constants.AttendanceStatusPresent),
				}, nil)
			},
			wantErr: false,
			assertFn: func(t *testing.T, resp *TodayStatusResponse) {
				assert.Equal(t, string(constants.AttendanceTypeCheckIn), resp.Type)
				assert.NotNil(t, resp.CheckInTime)
				assert.Nil(t, resp.CheckOutTime)
			},
		},
		{
			name:   "completed with checkout",
			userID: 1,
			setupMocks: func(r *mockRepo, u *mockUserProvider) {
				u.On("FindByID", mock.Anything, uint(1)).Return(&user.User{
					Employee: &user.Employee{ID: 1},
				}, nil)
				checkIn := time.Now().Add(-8 * time.Hour)
				checkOut := time.Now()
				r.On("GetTodayAttendance", mock.Anything, uint(1)).Return(&Attendance{
					ID:           1,
					CheckInTime:  checkIn,
					CheckOutTime: &checkOut,
					Status:       string(constants.AttendanceStatusPresent),
				}, nil)
			},
			wantErr: false,
			assertFn: func(t *testing.T, resp *TodayStatusResponse) {
				assert.Equal(t, string(constants.AttendanceTypeCompleted), resp.Type)
				assert.NotNil(t, resp.CheckInTime)
				assert.NotNil(t, resp.CheckOutTime)
				assert.NotEmpty(t, resp.WorkDuration)
			},
		},
		{
			name:   "employee not found",
			userID: 99,
			setupMocks: func(r *mockRepo, u *mockUserProvider) {
				u.On("FindByID", mock.Anything, uint(99)).Return(nil, errors.New("not found"))
			},
			wantErr: true,
			errMsg:  "employee not found",
		},
		{
			name:   "employee nil in user",
			userID: 1,
			setupMocks: func(r *mockRepo, u *mockUserProvider) {
				u.On("FindByID", mock.Anything, uint(1)).Return(&user.User{}, nil)
			},
			wantErr: true,
			errMsg:  "employee not found",
		},
		{
			name:   "repo error",
			userID: 1,
			setupMocks: func(r *mockRepo, u *mockUserProvider) {
				u.On("FindByID", mock.Anything, uint(1)).Return(&user.User{
					Employee: &user.Employee{ID: 1},
				}, nil)
				r.On("GetTodayAttendance", mock.Anything, uint(1)).Return(nil, errors.New("db error"))
			},
			wantErr: true,
			errMsg:  "db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, userProv, _, _, _, _ := newTestAttendanceService()
			tt.setupMocks(repo, userProv)

			resp, err := svc.GetTodayStatus(ctx, tt.userID)

			if tt.wantErr {
				require.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				if tt.assertFn != nil {
					tt.assertFn(t, resp)
				}
			}
		})
	}
}

func TestService_GetMyHistory(t *testing.T) {
	ctx := testutil.CtxWithTenant(1, 1, false)

	tests := []struct {
		name       string
		userID     uint
		cursor     string
		setupMocks func(*mockRepo, *mockUserProvider)
		wantErr    bool
		errMsg     string
		wantLen    int
	}{
		{
			name:   "success with data",
			userID: 1,
			cursor: "",
			setupMocks: func(r *mockRepo, u *mockUserProvider) {
				u.On("FindByID", mock.Anything, uint(1)).Return(&user.User{
					Employee: &user.Employee{ID: 1},
				}, nil)
				r.On("GetHistory", mock.Anything, uint(1), 5, 2026, 10, "").Return([]Attendance{
					{ID: 1, Status: string(constants.AttendanceStatusPresent)},
					{ID: 2, Status: string(constants.AttendanceStatusLate)},
				}, nil, nil)
			},
			wantErr: false,
			wantLen: 2,
		},
		{
			name:   "success empty",
			userID: 1,
			cursor: "",
			setupMocks: func(r *mockRepo, u *mockUserProvider) {
				u.On("FindByID", mock.Anything, uint(1)).Return(&user.User{
					Employee: &user.Employee{ID: 1},
				}, nil)
				r.On("GetHistory", mock.Anything, uint(1), 5, 2026, 10, "").Return([]Attendance{}, nil, nil)
			},
			wantErr: false,
			wantLen: 0,
		},
		{
			name:   "success with cursor",
			userID: 1,
			cursor: "some-cursor",
			setupMocks: func(r *mockRepo, u *mockUserProvider) {
				u.On("FindByID", mock.Anything, uint(1)).Return(&user.User{
					Employee: &user.Employee{ID: 1},
				}, nil)
				r.On("GetHistory", mock.Anything, uint(1), 5, 2026, 10, "some-cursor").Return([]Attendance{
					{ID: 3},
				}, &response.Cursor{ID: 3, SortValue: time.Now()}, nil)
			},
			wantErr: false,
			wantLen: 1,
		},
		{
			name:   "employee not found",
			userID: 99,
			cursor: "",
			setupMocks: func(r *mockRepo, u *mockUserProvider) {
				u.On("FindByID", mock.Anything, uint(99)).Return(nil, errors.New("not found"))
			},
			wantErr: true,
			errMsg:  "employee not found",
		},
		{
			name:   "repo error",
			userID: 1,
			cursor: "",
			setupMocks: func(r *mockRepo, u *mockUserProvider) {
				u.On("FindByID", mock.Anything, uint(1)).Return(&user.User{
					Employee: &user.Employee{ID: 1},
				}, nil)
				r.On("GetHistory", mock.Anything, uint(1), 5, 2026, 10, "").Return(nil, nil, errors.New("db error"))
			},
			wantErr: true,
			errMsg:  "db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, userProv, _, _, _, _ := newTestAttendanceService()
			tt.setupMocks(repo, userProv)

			logs, meta, err := svc.GetMyHistory(ctx, tt.userID, 5, 2026, 10, tt.cursor)

			if tt.wantErr {
				require.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				require.NoError(t, err)
				assert.Len(t, logs, tt.wantLen)
				if tt.wantLen > 0 {
					assert.NotNil(t, meta)
				}
			}
		})
	}
}

func TestService_GetAllRecap(t *testing.T) {
	ctx := testutil.CtxWithTenant(1, 1, false)

	tests := []struct {
		name       string
		filter     *FilterParams
		setupMocks func(*mockRepo)
		wantErr    bool
		wantLen    int
	}{
		{
			name:   "success with data",
			filter: &FilterParams{Limit: 10},
			setupMocks: func(r *mockRepo) {
				r.On("FindAll", mock.Anything, mock.AnythingOfType("*attendance.FilterParams")).Return([]Attendance{
					{
						ID:          1,
						CheckInTime: time.Now(),
						Status:      string(constants.AttendanceStatusPresent),
						Employee:    &user.Employee{FullName: "John", NIK: "001", Department: &department.Department{Name: "IT"}},
						Shift:       &master.Shift{Name: "Morning"},
					},
				}, nil, nil)
			},
			wantErr: false,
			wantLen: 1,
		},
		{
			name:   "success empty",
			filter: &FilterParams{Limit: 10},
			setupMocks: func(r *mockRepo) {
				r.On("FindAll", mock.Anything, mock.AnythingOfType("*attendance.FilterParams")).Return([]Attendance{}, nil, nil)
			},
			wantErr: false,
			wantLen: 0,
		},
		{
			name:   "success with checkout",
			filter: &FilterParams{Limit: 10},
			setupMocks: func(r *mockRepo) {
				checkOut := time.Now().Add(8 * time.Hour)
				r.On("FindAll", mock.Anything, mock.AnythingOfType("*attendance.FilterParams")).Return([]Attendance{
					{
						ID:           1,
						Date:         time.Now(),
						CheckInTime:  time.Now(),
						CheckOutTime: &checkOut,
						Status:       string(constants.AttendanceStatusPresent),
						Employee:     &user.Employee{FullName: "Jane", NIK: "002", Department: &department.Department{Name: "HR"}},
						Shift:        &master.Shift{Name: "Morning"},
					},
				}, nil, nil)
			},
			wantErr: false,
			wantLen: 1,
		},
		{
			name:   "repo error",
			filter: &FilterParams{Limit: 10},
			setupMocks: func(r *mockRepo) {
				r.On("FindAll", mock.Anything, mock.AnythingOfType("*attendance.FilterParams")).Return(nil, nil, errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _, _, _, _, _ := newTestAttendanceService()
			tt.setupMocks(repo)

			result, meta, err := svc.GetAllRecap(ctx, tt.filter)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Len(t, result, tt.wantLen)
				if tt.wantLen > 0 {
					assert.NotNil(t, meta)
				}
			}
		})
	}
}

func TestService_GenerateExcel(t *testing.T) {
	ctx := testutil.CtxWithTenant(1, 1, false)

	tests := []struct {
		name       string
		filter     *FilterParams
		setupMocks func(*mockRepo, *mockExcel)
		wantErr    bool
	}{
		{
			name:   "success with data",
			filter: &FilterParams{},
			setupMocks: func(r *mockRepo, e *mockExcel) {
				checkOut := time.Now().Add(8 * time.Hour)
				r.On("FindAll", mock.Anything, mock.AnythingOfType("*attendance.FilterParams")).Return([]Attendance{
					{
						ID:           1,
						Date:         time.Now(),
						CheckInTime:  time.Now(),
						CheckOutTime: &checkOut,
						Status:       string(constants.AttendanceStatusPresent),
						Employee:     &user.Employee{FullName: "John", NIK: "001", Department: &department.Department{Name: "IT"}},
						Shift:        &master.Shift{Name: "Morning"},
					},
				}, nil, nil)
				e.On("NewFile").Return(nil)
				e.On("WriteToBuffer", mock.Anything).Return([]byte("fake-excel"), nil)
			},
			wantErr: false,
		},
		{
			name:   "success empty data",
			filter: &FilterParams{},
			setupMocks: func(r *mockRepo, e *mockExcel) {
				r.On("FindAll", mock.Anything, mock.AnythingOfType("*attendance.FilterParams")).Return([]Attendance{}, nil, nil)
				e.On("NewFile").Return(nil)
				e.On("WriteToBuffer", mock.Anything).Return([]byte("fake-excel"), nil)
			},
			wantErr: false,
		},
		{
			name:   "success with leave status",
			filter: &FilterParams{},
			setupMocks: func(r *mockRepo, e *mockExcel) {
				r.On("FindAll", mock.Anything, mock.AnythingOfType("*attendance.FilterParams")).Return([]Attendance{
					{
						ID:          1,
						Date:        time.Now(),
						CheckInTime: time.Now(),
						Status:      "LEAVE",
						Employee:    &user.Employee{FullName: "John", NIK: "001", Department: &department.Department{Name: "IT"}},
						Shift:       &master.Shift{Name: "Morning"},
					},
				}, nil, nil)
				e.On("NewFile").Return(nil)
				e.On("WriteToBuffer", mock.Anything).Return([]byte("fake-excel"), nil)
			},
			wantErr: false,
		},
		{
			name:   "success with sick status",
			filter: &FilterParams{},
			setupMocks: func(r *mockRepo, e *mockExcel) {
				r.On("FindAll", mock.Anything, mock.AnythingOfType("*attendance.FilterParams")).Return([]Attendance{
					{
						ID:          1,
						Date:        time.Now(),
						CheckInTime: time.Now(),
						Status:      "SICK",
						Employee:    &user.Employee{FullName: "John", NIK: "001", Department: &department.Department{Name: "IT"}},
						Shift:       &master.Shift{Name: "Morning"},
					},
				}, nil, nil)
				e.On("NewFile").Return(nil)
				e.On("WriteToBuffer", mock.Anything).Return([]byte("fake-excel"), nil)
			},
			wantErr: false,
		},
		{
			name:   "repo error",
			filter: &FilterParams{},
			setupMocks: func(r *mockRepo, e *mockExcel) {
				r.On("FindAll", mock.Anything, mock.AnythingOfType("*attendance.FilterParams")).Return(nil, nil, errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, _, _, _, _, excel := newTestAttendanceService()
			tt.setupMocks(repo, excel)

			data, err := svc.GenerateExcel(ctx, tt.filter)

			if tt.wantErr {
				require.Error(t, err)
				assert.Nil(t, data)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, data)
			}
		})
	}
}

func TestService_GetDashboardStats(t *testing.T) {
	ctx := testutil.CtxWithTenant(1, 1, false)

	tests := []struct {
		name       string
		setupMocks func(*mockRepo, *mockUserProvider)
		wantErr    bool
		assertFn   func(*testing.T, *DashboardStatResponse)
	}{
		{
			name: "success",
			setupMocks: func(r *mockRepo, u *mockUserProvider) {
				u.On("CountActiveEmployee", mock.Anything).Return(int64(100), nil)
				r.On("CountAttendanceToday", mock.Anything, mock.Anything).Return(int64(80), nil)
				r.On("CountByStatus", mock.Anything, constants.AttendanceStatusLate, mock.Anything).Return(int64(10), nil)
			},
			wantErr: false,
			assertFn: func(t *testing.T, resp *DashboardStatResponse) {
				assert.Equal(t, int64(100), resp.TotalEmployees)
				assert.Equal(t, int64(80), resp.PresentToday)
				assert.Equal(t, int64(10), resp.LateToday)
				assert.Equal(t, int64(20), resp.AbsentToday)
			},
		},
		{
			name: "success all present",
			setupMocks: func(r *mockRepo, u *mockUserProvider) {
				u.On("CountActiveEmployee", mock.Anything).Return(int64(50), nil)
				r.On("CountAttendanceToday", mock.Anything, mock.Anything).Return(int64(50), nil)
				r.On("CountByStatus", mock.Anything, constants.AttendanceStatusLate, mock.Anything).Return(int64(0), nil)
			},
			wantErr: false,
			assertFn: func(t *testing.T, resp *DashboardStatResponse) {
				assert.Equal(t, int64(50), resp.TotalEmployees)
				assert.Equal(t, int64(50), resp.PresentToday)
				assert.Equal(t, int64(0), resp.AbsentToday)
			},
		},
		{
			name: "count active employee error",
			setupMocks: func(r *mockRepo, u *mockUserProvider) {
				u.On("CountActiveEmployee", mock.Anything).Return(int64(0), errors.New("db error"))
			},
			wantErr: true,
		},
		{
			name: "count attendance today error",
			setupMocks: func(r *mockRepo, u *mockUserProvider) {
				u.On("CountActiveEmployee", mock.Anything).Return(int64(100), nil)
				r.On("CountAttendanceToday", mock.Anything, mock.Anything).Return(int64(0), errors.New("db error"))
			},
			wantErr: true,
		},
		{
			name: "count late error",
			setupMocks: func(r *mockRepo, u *mockUserProvider) {
				u.On("CountActiveEmployee", mock.Anything).Return(int64(100), nil)
				r.On("CountAttendanceToday", mock.Anything, mock.Anything).Return(int64(80), nil)
				r.On("CountByStatus", mock.Anything, constants.AttendanceStatusLate, mock.Anything).Return(int64(0), errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo, userProv, _, _, _, _ := newTestAttendanceService()
			tt.setupMocks(repo, userProv)

			resp, err := svc.GetDashboardStats(ctx)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				if tt.assertFn != nil {
					tt.assertFn(t, resp)
				}
			}
		})
	}
}
