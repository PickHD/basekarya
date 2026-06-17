package seeder

import (
	"errors"
	"time"

	"basekarya-backend/internal/modules/bpjs"
	"basekarya-backend/internal/modules/tax"

	"gorm.io/gorm"
)

func seedTaxAndBPJSData(tx *gorm.DB) error {
	if err := seedPTKPConfigs(tx); err != nil {
		return err
	}
	if err := seedTERBrackets(tx); err != nil {
		return err
	}
	if err := seedBPJSRateConfigs(tx); err != nil {
		return err
	}
	return nil
}

func seedPTKPConfigs(tx *gorm.DB) error {
	ptkpConfigs := []tax.PTKPConfig{
		{Code: "TK/0", AnnualAmount: 54000000, EffectiveYear: 2026},
		{Code: "TK/1", AnnualAmount: 58500000, EffectiveYear: 2026},
		{Code: "TK/2", AnnualAmount: 63000000, EffectiveYear: 2026},
		{Code: "TK/3", AnnualAmount: 67500000, EffectiveYear: 2026},
		{Code: "K/0", AnnualAmount: 58500000, EffectiveYear: 2026},
		{Code: "K/1", AnnualAmount: 63000000, EffectiveYear: 2026},
		{Code: "K/2", AnnualAmount: 67500000, EffectiveYear: 2026},
		{Code: "K/3", AnnualAmount: 72000000, EffectiveYear: 2026},
	}

	for _, p := range ptkpConfigs {
		var existing tax.PTKPConfig
		err := tx.Where("code = ? AND effective_year = ?", p.Code, p.EffectiveYear).First(&existing).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := tx.Create(&p).Error; err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
	}
	return nil
}

