# Subscription Request Detail & History View — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a request detail view with proof image preview and restructure the subscription admin page from "pending only" to "all requests with status filter."

**Architecture:** Two small changes: (1) add `proof_payment_url` to existing raw SQL SELECT queries in the repository so the column reaches the API, (2) rework the frontend page to use `useAllRequests`, add a status filter, Status column, Detail button, and a Detail dialog with inline proof image.

**Tech Stack:** Go (GORM), TypeScript/React (React Query, Tailwind CSS, shadcn/ui)

---

## File Structure

| File | Action | Responsibility |
|------|--------|----------------|
| `backend/internal/modules/subscription/repository.go:69-91` | Modify | Add `proof_payment_url` to `FindAllPendingRequests` SELECT |
| `backend/internal/modules/subscription/repository.go:98-119` | Modify | Add `proof_payment_url` to `FindAllRequests` SELECT |
| `frontend/src/pages/admin/SubscriptionAdminPage.tsx` | Modify | Restructure page: all-requests view, status filter, Status column, Detail button, Detail dialog |

---

### Task 1: Repository — Add proof_payment_url to SELECT queries

**Files:**
- Modify: `backend/internal/modules/subscription/repository.go:69-91` (FindAllPendingRequests)
- Modify: `backend/internal/modules/subscription/repository.go:98-119` (FindAllRequests)

- [ ] **Step 1: Add proof_payment_url to FindAllPendingRequests SELECT**

In `FindAllPendingRequests` (around line 82), add this line after `subscription_requests.notes,`:

```go
subscription_requests.proof_payment_url,
```

The SELECT should go from:
```go
subscription_requests.notes,
subscription_requests.created_at`).
```
To:
```go
subscription_requests.notes,
subscription_requests.proof_payment_url,
subscription_requests.created_at`).
```

- [ ] **Step 2: Add proof_payment_url to FindAllRequests SELECT**

In `FindAllRequests` (around line 111), add the same line after `subscription_requests.notes,`:

```go
subscription_requests.proof_payment_url,
```

- [ ] **Step 3: Build and verify**

```bash
cd backend && go build ./internal/modules/subscription/
```

Expected: no errors.

- [ ] **Step 4: Commit**

```bash
git add backend/internal/modules/subscription/repository.go
git commit -m "feat(subscription): add proof_payment_url to request list queries"
```

---

### Task 2: Frontend Page — All-requests view with status filter and Detail dialog

**Files:**
- Modify: `frontend/src/pages/admin/SubscriptionAdminPage.tsx`

This is a significant rework. The implementer should read the current file first, then apply the changes described below in order.

- [ ] **Step 1: Update imports**

Add `Eye` to lucide-react import:
```typescript
import { Crown, Loader2, CheckCircle2, XCircle, Search, Eye } from "lucide-react";
```

Add `useAllRequests` to the hook import (keep `useReviewRequest`):
```typescript
import {
  useAllRequests,
  useReviewRequest,
} from "@/features/subscription/hooks/useSubscription";
```

Add Select components import:
```typescript
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
```

- [ ] **Step 2: Add new state variables**

Add after existing state declarations (around line 52):

```typescript
const [statusFilter, setStatusFilter] = useState<string>("ALL");
const [detailRequest, setDetailRequest] = useState<SubscriptionRequestItem | null>(null);
const [isDetailOpen, setIsDetailOpen] = useState(false);
```

- [ ] **Step 3: Switch data source**

Change line 38 from:
```typescript
const { data, isLoading } = usePendingRequests();
```
To:
```typescript
const { data, isLoading } = useAllRequests();
```

- [ ] **Step 4: Update filter logic**

Replace the `filtered` computation (lines 76-88):

```typescript
const allRequests: SubscriptionRequestItem[] = data?.data || [];

const filtered = allRequests.filter((r) => {
  const matchesSearch =
    !debouncedSearch ||
    r.company_name.toLowerCase().includes(debouncedSearch.toLowerCase()) ||
    r.requested_plan_name.toLowerCase().includes(debouncedSearch.toLowerCase());
  const matchesStatus = statusFilter === "ALL" || r.status === statusFilter;
  return matchesSearch && matchesStatus;
});
```

- [ ] **Step 5: Add openDetail handler**

```typescript
function openDetail(request: SubscriptionRequestItem) {
  setDetailRequest(request);
  setIsDetailOpen(true);
}
```

- [ ] **Step 6: Update header badge text**

Change the badge (line 141) from:
```tsx
{allRequests.length} permintaan menunggu
```
To:
```tsx
{filtered.length} permintaan
```

- [ ] **Step 7: Add status filter Select in CardHeader**

Add this between the CardTitle and the search input in the CardHeader (around line 152):

