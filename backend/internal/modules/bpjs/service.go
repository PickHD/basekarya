package bpjs

import (
	"basekarya-backend/pkg/utils"
	"context"
	"fmt"
	"math"
	"time"
)

type Service interface {
	CalculateAll(ctx context.Context, grossMonthlyIncome float64) ([]BPJSComponent, error)
	Create(ctx context.Context, req *BPJSRateConfigRequest) error
	GetByID(ctx context.Context, id uint) (*BPJSRateConfig, error)
	Update(ctx context.Context, id uint, req *BPJSRateConfigRequest) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, filter BPJSRateConfigFilter) ([]BPJSRateConfig, int64, error)
}

type service struct{ repo Repository }

func NewService(repo Repository) Service { return &service{repo} }

func (s *service) CalculateAll(ctx context.Context, grossMonthlyIncome float64) ([]BPJSComponent, error) {
	configs, err := s.repo.FindAllActive(ctx, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to get BPJS configs: %w", err)
	}

	seen := make(map[string]bool)
	var deduped []BPJSRateConfig
	for i := len(configs) - 1; i >= 0; i-- {
		if !seen[configs[i].Type] {
			seen[configs[i].Type] = true
			deduped = append([]BPJSRateConfig{configs[i]}, deduped...)
		}
	}

	var components []BPJSComponent
	for _, cfg := range deduped {
		if !cfg.IsActive {
			continue
		}
		effectiveSalary := grossMonthlyIncome
		if cfg.MaxSalaryCap != nil && *cfg.MaxSalaryCap > 0 {
			effectiveSalary = math.Min(grossMonthlyIncome, *cfg.MaxSalaryCap)
		}
		employeeAmount := math.Round(effectiveSalary * cfg.EmployeeRate)
		employerAmount := math.Round(effectiveSalary * cfg.EmployerRate)
		code := "BPJS_" + cfg.Type
		if cfg.EmployeeRate > 0 {
			components = append(components, BPJSComponent{
				Type:            cfg.Type,
				Code:            code + "_E",
				EmployeeAmount:  employeeAmount,
				EmployerAmount:  0,
				IsEmployerBorne: false,
				MaxCap:          effectiveSalary,
			})
		}
		if cfg.EmployerRate > 0 {
			components = append(components, BPJSComponent{
				Type:            cfg.Type,
				Code:            code + "_R",
				EmployeeAmount:  0,
				EmployerAmount:  employerAmount,
				IsEmployerBorne: true,
				MaxCap:          effectiveSalary,
			})
		}
	}
	return components, nil
}

func (s *service) Create(ctx context.Context, req *BPJSRateConfigRequest) error {
	from, err := time.Parse("2006-01-02", req.EffectiveFrom)
	if err != nil {
		return fmt.Errorf("invalid effective_from: %w", err)
	}
	companyID := utils.GetCompanyIDFromCtx(ctx)
	cfg := &BPJSRateConfig{
		Type:              req.Type,
		EmployeeRate:      req.EmployeeRate,
		EmployerRate:      req.EmployerRate,
		MaxSalaryCap:      req.MaxSalaryCap,
		IndustryRiskLevel: req.IndustryRiskLevel,
		IsActive:          req.IsActive,
		EffectiveFrom:     from,
	}
	if companyID > 0 {
		cfg.CompanyID = &companyID
	}
	if req.EffectiveUntil != nil {
		until, err := time.Parse("2006-01-02", *req.EffectiveUntil)
		if err != nil {
			return fmt.Errorf("invalid effective_until: %w", err)
		}
		cfg.EffectiveUntil = &until
	}
	return s.repo.Create(ctx, cfg)
}

func (s *service) GetByID(ctx context.Context, id uint) (*BPJSRateConfig, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *service) Update(ctx context.Context, id uint, req *BPJSRateConfigRequest) error {
	cfg, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	cfg.Type = req.Type
	cfg.EmployeeRate = req.EmployeeRate
	cfg.EmployerRate = req.EmployerRate
	cfg.MaxSalaryCap = req.MaxSalaryCap
	cfg.IndustryRiskLevel = req.IndustryRiskLevel
	cfg.IsActive = req.IsActive
	from, err := time.Parse("2006-01-02", req.EffectiveFrom)
	if err != nil {
		return fmt.Errorf("invalid effective_from: %w", err)
	}
	cfg.EffectiveFrom = from
	if req.EffectiveUntil != nil {
		until, err := time.Parse("2006-01-02", *req.EffectiveUntil)
		if err != nil {
			return fmt.Errorf("invalid effective_until: %w", err)
		}
		cfg.EffectiveUntil = &until
	} else {
		cfg.EffectiveUntil = nil
	}
	return s.repo.Update(ctx, cfg)
}

func (s *service) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

func (s *service) List(ctx context.Context, filter BPJSRateConfigFilter) ([]BPJSRateConfig, int64, error) {
	configs, _, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	seen := make(map[string]bool)
	var deduped []BPJSRateConfig
	for i := len(configs) - 1; i >= 0; i-- {
		if !seen[configs[i].Type] {
			seen[configs[i].Type] = true
			deduped = append([]BPJSRateConfig{configs[i]}, deduped...)
		}
	}

	return deduped, int64(len(deduped)), nil
}
