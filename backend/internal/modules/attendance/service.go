package attendance

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"hris-backend/internal/modules/user"
	"hris-backend/pkg/constants"
	"hris-backend/pkg/response"
	"hris-backend/pkg/utils"
	"time"

	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type Service interface {
	Clock(ctx context.Context, userID uint, req *ClockRequest) (*AttendanceResponse, error)
	GetTodayStatus(ctx context.Context, userID uint) (*TodayStatusResponse, error)
	GetMyHistory(ctx context.Context, userID uint, month, year, page, limit int) ([]Attendance, *response.Meta, error)
	GetAllRecap(ctx context.Context, filter *FilterParams) ([]RecapResponse, *response.Meta, error)
	GenerateExcel(ctx context.Context, filter *FilterParams) (*excelize.File, error)
	GetDashboardStats(ctx context.Context) (*DashboardStatResponse, error)
}

type service struct {
	repo         Repository
	userRepo     user.Repository
	storage      StorageProvider
	geocodeQueue chan<- GeocodeJob
}

func NewService(repo Repository, userRepo user.Repository, storage StorageProvider, geocodeQueue chan<- GeocodeJob) Service {
	return &service{repo, userRepo, storage, geocodeQueue}
}

func (s *service) Clock(ctx context.Context, userID uint, req *ClockRequest) (*AttendanceResponse, error) {
	u, err := s.userRepo.FindByID(userID)
	if err != nil || u.Employee == nil {
		return nil, errors.New("employee data not found")
	}

	employee := u.Employee

	if employee.Shift == nil {
		return nil, errors.New("employee shift not assigned")
	}

	imgBytes, err := utils.DecodeBase64Image(req.ImageBase64)
	if err != nil {
		return nil, errors.New("invalid image")
	}
	imageReader := bytes.NewReader(imgBytes)

	// set address temporary
	tempAddress := fmt.Sprintf("Processing location... (%f, %f)", req.Latitude, req.Longitude)

	now := time.Now()
	todayString := now.Format(constants.DefaultTimeFormat)
	fileName := fmt.Sprintf("attendance/%d/%s-%d.jpg", employee.ID, todayString, now.Unix())

	shiftStartToday, err := combineDateAndTime(now, employee.Shift.StartTime)
	if err != nil {
		return nil, errors.New("invalid shift time configuration")
	}

	expectedTime := time.Date(
		now.Year(), now.Month(), now.Day(),
		shiftStartToday.Hour(), shiftStartToday.Minute(), shiftStartToday.Second(), 0,
		now.Location(),
	)

	lateMinute := 0
	if now.After(expectedTime) {
		diff := now.Sub(expectedTime)
		lateMinute = int(diff.Minutes())
	}

	earliersAllowed := shiftStartToday.Add(-2 * time.Hour)
	if now.Before(earliersAllowed) {
		return nil, errors.New("cannot check-in, too early")
	}

	todayAtt, err := s.repo.GetTodayAttendance(employee.ID)
	// if today no data, its check-in of that employee
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// calculate status is LATE or PRESENT
		lateThreshold := shiftStartToday.Add(15 * time.Minute)
		status := string(constants.AttendanceStatusPresent)

		// compare current time with shift time, if more than late threshold, status will changed to LATE
		if now.After(lateThreshold) {
			status = string(constants.AttendanceStatusLate)
		}

		imgUrl, err := s.storage.UploadFileByte(ctx, fmt.Sprintf("in-%s", fileName), imageReader, int64(len(imgBytes)), "image/jpg")
		if err != nil {
			return nil, err
		}

		newAtt := &Attendance{
			EmployeeID:         employee.ID,
			ShiftID:            employee.ShiftID,
			Date:               time.Now(),
			CheckInTime:        now,
			CheckInLat:         req.Latitude,
			CheckInLong:        req.Longitude,
			CheckInImageURL:    imgUrl,
			CheckInAddress:     tempAddress,
			Status:             status,
			Notes:              req.Notes,
			LateDurationMinute: lateMinute,
			IsSuspicious:       false,
		}

		if err := s.repo.Create(newAtt); err != nil {
			return nil, err
		}

		// process this attendance (check-in) to geocode worker queue
		s.geocodeQueue <- GeocodeJob{
			AttendanceID: newAtt.ID,
			Latitude:     req.Latitude,
			Longitude:    req.Longitude,
			IsCheckout:   false,
		}

		return &AttendanceResponse{
			Type:    string(constants.AttendanceTypeCheckIn),
			Status:  status,
			Time:    now,
			Message: "Check-in succesful",
		}, nil
	}

	// if today already have attendance, but the checkout time is still null, its checkout of that employee
	if todayAtt != nil && todayAtt.CheckOutTime == nil {

		// Calculate Teleportation Check (to detect distance between location check-in & check-out employee, will get mark if suspicious )
		distanceMeters := utils.CalculateDistance(todayAtt.CheckInLat, todayAtt.CheckInLong, req.Latitude, req.Longitude)
		distanceKm := distanceMeters / 1000.0

		durationHours := time.Since(todayAtt.CheckInTime).Hours()

		isSuspicious := false
		notes := ""

		if durationHours > 0.01 {
			speedKmH := distanceKm / durationHours

			if speedKmH > 200 && distanceKm > 2.0 {
				isSuspicious = true
				notes = fmt.Sprintf("[SUSPICIOUS] Speed %.2f km/h detected. Teleportation check failed.", speedKmH)
			}
		}

		imgUrl, err := s.storage.UploadFileByte(ctx, fmt.Sprintf("out-%s", fileName), imageReader, int64(len(imgBytes)), "image/jpg")
		if err != nil {
			return nil, err
		}

		todayAtt.CheckOutTime = &now
		todayAtt.CheckOutLat = &req.Latitude
		todayAtt.CheckOutLong = &req.Longitude
		todayAtt.CheckOutImageURL = &imgUrl
		todayAtt.CheckOutAddress = &tempAddress

		if isSuspicious {
			todayAtt.IsSuspicious = true
			todayAtt.Notes = todayAtt.Notes + " " + notes
		}

		if err := s.repo.Update(todayAtt); err != nil {
			return nil, err
		}

		// process this attendance (check-out) to geocode worker queue
		s.geocodeQueue <- GeocodeJob{
			AttendanceID: todayAtt.ID,
			Latitude:     req.Latitude,
			Longitude:    req.Longitude,
			IsCheckout:   true,
		}

		return &AttendanceResponse{
			Type:    string(constants.AttendanceTypeCheckOut),
			Status:  todayAtt.Status,
			Time:    now,
			Message: "Check-out successful",
		}, nil
	}

	return nil, errors.New("you have already completed attendance for today")
}

