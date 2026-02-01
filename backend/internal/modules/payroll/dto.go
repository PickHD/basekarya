package payroll

import "time"

type GenerateRequest struct {
	Month int `json:"month" validate:"required,min=1,max=12"`
	Year  int `json:"year" validate:"required,min=2024"`
}

type GenerateResponse struct {
	SuccessCount int `json:"success_count"`
	Month        int `json:"month"`
	Year         int `json:"year"`
}

type PayrollFilter struct {
	Page    int    `json:"page"`
	Limit   int    `json:"limit"`
	Month   int    `json:"month"`
	Year    int    `json:"year"`
	Keyword string `json:"keyword"`
}

type PayrollListResponse struct {
	ID           uint      `json:"id"`
	EmployeeName string    `json:"employee_name"`
	EmployeeNIK  string    `json:"employee_nik"`
	PeriodDate   string    `json:"period_date"`
	NetSalary    float64   `json:"net_salary"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}