func seedTERBrackets(tx *gorm.DB) error {
	efFrom := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	catA := []tax.TERBracket{
		{Category: "A", BracketNumber: 1, MinMonthlySalary: 0, Rate: 0.0, EffectiveFrom: efFrom},
		{Category: "A", BracketNumber: 2, MinMonthlySalary: 5400000, Rate: 0.0025, EffectiveFrom: efFrom},
		{Category: "A", BracketNumber: 3, MinMonthlySalary: 5650000, Rate: 0.0050, EffectiveFrom: efFrom},
		{Category: "A", BracketNumber: 4, MinMonthlySalary: 5950000, Rate: 0.0075, EffectiveFrom: efFrom},
		{Category: "A", BracketNumber: 5, MinMonthlySalary: 6300000, Rate: 0.0100, EffectiveFrom: efFrom},
		{Category: "A", BracketNumber: 6, MinMonthlySalary: 6750000, Rate: 0.0125, EffectiveFrom: efFrom},
		{Category: "A", BracketNumber: 7, MinMonthlySalary: 7500000, Rate: 0.0150, EffectiveFrom: efFrom},
		{Category: "A", BracketNumber: 8, MinMonthlySalary: 8550000, Rate: 0.0175, EffectiveFrom: efFrom},
		{Category: "A", BracketNumber: 9, MinMonthlySalary: 9650000, Rate: 0.0200, EffectiveFrom: efFrom},
		{Category: "A", BracketNumber: 10, MinMonthlySalary: 10050000, Rate: 0.0225, EffectiveFrom: efFrom},
		{Category: "A", BracketNumber: 11, MinMonthlySalary: 10350000, Rate: 0.0250, EffectiveFrom: efFrom},
		{Category: "A", BracketNumber: 12, MinMonthlySalary: 10700000, Rate: 0.0300, EffectiveFrom: efFrom},
		{Category: "A", BracketNumber: 13, MinMonthlySalary: 11050000, Rate: 0.0350, EffectiveFrom: efFrom},
		{Category: "A", BracketNumber: 14, MinMonthlySalary: 11600000, Rate: 0.0400, EffectiveFrom: efFrom},
		{Category: "A", BracketNumber: 15, MinMonthlySalary: 12500000, Rate: 0.0500, EffectiveFrom: efFrom},
		{Category: "A", BracketNumber: 16, MinMonthlySalary: 13750000, Rate: 0.0600, EffectiveFrom: efFrom},
		{Category: "A", BracketNumber: 17, MinMonthlySalary: 15100000, Rate: 0.0700, EffectiveFrom: efFrom},
		{Category: "A", BracketNumber: 18, MinMonthlySalary: 16950000, Rate: 0.0800, EffectiveFrom: efFrom},
		{Category: "A", BracketNumber: 19, MinMonthlySalary: 19750000, Rate: 0.0900, EffectiveFrom: efFrom},
		{Category: "A", BracketNumber: 20, MinMonthlySalary: 24150000, Rate: 0.1000, EffectiveFrom: efFrom},
		{Category: "A", BracketNumber: 21, MinMonthlySalary: 26450000, Rate: 0.1100, EffectiveFrom: efFrom},
		{Category: "A", BracketNumber: 22, MinMonthlySalary: 28000000, Rate: 0.1200, EffectiveFrom: efFrom},
		{Category: "A", BracketNumber: 23, MinMonthlySalary: 30050000, Rate: 0.1300, EffectiveFrom: efFrom},
		{Category: "A", BracketNumber: 24, MinMonthlySalary: 32400000, Rate: 0.1400, EffectiveFrom: efFrom},
		{Category: "A", BracketNumber: 25, MinMonthlySalary: 35400000, Rate: 0.1500, EffectiveFrom: efFrom},
		{Category: "A", BracketNumber: 26, MinMonthlySalary: 39100000, Rate: 0.1600, EffectiveFrom: efFrom},
		{Category: "A", BracketNumber: 27, MinMonthlySalary: 43850000, Rate: 0.1700, EffectiveFrom: efFrom},
		{Category: "A", BracketNumber: 28, MinMonthlySalary: 47800000, Rate: 0.1800, EffectiveFrom: efFrom},
		{Category: "A", BracketNumber: 29, MinMonthlySalary: 51400000, Rate: 0.1900, EffectiveFrom: efFrom},
		{Category: "A", BracketNumber: 30, MinMonthlySalary: 56300000, Rate: 0.2000, EffectiveFrom: efFrom},
		{Category: "A", BracketNumber: 31, MinMonthlySalary: 62200000, Rate: 0.2100, EffectiveFrom: efFrom},
		{Category: "A", BracketNumber: 32, MinMonthlySalary: 68600000, Rate: 0.2200, EffectiveFrom: efFrom},
		{Category: "A", BracketNumber: 33, MinMonthlySalary: 77500000, Rate: 0.2300, EffectiveFrom: efFrom},
		{Category: "A", BracketNumber: 34, MinMonthlySalary: 89000000, Rate: 0.2400, EffectiveFrom: efFrom},
		{Category: "A", BracketNumber: 35, MinMonthlySalary: 103000000, Rate: 0.2500, EffectiveFrom: efFrom},
		{Category: "A", BracketNumber: 36, MinMonthlySalary: 125000000, Rate: 0.2600, EffectiveFrom: efFrom},
		{Category: "A", BracketNumber: 37, MinMonthlySalary: 157000000, Rate: 0.2700, EffectiveFrom: efFrom},
		{Category: "A", BracketNumber: 38, MinMonthlySalary: 206000000, Rate: 0.2800, EffectiveFrom: efFrom},
		{Category: "A", BracketNumber: 39, MinMonthlySalary: 337000000, Rate: 0.2900, EffectiveFrom: efFrom},
		{Category: "A", BracketNumber: 40, MinMonthlySalary: 454000000, Rate: 0.3000, EffectiveFrom: efFrom},
		{Category: "A", BracketNumber: 41, MinMonthlySalary: 550000000, Rate: 0.3100, EffectiveFrom: efFrom},
		{Category: "A", BracketNumber: 42, MinMonthlySalary: 695000000, Rate: 0.3200, EffectiveFrom: efFrom},
		{Category: "A", BracketNumber: 43, MinMonthlySalary: 910000000, Rate: 0.3300, EffectiveFrom: efFrom},
		{Category: "A", BracketNumber: 44, MinMonthlySalary: 1400000000, Rate: 0.3400, EffectiveFrom: efFrom},
	}

	for _, b := range catA {
		var existing tax.TERBracket
		err := tx.Where("category = ? AND bracket_number = ?", b.Category, b.BracketNumber).First(&existing).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := tx.Create(&b).Error; err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
	}

	return nil
}

func seedBPJSRateConfigs(tx *gorm.DB) error {
	efFrom := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	configs := []bpjs.BPJSRateConfig{
		{Type: "KESEHATAN", EmployeeRate: 0.0100, EmployerRate: 0.0400, IsActive: true, EffectiveFrom: efFrom},
		{Type: "JHT", EmployeeRate: 0.0200, EmployerRate: 0.0370, IsActive: true, EffectiveFrom: efFrom},
		{Type: "JKK", EmployeeRate: 0.0000, EmployerRate: 0.0054, IsActive: true, EffectiveFrom: efFrom},
		{Type: "JKM", EmployeeRate: 0.0000, EmployerRate: 0.0030, IsActive: true, EffectiveFrom: efFrom},
		{Type: "JP", EmployeeRate: 0.0100, EmployerRate: 0.0200, IsActive: true, EffectiveFrom: efFrom},
	}

	cap := 12000000.0
	configs[0].MaxSalaryCap = &cap

	for _, c := range configs {
		var existing bpjs.BPJSRateConfig
		err := tx.Where("type = ?", c.Type).First(&existing).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := tx.Create(&c).Error; err != nil {
				return err
			}
		} else if err != nil {
			return err
		} else {
			if err := tx.Model(&existing).Updates(map[string]interface{}{
				"employee_rate": c.EmployeeRate,
				"employer_rate": c.EmployerRate,
				"is_active":     c.IsActive,
			}).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
