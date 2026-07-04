# Subscription Proof of Payment Upload — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add required proof of payment image upload to the subscription upgrade approval dialog for superadmins.

**Architecture:** Add a `proof_payment_url` column to `subscription_requests` (GORM auto-migrated). The `PUT /admin/subscriptions/:id/review` endpoint switches from JSON to `multipart/form-data` to accept a file alongside `status` and `notes`. The service uploads images to MinIO using the existing `UploadFileMultipart` method. Frontend dialog adds a required file input with thumbnail preview when approving.

**Tech Stack:** Go (Echo + GORM + MinIO), TypeScript/React (React Query + Axios + Tailwind CSS)

---

## File Structure

| File | Action | Responsibility |
|------|--------|----------------|
| `backend/internal/modules/subscription/entity.go` | Modify | Add `ProofPaymentURL` GORM field |
| `backend/internal/modules/subscription/contract.go` | Modify | Add `StorageProvider` interface |
| `backend/internal/modules/subscription/dto.go` | Modify | Add `ProofPaymentURL` to response DTO |
| `backend/internal/modules/subscription/service.go` | Modify | Accept `StorageProvider`, upload proof in `ReviewRequest` |
| `backend/internal/modules/subscription/handler.go` | Modify | Parse multipart form in `ReviewRequest` |
| `backend/internal/bootstrap/container.go` | Modify | Wire `storage` into subscription service |
| `backend/internal/modules/subscription/mocks_test.go` | Modify | Add `mockStorageProvider`, update `newTestSubscriptionService` |
| `backend/internal/modules/subscription/service_test.go` | Modify | Add proof upload test cases to `TestService_ReviewRequest` |
| `frontend/src/features/subscription/types.ts` | Modify | Add `proof_payment_url` field |
| `frontend/src/features/subscription/hooks/useSubscription.ts` | Modify | Switch `useReviewRequest` to `FormData` |
| `frontend/src/pages/admin/SubscriptionAdminPage.tsx` | Modify | Add file input + preview to dialog |

---

### Task 1: Entity — Add `ProofPaymentURL` field

**Files:**
- Modify: `backend/internal/modules/subscription/entity.go`

- [ ] **Step 1: Add field to SubscriptionRequest struct**

```go
type SubscriptionRequest struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	CompanyID       uint       `gorm:"not null" json:"company_id"`
	CurrentPlanID   uint       `gorm:"not null" json:"current_plan_id"`
	RequestedPlanID uint       `gorm:"not null" json:"requested_plan_id"`
	Status          string     `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	RequestedBy     *uint      `json:"requested_by"`
	ReviewedBy      *uint      `json:"reviewed_by"`
	ReviewedAt      *time.Time `json:"reviewed_at"`
	Notes           string     `gorm:"type:text" json:"notes"`
	ProofPaymentURL *string    `gorm:"type:varchar(255)" json:"proof_payment_url"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}
```

- [ ] **Step 2: Commit**

```bash
git add backend/internal/modules/subscription/entity.go
git commit -m "feat(subscription): add ProofPaymentURL field to SubscriptionRequest entity"
```

---

### Task 2: Contract — Add StorageProvider interface

**Files:**
- Modify: `backend/internal/modules/subscription/contract.go`

- [ ] **Step 1: Add StorageProvider interface**

Add after the existing imports and interfaces:

```go
type StorageProvider interface {
	UploadFileMultipart(ctx context.Context, file *multipart.FileHeader, objectName string) (string, error)
}
```

The import for `"mime/multipart"` and `"context"` are already present. Verify `multipart` is imported (it might not be yet — if compile error, add `"mime/multipart"` to imports).

- [ ] **Step 2: Verify imports compile**

```bash
cd backend && go build ./internal/modules/subscription/
```

- [ ] **Step 3: Commit**

```bash
git add backend/internal/modules/subscription/contract.go
git commit -m "feat(subscription): add StorageProvider interface to contract"
```

---

### Task 3: DTO — Add ProofPaymentURL to response

**Files:**
- Modify: `backend/internal/modules/subscription/dto.go`

- [ ] **Step 1: Add field to SubscriptionRequestResponse**