```tsx
<Select value={statusFilter} onValueChange={setStatusFilter}>
  <SelectTrigger className="w-40">
    <SelectValue placeholder="Status" />
  </SelectTrigger>
  <SelectContent>
    <SelectItem value="ALL">Semua</SelectItem>
    <SelectItem value="PENDING">Menunggu</SelectItem>
    <SelectItem value="APPROVED">Disetujui</SelectItem>
    <SelectItem value="REJECTED">Ditolak</SelectItem>
  </SelectContent>
</Select>
```

Place this after the CardTitle div and before the search input div, inside the same `flex` row.

- [ ] **Step 8: Add "Status" column to desktop table**

In the TableHeader `<TableRow>`, add after `<TableHead>Tanggal</TableHead>`:

```tsx
<TableHead>Status</TableHead>
```

In each TableRow's `<TableCell>` section (desktop), add after the date cell and before the aksi cell:

```tsx
<TableCell>
  <Badge
    className={
      req.status === "APPROVED"
        ? "bg-emerald-100 text-emerald-800 hover:bg-emerald-100"
        : req.status === "REJECTED"
        ? "bg-red-100 text-red-800 hover:bg-red-100"
        : "text-amber-600"
    }
    variant={req.status === "PENDING" ? "outline" : undefined}
  >
    {req.status === "APPROVED"
      ? "Disetujui"
      : req.status === "REJECTED"
      ? "Ditolak"
      : "Menunggu"}
  </Badge>
</TableCell>
```

- [ ] **Step 9: Add Detail button + conditionally show approve/reject in desktop table**

Replace the "Aksi" cell in each TableRow. Currently it always shows Setujui/Tolak buttons. Change to:

```tsx
<TableCell className="text-right">
  <div className="flex justify-end gap-2">
    {req.status === "PENDING" && (
      <>
        <Button
          size="sm"
          variant="ghost"
          className="h-8 px-2 text-emerald-600 hover:text-emerald-700 hover:bg-emerald-50"
          onClick={() => handleReview(req, "APPROVED")}
          disabled={isReviewing}
        >
          <CheckCircle2 className="h-4 w-4 mr-1" />
          Setujui
        </Button>
        <Button
          size="sm"
          variant="ghost"
          className="h-8 px-2 text-red-600 hover:text-red-700 hover:bg-red-50"
          onClick={() => handleReview(req, "REJECTED")}
          disabled={isReviewing}
        >
          <XCircle className="h-4 w-4 mr-1" />
          Tolak
        </Button>
      </>
    )}
    <Button
      size="sm"
      variant="ghost"
      className="h-8 px-2 text-slate-500 hover:text-slate-700 hover:bg-slate-100"
      onClick={() => openDetail(req)}
    >
      <Eye className="h-4 w-4 mr-1" />
      Detail
    </Button>
  </div>
</TableCell>
```

- [ ] **Step 10: Update mobile card layout**

In the mobile view (`<div className="md:hidden space-y-4">`), each request card needs:
- Dynamic status badge (not just "Menunggu")
- Detail button always visible
- Approve/reject buttons only for PENDING

Replace each mobile card's content:

```tsx
<div className="flex items-center justify-between">
  <span className="font-medium">{req.company_name}</span>
  <Badge
    className={
      req.status === "APPROVED"
        ? "bg-emerald-100 text-emerald-800 hover:bg-emerald-100"
        : req.status === "REJECTED"
        ? "bg-red-100 text-red-800 hover:bg-red-100"
        : "text-amber-600"
    }
    variant={req.status === "PENDING" ? "outline" : undefined}
  >
    {req.status === "APPROVED"
      ? "Disetujui"
      : req.status === "REJECTED"
      ? "Ditolak"
      : "Menunggu"}
  </Badge>
</div>
```

And the button row at the bottom should become:
```tsx
<div className="flex gap-2 pt-2 border-t">
  {req.status === "PENDING" ? (
    <>
      <Button
        size="sm"
        className="flex-1"
        onClick={() => handleReview(req, "APPROVED")}
        disabled={isReviewing}
      >
        <CheckCircle2 className="w-4 h-4 mr-1" />
        Setujui
      </Button>
      <Button
        size="sm"
        variant="destructive"
        className="flex-1"
        onClick={() => handleReview(req, "REJECTED")}
        disabled={isReviewing}
      >
        <XCircle className="w-4 h-4 mr-1" />
        Tolak
      </Button>
    </>
  ) : (
    <div className="w-full" />
  )}
  <Button
    size="sm"
    variant="outline"
    className="flex-shrink-0"
    onClick={() => openDetail(req)}
  >
    <Eye className="w-4 h-4 mr-1" />
    Detail
  </Button>
</div>
```

- [ ] **Step 11: Add Detail dialog (after existing review dialog)**

Add this new Dialog component after the closing `</Dialog>` of the review dialog and before the closing `</div>` of the page wrapper:

