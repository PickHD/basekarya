import { describe, it, expect, beforeEach, vi } from "vitest";
import { renderHook, waitFor } from "@testing-library/react";
import { createQueryWrapper } from "@/test/utils";
import {
  useFinanceTransactions,
  useFinanceTransaction,
  useCreateFinanceTransaction,
  useFinanceAction,
} from "@/features/finance/hooks/useFinance";
import {
  useFinanceCategories,
  useCreateFinanceCategory,
  useUpdateFinanceCategory,
  useDeleteFinanceCategory,
} from "@/features/finance/hooks/useFinanceCategory";
import { useFinanceDashboard } from "@/features/finance/hooks/useFinanceDashboard";

vi.mock("@/lib/axios", () => ({
  api: { get: vi.fn(), post: vi.fn(), put: vi.fn(), delete: vi.fn() },
}));

vi.mock("sonner", () => ({
  toast: { success: vi.fn(), error: vi.fn() },
}));

import { api } from "@/lib/axios";
import { toast } from "sonner";

describe("useFinanceTransactions", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should fetch transactions with filter using infinite query", async () => {
    const mockData = { data: [{ id: 1, amount: 500000 }], meta: { limit: 10 } };
    (api.get as ReturnType<typeof vi.fn>).mockResolvedValue({ data: mockData });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(
      () => useFinanceTransactions({ limit: 10 }),
      { wrapper }
    );

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data?.pages[0]).toEqual(mockData);
  });
});

describe("useFinanceTransaction", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should fetch transaction detail", async () => {
    const mockDetail = { id: 1, amount: 500000, status: "pending" };
    (api.get as ReturnType<typeof vi.fn>).mockResolvedValue({ data: { data: mockDetail } });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useFinanceTransaction("1"), { wrapper });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data).toEqual(mockDetail);
  });

  it("should not fetch when id is empty", () => {
    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useFinanceTransaction(""), { wrapper });
    expect(result.current.fetchStatus).toBe("idle");
  });
});

describe("useCreateFinanceTransaction", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should create transaction and show toast", async () => {
    (api.post as ReturnType<typeof vi.fn>).mockResolvedValue({ data: { message: "ok" } });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useCreateFinanceTransaction(), { wrapper });

    result.current.mutate({
      category_id: 1,
      amount: 500000,
      date: "2024-01-01",
      description: "Test",
      type: "income",
    });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(toast.success).toHaveBeenCalledWith("Transaksi keuangan berhasil dibuat!");
  });
});

describe("useFinanceAction", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should approve finance transaction", async () => {
    (api.put as ReturnType<typeof vi.fn>).mockResolvedValue({ data: { message: "ok" } });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useFinanceAction(), { wrapper });

    result.current.mutate({ id: 1, action: "approve", rejection_reason: "" });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(toast.success).toHaveBeenCalledWith("Transaksi berhasil di-approve");
  });
});

describe("useFinanceCategories", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should fetch categories", async () => {
    const mockCats = [{ id: 1, name: "Salary" }];
    (api.get as ReturnType<typeof vi.fn>).mockResolvedValue({ data: { data: mockCats } });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useFinanceCategories(), { wrapper });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data).toEqual(mockCats);
  });

  it("should fetch categories with type filter", async () => {
    (api.get as ReturnType<typeof vi.fn>).mockResolvedValue({ data: { data: [] } });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useFinanceCategories("income"), { wrapper });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(api.get).toHaveBeenCalledWith("/finances/categories", {
      params: { type: "income" },
    });
  });
});

describe("useCreateFinanceCategory", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should create category", async () => {
    (api.post as ReturnType<typeof vi.fn>).mockResolvedValue({ data: { message: "ok" } });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useCreateFinanceCategory(), { wrapper });

    result.current.mutate({ name: "Bonus", type: "income" });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(toast.success).toHaveBeenCalledWith("Kategori berhasil ditambahkan!");
  });
});

describe("useUpdateFinanceCategory", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should update category", async () => {
    (api.put as ReturnType<typeof vi.fn>).mockResolvedValue({ data: { message: "ok" } });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useUpdateFinanceCategory(), { wrapper });

    result.current.mutate({ id: 1, name: "Bonus Updated", type: "income" });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(toast.success).toHaveBeenCalledWith("Kategori berhasil diperbarui!");
  });
});

describe("useDeleteFinanceCategory", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should delete category", async () => {
    (api.delete as ReturnType<typeof vi.fn>).mockResolvedValue({ data: { message: "ok" } });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useDeleteFinanceCategory(), { wrapper });

    result.current.mutate(1);

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(toast.success).toHaveBeenCalledWith("Kategori berhasil dihapus!");
  });
});

describe("useFinanceDashboard", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should fetch dashboard data", async () => {
    const mockDashboard = { total_income: 50000000, total_expense: 30000000 };
    (api.get as ReturnType<typeof vi.fn>).mockResolvedValue({ data: { data: mockDashboard } });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(
      () => useFinanceDashboard("2024-01-01", "2024-01-31"),
      { wrapper }
    );

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data).toEqual(mockDashboard);
  });

  it("should fetch dashboard without date params", async () => {
    (api.get as ReturnType<typeof vi.fn>).mockResolvedValue({ data: { data: {} } });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useFinanceDashboard(), { wrapper });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(api.get).toHaveBeenCalledWith("/finances/dashboard", { params: {} });
  });
});
