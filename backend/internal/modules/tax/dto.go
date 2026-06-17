package tax

type PPh21Result struct {
	TERCategory  string  `json:"ter_category"`
	PTKPCode     string  `json:"ptkp_code"`
	GrossMonthly float64 `json:"gross_monthly"`
	TERRate      float64 `json:"ter_rate"`
	MonthlyPPh21 float64 `json:"monthly_pph21"`
}

type AnnualSettlement struct {
	GrossAnnual  float64 `json:"gross_annual"`
	BiayaJabatan float64 `json:"biaya_jabatan"`
	PTKP         float64 `json:"ptkp"`
	PKP          float64 `json:"pkp"`
	TaxPayable   float64 `json:"tax_payable"`
	TERPaidYTD   float64 `json:"ter_paid_ytd"`
	Delta        float64 `json:"delta"`
}

type Form1721A1 struct {
	Year             int               `json:"year"`
	EmployeeName     string            `json:"employee_name"`
	NPWP             string            `json:"npwp"`
	PTKPCode         string            `json:"ptkp_code"`
	Settlement       AnnualSettlement  `json:"settlement"`
	MonthlyBreakdown []MonthlyTaxEntry `json:"monthly_breakdown"`
}

type MonthlyTaxEntry struct {
	Month       int     `json:"month"`
	GrossIncome float64 `json:"gross_income"`
	PPh21Paid   float64 `json:"pph21_paid"`
}

type TERBracketFilter struct {
	Category      string
	EffectiveDate string
	Page          int
	Limit         int
}

type TERBracketRequest struct {
	Category         string  `json:"category" validate:"required"`
	BracketNumber    int     `json:"bracket_number" validate:"required"`
	MinMonthlySalary float64 `json:"min_monthly_salary" validate:"required"`
	Rate             float64 `json:"rate" validate:"required"`
	EffectiveFrom    string  `json:"effective_from" validate:"required"`
	EffectiveUntil   *string `json:"effective_until"`
}

type PTKPConfigRequest struct {
	Code          string  `json:"code" validate:"required"`
	AnnualAmount  float64 `json:"annual_amount" validate:"required"`
	EffectiveYear int     `json:"effective_year" validate:"required"`
}
