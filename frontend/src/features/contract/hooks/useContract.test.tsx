import { describe, it, expect, beforeEach, vi } from "vitest";
import { renderHook, waitFor } from "@testing-library/react";
import { createQueryWrapper } from "@/test/utils";
import {
  useContracts,
  useContractByEmployee,
  useContractDetail,
  useUpsertContract,
  useDeleteContract,
} from "@/features/contract/hooks/useContract";

vi.mock("@/lib/axios", () => ({
  api: { get: vi.fn(), put: vi.fn(), delete: vi.fn() },
}));

vi.mock("sonner", () => ({
  toast: { success: vi.fn(), error: vi.fn() },
}));

import { api } from "@/lib/axios";
import { toast } from "sonner";

describe("useContracts", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should fetch contracts", async () => {
    const mockData = { data: [{ id: 1, type: "permanent" }] };
    (api.get as ReturnType<typeof vi.fn>).mockResolvedValue({ data: mockData });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(
      () => useContracts({ page: 1, limit: 10 }),
      { wrapper }
    );

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data).toEqual(mockData);
  });
});

describe("useContractByEmployee", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should fetch contract by employee", async () => {
    const mockData = { id: 1, type: "permanent", start_date: "2024-01-01" };
    (api.get as ReturnType<typeof vi.fn>).mockResolvedValue({ data: { data: mockData } });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useContractByEmployee(1), { wrapper });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data).toEqual(mockData);
  });

  it("should not fetch when employeeId is 0", () => {
    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useContractByEmployee(0), { wrapper });
    expect(result.current.fetchStatus).toBe("idle");
  });
});

describe("useContractDetail", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should fetch contract detail", async () => {
    const mockDetail = { id: 1, type: "contract", end_date: "2024-12-31" };
    (api.get as ReturnType<typeof vi.fn>).mockResolvedValue({ data: { data: mockDetail } });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useContractDetail(1), { wrapper });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data).toEqual(mockDetail);
  });

  it("should not fetch when id is null", () => {
    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useContractDetail(null), { wrapper });
    expect(result.current.fetchStatus).toBe("idle");
  });
});

describe("useUpsertContract", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should upsert contract", async () => {
    (api.put as ReturnType<typeof vi.fn>).mockResolvedValue({ data: { message: "ok" } });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useUpsertContract(), { wrapper });

    result.current.mutate({
      employee_id: 1,
      type: "permanent",
      start_date: "2024-01-01",
      end_date: "2024-12-31",
    });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(toast.success).toHaveBeenCalledWith("Contract saved successfully");
  });

  it("should show error on failure", async () => {
    (api.put as ReturnType<typeof vi.fn>).mockRejectedValue({
      response: { data: { message: "Invalid dates" } },
    });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useUpsertContract(), { wrapper });

    result.current.mutate({
      employee_id: 1,
      type: "permanent",
      start_date: "2024-01-01",
    });

    await waitFor(() => expect(result.current.isError).toBe(true));
    expect(toast.error).toHaveBeenCalled();
  });
});

describe("useDeleteContract", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should delete contract", async () => {
    (api.delete as ReturnType<typeof vi.fn>).mockResolvedValue({ data: { message: "ok" } });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useDeleteContract(), { wrapper });

    result.current.mutate(1);

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(toast.success).toHaveBeenCalledWith("Contract deleted successfully");
  });
});
