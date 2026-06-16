# BPJS & Tax Engine Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add production-grade BPJS (Kesehatan + Ketenagakerjaan) and PPh 21 (TER method) calculation to the payroll pipeline, with admin-configurable rates and two new config UI pages.

**Architecture:** Two new backend modules (`tax`, `bpjs`) following existing Clean Architecture pattern (entity/dto → repository → service → handler → routes), injected into the payroll service via interfaces. Two new admin frontend pages under `frontend/src/pages/admin/`.

**Tech Stack:** Go 1.25 + Echo + GORM + MySQL / React 19 + TypeScript + TanStack Query + React Hook Form + Zod + TailwindCSS

---

## File Map

```
Backend (new):
  backend/internal/modules/tax/
    entity.go, dto.go, repository.go, repository_test.go,
    service.go, service_test.go, mocks_test.go,
    handler.go, handler_test.go, routes.go
  backend/internal/modules/bpjs/
    entity.go, dto.go, repository.go, repository_test.go,
    service.go, service_test.go, mocks_test.go,
    handler.go, handler_test.go, routes.go

Backend (modified):
  backend/pkg/constants/permission_key.go       — new permissions
  backend/pkg/constants/marital_status.go       — new enum type
  backend/internal/modules/payroll/entity.go    — PayrollDetail new fields
  backend/internal/modules/payroll/contract.go  — TaxProvider, BPJSProvider interfaces
  backend/internal/modules/payroll/service.go   — integrate calc into GenerateAll
  backend/internal/bootstrap/container.go       — wire new modules
  backend/internal/routes/api.go                 — register routes
  backend/internal/modules/user/dto.go          — MaritalStatus + Dependents
  backend/internal/modules/user/service.go      — persist new fields

Database:
  backend/migrations/000027_add_tax_bpjs_fields.up.sql
  backend/migrations/000027_add_tax_bpjs_fields.down.sql

Frontend (new):
  frontend/src/features/bpjs-config/
    types.ts, hooks/useBpjsConfig.ts, components/BpjsConfigCard.tsx
  frontend/src/pages/admin/BpjsConfigPage.tsx
  frontend/src/features/pph21-config/
    types.ts, hooks/usePph21Config.ts,
    components/TerBracketsTable.tsx, components/PtkpTable.tsx
  frontend/src/pages/admin/Pph21ConfigPage.tsx

Frontend (modified):
  frontend/src/config/permissions.ts           — new permission keys
  frontend/src/config/menu.ts                  — new menu items
  frontend/src/router.tsx (or App.tsx)         — new routes
```

---
## Task 1: DB Migration

**Files:**
- Create: `backend/migrations/000027_add_tax_bpjs_fields.up.sql`
- Create: `backend/migrations/000027_add_tax_bpjs_fields.down.sql`

- [ ] **Step 1: Create migration files**

```bash
make migrate-create NAME=add_tax_bpjs_fields
```

- [ ] **Step 2: Write up migration**

File: `backend/migrations/000027_add_tax_bpjs_fields.up.sql`

```sql
ALTER TABLE employees
  ADD COLUMN marital_status ENUM('TK','K') NULL AFTER npwp,
  ADD COLUMN dependents_count TINYINT UNSIGNED DEFAULT 0 AFTER marital_status;

ALTER TABLE companies
  ADD COLUMN bpjs_kesehatan_number VARCHAR(50) NULL AFTER tax_number,
  ADD COLUMN bpjs_ketenagakerjaan_number VARCHAR(50) NULL AFTER bpjs_kesehatan_number;

ALTER TABLE payroll_details
  ADD COLUMN code VARCHAR(30) NULL AFTER title,
  ADD COLUMN `group` VARCHAR(30) NULL AFTER code,
  ADD COLUMN is_employer_borne TINYINT(1) DEFAULT 0 AFTER `group`,
  ADD INDEX idx_payroll_details_group (`group`);

CREATE TABLE bpjs_rate_configs (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL,
  company_id BIGINT NULL,
  type ENUM('KESEHATAN','JHT','JKK','JKM','JP') NOT NULL,
  employee_rate DECIMAL(5,4) NOT NULL DEFAULT 0,
  employer_rate DECIMAL(5,4) NOT NULL DEFAULT 0,
  max_salary_cap DECIMAL(15,2) NULL,
  industry_risk_level ENUM('VERY_LOW','LOW','MEDIUM','HIGH','VERY_HIGH') NULL,
  is_active TINYINT(1) DEFAULT 1,
  effective_from DATE NOT NULL,
  effective_until DATE NULL,
  INDEX idx_bpjs_config_company (company_id),
  INDEX idx_bpjs_config_type (type),
  INDEX idx_bpjs_config_active (is_active),
  INDEX idx_bpjs_config_deleted (deleted_at),
  CONSTRAINT fk_bpjs_config_company FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE pph21_term_configs (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL,
  company_id BIGINT NULL,
  category CHAR(1) NOT NULL,
  bracket_number INT NOT NULL,
  min_monthly_salary DECIMAL(15,2) NOT NULL,
  rate DECIMAL(5,4) NOT NULL,
  effective_from DATE NOT NULL,
  effective_until DATE NULL,
  INDEX idx_term_category (category),
  INDEX idx_term_effective (effective_from),
  CONSTRAINT fk_term_config_company FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE ptkp_configs (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL,
  code VARCHAR(10) NOT NULL,
  annual_amount DECIMAL(15,2) NOT NULL,
  effective_year INT NOT NULL,
  INDEX idx_ptkp_code (code),
  INDEX idx_ptkp_year (effective_year)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

- [ ] **Step 3: Write down migration**

File: `backend/migrations/000027_add_tax_bpjs_fields.down.sql`

```sql
ALTER TABLE payroll_details DROP COLUMN is_employer_borne, DROP COLUMN `group`, DROP COLUMN code, DROP INDEX idx_payroll_details_group;
ALTER TABLE companies DROP COLUMN bpjs_ketenagakerjaan_number, DROP COLUMN bpjs_kesehatan_number;
ALTER TABLE employees DROP COLUMN dependents_count, DROP COLUMN marital_status;
DROP TABLE IF EXISTS ptkp_configs;
DROP TABLE IF EXISTS pph21_term_configs;
DROP TABLE IF EXISTS bpjs_rate_configs;
```

- [ ] **Step 4: Run migration**

```bash
make migrate-up
```

Expected: migration applied without errors.

- [ ] **Step 5: Commit**

```bash
git add backend/migrations/000027_add_tax_bpjs_fields.up.sql backend/migrations/000027_add_tax_bpjs_fields.down.sql
git commit -m "feat: add tax and BPJS database schema"
```


## Task 2: Backend Constants & Permissions

**Files:**
- Modify: `backend/pkg/constants/permission_key.go`
- Create: `backend/pkg/constants/marital_status.go`

- [ ] **Step 1: Add BPJS and Tax permission constants**

File: `backend/pkg/constants/permission_key.go` — append before last line:

```go
	// bpjs
	VIEW_BPJS_CONFIG   = "VIEW_BPJS_CONFIG"
	MANAGE_BPJS_CONFIG = "MANAGE_BPJS_CONFIG"

	// tax
	VIEW_TAX_CONFIG   = "VIEW_TAX_CONFIG"
	MANAGE_TAX_CONFIG = "MANAGE_TAX_CONFIG"
