# Subscription Proof of Payment Upload — Design Spec

**Date:** 2026-07-04  
**Status:** Approved, pending implementation

## Overview

Add proof of payment image upload to the subscription upgrade approval flow. When a superadmin approves a subscription upgrade request, they must upload a proof of payment image (transfer receipt). The image is uploaded to MinIO using the existing `UploadFileMultipart` infrastructure and stored as a URL on the `subscription_requests` table.

## Who

- **Superadmin** uploads the proof of payment image during the approval step.
- Upload happens in the existing "Setujui" dialog, as part of the approval action.

## What Changes (by layer)

### Database

New migration: add `proof_payment_url` column to `subscription_requests` table.

| Column | Type | Nullable | Notes |
|--------|------|----------|-------|
| `proof_payment_url` | VARCHAR(255) | YES | MinIO URL; set only on APPROVED requests |

### Backend (Go)

**Files changed: 6**

| File | Change |
|------|--------|
| `backend/internal/modules/subscription/entity.go` | Add `ProofPaymentURL *string` field on `SubscriptionRequest` struct with GORM tag |
| `backend/internal/modules/subscription/contract.go` | Add `StorageProvider` interface with `UploadFileMultipart(ctx, file, objectName) (string, error)` |
| `backend/internal/modules/subscription/dto.go` | Add `ProofPaymentURL *string` to `SubscriptionRequestResponse` |
| `backend/internal/modules/subscription/service.go` | Add `storage StorageProvider` field to `service` struct; wire into `NewService`; in `ReviewRequest`, when status is `APPROVED`: validate file present, generate object name `"subscriptions/requests/{id}/proof-payment-{timestamp}.jpg"`, upload to MinIO, store URL in entity before saving |
| `backend/internal/modules/subscription/handler.go` | Change `ReviewRequest` handler: read `status` and `notes` from `ctx.FormValue(...)`, read file from `ctx.FormFile("proof_payment")`, validate max 5MB, require file when status is `APPROVED` |
| `backend/internal/bootstrap/container.go` | Pass `storage` to `subscription.NewService(...)` |

**Key behavior:**
- When `status == "APPROVED"`: `proof_payment` file is **required** and uploaded to MinIO. URL stored in `ProofPaymentURL`.
- When `status == "REJECTED"`: no file needed, no upload. `ProofPaymentURL` stays nil.
- File max size: 5MB (same limit as reimbursement).
- Image is auto-resized to 800px width, JPEG quality 75% by the MinIO provider (existing behavior).

### Frontend (TypeScript/React)

**Files changed: 3**

| File | Change |
|------|--------|
| `frontend/src/features/subscription/types.ts` | Add `proof_payment_url?: string` to `SubscriptionRequestItem` |
| `frontend/src/features/subscription/hooks/useSubscription.ts` | Change `useReviewRequest` mutation to send `FormData` instead of JSON (append `status`, `notes`, and `proof_payment` file) |
| `frontend/src/pages/admin/SubscriptionAdminPage.tsx` | Add file input (hidden, `accept="image/*"`) + `useRef`, file state, and preview to the approval dialog. Validate file is selected for APPROVED. Show image thumbnail on file select. Client-side 5MB validation. |

**Dialog behavior:**
- APPROVED mode: shows a file upload area + thumbnail preview below the notes textarea. File is required.
- REJECTED mode: no file input shown (unchanged behavior).

### Routes

No new routes. The existing `PUT /admin/subscriptions/:id/review` is changed from `application/json` to `multipart/form-data`.

## Migration Strategy

- Backward compatible: old rows have `proof_payment_url = NULL`.
- The new column is added with `DEFAULT NULL` — no data migration needed for existing requests.
- The frontend is updated atomically with the backend — both must deploy together since the review endpoint changes content type.

## Testing

Backend unit tests (service layer):
- Test `ReviewRequest` with APPROVED status and valid file → proof URL is set, upload is called
- Test `ReviewRequest` with APPROVED status and no file → returns error
- Test `ReviewRequest` with REJECTED status → no upload call, proof URL stays nil
- Test file > 5MB → handler returns error

Frontend: manual testing of dialog, file select, preview, and submit.
