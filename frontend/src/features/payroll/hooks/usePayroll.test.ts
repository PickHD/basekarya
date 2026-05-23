import { describe, it, expect, beforeEach, vi } from "vitest";
import { renderHook, waitFor } from "@testing-library/react";
import { createQueryWrapper } from "@/test/utils";
import {
  usePayrolls,
  usePayroll,
  useGeneratePayroll,
  useMarkAsPaid,
  useSendPayslipEmail,
} from "@/features/payroll/hooks/usePayroll";

vi.mock("@/lib/axios", () => ({
  api: { get: vi.fn(), post: vi.fn(), put: vi.fn() },
}));

vi.mock("sonner", () => ({
  toast: { success: vi.fn(), error: vi.fn() },
}));

import { api } from "@/lib/axios";
import { toast } from "sonner";

describe("usePayrolls", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should fetch payrolls with filter", async () => {
    const mockData = { data: [{ id: 1, status: "pending" }], meta: { limit: 10 } };
    (api.get as ReturnType<typeof vi.fn>).mockResolvedValue({ data: mockData });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => usePayrolls({ page: 1, limit: 10 }), { wrapper });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data).toEqual(mockData);
  });
});

describe("usePayroll", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should fetch payroll detail", async () => {
    const mockDetail = { id: 1, status: "paid", net_salary: 5000000 };
    (api.get as ReturnType<typeof vi.fn>).mockResolvedValue({ data: { data: mockDetail } });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => usePayroll(1), { wrapper });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data).toEqual(mockDetail);
  });

  it("should not fetch when id is null", () => {
    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => usePayroll(null), { wrapper });
    expect(result.current.fetchStatus).toBe("idle");
  });
});

describe("useGeneratePayroll", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should generate payroll and show toast", async () => {
    (api.post as ReturnType<typeof vi.fn>).mockResolvedValue({
      data: { message: "ok", data: { success_count: 10 } },
    });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useGeneratePayroll(), { wrapper });

    result.current.mutate({ month: 1, year: 2024 });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(toast.success).toHaveBeenCalled();
  });

  it("should show error on failure", async () => {
    (api.post as ReturnType<typeof vi.fn>).mockRejectedValue({
      response: { data: { message: "Already generated" } },
    });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useGeneratePayroll(), { wrapper });

    result.current.mutate({ month: 1, year: 2024 });

    await waitFor(() => expect(result.current.isError).toBe(true));
    expect(toast.error).toHaveBeenCalled();
  });
});

describe("useMarkAsPaid", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should mark as paid and show toast", async () => {
    (api.put as ReturnType<typeof vi.fn>).mockResolvedValue({ data: { message: "ok" } });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useMarkAsPaid(), { wrapper });

    result.current.mutate(1);

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(toast.success).toHaveBeenCalledWith("Status berhasil diubah menjadi PAID");
  });
});

describe("useSendPayslipEmail", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should send payslip email", async () => {
    (api.post as ReturnType<typeof vi.fn>).mockResolvedValue({ data: { message: "ok" } });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useSendPayslipEmail(), { wrapper });

    result.current.mutate(1);

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(toast.success).toHaveBeenCalledWith("Email sedang dikirim...");
  });

  it("should show error on failure", async () => {
    (api.post as ReturnType<typeof vi.fn>).mockRejectedValue({
      response: { data: { message: "SMTP error" } },
    });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useSendPayslipEmail(), { wrapper });

    result.current.mutate(1);

    await waitFor(() => expect(result.current.isError).toBe(true));
    expect(toast.error).toHaveBeenCalled();
  });
});