```go
type SubscriptionRequestResponse struct {
	ID                 uint    `json:"id"`
	CompanyName        string  `json:"company_name"`
	CurrentPlanName    string  `json:"current_plan_name"`
	RequestedPlanName  string  `json:"requested_plan_name"`
	CurrentPlanPrice   float64 `json:"current_plan_price"`
	RequestedPlanPrice float64 `json:"requested_plan_price"`
	PriceDifference    float64 `json:"price_difference"`
	Status             string  `json:"status"`
	RequestedByName    string  `json:"requested_by_name"`
	RequestedByEmail   string  `json:"requested_by_email"`
	Notes              string  `json:"notes"`
	ProofPaymentURL    *string `json:"proof_payment_url"`
	CreatedAt          string  `json:"created_at"`
}
```

- [ ] **Step 2: Commit**

```bash
git add backend/internal/modules/subscription/dto.go
git commit -m "feat(subscription): add ProofPaymentURL to response DTO"
```

---

### Task 4: Service — Wire storage and upload proof

**Files:**
- Modify: `backend/internal/modules/subscription/service.go`

- [ ] **Step 1: Add `storage` field to service struct**

```go
type service struct {
	repo    Repository
	company company.Repository
	role    RoleProvider
	user    UserProvider
	cache   CacheProvider
	storage StorageProvider
}
```

- [ ] **Step 2: Update NewService signature and constructor**

```go
func NewService(repo Repository, companyRepo company.Repository, role RoleProvider, user UserProvider, cache CacheProvider, storage StorageProvider) Service {
	return &service{repo: repo, company: companyRepo, role: role, user: user, cache: cache, storage: storage}
}
```

- [ ] **Step 3: Add import for "mime/multipart" and "fmt" at top (alphabetize)**

Add `"mime/multipart"` to imports if not present. `"fmt"` is already imported.

- [ ] **Step 4: Update ReviewRequest interface method signature on Service**

```go
type Service interface {
	// ... existing methods ...
	ReviewRequest(ctx context.Context, requestID uint, req *ReviewRequest, proofPayment *multipart.FileHeader) error
	// ... rest ...
}
```

- [ ] **Step 5: Update ReviewRequest implementation to accept file and upload proof**

Replace the existing `ReviewRequest` method body. Insert the proof upload logic after the request status/meta update and before the approval block:

```go
func (s *service) ReviewRequest(ctx context.Context, requestID uint, req *ReviewRequest, proofPayment *multipart.FileHeader) error {
	subReq, err := s.repo.FindRequestByID(ctx, requestID)
	if err != nil {
		return errors.New("request not found")
	}

	if subReq.Status != constants.SubReqStatusPending {
		return errors.New("request already reviewed")
	}

	if req.Status == constants.SubReqStatusApproved {
		if proofPayment == nil {
			return errors.New("proof of payment is required")
		}

		objectName := fmt.Sprintf("subscriptions/requests/%d/proof-payment-%d.jpg", requestID, time.Now().UnixNano())
		proofURL, err := s.storage.UploadFileMultipart(ctx, proofPayment, objectName)
		if err != nil {
			return fmt.Errorf("failed to upload proof of payment: %w", err)
		}
		subReq.ProofPaymentURL = &proofURL
	}

	reviewerID := utils.GetUserIDFromCtx(ctx)
	now := time.Now()

	subReq.Status = req.Status
	subReq.ReviewedBy = &reviewerID
	subReq.ReviewedAt = &now
	subReq.Notes = req.Notes

	if err := s.repo.UpdateRequest(ctx, subReq); err != nil {
		return fmt.Errorf("failed to update request: %w", err)
	}

	if req.Status == constants.SubReqStatusApproved {
		comp, err := s.company.FindByID(ctx, subReq.CompanyID)
		if err != nil {
			return errors.New("company not found")
		}

		comp.SubscriptionPlanID = &subReq.RequestedPlanID
		comp.SubscriptionStatus = constants.SubStatusActive
		expiresAt := time.Now().Add(30 * 24 * time.Hour)
		comp.SubscriptionExpiresAt = &expiresAt

		if err := s.company.Update(ctx, comp); err != nil {
			return fmt.Errorf("failed to update company plan: %w", err)
		}

		plan, err := s.repo.FindPlanByID(ctx, subReq.RequestedPlanID)
		if err == nil && plan != nil {
			allowedGroups := buildAllowedGroups(plan.Slug)
			permissionIDs, err := s.role.FindPermissionIDsByGroupNames(ctx, allowedGroups)
			if err == nil && len(permissionIDs) > 0 {
				roleIDs, _ := s.role.FindRoleIDsByCompanyID(ctx, subReq.CompanyID)
				for _, roleID := range roleIDs {
					_ = s.role.AssignPermissions(ctx, roleID, permissionIDs, subReq.CompanyID)
				}
			}
		}

		_ = s.cache.Del(ctx, fmt.Sprintf(constants.SUBSCRIPTION_FEATURES_CACHE_KEY, subReq.CompanyID))
		_ = s.cache.Del(ctx, fmt.Sprintf(constants.COMPANY_PROFILE_CACHE_KEY, subReq.CompanyID))

		_ = s.user.ForceResetPasswordByCompanyID(ctx, subReq.CompanyID)
	}

	return nil
}
```

