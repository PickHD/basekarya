# Subscription Request Detail & History View — Design Spec

**Date:** 2026-07-04  
**Status:** Approved, pending implementation

## Overview

Add a request detail view to the subscription admin page so superadmins can recheck approved/rejected requests and view the proof of payment image. The page is restructured from "pending only" to "all requests with status filter," and each row gets a Detail button that opens a dialog showing full request info including inline proof image preview.

## What Changes

### Page Restructure

- Switch from `usePendingRequests()` → `useAllRequests()`
- Add a `Select` status filter: Semua / PENDING / APPROVED / REJECTED (default: Semua)
- Add "Status" column to the table with colored badges
- "Aksi" column: Setujui/Tolak buttons for PENDING rows only; **Detail** button for all rows

### Detail Dialog

A new dialog (separate from the review dialog) triggered by the Detail button. Shows:

| Field | Source |
|-------|--------|
| Perusahaan | `company_name` |
| Paket Saat Ini → Diminta | `current_plan_name` → `requested_plan_name` |
| Harga Saat Ini | `current_plan_price` formatted as IDR |
| Harga Diminta | `requested_plan_price` formatted as IDR |
| Selisih Harga | `price_difference` formatted as IDR |
| Status | colored badge |
| Diminta Oleh | `requested_by_name` |
| Tanggal | `created_at` formatted date |
| Catatan | `notes` (or "-" if empty) |
| Bukti Pembayaran | inline image preview if `proof_payment_url` exists (max-height 300px, object-fit contain), with "Lihat Gambar Penuh" link opening in new tab. If no proof: "Tidak ada bukti pembayaran." |

Dialog has a single "Tutup" button.

### Backend

| File | Change |
|------|--------|
| `backend/internal/modules/subscription/repository.go` | Add `subscription_requests.proof_payment_url` to SELECT clause in `FindAllRequests` and `FindAllPendingRequests` |
| `backend/migrations/000028_add_column_proof_payment_url.up.sql` | Already created by user — no change |

No handler/service/route changes needed. The `SubscriptionRequestResponse` DTO already has `ProofPaymentURL *string`.

### Frontend

| File | Change |
|------|--------|
| `frontend/src/pages/admin/SubscriptionAdminPage.tsx` | Major rework: switch to `useAllRequests`, add status filter, Status column, Detail button, Detail dialog with proof image |

## Migration

The SQL migration (`000028`) is already created. The GORM entity field (`ProofPaymentURL *string`) was added in previous work. Only the repository SELECT queries need updating to include the column in the raw SQL queries.

## Testing

- Manual: verify status filter works, Detail dialog shows correct info, proof image renders for approved requests, click-to-enlarge opens new tab
- Unit: no new backend logic to test (only SELECT column additions)