```

- [ ] **Step 2: Create MaritalStatus type**

File: `backend/pkg/constants/marital_status.go`

```go
package constants

type MaritalStatus string

const (
	MaritalStatusSingle  MaritalStatus = "TK"
	MaritalStatusMarried MaritalStatus = "K"
)
```

- [ ] **Step 3: Commit**

```bash
git add backend/pkg/constants/permission_key.go backend/pkg/constants/marital_status.go
git commit -m "feat: add tax and BPJS permission keys and marital status constant"
```

## Task 3: Tax Module — Entity & DTO

**Files:**
- Create: `backend/internal/modules/tax/entity.go`
- Create: `backend/internal/modules/tax/dto.go`

- [ ] **Step 1: Create entity with GORM models for TERBracket and PTKPConfig**

File: `backend/internal/modules/tax/entity.go`

```go
package tax

import (
	"time"
	"gorm.io/gorm"
)

type TERBracket struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
	CompanyID        *uint          `gorm:"index" json:"company_id"`
	Category         string         `gorm:"type:char(1);not null" json:"category"`
	BracketNumber    int            `gorm:"not null" json:"bracket_number"`
	MinMonthlySalary float64        `gorm:"type:decimal(15,2);not null" json:"min_monthly_salary"`
	Rate             float64        `gorm:"type:decimal(5,4);not null" json:"rate"`
	EffectiveFrom    time.Time      `gorm:"type:date;not null" json:"effective_from"`
	EffectiveUntil   *time.Time     `gorm:"type:date" json:"effective_until"`
}

func (TERBracket) TableName() string { return "pph21_term_configs" }

type PTKPConfig struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	Code          string         `gorm:"type:varchar(10);not null" json:"code"`
	AnnualAmount  float64        `gorm:"type:decimal(15,2);not null" json:"annual_amount"`
	EffectiveYear int            `gorm:"not null" json:"effective_year"`
}

