# Design: Extract Department into Separate Module

**Date:** 2026-06-12
**Status:** Approved

## Problem

The `master` module bundles three unrelated data types: Department, Shift, and LeaveType.
Department has full CRUD, while Shift and LeaveType are read-only lookups. This mixing
violates single-responsibility and makes the module grow unnecessarily.

## Solution

Extract Department into its own `department` module, leaving Shift and LeaveType in `master`.

### Architecture

New module `backend/internal/modules/department/` following the existing module pattern:

| File | Purpose |
|------|---------|
| `entity.go` | `Department` struct, GORM tags, `TableName() = "ref_departments"` |
| `dto.go` | `LookupResponse`, `CreateDepartmentRequest`, `UpdateDepartmentRequest` |
| `contract.go` | `CacheProvider` interface (same signature as master's) |
| `repository.go` | `Repository` interface + GORM implementation (CRUD, finders, existence check, employee count) |
| `service.go` | `Service` interface + implementation (business logic, caching, validation) |
| `handler.go` | Echo HTTP handlers for all department endpoints |

Deleted from `master`: `Department` entity, all department methods in repo/service/handler/DTOs,
mocks, and tests. SeedDefaults loses its department creation line.

### Routes

New `routes/department.go`, registered at `/api/v1/departments`:

| Method | Path | Permission |
|--------|------|------------|
| GET | `/departments` | `VIEW_MASTER` |
| GET | `/departments/:id` | `VIEW_MASTER` |
| POST | `/departments` | `MANAGE_MASTER` |
| PUT | `/departments/:id` | `MANAGE_MASTER` |
| DELETE | `/departments/:id` | `MANAGE_MASTER` |

Permissions reuse existing `VIEW_MASTER` / `MANAGE_MASTER` constants — no new permission keys.

Master routes lose 5 department endpoints, keep only shifts and leave types under `/masters`.

### Container Wiring

- `departmentRepo := department.NewRepository(db.GetDB())`
- `departmentSvc := department.NewService(departmentRepo, redis)`
- `departmentHandler := department.NewHandler(departmentSvc)`
- `container.DepartmentHandler = departmentHandler`
- `onboardingSvc` wired with `departmentRepo` instead of `masterRepo`

### Consumer Updates

**auth module:** `SeedDefaults` keeps shift and leave type seeding, drops department seeding.
The auth module's `MasterProvider` interface and wiring are unchanged (it calls `SeedDefaults`
on the master repo, which still exists).

**onboarding module:** Uses `department.Repository.FindDepartmentByName` instead of
`master.Repository.FindDepartmentByName`. A department-specific contract interface
is defined in onboarding's `contract.go`.

### Frontend

- Consolidate duplicate `useDepartments` hooks (`useAdmin.tsx` and `useMasterData.tsx`)
  into a single definition in `useAdmin.tsx`. `useMasterData.tsx` imports from there.
- Update all department API paths from `/masters/departments` to `/departments`.

Affected files:
- `frontend/src/features/admin/hooks/useAdmin.tsx`
- `frontend/src/features/admin/hooks/useMasterData.tsx`

All component consumers of `useDepartments` work without changes.

### Testing

New `department/` module gets `service_test.go`, `handler_test.go`, `repository_test.go`,
`mocks_test.go` — ported from current master tests, adapted to new package.

Master tests are stripped of all department test cases.

### What Does NOT Change

- No DB schema changes — `ref_departments` table stays as-is
- No migration files
- No new permission constants
- No seeder changes
- No changes to other consumer modules (attendance, payroll, etc.)

## Change Summary

| Layer | Add | Modify | Delete (from master) |
|-------|-----|--------|---------------------|
| Module | `department/` (6 files + tests + mocks) | — | Department entity, CRUD, tests, mocks |
| Routes | `routes/department.go` | `routes/master.go`, `routes/api.go` | 5 department endpoints |
| Container | — | `container.go` (repo, svc, handler wiring) | — |
| Consumers | — | `auth` (remove dept from seed), `onboarding` (new dep provider) | Cross-module refs to master departments |
| Frontend | — | Consolidate `useDepartments`, update URLs | — |
