import { describe, it, expect, beforeEach, vi } from "vitest";
import { renderHook, waitFor } from "@testing-library/react";
import { createQueryWrapper } from "@/test/utils";
import {
  useAssets,
  useAsset,
  useCreateAsset,
  useAssetAssignments,
  useAssetAssignment,
  useCreateAssetAssignment,
  useAssetAssignmentAction,
  useReturnAssetAssignment,
} from "@/features/asset/hooks/useAsset";

vi.mock("@/lib/axios", () => ({
  api: { get: vi.fn(), post: vi.fn(), put: vi.fn(), delete: vi.fn() },
}));

vi.mock("sonner", () => ({
  toast: { success: vi.fn(), error: vi.fn() },
}));

import { api } from "@/lib/axios";
import { toast } from "sonner";

describe("useAssets", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should fetch assets with filter", async () => {
    const mockData = {
      data: [{ id: 1, name: "Laptop", status: "AVAILABLE", condition: "GOOD" }],
      meta: { limit: 10 },
    };
    (api.get as ReturnType<typeof vi.fn>).mockResolvedValue({ data: mockData });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useAssets({ page: 1, limit: 10 }), { wrapper });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data).toEqual(mockData);
  });
});

describe("useAsset", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should fetch asset detail", async () => {
    const mockDetail = { id: 1, name: "Laptop", status: "AVAILABLE", condition: "GOOD" };
    (api.get as ReturnType<typeof vi.fn>).mockResolvedValue({ data: { data: mockDetail } });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useAsset("1"), { wrapper });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data).toEqual(mockDetail);
  });

  it("should not fetch when id is empty", () => {
    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useAsset(""), { wrapper });
    expect(result.current.fetchStatus).toBe("idle");
  });
});

describe("useCreateAsset", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should create asset and show toast", async () => {
    (api.post as ReturnType<typeof vi.fn>).mockResolvedValue({ data: { message: "ok" } });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useCreateAsset(), { wrapper });

    result.current.mutate({
      asset_category_id: 1,
      name: "MacBook Pro",
      description: "",
      serial_number: "SN001",
      condition: "GOOD",
    });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(toast.success).toHaveBeenCalledWith("Aset berhasil dibuat!");
  });

  it("should show error toast on failure", async () => {
    (api.post as ReturnType<typeof vi.fn>).mockRejectedValue({
      response: { data: { message: "Failed" } },
    });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useCreateAsset(), { wrapper });

    result.current.mutate({
      asset_category_id: 1,
      name: "",
      description: "",
      serial_number: "",
      condition: "GOOD",
    });

    await waitFor(() => expect(result.current.isError).toBe(true));
    expect(toast.error).toHaveBeenCalled();
  });
});

describe("useAssetAssignments", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should fetch assignments with filter", async () => {
    const mockData = {
      data: [{ id: 1, asset_id: 1, employee_name: "John", status: "PENDING" }],
      meta: { limit: 10 },
    };
    (api.get as ReturnType<typeof vi.fn>).mockResolvedValue({ data: mockData });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useAssetAssignments({ page: 1, limit: 10 }), { wrapper });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data).toEqual(mockData);
  });
});

describe("useAssetAssignment", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should fetch assignment detail", async () => {
    const mockDetail = { id: 1, asset_id: 1, employee_name: "John", status: "ACTIVE" };
    (api.get as ReturnType<typeof vi.fn>).mockResolvedValue({ data: { data: mockDetail } });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useAssetAssignment("1"), { wrapper });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data).toEqual(mockDetail);
  });
});

describe("useCreateAssetAssignment", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should create assignment and show toast", async () => {
    (api.post as ReturnType<typeof vi.fn>).mockResolvedValue({ data: { message: "ok" } });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useCreateAssetAssignment(), { wrapper });

    result.current.mutate({
      asset_id: 1,
      purpose: "Need laptop for presentation",
    });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(toast.success).toHaveBeenCalledWith("Permintaan aset berhasil dikirim!");
  });
});

describe("useAssetAssignmentAction", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should approve assignment", async () => {
    (api.put as ReturnType<typeof vi.fn>).mockResolvedValue({ data: { message: "ok" } });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useAssetAssignmentAction(), { wrapper });

    result.current.mutate({ id: 1, action: "APPROVE" });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(toast.success).toHaveBeenCalled();
  });

  it("should reject assignment with reason", async () => {
    (api.put as ReturnType<typeof vi.fn>).mockResolvedValue({ data: { message: "ok" } });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useAssetAssignmentAction(), { wrapper });

    result.current.mutate({ id: 1, action: "REJECT", rejection_reason: "Not available" });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(toast.success).toHaveBeenCalled();
  });
});

describe("useReturnAssetAssignment", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should return assignment", async () => {
    (api.put as ReturnType<typeof vi.fn>).mockResolvedValue({ data: { message: "ok" } });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useReturnAssetAssignment(), { wrapper });

    result.current.mutate(1);

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(toast.success).toHaveBeenCalledWith("Aset berhasil dikembalikan!");
  });
});
