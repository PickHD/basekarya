import { describe, it, expect, beforeEach, vi } from "vitest";
import { renderHook, waitFor } from "@testing-library/react";
import { createQueryWrapper } from "@/test/utils";
import {
  useReimbursements,
  useReimbursement,
  useCreateReimbursement,
  useReimbursementAction,
} from "@/features/reimbursement/hooks/useReimbursement";

vi.mock("@/lib/axios", () => ({
  api: { get: vi.fn(), post: vi.fn(), put: vi.fn() },
}));

vi.mock("sonner", () => ({
  toast: { success: vi.fn(), error: vi.fn() },
}));

vi.mock("react-router-dom", () => ({
  useNavigate: () => vi.fn(),
}));

import { api } from "@/lib/axios";
import { toast } from "sonner";

describe("useReimbursements", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should fetch reimbursements with filter", async () => {
    const mockData = { data: [{ id: 1, amount: 500000 }], meta: { limit: 10 } };
    (api.get as ReturnType<typeof vi.fn>).mockResolvedValue({ data: mockData });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(
      () => useReimbursements({ page: 1, limit: 10 }),
      { wrapper }
    );

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data).toEqual(mockData);
  });
});

describe("useReimbursement", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should fetch reimbursement detail", async () => {
    const mockDetail = { id: 1, amount: 500000, status: "pending" };
    (api.get as ReturnType<typeof vi.fn>).mockResolvedValue({ data: { data: mockDetail } });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useReimbursement("1"), { wrapper });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data).toEqual(mockDetail);
  });

  it("should not fetch when id is empty", () => {
    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useReimbursement(""), { wrapper });
    expect(result.current.fetchStatus).toBe("idle");
  });
});

describe("useCreateReimbursement", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should create reimbursement with FormData", async () => {
    (api.post as ReturnType<typeof vi.fn>).mockResolvedValue({ data: { message: "ok" } });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useCreateReimbursement(), { wrapper });

    result.current.mutate({
      title: "Transport",
      description: "Taxi fare",
      amount: 50000,
      date: "2024-01-15",
      proof_file: null,
    });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(api.post).toHaveBeenCalledWith("/reimbursements", expect.any(FormData));
  });

  it("should handle file upload", async () => {
    (api.post as ReturnType<typeof vi.fn>).mockResolvedValue({ data: { message: "ok" } });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useCreateReimbursement(), { wrapper });

    const mockFile = new File(["test"], "receipt.jpg", { type: "image/jpeg" });
    result.current.mutate({
      title: "Transport",
      description: "Taxi fare",
      amount: 50000,
      date: "2024-01-15",
      proof_file: [mockFile],
    });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    const calledFormData = (api.post as ReturnType<typeof vi.fn>).mock.calls[0][1] as FormData;
    expect(calledFormData.get("file")).toBeTruthy();
  });
});

describe("useReimbursementAction", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should approve reimbursement", async () => {
    (api.put as ReturnType<typeof vi.fn>).mockResolvedValue({ data: { message: "ok" } });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useReimbursementAction(), { wrapper });

    result.current.mutate({ id: 1, action: "approve", rejection_reason: "" });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(toast.success).toHaveBeenCalledWith("Reimbursement berhasil di-approve");
  });

  it("should reject reimbursement", async () => {
    (api.put as ReturnType<typeof vi.fn>).mockResolvedValue({ data: { message: "ok" } });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useReimbursementAction(), { wrapper });

    result.current.mutate({ id: 1, action: "reject", rejection_reason: "Invalid" });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(toast.success).toHaveBeenCalledWith("Reimbursement berhasil di-reject");
  });
});
