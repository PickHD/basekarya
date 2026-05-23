import { describe, it, expect, beforeEach, vi } from "vitest";
import { renderHook, waitFor } from "@testing-library/react";
import { createQueryWrapper } from "@/test/utils";
import {
  useRequestUpgrade,
  usePendingRequests,
  useReviewRequest,
  useCompanies,
  useDashboardStats,
  useCompanyDetail,
  useUpdateCompanyStatus,
} from "@/features/subscription/hooks/useSubscription";

vi.mock("@/lib/axios", () => ({
  api: { get: vi.fn(), post: vi.fn(), put: vi.fn() },
}));

vi.mock("sonner", () => ({
  toast: { success: vi.fn(), error: vi.fn() },
}));

import { api } from "@/lib/axios";
import { toast } from "sonner";

describe("useRequestUpgrade", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should request upgrade and show success toast", async () => {
    (api.post as ReturnType<typeof vi.fn>).mockResolvedValue({
      data: { message: "ok" },
    });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useRequestUpgrade(), { wrapper });

    result.current.mutate({ plan_slug: "professional" });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(toast.success).toHaveBeenCalledWith("Permintaan upgrade terkirim", {
      description:
        "Tim kami akan menghubungi Anda untuk konfirmasi pembayaran.",
    });
  });

  it("should show error toast on failure", async () => {
    (api.post as ReturnType<typeof vi.fn>).mockRejectedValue({
      response: { data: { message: "Already pending" } },
    });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useRequestUpgrade(), { wrapper });

    result.current.mutate({ plan_slug: "professional" });

    await waitFor(() => expect(result.current.isError).toBe(true));
    expect(toast.error).toHaveBeenCalled();
  });
});

describe("usePendingRequests", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should fetch pending requests", async () => {
    const mockData = { data: [{ id: 1, status: "pending" }] };
    (api.get as ReturnType<typeof vi.fn>).mockResolvedValue({ data: mockData });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => usePendingRequests(), { wrapper });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data).toEqual(mockData);
  });
});

describe("useReviewRequest", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should review request and show success toast", async () => {
    (api.put as ReturnType<typeof vi.fn>).mockResolvedValue({
      data: { message: "ok" },
    });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useReviewRequest(), { wrapper });

    result.current.mutate({ id: 1, payload: { status: "APPROVED" } });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(toast.success).toHaveBeenCalledWith("Berhasil", {
      description: "Permintaan berhasil diproses.",
    });
  });
});

describe("useCompanies", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should fetch companies with search", async () => {
    const mockData = { data: [{ id: 1, name: "Test Co" }] };
    (api.get as ReturnType<typeof vi.fn>).mockResolvedValue({ data: mockData });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useCompanies("test"), { wrapper });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(api.get).toHaveBeenCalledWith("/admin/subscriptions/companies", {
      params: { search: "test" },
    });
  });

  it("should fetch companies without search", async () => {
    const mockData = { data: [] };
    (api.get as ReturnType<typeof vi.fn>).mockResolvedValue({ data: mockData });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useCompanies(), { wrapper });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(api.get).toHaveBeenCalledWith("/admin/subscriptions/companies", {
      params: undefined,
    });
  });
});

describe("useCompanyDetail", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should fetch company detail", async () => {
    const mockData = { data: { id: 1, name: "Test" } };
    (api.get as ReturnType<typeof vi.fn>).mockResolvedValue({ data: mockData });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useCompanyDetail(1), { wrapper });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data).toEqual(mockData);
  });

  it("should not fetch when id is 0", () => {
    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useCompanyDetail(0), { wrapper });
    expect(result.current.fetchStatus).toBe("idle");
  });
});

describe("useUpdateCompanyStatus", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should update company status", async () => {
    (api.put as ReturnType<typeof vi.fn>).mockResolvedValue({
      data: { message: "ok" },
    });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useUpdateCompanyStatus(), { wrapper });

    result.current.mutate({ id: 1, status: "active" });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(toast.success).toHaveBeenCalledWith("Berhasil", {
      description: "Status perusahaan diperbarui.",
    });
  });
});

describe("useDashboardStats", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should fetch dashboard stats", async () => {
    const mockData = { data: { total_companies: 10, active_subscriptions: 8 } };
    (api.get as ReturnType<typeof vi.fn>).mockResolvedValue({ data: mockData });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useDashboardStats(), { wrapper });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data).toEqual(mockData);
  });
});
