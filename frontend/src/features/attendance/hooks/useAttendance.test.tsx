import { describe, it, expect, beforeEach, vi } from "vitest";
import { renderHook, waitFor } from "@testing-library/react";
import { createQueryWrapper } from "@/test/utils";
import { useClock, useTodayAttendance } from "@/features/attendance/hooks/useAttendance";

vi.mock("@/lib/axios", () => ({
  api: { post: vi.fn(), get: vi.fn() },
}));

vi.mock("sonner", () => ({
  toast: { success: vi.fn(), error: vi.fn() },
}));

import { api } from "@/lib/axios";
import { toast } from "sonner";

describe("useClock", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should submit clock and show success toast", async () => {
    (api.post as ReturnType<typeof vi.fn>).mockResolvedValue({
      data: { message: "Clock in successful", data: {} },
    });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useClock(), { wrapper });

    result.current.mutate({ type: "clock_in", latitude: -6.2, longitude: 106.8 });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(toast.success).toHaveBeenCalledWith("Clock in successful");
  });

  it("should show error toast on failure", async () => {
    (api.post as ReturnType<typeof vi.fn>).mockRejectedValue({
      response: { data: { message: "Already clocked in" } },
    });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useClock(), { wrapper });

    result.current.mutate({ type: "clock_in", latitude: 0, longitude: 0 });

    await waitFor(() => expect(result.current.isError).toBe(true));
    expect(toast.error).toHaveBeenCalled();
  });
});

describe("useTodayAttendance", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should fetch today attendance", async () => {
    const mockData = { clock_in: "08:00", clock_out: null, status: "present" };
    (api.get as ReturnType<typeof vi.fn>).mockResolvedValue({
      data: { data: mockData },
    });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useTodayAttendance(), { wrapper });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data).toEqual(mockData);
  });
});