func (PTKPConfig) TableName() string { return "ptkp_configs" }
```

- [ ] **Step 2: Create DTOs**

File: `backend/internal/modules/tax/dto.go`

```go
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
```

- [ ] **Step 3: Commit**

```bash
git add backend/internal/modules/tax/entity.go backend/internal/modules/tax/dto.go
git commit -m "feat: add tax module entity and DTOs"
```

## Task 4: Tax Module — Repository + Tests

**Files:**
- Create: `backend/internal/modules/tax/repository.go`
- Create: `backend/internal/modules/tax/repository_test.go`

- [ ] **Step 1: Write repository test** (TDD — test first)

```bash
# Use the test code from the conversation design. Key test cases:
# - FindTERBrackets returns valid brackets for category+date
# - FindTERBrackets returns empty when future date with no match
# - FindPTKPByYear returns configs for given year
# - FindPTKPByYear returns empty for year with no data
# - CRUD CreateTERBracket/FindByID/Update/Delete (soft delete)
# - CRUD CreatePTKPConfig/FindByID/Update/Delete
```

Write: `backend/internal/modules/tax/repository_test.go` — use `testutil.NewTestDB(&TERBracket{}, &PTKPConfig{})` pattern, see loan repository_test.go for reference. Tests are written with the full code in the conversations's Task 4 design.

- [ ] **Step 2: Run test to verify it fails**

```bash
go test ./backend/internal/modules/tax/... -v -run TestRepository
```

Expected: compilation errors (types not defined yet).

- [ ] **Step 3: Implement repository**

File: `backend/internal/modules/tax/repository.go`

The repository interface includes: FindTERBrackets, CreateTERBracket, FindTERBracketByID, UpdateTERBracket, DeleteTERBracket, ListTERBrackets, FindPTKPByYear, CreatePTKPConfig, FindPTKPConfigByID, UpdatePTKPConfig, DeletePTKPConfig, ListPTKPConfigs.

Implementation follows standard pattern: `utils.GetDBFromContext(ctx, r.db)`, soft delete via `deleted_at = now()`, pagination with offset/limit.

Full code reference: see conversation Task 4 Step 3.

- [ ] **Step 4: Run tests to verify they pass**

```bash
go test ./backend/internal/modules/tax/... -v -run TestRepository
```

Expected: ALL PASS.

- [ ] **Step 5: Commit**

```bash
git add backend/internal/modules/tax/repository.go backend/internal/modules/tax/repository_test.go
git commit -m "feat: add tax repository with TER bracket and PTKP lookups"
```

## Task 5: Tax Module — Service + Tests

**Files:**
- Create: `backend/internal/modules/tax/service.go`
- Create: `backend/internal/modules/tax/service_test.go`
- Create: `backend/internal/modules/tax/mocks_test.go`

- [ ] **Step 1: Create mocks** (testify/mock for Repository interface)

File: `backend/internal/modules/tax/mocks_test.go` — one mock per repository method. See conversation Task 5 Step 1 for full code.

- [ ] **Step 2: Write service tests** (TDD)

File: `backend/internal/modules/tax/service_test.go`

Key test cases:
- CalculateTER for Category A (single employee, TK/0) at 10M salary → 1% rate
- CalculateTER highest bracket fallback
- CalculateTER exact bracket boundary
- CalculateTER for Category B (married K/1) 
- CalculateTER for Category C (married K/3)
- CalculateTER error when no brackets found
- derivePTKPCode mapping tests (TK/0 → "TK/0", K/3 → "K/3", cap at 3)
- deriveTERCategory mapping tests (A/B/C)
- findTERRate returns correct rate for salary ranges
- calculateBiayaJabatanMonthly (5% cap 500K/month)
- calculateProgressiveTax (50M → 2.5M, 100M → 9M)
- ReconcileAnnual underpaid scenario
- CRUD CreateTERBracket/CreatePTKPConfig

Full code: see conversation Task 5 Step 2.

- [ ] **Step 3: Run test to verify it fails**

```bash
go test ./backend/internal/modules/tax/... -v -run TestService
```

Expected: compilation errors.

- [ ] **Step 4: Implement service**

File: `backend/internal/modules/tax/service.go`

Key functions:
- `derivePTKPCode(maritalStatus, dependents)` → "TK/0", "K/1", etc (cap at 3)
- `deriveTERCategory(ptkpCode)` → "A", "B", "C" per regulation mapping
- `findTERRate(brackets, grossIncome)` → highest matching bracket rate
- `calculateBiayaJabatanMonthly(gross)` → min(5%×gross, 500000)
- `calculateProgressiveTax(pkp)` → layer 1: 5% up to 60M, layer 2: 15% up to 250M, layer 3: 25% up to 500M, layer 4: 30% up to 5B, layer 5: 35%
- `CalculateTER(ctx, grossMonthly, maritalStatus, dependents)` → PPh21Result
- `ReconcileAnnual(ctx, year, ptkpCode, grossAnnual, monthlyDetails)` → AnnualSettlement

Full code: see conversation Task 5 Step 4.

- [ ] **Step 5: Run tests to verify they pass**

```bash
go test ./backend/internal/modules/tax/... -v -run TestService
```

Expected: ALL PASS.

- [ ] **Step 6: Commit**

```bash
git add backend/internal/modules/tax/service.go backend/internal/modules/tax/service_test.go backend/internal/modules/tax/mocks_test.go
git commit -m "feat: add tax service with TER calculation and progressive tax"
```

## Task 6: Tax Module — Handler, Routes + Tests

**Files:**
- Create: `backend/internal/modules/tax/handler.go`
- Create: `backend/internal/modules/tax/handler_test.go`
- Create: `backend/internal/modules/tax/routes.go`
- Modify: `backend/internal/routes/api.go`

- [ ] **Step 1: Write handler tests** (TDD)

File: `backend/internal/modules/tax/handler_test.go`

Test cases:
- ListTERBrackets returns 200 with bracket data
- CreateTERBracket returns 201 with valid payload
- ListPTKPConfigs returns 200 with config data

Use `testutil.NewAPITest`, `testutil.WithAuthContext`. Full code: see conversation Task 6 Step 1.

- [ ] **Step 2: Run test to verify it fails**

```bash
go test ./backend/internal/modules/tax/... -v -run TestHandler
```

- [ ] **Step 3: Implement handler**

File: `backend/internal/modules/tax/handler.go`

CRUD handlers for TER brackets (List/Create/GetByID/Update/Delete) and PTKP configs (List/Create/Update/Delete). Standard Echo handler pattern with `utils.GetUserContext`, `ctx.Bind`, `ctx.Validate`, pagination via query params, response via `response.NewResponses`.

Full code: see conversation Task 6 Step 3.

- [ ] **Step 4: Implement routes**

File: `backend/internal/modules/tax/routes.go`

```go
package tax

