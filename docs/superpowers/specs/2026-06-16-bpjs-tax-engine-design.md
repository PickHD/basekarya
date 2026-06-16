# BPJS & Tax Engine Design

**Date:** 2026-06-16  
**Status:** Approved  
**Scope:** Production-grade BPJS Kesehatan, BPJS Ketenagakerjaan, and PPh 21 (TER) calculation engine for BaseKarya HRIS.

## Overview

Add tax and BPJS calculation capabilities to the existing payroll pipeline. Two new modules (`tax`, `bpjs`) follow the existing Clean Architecture pattern. Admins configure rates via new UI pages. Payroll service orchestrates by calling calculator interfaces.

**Key decisions:**
- PPh 21: TER method only (PP 58/2023), no older progressive-monthly method
- BPJS rates: admin-configurable per company with global defaults
- Contributions: track both employee-borne and employer-borne
- PTKP: derived from Employee marital status + dependents count

---

## 1. Data Model

### 1.1 Employee entity additions

| Field | Type | Notes |
|-------|------|-------|
| `MaritalStatus` | enum(`TK`,`K`) | PTKP base category |
| `DependentsCount` | int (0–3) | PTKP dependents, capped at 3 |

`NPWP` already exists on Employee. PTKP code (`TK/0`, `K/2`, etc.) is **derived** at calculation time from `MaritalStatus + DependentsCount`.

### 1.2 Company entity additions

| Field | Type |
|-------|------|
| `BPJSKesehatanNumber` | string |
| `BPJSKetenagakerjaanNumber` | string |

### 1.3 New tables

**`bpjs_rate_configs`** — rates per component, per company override:

| Column | Type | Notes |
|--------|------|-------|
| id | uint PK | |
| company_id | uint FK nullable | NULL = global default |
| type | enum(`KESEHATAN`,`JHT`,`JKK`,`JKM`,`JP`) | Component |
| employee_rate | decimal(5,4) | e.g. 0.0100 = 1% |
| employer_rate | decimal(5,4) | |
| max_salary_cap | decimal(15,2) nullable | Cap for Kesehatan, JP |
| industry_risk_level | enum(`VERY_LOW`,`LOW`,`MEDIUM`,`HIGH`,`VERY_HIGH`) | JKK only |
| is_active | bool | |
| effective_from | date | |
| effective_until | date nullable | |
| created_at / updated_at / deleted_at | timestamps | |

**`pph21_term_configs`** — TER rate brackets by PTKP category:

| Column | Type |
|--------|------|
| id | uint PK |
| company_id | uint FK nullable |
| category | char(1) (`A`,`B`,`C`) |
| bracket_number | int |
| min_monthly_salary | decimal(15,2) |
| rate | decimal(5,4) |
| effective_from | date |
| effective_until | date nullable |
| created_at / updated_at / deleted_at | |

**`ptkp_configs`** — annual PTKP thresholds:

| Column | Type |
|--------|------|
| id | uint PK |
| code | string (`TK/0`,`TK/1`,`TK/2`,`TK/3`,`K/0`,`K/1`,`K/2`,`K/3`) |
| annual_amount | decimal(15,2) |
| effective_year | int |
| created_at / updated_at / deleted_at | |

### 1.4 PayrollDetail extensions

Add these columns to existing `payroll_details`:

| Field | Type | Purpose |
|-------|------|---------|
| `code` | string | Machine-readable key (`PPH21`, `BPJS_KES`, `BPJS_JHT_E`) |
| `group` | string | Grouping for payslip rendering (`TAX`, `BPJS`, `SALARY`, `PENALTY`) |
| `is_employer_borne` | bool | True for employer contributions (informational, excluded from net) |

---

## 2. Calculation Engine

### 2.1 PPh 21 TER (Monthly)

```
Monthly PPh 21 = Gross Monthly Income × TER Rate
```

1. Gross income = base salary + regular fixed allowances (not overtime/reimbursement)
2. Derive PTKP code → TER category (A/B/C)
3. Lookup TER rate by category + gross income bracket
4. Result: PayrollDetail with code=`PPH21`, group=`TAX`, type=`DEDUCTION`

TER categories mapping:
- Category A: TK/0, TK/1, K/0
- Category B: TK/2, TK/3, K/1, K/2
- Category C: K/3

### 2.2 BPJS Calculations

All percentages from `bpjs_rate_configs`. Employee-borne → DEDUCTION. Employer-borne → ALLOWANCE with `is_employer_borne=true`.

| Component | Employee | Employer | Cap | Code |
|-----------|----------|----------|-----|------|
| Kesehatan | 1% | 4% | Rp 12.000.000 | BPJS_KES |
| JHT | 2% | 3.7% | None | BPJS_JHT |
| JKK | — | 0.24%–1.74% | None | BPJS_JKK |
| JKM | — | 0.3% | None | BPJS_JKM |
| JP | 1% | 2% | Configurable | BPJS_JP |

Cap is applied: `min(actual_salary, max_salary_cap) × rate`.

### 2.3 Module Structure

Two new modules under `backend/internal/modules/`:

