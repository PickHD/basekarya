import { describe, it, expect, beforeEach, vi } from "vitest";
import { renderHook, waitFor } from "@testing-library/react";
import { createQueryWrapper } from "@/test/utils";
import { useLoans, useCreateLoan, useLoanAction, useLoan } from "@/features/loan/hooks/useLoan";

vi.mock("@/lib/axios", () => ({
  api: { get: vi.fn(), post: vi.fn(), put: vi.fn() },
}));

vi.mock("sonner", () => ({
  toast: { success: vi.fn(), error: vi.fn() },
}));

import { api } from "@/lib/axios";
import { toast } from "sonner";

describe("useLoans", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should fetch loans with filter", async () => {
    const mockData = { data: [{ id: 1, amount: 500000, status: "pending" }], meta: { limit: 10 } };
    (api.get as ReturnType<typeof vi.fn>).mockResolvedValue({ data: mockData });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useLoans({ page: 1, limit: 10 }), { wrapper });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data).toEqual(mockData);
  });
});

describe("useLoan", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should fetch loan detail", async () => {
    const mockDetail = { id: 1, amount: 500000, status: "approved" };
    (api.get as ReturnType<typeof vi.fn>).mockResolvedValue({ data: { data: mockDetail } });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useLoan("1"), { wrapper });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data).toEqual(mockDetail);
  });

  it("should not fetch when id is empty", () => {
    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useLoan(""), { wrapper });
    expect(result.current.fetchStatus).toBe("idle");
  });
});

describe("useCreateLoan", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should create loan and show toast", async () => {
    (api.post as ReturnType<typeof vi.fn>).mockResolvedValue({ data: { message: "ok" } });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useCreateLoan(), { wrapper });

    result.current.mutate({ amount: 500000, reason: "Emergency", installments: 6 });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(toast.success).toHaveBeenCalledWith("Pengajuan kasbon berhasil dikirim!");
  });

  it("should show error toast on failure", async () => {
    (api.post as ReturnType<typeof vi.fn>).mockRejectedValue({
      response: { data: { message: "Insufficient tenure" } },
    });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useCreateLoan(), { wrapper });

    result.current.mutate({ amount: 500000, reason: "Emergency", installments: 6 });

    await waitFor(() => expect(result.current.isError).toBe(true));
    expect(toast.error).toHaveBeenCalled();
  });
});

describe("useLoanAction", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should approve loan", async () => {
    (api.put as ReturnType<typeof vi.fn>).mockResolvedValue({ data: { message: "ok" } });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useLoanAction(), { wrapper });

    result.current.mutate({ id: 1, action: "approve", rejection_reason: "" });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(toast.success).toHaveBeenCalled();
  });
});