import (
	"basekarya-backend/internal/middleware"
	"basekarya-backend/pkg/constants"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Group, handler *Handler, auth *middleware.AuthMiddleware) {
	g := e.Group("/admin/tax")
	g.GET("/ter-brackets", handler.ListTERBrackets, auth.GrantPermission(constants.VIEW_TAX_CONFIG))
	g.POST("/ter-brackets", handler.CreateTERBracket, auth.GrantPermission(constants.MANAGE_TAX_CONFIG))
	g.GET("/ter-brackets/:id", handler.GetTERBracketByID, auth.GrantPermission(constants.VIEW_TAX_CONFIG))
	g.PUT("/ter-brackets/:id", handler.UpdateTERBracket, auth.GrantPermission(constants.MANAGE_TAX_CONFIG))
	g.DELETE("/ter-brackets/:id", handler.DeleteTERBracket, auth.GrantPermission(constants.MANAGE_TAX_CONFIG))
	g.GET("/ptkp-configs", handler.ListPTKPConfigs, auth.GrantPermission(constants.VIEW_TAX_CONFIG))
	g.POST("/ptkp-configs", handler.CreatePTKPConfig, auth.GrantPermission(constants.MANAGE_TAX_CONFIG))
	g.PUT("/ptkp-configs/:id", handler.UpdatePTKPConfig, auth.GrantPermission(constants.MANAGE_TAX_CONFIG))
	g.DELETE("/ptkp-configs/:id", handler.DeletePTKPConfig, auth.GrantPermission(constants.MANAGE_TAX_CONFIG))
}
```

- [ ] **Step 5: Register routes in api.go**

Add import and route call in `backend/internal/routes/api.go`:
```go
import "basekarya-backend/internal/modules/tax"
// inside setupRoutes():
tax.RegisterRoutes(protected, r.container.TaxHandler, r.container.AuthMiddleware)
```

- [ ] **Step 6: Run tests**

```bash
go test ./backend/internal/modules/tax/... -v -run TestHandler
```

Expected: handlers compile and tests PASS.

- [ ] **Step 7: Commit**

```bash
git add backend/internal/modules/tax/handler.go backend/internal/modules/tax/handler_test.go backend/internal/modules/tax/routes.go backend/internal/routes/api.go
git commit -m "feat: add tax handler, routes and admin API endpoints"
```

## Task 7: BPJS Module — Entity & DTO

**Files:**
- Create: `backend/internal/modules/bpjs/entity.go`
- Create: `backend/internal/modules/bpjs/dto.go`

- [ ] **Step 1: Create entity**

File: `backend/internal/modules/bpjs/entity.go`

```go
package bpjs

import (
	"time"
	"gorm.io/gorm"
)

