package bpjs

type BPJSComponent struct {
	Type            string  `json:"type"`
	Code            string  `json:"code"`
	EmployeeAmount  float64 `json:"employee_amount"`
	EmployerAmount  float64 `json:"employer_amount"`
	IsEmployerBorne bool    `json:"is_employer_borne"`
	MaxCap          float64 `json:"max_cap,omitempty"`
}

type BPJSRateConfigRequest struct {
	Type              string   `json:"type" validate:"required"`
	EmployeeRate      float64  `json:"employee_rate" validate:"required"`
	EmployerRate      float64  `json:"employer_rate" validate:"required"`
	MaxSalaryCap      *float64 `json:"max_salary_cap"`
	IndustryRiskLevel *string  `json:"industry_risk_level"`
	IsActive          bool     `json:"is_active"`
	EffectiveFrom     string   `json:"effective_from" validate:"required"`
	EffectiveUntil    *string  `json:"effective_until"`
}

type BPJSRateConfigFilter struct {
	Type     string
	IsActive *bool
	Page     int
	Limit    int
}
