import { describe, it, expect, beforeEach, vi } from "vitest";
import { renderHook, waitFor } from "@testing-library/react";
import { createQueryWrapper } from "@/test/utils";
import { useLeaves, useApplyLeave, useLeaveAction, useLeaveDetail, useLeaveTypes } from "@/features/leave/hooks/useLeave";

vi.mock("@/lib/axios", () => ({
  api: {
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
  },
}));

vi.mock("sonner", () => ({
  toast: { success: vi.fn(), error: vi.fn() },
}));

import { api } from "@/lib/axios";
import { toast } from "sonner";

const mockLeaves = {
  data: [
    {
      id: 1,
      employee_id: 1,
      employee_name: "John",
      leave_type: "annual",
      start_date: "2024-01-01",
      end_date: "2024-01-03",
      status: "pending",
      reason: "Vacation",
    },
  ],
  meta: { limit: 10, page: 1, total_page: 1, total_data: 1 },
};

describe("useLeaves", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("should fetch leaves with params", async () => {
    (api.get as ReturnType<typeof vi.fn>).mockResolvedValue({ data: mockLeaves });

    const params = { page: 1, limit: 10 };
    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useLeaves(params), { wrapper });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data).toEqual(mockLeaves);
    expect(api.get).toHaveBeenCalledWith("/leaves", { params });
  });
});

describe("useApplyLeave", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("should apply leave and invalidate queries", async () => {
    (api.post as ReturnType<typeof vi.fn>).mockResolvedValue({ data: { message: "ok" } });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useApplyLeave(), { wrapper });

    result.current.mutate({
      leave_type_id: 1,
      start_date: "2024-01-01",
      end_date: "2024-01-03",
      reason: "Vacation",
    });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(toast.success).toHaveBeenCalled();
  });

  it("should show error toast on failure", async () => {
    (api.post as ReturnType<typeof vi.fn>).mockRejectedValue({
      response: { data: { message: "Insufficient balance" } },
    });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useApplyLeave(), { wrapper });

    result.current.mutate({
      leave_type_id: 1,
      start_date: "2024-01-01",
      end_date: "2024-01-03",
      reason: "Vacation",
    });

    await waitFor(() => expect(result.current.isError).toBe(true));
    expect(toast.error).toHaveBeenCalled();
  });
});

describe("useLeaveAction", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("should approve leave and show success toast", async () => {
    (api.put as ReturnType<typeof vi.fn>).mockResolvedValue({ data: { message: "ok" } });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useLeaveAction(), { wrapper });

    result.current.mutate({ id: 1, action: "approve", rejection_reason: "" });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(toast.success).toHaveBeenCalled();
    expect(api.put).toHaveBeenCalledWith("/leaves/1/action", {
      action: "approve",
      rejection_reason: "",
    });
  });

  it("should reject leave with reason", async () => {
    (api.put as ReturnType<typeof vi.fn>).mockResolvedValue({ data: { message: "ok" } });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useLeaveAction(), { wrapper });

    result.current.mutate({ id: 2, action: "reject", rejection_reason: "Not approved" });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(api.put).toHaveBeenCalledWith("/leaves/2/action", {
      action: "reject",
      rejection_reason: "Not approved",
    });
  });
});

describe("useLeaveDetail", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("should fetch leave detail by id", async () => {
    const mockDetail = { data: mockLeaves.data[0] };
    (api.get as ReturnType<typeof vi.fn>).mockResolvedValue({ data: mockDetail });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useLeaveDetail("1"), { wrapper });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data).toEqual(mockLeaves.data[0]);
  });

  it("should not fetch when id is empty", () => {
    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useLeaveDetail(""), { wrapper });

    expect(result.current.fetchStatus).toBe("idle");
  });
});

describe("useLeaveTypes", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("should fetch leave types", async () => {
    const mockTypes = [{ id: 1, name: "Annual", quota: 12 }];
    (api.get as ReturnType<typeof vi.fn>).mockResolvedValue({
      data: { data: mockTypes },
    });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useLeaveTypes(), { wrapper });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data).toEqual(mockTypes);
  });
});