func (s *service) GetTodayStatus(ctx context.Context, userID uint) (*TodayStatusResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil || user.Employee == nil {
		return nil, errors.New("employee not found")
	}

	att, err := s.repo.GetTodayAttendance(user.Employee.ID)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &TodayStatusResponse{
			Status: string(constants.AttendanceStatusAbsent),
			Type:   string(constants.AttendanceTypeNone),
		}, nil
	}
	if err != nil {
		return nil, err
	}

	if att.CheckOutTime == nil {
		return &TodayStatusResponse{
			Status:      att.Status,
			Type:        string(constants.AttendanceTypeCheckIn),
			CheckInTime: &att.CheckInTime,
		}, nil
	}

	duration := att.CheckOutTime.Sub(att.CheckInTime).String()

	return &TodayStatusResponse{
		Status:       att.Status,
		Type:         string(constants.AttendanceTypeCompleted),
		CheckInTime:  &att.CheckInTime,
		CheckOutTime: att.CheckOutTime,
		WorkDuration: duration,
	}, nil
}

func (s *service) GetMyHistory(ctx context.Context, userID uint, month, year, page, limit int) ([]Attendance, *response.Meta, error) {
	u, err := s.userRepo.FindByID(userID)
	if err != nil || u.Employee == nil {
		return nil, nil, errors.New("employee not found")
	}

	logs, total, err := s.repo.GetHistory(u.Employee.ID, month, year, page, limit)
	if err != nil {
		return nil, nil, err
	}

	meta := response.NewMeta(page, limit, total)

	return logs, meta, nil
}

func (s *service) GetAllRecap(ctx context.Context, filter *FilterParams) ([]RecapResponse, *response.Meta, error) {
	data, total, err := s.repo.FindAll(filter)
	if err != nil {
		return nil, nil, err
	}

	if len(data) == 0 {
		return []RecapResponse{}, nil, nil
	}

	var result []RecapResponse
	for _, item := range data {
		cIn := "-"
		cOut := "-"
		duration := "-"

		if !item.CheckInTime.IsZero() {
			cIn = item.CheckInTime.Format(constants.ShiftHourFormat)
		}
		if item.CheckOutTime != nil {
			cOut = item.CheckOutTime.Format(constants.ShiftHourFormat)
			dur := item.CheckOutTime.Sub(item.CheckInTime)
			duration = fmt.Sprintf("%.1f Hours", dur.Hours())
		}

		result = append(result, RecapResponse{
			ID:           item.ID,
			Date:         item.Date.Format(constants.DefaultTimeFormat),
			EmployeeName: item.Employee.FullName,
			NIK:          item.Employee.NIK,
			Department:   item.Employee.Department.Name,
			Shift:        item.Shift.Name,
			CheckInTime:  cIn,
			CheckOutTime: cOut,
			Status:       item.Status,
			WorkDuration: duration,
		})
	}

	meta := response.NewMeta(filter.Page, filter.Limit, total)
	return result, meta, nil
}

