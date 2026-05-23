import { describe, it, expect, beforeEach, vi } from "vitest";
import { renderHook, waitFor } from "@testing-library/react";
import { createQueryWrapper } from "@/test/utils";
import {
  useRequisitions,
  useRequisitionDetail,
  useCreateRequisition,
  useSubmitRequisition,
  useRequisitionAction,
  useCloseRequisition,
  useDeleteRequisition,
} from "@/features/recruitment/hooks/useRequisition";
import {
  useApplicants,
  useApplicantDetail,
  useAddApplicant,
  useUpdateApplicantStage,
} from "@/features/recruitment/hooks/useApplicant";

vi.mock("@/lib/axios", () => ({
  api: { get: vi.fn(), post: vi.fn(), put: vi.fn(), delete: vi.fn() },
}));

vi.mock("sonner", () => ({
  toast: { success: vi.fn(), error: vi.fn() },
}));

import { api } from "@/lib/axios";
import { toast } from "sonner";

describe("useRequisitions", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should fetch requisitions", async () => {
    const mockData = { data: [{ id: 1, title: "Dev" }] };
    (api.get as ReturnType<typeof vi.fn>).mockResolvedValue({ data: mockData });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useRequisitions({ page: 1 }), { wrapper });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data).toEqual(mockData);
  });
});

describe("useRequisitionDetail", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should fetch requisition detail", async () => {
    const mockDetail = { id: 1, title: "Dev", status: "open" };
    (api.get as ReturnType<typeof vi.fn>).mockResolvedValue({ data: { data: mockDetail } });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useRequisitionDetail(1), { wrapper });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data).toEqual(mockDetail);
  });

  it("should not fetch when id is null", () => {
    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useRequisitionDetail(null), { wrapper });
    expect(result.current.fetchStatus).toBe("idle");
  });
});

describe("useCreateRequisition", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should create requisition", async () => {
    (api.post as ReturnType<typeof vi.fn>).mockResolvedValue({ data: { message: "ok" } });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useCreateRequisition(), { wrapper });

    result.current.mutate({ title: "Dev", department_id: 1, quantity: 1, description: "Test" });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(toast.success).toHaveBeenCalledWith("Requisition created successfully");
  });
});

describe("useSubmitRequisition", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should submit requisition", async () => {
    (api.put as ReturnType<typeof vi.fn>).mockResolvedValue({ data: { message: "ok" } });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useSubmitRequisition(), { wrapper });

    result.current.mutate(1);

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(toast.success).toHaveBeenCalledWith("Requisition submitted for approval");
  });
});

describe("useRequisitionAction", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should process action on requisition", async () => {
    (api.put as ReturnType<typeof vi.fn>).mockResolvedValue({ data: { message: "ok" } });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useRequisitionAction(), { wrapper });

    result.current.mutate({ id: 1, payload: { action: "approve", rejection_reason: "" } });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(toast.success).toHaveBeenCalledWith("Action processed successfully");
  });
});

describe("useCloseRequisition", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should close requisition", async () => {
    (api.put as ReturnType<typeof vi.fn>).mockResolvedValue({ data: { message: "ok" } });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useCloseRequisition(), { wrapper });

    result.current.mutate(1);

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(toast.success).toHaveBeenCalledWith("Requisition closed");
  });
});

describe("useDeleteRequisition", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should delete requisition", async () => {
    (api.delete as ReturnType<typeof vi.fn>).mockResolvedValue({ data: { message: "ok" } });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useDeleteRequisition(), { wrapper });

    result.current.mutate(1);

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(toast.success).toHaveBeenCalledWith("Requisition deleted");
  });
});

describe("useApplicants", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should fetch applicants for requisition", async () => {
    const mockData = { data: [{ id: 1, name: "Alice" }] };
    (api.get as ReturnType<typeof vi.fn>).mockResolvedValue({ data: mockData });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useApplicants(1), { wrapper });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(api.get).toHaveBeenCalledWith("/recruitments/requisitions/1/applicants");
  });

  it("should not fetch when requisitionId is null", () => {
    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useApplicants(null), { wrapper });
    expect(result.current.fetchStatus).toBe("idle");
  });
});

describe("useApplicantDetail", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should fetch applicant detail", async () => {
    const mockDetail = { id: 1, name: "Alice", stage: "interview" };
    (api.get as ReturnType<typeof vi.fn>).mockResolvedValue({ data: { data: mockDetail } });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useApplicantDetail(1), { wrapper });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data).toEqual(mockDetail);
  });

  it("should not fetch when id is null", () => {
    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useApplicantDetail(null), { wrapper });
    expect(result.current.fetchStatus).toBe("idle");
  });
});

describe("useAddApplicant", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should add applicant", async () => {
    (api.post as ReturnType<typeof vi.fn>).mockResolvedValue({ data: { message: "ok" } });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useAddApplicant(1), { wrapper });

    result.current.mutate({ name: "Alice", email: "alice@test.com" });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(toast.success).toHaveBeenCalledWith("Applicant added successfully");
  });
});

describe("useUpdateApplicantStage", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should update applicant stage", async () => {
    (api.put as ReturnType<typeof vi.fn>).mockResolvedValue({ data: { message: "ok" } });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useUpdateApplicantStage(), { wrapper });

    result.current.mutate({ id: 1, payload: { stage: "interview" } });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(toast.success).toHaveBeenCalledWith("Stage updated");
  });
});