- [ ] **Step 6: Verify compilation**

```bash
cd backend && go build ./internal/modules/subscription/
```

- [ ] **Step 7: Commit**

```bash
git add backend/internal/modules/subscription/service.go
git commit -m "feat(subscription): add proof of payment upload to ReviewRequest"
```

---

### Task 5: Handler — Switch ReviewRequest to multipart form

**Files:**
- Modify: `backend/internal/modules/subscription/handler.go`

- [ ] **Step 1: Add import for "mime/multipart" (if needed) and "constants"**

The imports need `"mime/multipart"` and `"basekarya-backend/pkg/constants"`. Check current imports and add if missing.

- [ ] **Step 2: Replace ReviewRequest handler to parse multipart form**

Replace the entire `ReviewRequest` method:

```go
func (h *Handler) ReviewRequest(ctx echo.Context) error {
	id := ctx.Param("id")
	var requestID uint
	if _, err := fmt.Sscanf(id, "%d", &requestID); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid request ID", nil, err, nil)
	}

	status := ctx.FormValue("status")
	notes := ctx.FormValue("notes")

	req := &ReviewRequest{Status: status, Notes: notes}
	if err := ctx.Validate(req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}

	var proofPayment *multipart.FileHeader
	if status == constants.SubReqStatusApproved {
		fileHeader, err := ctx.FormFile("proof_payment")
		if err != nil {
			return response.NewResponses[any](ctx, http.StatusBadRequest, "Proof of payment is required", nil, err, nil)
		}
		if fileHeader.Size > 5*1024*1024 {
			return response.NewResponses[any](ctx, http.StatusBadRequest, "File size exceeds 5MB limit", nil, nil, nil)
		}
		proofPayment = fileHeader
	}

	if err := h.service.ReviewRequest(ctx.Request().Context(), requestID, req, proofPayment); err != nil {
		logger.Errorw("Failed to review request: ", err)
		return response.NewResponses[any](ctx, http.StatusBadRequest, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Request reviewed successfully", nil, nil, nil)
}
```

- [ ] **Step 3: Verify compilation**

```bash
cd backend && go build ./internal/modules/subscription/
```

- [ ] **Step 4: Commit**

```bash
git add backend/internal/modules/subscription/handler.go
git commit -m "feat(subscription): switch ReviewRequest handler to multipart form upload"
```

---

### Task 6: Container — Wire storage into subscription service

**Files:**
- Modify: `backend/internal/bootstrap/container.go`

- [ ] **Step 1: Pass storage to subscription.NewService**

On the line that creates `subscriptionSvc`:

```go
subscriptionSvc := subscription.NewService(subscriptionRepo, companyRepo, rbacRepo, userRepo, planCache, storage)
```

