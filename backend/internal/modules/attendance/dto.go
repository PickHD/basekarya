package attendance

import "time"

type ClockRequest struct {
	Latitude    float64 `json:"latitude" validate:"required"`
	Longitude   float64 `json:"longitude" validate:"required"`
	ImageBase64 string  `json:"image_base64" validate:"required"`
	Address     string  `json:"address"`
	Notes       string  `json:"notes"`
}

type AttendanceResponse struct {
	Type    string    `json:"type"`
	Status  string    `json:"status"`
	Time    time.Time `json:"time"`
	Message string    `json:"message"`
}

type TodayStatusResponse struct {
	Status       string     `json:"status"`
	Type         string     `json:"type"`
	CheckInTime  *time.Time `json:"check_in_time"`
	CheckOutTime *time.Time `json:"check_out_time"`
	WorkDuration string     `json:"work_duration"`
}

type FilterParams struct {
	StartDate    string
	EndDate      string
	DepartmentID uint
	Search       string
	Page         int
	Limit        int
	Timezone     string
}

type RecapResponse struct {
	ID           uint   `json:"id"`
	Date         string `json:"date"`
	EmployeeName string `json:"employee_name"`
	NIK          string `json:"nik"`
	Department   string `json:"department"`
	Shift        string `json:"shift"`
	CheckInTime  string `json:"check_in_time"`
	CheckOutTime string `json:"check_out_time"`
	Status       string `json:"status"`
	WorkDuration string `json:"work_duration"`
}

type DashboardStatResponse struct {
	TotalEmployees int64 `json:"total_employees"`
	PresentToday   int64 `json:"present_today"`
	LateToday      int64 `json:"late_today"`
	AbsentToday    int64 `json:"absent_today"`
}