func (s *service) GenerateExcel(ctx context.Context, filter *FilterParams) (*excelize.File, error) {
	filter.Limit = 0
	data, _, err := s.repo.FindAll(filter)
	if err != nil {
		return nil, err
	}

	f := excelize.NewFile()
	sheet := "Sheet1"
	f.SetSheetName("Sheet1", "Attendance Recap")
	sheet = "Attendance Recap"

	styleLeave, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#DDEBF7"}, Pattern: 1},
		Font: &excelize.Font{Color: "#1F4E78", Bold: true},
	})
	styleSick, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#FCE4D6"}, Pattern: 1},
		Font: &excelize.Font{Color: "#833C0C", Bold: true},
	})
	styleLate, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Color: "#FF0000", Bold: true},
	})

	headers := []string{"Date", "NIK", "Name", "Department", "Shift", "Check In", "Check Out", "Status", "Notes"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
		f.SetCellStyle(sheet, cell, cell, styleLeave)
	}

	for i, item := range data {
		row := i + 2

		cIn := ""
		if !item.CheckInTime.IsZero() {
			cIn = item.CheckInTime.Format(constants.ShiftHourFormat)
		}
		cOut := ""
		if item.CheckOutTime != nil && !item.CheckOutTime.IsZero() {
			cOut = item.CheckOutTime.Format(constants.ShiftHourFormat)
		}

		displayStatus := item.Status
		var styleID int

		switch item.Status {
		case "LEAVE":
			displayStatus = "CUTI"
			styleID = styleLeave
		case "SICK":
			displayStatus = "SAKIT"
			styleID = styleSick
		case "LATE":
			displayStatus = "TERLAMBAT"
			styleID = styleLate
		case "ABSENT":
			displayStatus = "ALPHA"
			styleID = styleLate
		case "PRESENT":
			displayStatus = "HADIR"
		}

		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), item.Date.Format("2006-01-02"))
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), item.Employee.NIK)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), item.Employee.FullName)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), item.Employee.Department.Name)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), item.Shift.Name)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), cIn)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", row), cOut)

		cellStatus := fmt.Sprintf("H%d", row)
		f.SetCellValue(sheet, cellStatus, displayStatus)

		if styleID != 0 {
			f.SetCellStyle(sheet, cellStatus, cellStatus, styleID)
		}

		f.SetCellValue(sheet, fmt.Sprintf("I%d", row), item.Notes)
	}

	f.SetColWidth(sheet, "A", "A", 12)
	f.SetColWidth(sheet, "B", "B", 15)
	f.SetColWidth(sheet, "C", "C", 25)
	f.SetColWidth(sheet, "D", "E", 15)
	f.SetColWidth(sheet, "H", "H", 12)
	f.SetColWidth(sheet, "I", "I", 30)

	return f, nil
}

func (s *service) GetDashboardStats(ctx context.Context) (*DashboardStatResponse, error) {
	todayDate := time.Now().Format(constants.DefaultTimeFormat)

	totalActiveEmployee, err := s.userRepo.CountActiveEmployee()
	if err != nil {
		return nil, err
	}

	totalPresentToday, err := s.repo.CountAttendanceToday(todayDate)
	if err != nil {
		return nil, err
	}

	totalLateToday, err := s.repo.CountByStatus(constants.AttendanceStatusLate, todayDate)
	if err != nil {
		return nil, err
	}

	stats := &DashboardStatResponse{
		TotalEmployees: totalActiveEmployee,
		PresentToday:   totalPresentToday,
		LateToday:      totalLateToday,
	}

	if stats.TotalEmployees >= stats.PresentToday {
		stats.AbsentToday = stats.TotalEmployees - stats.PresentToday
	} else {
		stats.AbsentToday = 0
	}

	return stats, nil
}

func combineDateAndTime(date time.Time, timeStr string) (time.Time, error) {
	parsedTime, err := time.Parse(constants.AttendanceTimeFormat, timeStr)
	if err != nil {
		return time.Time{}, err
	}

	fullTime := time.Date(
		date.Year(), date.Month(), date.Day(),
		parsedTime.Hour(), parsedTime.Minute(), parsedTime.Second(), 0,
		date.Location(),
	)

	return fullTime, nil
}
