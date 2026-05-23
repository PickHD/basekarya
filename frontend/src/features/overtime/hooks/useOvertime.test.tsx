import { describe, it, expect, beforeEach, vi } from "vitest";
import { renderHook, waitFor } from "@testing-library/react";
import { createQueryWrapper } from "@/test/utils";
import { useOvertimes, useCreateOvertime, useOvertimeAction, useOvertime } from "@/features/overtime/hooks/useOvertime";

vi.mock("@/lib/axios", () => ({
  api: { get: vi.fn(), post: vi.fn(), put: vi.fn() },
}));

vi.mock("sonner", () => ({
  toast: { success: vi.fn(), error: vi.fn() },
}));

import { api } from "@/lib/axios";
import { toast } from "sonner";

describe("useOvertimes", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should fetch overtimes with filter", async () => {
    const mockData = { data: [{ id: 1, status: "pending" }], meta: { limit: 10 } };
    (api.get as ReturnType<typeof vi.fn>).mockResolvedValue({ data: mockData });

    const filter = { page: 1, limit: 10 };
    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useOvertimes(filter), { wrapper });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data).toEqual(mockData);
  });
});

describe("useOvertime", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should fetch overtime detail", async () => {
    const mockDetail = { id: 1, status: "approved" };
    (api.get as ReturnType<typeof vi.fn>).mockResolvedValue({
      data: { data: mockDetail },
    });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useOvertime("1"), { wrapper });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data).toEqual(mockDetail);
  });

  it("should not fetch when id is empty", () => {
    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useOvertime(""), { wrapper });
    expect(result.current.fetchStatus).toBe("idle");
  });
});

describe("useCreateOvertime", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should create overtime and show toast", async () => {
    (api.post as ReturnType<typeof vi.fn>).mockResolvedValue({ data: { message: "ok" } });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useCreateOvertime(), { wrapper });

    result.current.mutate({ date: "2024-01-01", start_time: "18:00", end_time: "20:00", reason: "Project" });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(toast.success).toHaveBeenCalled();
  });

  it("should show error toast on failure", async () => {
    (api.post as ReturnType<typeof vi.fn>).mockRejectedValue({
      response: { data: { message: "Overlap" } },
    });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useCreateOvertime(), { wrapper });

    result.current.mutate({ date: "2024-01-01", start_time: "18:00", end_time: "20:00", reason: "Project" });

    await waitFor(() => expect(result.current.isError).toBe(true));
    expect(toast.error).toHaveBeenCalled();
  });
});

describe("useOvertimeAction", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should approve overtime", async () => {
    (api.put as ReturnType<typeof vi.fn>).mockResolvedValue({ data: { message: "ok" } });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useOvertimeAction(), { wrapper });

    result.current.mutate({ id: 1, action: "APPROVE", rejection_reason: "" });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(toast.success).toHaveBeenCalledWith("Lembur berhasil disetujui");
  });

  it("should reject overtime", async () => {
    (api.put as ReturnType<typeof vi.fn>).mockResolvedValue({ data: { message: "ok" } });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useOvertimeAction(), { wrapper });

    result.current.mutate({ id: 1, action: "REJECT", rejection_reason: "Not needed" });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(toast.success).toHaveBeenCalledWith("Lembur berhasil ditolak");
  });
});