```
tax/
  entity.go        — PPh21Result, TERBracket, PTKPConfig
  repository.go    — read term_configs, ptkp_configs
  service.go       — CalculateTER, ReconcileAnnual, Generate1721A1
  handler.go       — admin CRUD for TER/PTKP configuration
  routes.go

bpjs/
  entity.go        — BPJSComponent, BPJSRateConfig
  repository.go    — read bpjs_rate_configs
  service.go       — CalculateAll
  handler.go       — admin CRUD for BPJS rate configuration
  routes.go
```

### 2.4 Interfaces

```go
type TaxCalculator interface {
    CalculateTER(ctx context.Context, grossIncome float64, employee Employee) (PPh21Result, error)
    ReconcileAnnual(ctx context.Context, employeeID uint, year uint) (AnnualSettlement, error)
    Generate1721A1(ctx context.Context, employeeID uint, year uint) (Form1721A1, error)
}

type BPJSCalculator interface {
    CalculateAll(ctx context.Context, grossIncome float64, companyID uint) ([]BPJSComponent, error)
}
```

---

## 3. Payroll Integration

### 3.1 Modified GenerateAll flow

```
For each employee:
  1. Gross income = baseSalary + fixed allowances
  2. Existing: overtime, reimbursement calculations
  3. bpjsCalculator.CalculateAll(grossIncome, companyID) → []BPJSComponent
  4. taxCalculator.CalculateTER(grossIncome, employee) → PPh21Result
  5. Existing: late penalty, loan installment
  6. Compose PayrollDetail entries
  7. Net = TotalAllowance - TotalDeduction (employer BPJS excluded)
```

### 3.2 DI Container

```go
container.Provide(tax.NewService)
container.Provide(tax.NewRepository)
container.Provide(bpjs.NewService)
container.Provide(bpjs.NewRepository)

container.Provide(func(
    // existing deps...
    taxCalc tax.TaxCalculator,
    bpjsCalc bpjs.BPJSCalculator,
) payroll.PayrollService {
    return payroll.NewService(..., taxCalc, bpjsCalc)
})
```

### 3.3 Payslip Layout

```
INCOME
  Base Salary              Rp 10.000.000
  Overtime                 Rp  1.500.000
  Reimbursement            Rp    500.000
  ─────────────────────────────────────
  Gross Income             Rp 12.000.000

DEDUCTIONS
  PPh 21                   Rp    250.000
  BPJS Kesehatan           Rp    120.000
  BPJS JHT                 Rp    240.000
  BPJS JP                  Rp    120.000
  Late Penalty             Rp     50.000
  ─────────────────────────────────────
  Total Deductions         Rp    780.000

NET SALARY                 Rp 11.220.000

EMPLOYER CONTRIBUTIONS
  BPJS Kesehatan           Rp    480.000
  BPJS JHT                 Rp    444.000
  BPJS JKK                 Rp     64.800
  BPJS JKM                 Rp     36.000
  BPJS JP                  Rp    240.000
```

---

## 4. Configuration UI

Two new admin pages under existing `frontend/src/pages/admin/` pattern.

### 4.1 BPJS Configuration (`/admin/bpjs-config`)

Per-component cards: employee rate (%), employer rate (%), max salary cap, effective dates. JKK has industry risk level dropdown. Global defaults vs company override indicator.

### 4.2 PPh 21 Configuration (`/admin/pph21-config`)

Tabbed: TER Rate Brackets (A/B/C editable tables) and PTKP Thresholds (code + annual amount). Import default TER table, export for audit.

### 4.3 Frontend structure

```
frontend/src/pages/admin/bpjs-config/
  types.ts
  hooks/useBpjsConfig.ts
  components/BpjsConfigForm.tsx
  index.tsx
frontend/src/pages/admin/pph21-config/
  types.ts
  hooks/usePph21Config.ts
  components/TerBracketsTable.tsx
  components/PtkpTable.tsx
  index.tsx
```

---

## 5. Year-End Reconciliation

### 5.1 Annual PPh 21 Settlement

```
Annual liability =
  Gross annual income (base salaries + regular allowances)
  - Biaya Jabatan (5% of gross, max Rp 6.000.000/year)
  - PTKP (annual threshold from ptkp_configs)
  = PKP (Taxable Income)

Progressive brackets on PKP:
  0 – 60.000.000            5%
  60.000.001 – 250.000.000  15%
  250.000.001 – 500.000.000 25%
  500.000.001 – 5.000.000.000 30%
  > 5.000.000.000          35%

Delta = TaxPayable - Σ TER paid YTD
  Delta > 0 → additional deduction
  Delta < 0 → refund
```

### 5.2 Trigger Points

| Trigger | Behavior |
|---------|----------|
| December payroll | PPh 21 entry = settlement amount (not TER). Delta from Jan–Nov TER paid. |
| Mid-year resignation | Reconcile Jan–termination month on final payslip. |
| Form 1721-A1 | Generated per employee after December settlement. |

### 5.3 Service Types

```go
type AnnualSettlement struct {
    GrossAnnual  float64
    BiayaJabatan float64
    PTKP         float64
    PKP          float64
    TaxPayable   float64
    TERPaidYTD   float64
    Delta        float64 // >0 underpaid, <0 overpaid
}
```

---

## 6. Out of Scope

- THR (religious holiday allowance) calculation
- Non-TER (progressive monthly) PPh 21 method
- PPh 21 for non-permanent employees (freelancers, consultants)
- BPJS contribution payment/disbursement automation (just calculation)
- Multi-country tax support