type BPJSRateConfig struct {
	ID                uint           `gorm:"primaryKey" json:"id"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
	CompanyID         *uint          `gorm:"index" json:"company_id"`
	Type              string         `gorm:"type:varchar(20);not null" json:"type"`
	EmployeeRate      float64        `gorm:"type:decimal(5,4);not null;default:0" json:"employee_rate"`
	EmployerRate      float64        `gorm:"type:decimal(5,4);not null;default:0" json:"employer_rate"`
	MaxSalaryCap      *float64       `gorm:"type:decimal(15,2)" json:"max_salary_cap"`
	IndustryRiskLevel *string        `gorm:"type:varchar(10)" json:"industry_risk_level"`
	IsActive          bool           `gorm:"default:true" json:"is_active"`
	EffectiveFrom     time.Time      `gorm:"type:date;not null" json:"effective_from"`
	EffectiveUntil    *time.Time     `gorm:"type:date" json:"effective_until"`
}

func (BPJSRateConfig) TableName() string { return "bpjs_rate_configs" }
```

- [ ] **Step 2: Create DTOs**

File: `backend/internal/modules/bpjs/dto.go`

```go
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
```

- [ ] **Step 3: Commit**

```bash
git add backend/internal/modules/bpjs/entity.go backend/internal/modules/bpjs/dto.go
git commit -m "feat: add BPJS module entity and DTOs"
```

## Task 8: BPJS Module — Repository + Tests

**Files:**
- Create: `backend/internal/modules/bpjs/repository.go`
- Create: `backend/internal/modules/bpjs/repository_test.go`

- [ ] **Step 1: Write repository tests** (TDD)

File: `backend/internal/modules/bpjs/repository_test.go`

Use `testutil.NewTestDB(&BPJSRateConfig{})`. Test: FindActiveByType, FindActiveByType with inactive record, FindAllActive, CRUD cycle. See conversation Task 8 for full test code.

- [ ] **Step 2: Run test (expect failure)**

```bash
go test ./backend/internal/modules/bpjs/... -v -run TestRepository
```

- [ ] **Step 3: Implement repository**

File: `backend/internal/modules/bpjs/repository.go`

Interface: FindActiveByType, FindAllActive, Create, FindByID, Update, Delete, List. Follow `utils.GetDBFromContext` pattern. See conversation Task 8 Step 3 for full code.

- [ ] **Step 4: Run tests (expect pass)**

```bash
go test ./backend/internal/modules/bpjs/... -v -run TestRepository
```

- [ ] **Step 5: Commit**

```bash
git add backend/internal/modules/bpjs/repository.go backend/internal/modules/bpjs/repository_test.go
git commit -m "feat: add BPJS repository with rate config lookups"
```

## Task 9: BPJS Module — Service + Tests

**Files:**
- Create: `backend/internal/modules/bpjs/service.go`
- Create: `backend/internal/modules/bpjs/service_test.go`
- Create: `backend/internal/modules/bpjs/mocks_test.go`

- [ ] **Step 1: Create mocks**

File: `backend/internal/modules/bpjs/mocks_test.go` — mock Repository. See conversation Task 9 Step 1.

- [ ] **Step 2: Write service tests** (TDD)

File: `backend/internal/modules/bpjs/service_test.go`

Test: CalculateAll returns 7 components (3 employee + 5 employer minus 2 overlap), Kesehatan 1%/4% at 10M, max salary cap applied, empty when no active configs, inactive configs excluded. See conversation Task 9 Step 2 for full test code.

- [ ] **Step 3: Run test (expect failure)**

```bash
go test ./backend/internal/modules/bpjs/... -v -run TestService
```

- [ ] **Step 4: Implement service**

File: `backend/internal/modules/bpjs/service.go`

`CalculateAll` iterates all active configs, applies `min(salary, cap) * rate`, produces separate employee/employer components with `IsEmployerBorne` flag. Standard CRUD methods. See conversation Task 9 Step 4 for full code.

- [ ] **Step 5: Run tests (expect pass)**

```bash
go test ./backend/internal/modules/bpjs/... -v -run TestService
```

- [ ] **Step 6: Commit**

```bash
git add backend/internal/modules/bpjs/service.go backend/internal/modules/bpjs/service_test.go backend/internal/modules/bpjs/mocks_test.go
git commit -m "feat: add BPJS service with contribution calculation"
```

## Task 10: BPJS Module — Handler, Routes + Tests

**Files:**
- Create: `backend/internal/modules/bpjs/handler.go`
- Create: `backend/internal/modules/bpjs/handler_test.go`
- Create: `backend/internal/modules/bpjs/routes.go`
- Modify: `backend/internal/routes/api.go`

- [ ] **Step 1: Write handler tests** (TDD)

File: `backend/internal/modules/bpjs/handler_test.go`

Test: List returns 200, Create returns 201. Use testutil.NewAPITest pattern. See conversation Task 10 Step 1.

- [ ] **Step 2: Run test (expect failure)**

```bash
go test ./backend/internal/modules/bpjs/... -v -run TestHandler
```

- [ ] **Step 3: Implement handler**

File: `backend/internal/modules/bpjs/handler.go`

Standard CRUD handler: List, Create, GetByID, Update, Delete. See conversation Task 10 Step 3 for full code.

- [ ] **Step 4: Implement routes**

File: `backend/internal/modules/bpjs/routes.go`

```go
package bpjs

import (
	"basekarya-backend/internal/middleware"
	"basekarya-backend/pkg/constants"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Group, handler *Handler, auth *middleware.AuthMiddleware) {
	g := e.Group("/admin/bpjs")
	g.GET("/configs", handler.List, auth.GrantPermission(constants.VIEW_BPJS_CONFIG))
	g.POST("/configs", handler.Create, auth.GrantPermission(constants.MANAGE_BPJS_CONFIG))
	g.GET("/configs/:id", handler.GetByID, auth.GrantPermission(constants.VIEW_BPJS_CONFIG))
	g.PUT("/configs/:id", handler.Update, auth.GrantPermission(constants.MANAGE_BPJS_CONFIG))
	g.DELETE("/configs/:id", handler.Delete, auth.GrantPermission(constants.MANAGE_BPJS_CONFIG))
}
```

- [ ] **Step 5: Register routes in api.go**

Add import: `"basekarya-backend/internal/modules/bpjs"` and route call:
```go
bpjs.RegisterRoutes(protected, r.container.BpjsHandler, r.container.AuthMiddleware)
```

- [ ] **Step 6: Run tests**

```bash
go test ./backend/internal/modules/bpjs/... -v -run TestHandler
```

- [ ] **Step 7: Commit**

```bash
git add backend/internal/modules/bpjs/handler.go backend/internal/modules/bpjs/handler_test.go backend/internal/modules/bpjs/routes.go backend/internal/routes/api.go
git commit -m "feat: add BPJS handler, routes and admin API endpoints"
```

## Task 11: Payroll Integration — Wire Everything Together

**Files:**
- Modify: `backend/internal/modules/payroll/contract.go`
- Modify: `backend/internal/modules/payroll/entity.go`
- Modify: `backend/internal/modules/payroll/service.go`
- Modify: `backend/internal/bootstrap/container.go`

- [ ] **Step 1: Add TaxProvider and BPJSProvider interfaces to payroll contract**

File: `backend/internal/modules/payroll/contract.go` — append:

```go
type TaxProvider interface {
	CalculateTER(ctx context.Context, grossMonthlyIncome float64, maritalStatus constants.MaritalStatus, dependentsCount int) (PPh21Result, error)
}

type PPh21Result struct {
	TERCategory  string  `json:"ter_category"`
	PTKPCode     string  `json:"ptkp_code"`
	GrossMonthly float64 `json:"gross_monthly"`
	TERRate      float64 `json:"ter_rate"`
	MonthlyPPh21 float64 `json:"monthly_pph21"`
}

type BPJSProvider interface {
	CalculateAll(ctx context.Context, grossMonthlyIncome float64) ([]BPJSComponent, error)
}

type BPJSComponent struct {
	Type            string  `json:"type"`
	Code            string  `json:"code"`
	EmployeeAmount  float64 `json:"employee_amount"`
	EmployerAmount  float64 `json:"employer_amount"`
	IsEmployerBorne bool    `json:"is_employer_borne"`
}
```

- [ ] **Step 2: Update PayrollDetail entity with new fields**

File: `backend/internal/modules/payroll/entity.go` — update PayrollDetail:

```go
type PayrollDetail struct {
	ID              uint                       `gorm:"primaryKey" json:"id"`
	PayrollID       uint                       `json:"payroll_id"`
	CompanyID       uint                       `gorm:"index;not null" json:"company_id"`
	Title           string                     `gorm:"type:varchar(150);not null" json:"title"`
	Code            *string                    `gorm:"type:varchar(30)" json:"code"`
	Group           *string                    `gorm:"type:varchar(30)" json:"group"`
	IsEmployerBorne bool                       `gorm:"default:false" json:"is_employer_borne"`
	Type            constants.PayrollDetailType `gorm:"type:varchar(20);not null" json:"type"`
	Amount          float64                    `gorm:"type:decimal(15,2);not null" json:"amount"`
}
```

- [ ] **Step 3: Update payroll service to accept and use tax/bpjs providers**

File: `backend/internal/modules/payroll/service.go`

Add fields to service struct:
```go
taxProv  TaxProvider
bpjsProv BPJSProvider
```

Update `NewService` signature to accept the two new dependencies.

In `GenerateAll`, inside the employee loop, after overtime calculation and before penalties:

```go
// Calculate PPh 21 TER
var pph21Amount float64
if s.taxProv != nil {
    result, err := s.taxProv.CalculateTER(ctx, baseSalary, emp.MaritalStatus, emp.DependentsCount)
    if err != nil {
        logger.Warn("failed to calculate PPh 21 for employee %d: %v", emp.ID, err)
    } else {
        pph21Amount = result.MonthlyPPh21
    }
}

// Calculate BPJS
var bpjsEmployeeTotal float64
var bpjsComponents []BPJSComponent
if s.bpjsProv != nil {
    components, err := s.bpjsProv.CalculateAll(ctx, baseSalary)
    if err != nil {
        logger.Warn("failed to calculate BPJS for employee %d: %v", emp.ID, err)
    } else {
        bpjsComponents = components
        for _, c := range components {
            if !c.IsEmployerBorne {
                bpjsEmployeeTotal += c.EmployeeAmount
            }
        }
    }
}
```

Update PayrollDetail creation in the loop to include tax and BPJS entries:

```go
// BPJS employee deductions
for _, c := range bpjsComponents {
    if !c.IsEmployerBorne && c.EmployeeAmount > 0 {
        code := "BPJS_" + c.Type + "_E"
        payroll.Details = append(payroll.Details, PayrollDetail{
            Title: "BPJS " + c.Type, Type: constants.DetailTypeDeduction,
            Amount: c.EmployeeAmount, Code: &code, Group: strPtr("BPJS"),
            CompanyID: emp.CompanyID,
        })
    }
}

// PPh 21
if pph21Amount > 0 {
    code := "PPH21"
    payroll.Details = append(payroll.Details, PayrollDetail{
        Title: "PPh 21", Type: constants.DetailTypeDeduction,
        Amount: pph21Amount, Code: &code, Group: strPtr("TAX"),
        CompanyID: emp.CompanyID,
    })
}

// BPJS employer contributions (info only)
for _, c := range bpjsComponents {
    if c.IsEmployerBorne && c.EmployerAmount > 0 {
        code := "BPJS_" + c.Type + "_R"
        payroll.Details = append(payroll.Details, PayrollDetail{
            Title: "BPJS " + c.Type + " (Employer)", Type: constants.DetailTypeAllowance,
            Amount: c.EmployerAmount, Code: &code, Group: strPtr("BPJS"),
            IsEmployerBorne: true, CompanyID: emp.CompanyID,
        })
    }
}
```

Update net salary: `totalDeduction += pph21Amount + bpjsEmployeeTotal`

The full integration code is in the conversation Task 11 Step 3 design.

- [ ] **Step 4: Wire up in DI container**

File: `backend/internal/bootstrap/container.go`

Add imports for `tax` and `bpjs`. Add after existing repo initializations:
```go
taxRepo := tax.NewRepository(db.GetDB())
bpjsRepo := bpjs.NewRepository(db.GetDB())
taxSvc := tax.NewService(taxRepo)
bpjsSvc := bpjs.NewService(bpjsRepo)
```

Add handlers:
```go
taxHandler := tax.NewHandler(taxSvc)
bpjsHandler := bpjs.NewHandler(bpjsSvc)
```

Update payroll service construction:
```go
payrollSvc := payroll.NewService(payrollRepo, userRepo, reimburseRepo, attendanceRepo, companyRepo, notificationSvc, transactionManager, httpClient.GetClient(), email, loanRepo, overtimeRepo, taxSvc, bpjsSvc)
```

Add to Container struct and return value: `TaxHandler *tax.Handler`, `BpjsHandler *bpjs.Handler`.

- [ ] **Step 5: Verify compilation**

```bash
cd backend && go build ./...
```

Expected: no compilation errors.

- [ ] **Step 6: Commit**

```bash
git add backend/internal/modules/payroll/contract.go backend/internal/modules/payroll/entity.go backend/internal/modules/payroll/service.go backend/internal/bootstrap/container.go
git commit -m "feat: integrate tax and BPJS calculators into payroll pipeline"
```

## Task 12: Employee Entity — Add MaritalStatus & Dependents fields

**Files:**
- Modify: `backend/internal/modules/user/entity.go`
- Modify: `backend/internal/modules/user/dto.go`
- Modify: `backend/internal/modules/user/service.go`

- [ ] **Step 1: Add fields to Employee GORM entity**

File: `backend/internal/modules/user/entity.go` — add to Employee struct:

```go
MaritalStatus   constants.MaritalStatus `gorm:"type:enum('TK','K')" json:"marital_status"`
DependentsCount int                     `gorm:"type:tinyint;default:0" json:"dependents_count"`
```

- [ ] **Step 2: Add fields to request DTOs**

Find the employee create/update request struct (typically in `dto.go` or as inline struct) and add:
```go
MaritalStatus   *string `json:"marital_status" validate:"omitempty,oneof=TK K"`
DependentsCount *int    `json:"dependents_count" validate:"omitempty,min=0,max=3"`
```

- [ ] **Step 3: Map fields in service**

In the employee update/create method, map the new fields:
```go
if req.MaritalStatus != nil {
    employee.MaritalStatus = constants.MaritalStatus(*req.MaritalStatus)
}
if req.DependentsCount != nil {
    employee.DependentsCount = *req.DependentsCount
}
```

- [ ] **Step 4: Commit**

```bash
git add backend/internal/modules/user/entity.go backend/internal/modules/user/dto.go backend/internal/modules/user/service.go
git commit -m "feat: add marital status and dependents to employee"
```

## Task 13: Seed Data — TER Brackets & PTKP Defaults

**Files:**
- Create or modify seed files to insert default TER and PTKP data.

- [ ] **Step 1: Create seed SQL for TER brackets (2026 defaults)**

Insert default TER brackets for all categories A, B, C. Key brackets (Category A example):
- No 1: min 0, rate 0.00
- No 2: min 5,400,000, rate 0.25%
- No 3: min 5,650,000, rate 0.50%
- ...continuing per PP 58/2023 schedule

Insert into `pph21_term_configs` with company_id=NULL and effective_from='2026-01-01'.

- [ ] **Step 2: Create seed SQL for PTKP configs (2026)**

```sql
INSERT INTO ptkp_configs (code, annual_amount, effective_year) VALUES
('TK/0', 54000000, 2026),
('TK/1', 58500000, 2026),
('TK/2', 63000000, 2026),
('TK/3', 67500000, 2026),
('K/0',  58500000, 2026),
('K/1',  63000000, 2026),
('K/2',  67500000, 2026),
('K/3',  72000000, 2026);
```

- [ ] **Step 3: Create seed SQL for BPJS defaults (2026)**

```sql
INSERT INTO bpjs_rate_configs (type, employee_rate, employer_rate, max_salary_cap, is_active, effective_from) VALUES
('KESEHATAN', 0.0100, 0.0400, 12000000, 1, '2026-01-01'),
('JHT',       0.0200, 0.0370, NULL,      1, '2026-01-01'),
('JKK',       0.0000, 0.0054, NULL,      1, '2026-01-01'),
('JKM',       0.0000, 0.0030, NULL,      1, '2026-01-01'),
('JP',        0.0100, 0.0200, NULL,      1, '2026-01-01');
```

- [ ] **Step 4: Run seeds**

```bash
make seed
```

- [ ] **Step 5: Commit**

```bash
git add backend/seeds/
git commit -m "feat: add seed data for TER brackets, PTKP, and BPJS defaults"
```

## Task 14: Frontend — Permissions & Menu

**Files:**
- Modify: `frontend/src/config/permissions.ts`
- Modify: `frontend/src/config/menu.ts`

- [ ] **Step 1: Add permission keys**

File: `frontend/src/config/permissions.ts` — append before `} as const;`:

```ts
  // bpjs config
  VIEW_BPJS_CONFIG: "VIEW_BPJS_CONFIG",
  MANAGE_BPJS_CONFIG: "MANAGE_BPJS_CONFIG",

  // tax config
  VIEW_TAX_CONFIG: "VIEW_TAX_CONFIG",
  MANAGE_TAX_CONFIG: "MANAGE_TAX_CONFIG",
```

- [ ] **Step 2: Add menu items**

File: `frontend/src/config/menu.ts`

Add imports: `import { HeartPulse, Percent } from "lucide-react";`

Add menu items in "Pengaturan" group (before Departments):

```ts
  {
    title: "BPJS Config",
    href: "/admin/bpjs-config",
    icon: HeartPulse,
    permission: PERMISSIONS.VIEW_BPJS_CONFIG,
    group: "Pengaturan",
    hideForPlatformAdmin: true,
  },
  {
    title: "PPh 21 Config",
    href: "/admin/pph21-config",
    icon: Percent,
    permission: PERMISSIONS.VIEW_TAX_CONFIG,
    group: "Pengaturan",
    hideForPlatformAdmin: true,
  },
```

- [ ] **Step 3: Run frontend tests**

```bash
cd frontend && pnpm test -- --run
```

Expected: existing permission tests still pass.

- [ ] **Step 4: Commit**

```bash
git add frontend/src/config/permissions.ts frontend/src/config/menu.ts
git commit -m "feat: add BPJS and tax config menu items and permissions"
```

## Task 15: Frontend — BPJS Config Page

**Files:**
- Create: `frontend/src/features/bpjs-config/types.ts`
- Create: `frontend/src/features/bpjs-config/hooks/useBpjsConfig.ts`
- Create: `frontend/src/features/bpjs-config/components/BpjsConfigCard.tsx`
- Create: `frontend/src/pages/admin/BpjsConfigPage.tsx`
- Modify: Router to add `/admin/bpjs-config` route

- [ ] **Step 1: Create types**

File: `frontend/src/features/bpjs-config/types.ts`

Types: `BPJSComponentType`, `IndustryRiskLevel`, `BPJSRateConfig`, `BPJSRateConfigPayload`. See conversation Task 15 Step 1 for full type definitions.

- [ ] **Step 2: Create React Query hooks**

File: `frontend/src/features/bpjs-config/hooks/useBpjsConfig.ts`

Hooks: `useBpjsConfigs` (GET), `useCreateBpjsConfig` (POST), `useUpdateBpjsConfig` (PUT), `useDeleteBpjsConfig` (DELETE). Standard pattern with query invalidation and toast notifications. See conversation Task 15 Step 2.

- [ ] **Step 3: Create config card component**

File: `frontend/src/features/bpjs-config/components/BpjsConfigCard.tsx`

React Hook Form + Zod form with: employee rate (%), employer rate (%), max salary cap (Rp), industry risk level (JKK only), effective dates, active toggle. Shows "Configured"/"Default" badge based on existing config. See conversation Task 15 Step 3 for full component code.

- [ ] **Step 4: Create page**

File: `frontend/src/pages/admin/BpjsConfigPage.tsx`

Renders 5 BpjsConfigCard components (Kesehatan, JHT, JKK, JKM, JP) in a 2-column grid. Permission gate with `hasPermission(VIEW_BPJS_CONFIG)`. See conversation Task 15 Step 4.

- [ ] **Step 5: Add route**

Find the router file (likely `frontend/src/App.tsx` or `frontend/src/router.tsx`) and add:
```tsx
import BpjsConfigPage from "@/pages/admin/BpjsConfigPage";
// in routes:
{ path: "/admin/bpjs-config", element: <BpjsConfigPage /> }
```

- [ ] **Step 6: Commit**

```bash
git add frontend/src/features/bpjs-config/ frontend/src/pages/admin/BpjsConfigPage.tsx
git commit -m "feat: add BPJS configuration admin page"
```

## Task 16: Frontend — PPh 21 Config Page

**Files:**
- Create: `frontend/src/features/pph21-config/types.ts`
- Create: `frontend/src/features/pph21-config/hooks/usePph21Config.ts`
- Create: `frontend/src/features/pph21-config/components/TerBracketsTable.tsx`
- Create: `frontend/src/features/pph21-config/components/PtkpTable.tsx`
- Create: `frontend/src/pages/admin/Pph21ConfigPage.tsx`
- Modify: Router to add `/admin/pph21-config` route

- [ ] **Step 1: Create types**

File: `frontend/src/features/pph21-config/types.ts`

Types: `TERBracket`, `TERBracketPayload`, `PTKPConfig`, `PTKPConfigPayload`.

- [ ] **Step 2: Create hooks**

File: `frontend/src/features/pph21-config/hooks/usePph21Config.ts`

Hooks: `useTerBrackets(category)` (GET with filter), `useCreateTerBracket`, `useUpdateTerBracket`, `useDeleteTerBracket`, `usePtkpConfigs(year)`, `useCreatePtkpConfig`, `useUpdatePtkpConfig`, `useDeletePtkpConfig`. Standard pattern with toast + query invalidation.

- [ ] **Step 3: Create TER brackets table component**

File: `frontend/src/features/pph21-config/components/TerBracketsTable.tsx`

Tabbed table with Category A/B/C tabs. Each tab shows an editable table: bracket number, min monthly salary (Rp), rate (%). Inline add/edit form rows. Save per-category button. See conversation design Section 3 wireframe.

- [ ] **Step 4: Create PTKP table component**

File: `frontend/src/features/pph21-config/components/PtkpTable.tsx`

Simple table: PTKP Code (TK/0, K/1, etc.), Annual Amount (Rp), Effective Year. Edit dialog/inline for each row. Supported PTKP codes: TK/0, TK/1, TK/2, TK/3, K/0, K/1, K/2, K/3.

- [ ] **Step 5: Create page**

File: `frontend/src/pages/admin/Pph21ConfigPage.tsx`

Tab layout: "TER Brackets" tab (TerBracketsTable) and "PTKP" tab (PtkpTable). Permission gate with `hasPermission(VIEW_TAX_CONFIG)`.

- [ ] **Step 6: Add route**

In the router file, add:
```tsx
import Pph21ConfigPage from "@/pages/admin/Pph21ConfigPage";
{ path: "/admin/pph21-config", element: <Pph21ConfigPage /> }
```

- [ ] **Step 7: Commit**

```bash
git add frontend/src/features/pph21-config/ frontend/src/pages/admin/Pph21ConfigPage.tsx
git commit -m "feat: add PPh 21 configuration admin page"
```

## Task 17: Final Verification

- [ ] **Step 1: Run all backend tests**

```bash
cd backend && go test ./... -count=1
```

Expected: ALL tests pass, including new tax and bpjs module tests.

- [ ] **Step 2: Run all frontend tests**

```bash
cd frontend && pnpm test -- --run
```

Expected: ALL tests pass.

- [ ] **Step 3: Build backend**

```bash
make build-be
```

Expected: builds without errors.

- [ ] **Step 4: Build frontend**

```bash
make build-fe
```

Expected: builds without TypeScript or Vite errors.

- [ ] **Step 5: Verify the full test pipeline**

```bash
make test
```

- [ ] **Step 6: Final commit (if any cleanup changes)**

```bash
git add -A
git commit -m "chore: final cleanup and verification for BPJS and tax engine"
```

## Out of Scope (Not in this plan)

- THR (holiday allowance) calculation
- Non-TER (progressive monthly) PPh 21 method
- Form 1721-A1 PDF generation
- December auto-reconciliation in payroll
- Mid-year resignation termination settlement
- Company-specific BPJS/Tax rate overrides (only global defaults)
- Multi-country tax support