The `storage` variable is already declared and in scope (it's the `*infrastructure.MinioStorageProvider` used by other services).

- [ ] **Step 2: Verify compilation**

```bash
cd backend && go build ./...
```

- [ ] **Step 3: Commit**

```bash
git add backend/internal/bootstrap/container.go
git commit -m "feat(subscription): wire storage provider into subscription service"
```

---

### Task 7: Mock — Add mockStorageProvider for tests

**Files:**
- Modify: `backend/internal/modules/subscription/mocks_test.go`

- [ ] **Step 1: Add mockStorageProvider struct**

Add at end of file, before `newTestSubscriptionService`:

```go
type mockStorage struct{ mock.Mock }

func (m *mockStorage) UploadFileMultipart(ctx context.Context, file *multipart.FileHeader, objectName string) (string, error) {
	args := m.Called(ctx, file, objectName)
	return args.String(0), args.Error(1)
}
```

This requires adding `"mime/multipart"` to imports.

- [ ] **Step 2: Update newTestSubscriptionService to return the mock storage**

```go
func newTestSubscriptionService() (Service, *mockRepo, *mockCompanyRepo, *mockRole, *mockUser, *mockCache, *mockStorage) {
	repo := new(mockRepo)
	companyRepo := new(mockCompanyRepo)
	role := new(mockRole)
	user := new(mockUser)
	cache := new(mockCache)
	storage := new(mockStorage)

	svc := NewService(repo, companyRepo, role, user, cache, storage)
	return svc, repo, companyRepo, role, user, cache, storage
}
```

- [ ] **Step 3: Commit**

```bash
git add backend/internal/modules/subscription/mocks_test.go
git commit -m "test(subscription): add mockStorageProvider and update test constructor"
```

---

### Task 8: Service Tests — Add proof upload test cases

**Files:**
- Modify: `backend/internal/modules/subscription/service_test.go`

- [ ] **Step 1: Add import for "mime/multipart" at top**

```go
import (
	stdcontext "context"
	"errors"
	"mime/multipart"
	"testing"
	// ... rest of imports
)
```

- [ ] **Step 2: Update test table struct to include file field**

Add `proofPayment *multipart.FileHeader` field to the test struct inside `TestService_ReviewRequest`:

```go
tests := []struct {
	name         string
	requestID    uint
	ctx          func() stdcontext.Context
	req          *ReviewRequest
	proofPayment *multipart.FileHeader
	setupMocks   func(*mockRepo, *mockCompanyRepo, *mockRole, *mockUser, *mockCache, *mockStorage)
	wantErr      bool
	errMsg       string
}{
```

- [ ] **Step 3: Add new test cases for proof upload**

Insert these into the `tests` slice:

```go
{
	name:      "approve success with proof upload",
	requestID: 1,
	ctx: func() stdcontext.Context {
		return testutil.CtxWithTenant(1, reviewerID, true)
	},
	req:          &ReviewRequest{Status: constants.SubReqStatusApproved, Notes: "paid confirmed"},
	proofPayment: &multipart.FileHeader{Filename: "proof.jpg", Size: 5000},
	setupMocks: func(repo *mockRepo, companyRepo *mockCompanyRepo, role *mockRole, user *mockUser, cache *mockCache, storage *mockStorage) {
		repo.On("FindRequestByID", mock.Anything, uint(1)).Return(&SubscriptionRequest{
			ID: 1, CompanyID: 1, CurrentPlanID: 1, RequestedPlanID: proPlanID, Status: constants.SubReqStatusPending,
		}, nil)
		storage.On("UploadFileMultipart", mock.Anything, mock.Anything, mock.Anything).Return("https://storage.example.com/subscriptions/requests/1/proof-payment-123.jpg", nil)
		repo.On("UpdateRequest", mock.Anything, mock.MatchedBy(func(rq *SubscriptionRequest) bool {
			return rq.ProofPaymentURL != nil && *rq.ProofPaymentURL == "https://storage.example.com/subscriptions/requests/1/proof-payment-123.jpg"
		})).Return(nil)
		companyRepo.On("FindByID", mock.Anything, uint(1)).Return(&company.Company{ID: 1, SubscriptionPlanID: nil}, nil)
		companyRepo.On("Update", mock.Anything, mock.AnythingOfType("*company.Company")).Return(nil)
		repo.On("FindPlanByID", mock.Anything, proPlanID).Return(&SubscriptionPlan{ID: proPlanID, Slug: "pro"}, nil)
		role.On("FindPermissionIDsByGroupNames", mock.Anything, mock.Anything).Return([]uint{1, 2}, nil)
		role.On("FindRoleIDsByCompanyID", mock.Anything, uint(1)).Return([]uint{1}, nil)
		role.On("AssignPermissions", mock.Anything, uint(1), []uint{1, 2}, uint(1)).Return(nil)
		cache.On("Del", mock.Anything, "subscription:features:1").Return(nil)
		cache.On("Del", mock.Anything, "company:profile:1").Return(nil)
		user.On("ForceResetPasswordByCompanyID", mock.Anything, uint(1)).Return(nil)
	},
	wantErr: false,
},
{
	name:         "approve error proof payment required",
	requestID:    1,
	ctx:          func() stdcontext.Context { return testutil.CtxWithTenant(1, reviewerID, true) },
	req:          &ReviewRequest{Status: constants.SubReqStatusApproved},
	proofPayment: nil,
	setupMocks: func(repo *mockRepo, companyRepo *mockCompanyRepo, role *mockRole, user *mockUser, cache *mockCache, storage *mockStorage) {
		repo.On("FindRequestByID", mock.Anything, uint(1)).Return(&SubscriptionRequest{
			ID: 1, CompanyID: 1, Status: constants.SubReqStatusPending,
		}, nil)
	},
	wantErr: true,
	errMsg:  "proof of payment is required",
},
{
	name:         "approve error upload fails",
	requestID:    1,
	ctx:          func() stdcontext.Context { return testutil.CtxWithTenant(1, reviewerID, true) },
	req:          &ReviewRequest{Status: constants.SubReqStatusApproved},
	proofPayment: &multipart.FileHeader{Filename: "proof.jpg", Size: 5000},
	setupMocks: func(repo *mockRepo, companyRepo *mockCompanyRepo, role *mockRole, user *mockUser, cache *mockCache, storage *mockStorage) {
		repo.On("FindRequestByID", mock.Anything, uint(1)).Return(&SubscriptionRequest{
			ID: 1, CompanyID: 1, CurrentPlanID: 1, RequestedPlanID: proPlanID, Status: constants.SubReqStatusPending,
		}, nil)
		storage.On("UploadFileMultipart", mock.Anything, mock.Anything, mock.Anything).Return("", errors.New("upload failed"))
	},
	wantErr: true,
	errMsg:  "failed to upload proof of payment: upload failed",
},
```

- [ ] **Step 4: Update all existing test cases in TestService_ReviewRequest**

Update every existing test case's:
- `setupMocks` function signature: add `*mockStorage` parameter (unused for reject/not-found cases)
- Call to `newTestSubscriptionService`: capture 7 return values instead of 6
- Call to `svc.ReviewRequest`: pass `tt.proofPayment` as the file argument (nil for existing cases)

For example, the "reject success" test case becomes:

```go
{
	name:      "reject success",
	requestID: 1,
	ctx: func() stdcontext.Context {
		return testutil.CtxWithTenant(1, reviewerID, true)
	},
	req:          &ReviewRequest{Status: constants.SubReqStatusRejected, Notes: "not eligible"},
	proofPayment: nil,
	setupMocks: func(repo *mockRepo, companyRepo *mockCompanyRepo, role *mockRole, user *mockUser, cache *mockCache, storage *mockStorage) {
		repo.On("FindRequestByID", mock.Anything, uint(1)).Return(&SubscriptionRequest{
			ID: 1, CompanyID: 1, Status: constants.SubReqStatusPending,
		}, nil)
		repo.On("UpdateRequest", mock.Anything, mock.AnythingOfType("*subscription.SubscriptionRequest")).Return(nil)
	},
	wantErr: false,
},
```

And the test loop becomes:

```go
for _, tt := range tests {
	t.Run(tt.name, func(t *testing.T) {
		svc, repo, companyRepo, role, user, cache, storage := newTestSubscriptionService()
		tt.setupMocks(repo, companyRepo, role, user, cache, storage)

		err := svc.ReviewRequest(tt.ctx(), tt.requestID, tt.req, tt.proofPayment)

		if tt.wantErr {
			require.Error(t, err)
			assert.Equal(t, tt.errMsg, err.Error())
		} else {
			require.NoError(t, err)
		}
	})
}
```

Also update the "approve success" (original, without proof) test case: change its setupMocks signature to include `*mockStorage`, add `storage` as the 7th capture, add `proofPayment: nil` (since it tests the case with no proof for APPROVED — but wait, that would now fail with "proof of payment is required"). Either remove this case or add a file to it.

Actually: the original "approve success" test case without a proof file should now expect an error (`proof of payment is required`), because proof is required for APPROVED. Convert it to the error case above or add proof payment to make it succeed. The new test case "approve success with proof upload" replaces the old one. You can remove the old "approve success" case since it's now covered by the new "approve success with proof upload" case.

- [ ] **Step 5: Run tests to verify**

```bash
cd backend && go test ./internal/modules/subscription/ -v -run TestService_ReviewRequest
```

Expected: all ReviewRequest tests pass including new proof payment cases.

- [ ] **Step 6: Run full test suite**

```bash
cd backend && go test ./internal/modules/subscription/ -v
```

Expected: all tests pass.

- [ ] **Step 7: Commit**

```bash
git add backend/internal/modules/subscription/service_test.go
git commit -m "test(subscription): add proof payment upload test cases to ReviewRequest"
```

---

### Task 9: Frontend Types — Add proof_payment_url field

**Files:**
- Modify: `frontend/src/features/subscription/types.ts`

- [ ] **Step 1: Add field to SubscriptionRequestItem**

```typescript
export interface SubscriptionRequestItem {
  id: number;
  company_name: string;
  current_plan_name: string;
  requested_plan_name: string;
  current_plan_price: number;
  requested_plan_price: number;
  price_difference: number;
  status: string;
  requested_by_name: string;
  requested_by_email: string;
  notes: string;
  proof_payment_url?: string;
  created_at: string;
}
```

- [ ] **Step 2: Commit**

```bash
git add frontend/src/features/subscription/types.ts
git commit -m "feat(subscription): add proof_payment_url to frontend types"
```

---

### Task 10: Frontend Hook — Switch to FormData

**Files:**
- Modify: `frontend/src/features/subscription/hooks/useSubscription.ts`

- [ ] **Step 1: Change useReviewRequest to send FormData**

Replace the `useReviewRequest` hook:

```typescript
export const useReviewRequest = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({
      id,
      payload,
      file,
    }: {
      id: number;
      payload: ReviewPayload;
      file?: File;
    }) => {
      const formData = new FormData();
      formData.append("status", payload.status);
      formData.append("notes", payload.notes || "");
      if (file) {
        formData.append("proof_payment", file);
      }
      const { data } = await api.put(
        `/admin/subscriptions/${id}/review`,
        formData
      );
      return data;
    },
    onSuccess: () => {
      toast.success("Berhasil", {
        description: "Permintaan berhasil diproses.",
      });
      queryClient.invalidateQueries({
        queryKey: ["pending-subscription-requests"],
      });
      queryClient.invalidateQueries({
        queryKey: ["all-subscription-requests"],
      });
      queryClient.invalidateQueries({ queryKey: ["admin-companies"] });
      queryClient.invalidateQueries({ queryKey: ["admin-dashboard-stats"] });
    },
    onError: (error: any) => {
      const msg =
        error.response?.data?.message || "Gagal memproses permintaan";
      toast.error("Gagal", { description: msg });
    },
  });
};
```

- [ ] **Step 2: Commit**

```bash
git add frontend/src/features/subscription/hooks/useSubscription.ts
git commit -m "feat(subscription): switch useReviewRequest to FormData with file upload"
```

---

### Task 11: Frontend Page — Add file input and preview to dialog

**Files:**
- Modify: `frontend/src/pages/admin/SubscriptionAdminPage.tsx`

- [ ] **Step 1: Add file ref, file state, and preview state**

Add inside the component, after existing state declarations:

```typescript
const fileInputRef = useRef<HTMLInputElement>(null);
const [selectedFile, setSelectedFile] = useState<File | null>(null);
const [previewUrl, setPreviewUrl] = useState<string | null>(null);
```

Add `useRef` to the React import at top: `import { useState, useRef } from "react";`

- [ ] **Step 2: Add file selection handler**

```typescript
function handleFileSelect(e: React.ChangeEvent<HTMLInputElement>) {
  const file = e.target.files?.[0];
  if (!file) return;

  if (file.size > 5 * 1024 * 1024) {
    toast.error("File terlalu besar", {
      description: "Ukuran maksimal 5MB.",
    });
    return;
  }

  setSelectedFile(file);
  setPreviewUrl(URL.createObjectURL(file));
}
```

- [ ] **Step 3: Add cleanup on dialog close**

In `handleReview`, reset file state:

```typescript
function handleReview(
  request: SubscriptionRequestItem,
  action: "APPROVED" | "REJECTED"
) {
  setSelectedRequest(request);
  setReviewAction(action);
  setNotes("");
  setSelectedFile(null);
  setPreviewUrl(null);
  setIsDialogOpen(true);
}
```

- [ ] **Step 4: Update confirmReview to validate and send file**

```typescript
function confirmReview() {
  if (!selectedRequest) return;

  if (reviewAction === "APPROVED" && !selectedFile) {
    toast.error("Bukti pembayaran wajib diunggah", {
      description: "Silakan unggah bukti pembayaran sebelum menyetujui.",
    });
    return;
  }

  reviewRequest(
    {
      id: selectedRequest.id,
      payload: { status: reviewAction, notes },
      file: selectedFile || undefined,
    },
    {
      onSettled: () => {
        setIsDialogOpen(false);
        setSelectedRequest(null);
        setSelectedFile(null);
        setPreviewUrl(null);
      },
    }
  );
}
```

- [ ] **Step 5: Add file input and preview in the dialog JSX**

After the existing `<Textarea>` block (the notes field), add:

```tsx
{reviewAction === "APPROVED" && (
  <div className="space-y-2">
    <label className="text-sm font-medium text-slate-700">
      Bukti Pembayaran <span className="text-red-500">*</span>
    </label>
    <input
      ref={fileInputRef}
      type="file"
      accept="image/*"
      onChange={handleFileSelect}
      className="hidden"
    />
    {!selectedFile ? (
      <button
        type="button"
        onClick={() => fileInputRef.current?.click()}
        className="w-full rounded-lg border-2 border-dashed border-slate-300 p-6 text-center text-sm text-slate-500 hover:border-blue-400 hover:text-blue-600 transition-colors"
      >
        Klik untuk unggah bukti pembayaran
        <br />
        <span className="text-xs text-slate-400">Maks. 5MB (JPG, PNG)</span>
      </button>
    ) : (
      <div className="relative rounded-lg border border-slate-200 overflow-hidden">
        <img
          src={previewUrl!}
          alt="Preview bukti pembayaran"
          className="w-full h-40 object-cover"
        />
        <button
          type="button"
          onClick={() => {
            setSelectedFile(null);
            setPreviewUrl(null);
            if (fileInputRef.current) {
              fileInputRef.current.value = "";
            }
          }}
          className="absolute top-2 right-2 rounded-full bg-white/80 p-1 text-slate-600 hover:bg-white hover:text-red-600 shadow-sm"
        >
          <XCircle className="h-4 w-4" />
        </button>
      </div>
    )}
  </div>
)}
```

Place this right after the `</Textarea>` closing tag and before the closing `</div>` of the `space-y-4 py-2` container.

- [ ] **Step 6: Verify TypeScript compiles**

```bash
cd frontend && npx tsc --noEmit
```

- [ ] **Step 7: Commit**

```bash
git add frontend/src/pages/admin/SubscriptionAdminPage.tsx
git commit -m "feat(subscription): add proof payment file input and preview to approval dialog"
```

---

### Task 12: Final verification and build

**Files:**
- None created/modified — this is a verification task

- [ ] **Step 1: Build backend**

```bash
cd backend && go build ./...
```

Expected: no errors.

- [ ] **Step 2: Run all backend tests**

```bash
cd backend && go test ./... -count=1
```

Expected: all tests pass.

- [ ] **Step 3: Build frontend**

```bash
cd frontend && npm run build
```

Expected: build succeeds.

- [ ] **Step 4: Confirm all changes are committed**

```bash
git status
```

Expected: clean working tree.