```tsx
<Dialog open={isDetailOpen} onOpenChange={setIsDetailOpen}>
  <DialogContent className="max-w-lg">
    <DialogHeader>
      <DialogTitle>Detail Permintaan</DialogTitle>
    </DialogHeader>

    {detailRequest && (
      <div className="space-y-4 py-2">
        <div className="grid grid-cols-2 gap-3 text-sm">
          <div>
            <p className="text-slate-500">Perusahaan</p>
            <p className="font-medium">{detailRequest.company_name}</p>
          </div>
          <div>
            <p className="text-slate-500">Status</p>
            <Badge
              className={
                detailRequest.status === "APPROVED"
                  ? "bg-emerald-100 text-emerald-800 hover:bg-emerald-100"
                  : detailRequest.status === "REJECTED"
                  ? "bg-red-100 text-red-800 hover:bg-red-100"
                  : "text-amber-600"
              }
              variant={detailRequest.status === "PENDING" ? "outline" : undefined}
            >
              {detailRequest.status === "APPROVED"
                ? "Disetujui"
                : detailRequest.status === "REJECTED"
                ? "Ditolak"
                : "Menunggu"}
            </Badge>
          </div>
          <div>
            <p className="text-slate-500">Paket Saat Ini</p>
            <p className="font-medium">{detailRequest.current_plan_name}</p>
          </div>
          <div>
            <p className="text-slate-500">Paket Diminta</p>
            <p className="font-medium text-blue-700">{detailRequest.requested_plan_name}</p>
          </div>
          <div>
            <p className="text-slate-500">Harga Saat Ini</p>
            <p className="font-medium">
              {detailRequest.current_plan_price.toLocaleString("id-ID", {
                style: "currency",
                currency: "IDR",
                maximumFractionDigits: 0,
              })}
            </p>
          </div>
          <div>
            <p className="text-slate-500">Harga Diminta</p>
            <p className="font-medium text-blue-700">
              {detailRequest.requested_plan_price.toLocaleString("id-ID", {
                style: "currency",
                currency: "IDR",
                maximumFractionDigits: 0,
              })}
            </p>
          </div>
          <div>
            <p className="text-slate-500">Selisih</p>
            <p className="font-medium text-blue-700">
              +
              {detailRequest.price_difference.toLocaleString("id-ID", {
                style: "currency",
                currency: "IDR",
                maximumFractionDigits: 0,
              })}
              /bln
            </p>
          </div>
          <div>
            <p className="text-slate-500">Tanggal</p>
            <p className="font-medium">
              {new Date(detailRequest.created_at).toLocaleDateString("id-ID", {
                day: "numeric",
                month: "short",
                year: "numeric",
              })}
            </p>
          </div>
          <div>
            <p className="text-slate-500">Diminta Oleh</p>
            <p className="font-medium">{detailRequest.requested_by_name}</p>
          </div>
        </div>

        <div>
          <p className="text-slate-500 text-sm mb-1">Catatan</p>
          <p className="text-sm bg-slate-50 rounded-lg p-3">
            {detailRequest.notes || "-"}
          </p>
        </div>

        <div>
          <p className="text-slate-500 text-sm mb-2">Bukti Pembayaran</p>
          {detailRequest.proof_payment_url ? (
            <div className="space-y-2">
              <a
                href={detailRequest.proof_payment_url}
                target="_blank"
                rel="noopener noreferrer"
              >
                <img
                  src={detailRequest.proof_payment_url}
                  alt="Bukti Pembayaran"
                  className="w-full max-h-[300px] object-contain rounded-lg border border-slate-200 hover:opacity-90 transition-opacity cursor-pointer"
                />
              </a>
              <a
                href={detailRequest.proof_payment_url}
                target="_blank"
                rel="noopener noreferrer"
                className="text-xs text-blue-600 hover:text-blue-800 underline inline-block"
              >
                Lihat Gambar Penuh
              </a>
            </div>
          ) : (
            <p className="text-sm text-slate-400 italic">
              Tidak ada bukti pembayaran.
            </p>
          )}
        </div>
      </div>
    )}

    <DialogFooter>
      <Button variant="outline" onClick={() => setIsDetailOpen(false)}>
        Tutup
      </Button>
    </DialogFooter>
  </DialogContent>
</Dialog>
```

- [ ] **Step 12: Verify TypeScript**

```bash
cd frontend && npx tsc --noEmit
```

Expected: no errors.

- [ ] **Step 13: Verify frontend build**

```bash
cd frontend && npm run build
```

Expected: build succeeds.

- [ ] **Step 14: Commit**

```bash
git add frontend/src/pages/admin/SubscriptionAdminPage.tsx
git commit -m "feat(subscription): add all-requests view with status filter and detail dialog"
```

---

### Task 3: Final verification

- [ ] **Step 1: Build backend**

```bash
cd backend && go build ./...
```

- [ ] **Step 2: Run backend tests**

```bash
cd backend && go test ./internal/modules/subscription/ -v -count=1
```

- [ ] **Step 3: Verify git status clean**

```bash
git status
```
