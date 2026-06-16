package tax

import (
	"basekarya-backend/pkg/constants"
	"context"
	"errors"
	"fmt"
	"math"
	"time"
)

type Service interface {
	CalculateTER(ctx context.Context, grossMonthlyIncome float64, maritalStatus constants.MaritalStatus, dependentsCount int) (*PPh21Result, error)
	ReconcileAnnual(ctx context.Context, year int, ptkpCode string, grossAnnual float64, monthlyDetails []MonthlyTaxEntry) (*AnnualSettlement, error)
	CreateTERBracket(ctx context.Context, req *TERBracketRequest) error
	GetTERBracketByID(ctx context.Context, id uint) (*TERBracket, error)
	UpdateTERBracket(ctx context.Context, id uint, req *TERBracketRequest) error
	DeleteTERBracket(ctx context.Context, id uint) error
	ListTERBrackets(ctx context.Context, filter TERBracketFilter) ([]TERBracket, int64, error)
	CreatePTKPConfig(ctx context.Context, req *PTKPConfigRequest) error
	GetPTKPConfigByID(ctx context.Context, id uint) (*PTKPConfig, error)
	UpdatePTKPConfig(ctx context.Context, id uint, req *PTKPConfigRequest) error
	DeletePTKPConfig(ctx context.Context, id uint) error
	ListPTKPConfigs(ctx context.Context, year int) ([]PTKPConfig, int64, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) CalculateTER(ctx context.Context, grossMonthlyIncome float64, maritalStatus constants.MaritalStatus, dependentsCount int) (*PPh21Result, error) {
	ptkpCode := derivePTKPCode(maritalStatus, dependentsCount)
	terCategory := deriveTERCategory(ptkpCode)

	brackets, err := s.repo.FindTERBrackets(ctx, terCategory, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to find TER brackets: %w", err)
	}
	if len(brackets) == 0 {
		return nil, errors.New("no TER brackets found for category " + terCategory)
	}

	rate := findTERRate(brackets, grossMonthlyIncome)
	monthlyPPh21 := math.Round(grossMonthlyIncome * rate)

	return &PPh21Result{
		TERCategory:  terCategory,
		PTKPCode:     ptkpCode,
		GrossMonthly: grossMonthlyIncome,
		TERRate:      rate,
		MonthlyPPh21: monthlyPPh21,
	}, nil
}

func (s *service) ReconcileAnnual(ctx context.Context, year int, ptkpCode string, grossAnnual float64, monthlyDetails []MonthlyTaxEntry) (*AnnualSettlement, error) {
	ptkpConfigs, err := s.repo.FindPTKPByYear(ctx, year)
	if err != nil {
		return nil, fmt.Errorf("failed to find PTKP configs: %w", err)
	}
	var ptkpAnnual float64
	for _, p := range ptkpConfigs {
		if p.Code == ptkpCode {
			ptkpAnnual = p.AnnualAmount
			break
		}
	}
	biayaJabatanAnnual := math.Min(grossAnnual*0.05, 6000000)
	pkp := math.Max(0, grossAnnual-biayaJabatanAnnual-ptkpAnnual)
	pkp = math.Floor(pkp/1000) * 1000
	taxPayable := calculateProgressiveTax(pkp)
	var terPaidYTD float64
	for _, m := range monthlyDetails {
		terPaidYTD += m.PPh21Paid
	}
	delta := taxPayable - terPaidYTD
	return &AnnualSettlement{
		GrossAnnual:  grossAnnual,
		BiayaJabatan: biayaJabatanAnnual,
		PTKP:         ptkpAnnual,
		PKP:          pkp,
		TaxPayable:   taxPayable,
		TERPaidYTD:   terPaidYTD,
		Delta:        delta,
	}, nil
}

func (s *service) CreateTERBracket(ctx context.Context, req *TERBracketRequest) error {
	from, err := time.Parse("2006-01-02", req.EffectiveFrom)
	if err != nil {
		return fmt.Errorf("invalid effective_from: %w", err)
	}
	b := &TERBracket{
		Category: req.Category, BracketNumber: req.BracketNumber,
		MinMonthlySalary: req.MinMonthlySalary, Rate: req.Rate,
		EffectiveFrom: from,
	}
	if req.EffectiveUntil != nil {
		until, err := time.Parse("2006-01-02", *req.EffectiveUntil)
		if err != nil {
			return fmt.Errorf("invalid effective_until: %w", err)
		}
		b.EffectiveUntil = &until
	}
	return s.repo.CreateTERBracket(ctx, b)
}

func (s *service) GetTERBracketByID(ctx context.Context, id uint) (*TERBracket, error) {
	return s.repo.FindTERBracketByID(ctx, id)
}

func (s *service) UpdateTERBracket(ctx context.Context, id uint, req *TERBracketRequest) error {
	bracket, err := s.repo.FindTERBracketByID(ctx, id)
	if err != nil {
		return err
	}
	bracket.Category = req.Category
	bracket.BracketNumber = req.BracketNumber
	bracket.MinMonthlySalary = req.MinMonthlySalary
	bracket.Rate = req.Rate
	from, err := time.Parse("2006-01-02", req.EffectiveFrom)
	if err != nil {
		return fmt.Errorf("invalid effective_from: %w", err)
	}
	bracket.EffectiveFrom = from
	if req.EffectiveUntil != nil {
		until, err := time.Parse("2006-01-02", *req.EffectiveUntil)
		if err != nil {
			return fmt.Errorf("invalid effective_until: %w", err)
		}
		bracket.EffectiveUntil = &until
	} else {
		bracket.EffectiveUntil = nil
	}
	return s.repo.UpdateTERBracket(ctx, bracket)
}

func (s *service) DeleteTERBracket(ctx context.Context, id uint) error {
	return s.repo.DeleteTERBracket(ctx, id)
}

func (s *service) ListTERBrackets(ctx context.Context, filter TERBracketFilter) ([]TERBracket, int64, error) {
	return s.repo.ListTERBrackets(ctx, filter)
}

func (s *service) CreatePTKPConfig(ctx context.Context, req *PTKPConfigRequest) error {
	ptkp := &PTKPConfig{Code: req.Code, AnnualAmount: req.AnnualAmount, EffectiveYear: req.EffectiveYear}
	return s.repo.CreatePTKPConfig(ctx, ptkp)
}

func (s *service) GetPTKPConfigByID(ctx context.Context, id uint) (*PTKPConfig, error) {
	return s.repo.FindPTKPConfigByID(ctx, id)
}

func (s *service) UpdatePTKPConfig(ctx context.Context, id uint, req *PTKPConfigRequest) error {
	ptkp, err := s.repo.FindPTKPConfigByID(ctx, id)
	if err != nil {
		return err
	}
	ptkp.Code = req.Code
	ptkp.AnnualAmount = req.AnnualAmount
	ptkp.EffectiveYear = req.EffectiveYear
	return s.repo.UpdatePTKPConfig(ctx, ptkp)
}

func (s *service) DeletePTKPConfig(ctx context.Context, id uint) error {
	return s.repo.DeletePTKPConfig(ctx, id)
}

func (s *service) ListPTKPConfigs(ctx context.Context, year int) ([]PTKPConfig, int64, error) {
	return s.repo.ListPTKPConfigs(ctx, year)
}

func derivePTKPCode(maritalStatus constants.MaritalStatus, dependents int) string {
	capped := dependents
	if capped > 3 {
		capped = 3
	}
	return fmt.Sprintf("%s/%d", maritalStatus, capped)
}

func deriveTERCategory(ptkpCode string) string {
	switch ptkpCode {
	case "TK/0", "TK/1", "K/0":
		return "A"
	case "TK/2", "TK/3", "K/1", "K/2":
		return "B"
	case "K/3":
		return "C"
	default:
		return "A"
	}
}

func findTERRate(brackets []TERBracket, grossIncome float64) float64 {
	var matchedRate float64
	for _, b := range brackets {
		if grossIncome >= b.MinMonthlySalary {
			matchedRate = b.Rate
		} else {
			break
		}
	}
	return matchedRate
}

func calculateBiayaJabatanMonthly(grossMonthly float64) float64 {
	return math.Min(grossMonthly*0.05, 500000)
}

func calculateProgressiveTax(pkp float64) float64 {
	if pkp <= 0 {
		return 0
	}
	type bracket struct {
		limit, rate float64
	}
	brackets := []bracket{
		{60000000, 0.05},
		{250000000, 0.15},
		{500000000, 0.25},
		{5000000000, 0.30},
	}
	var totalTax, prevLimit float64
	remaining := pkp
	for _, b := range brackets {
		if remaining <= 0 {
			break
		}
		taxable := math.Min(remaining, b.limit-prevLimit)
		totalTax += taxable * b.rate
		remaining -= taxable
		prevLimit = b.limit
	}
	if remaining > 0 {
		totalTax += remaining * 0.35
	}
	return math.Floor(totalTax)
}
