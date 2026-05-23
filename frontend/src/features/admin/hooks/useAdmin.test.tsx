import { describe, it, expect, beforeEach, vi } from "vitest";
import { renderHook, waitFor } from "@testing-library/react";
import { createQueryWrapper } from "@/test/utils";
import {
  useAllEmployees,
  useEmployeeMutations,
  useDashboardStats,
} from "@/features/admin/hooks/useAdmin";
import {
  useDepartments,
  useShifts,
} from "@/features/admin/hooks/useMasterData";

vi.mock("@/lib/axios", () => ({
  api: { get: vi.fn(), post: vi.fn(), put: vi.fn(), delete: vi.fn() },
}));

vi.mock("sonner", () => ({
  toast: { success: vi.fn(), error: vi.fn() },
}));

import { api } from "@/lib/axios";
import { toast } from "sonner";

describe("useAllEmployees", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should fetch employees with page and search", async () => {
    const mockData = { data: [{ id: 1, name: "John" }], meta: { limit: 10 } };
    (api.get as ReturnType<typeof vi.fn>).mockResolvedValue({ data: mockData });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useAllEmployees(1, "john"), {
      wrapper,
    });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(api.get).toHaveBeenCalledWith("/employees", {
      params: { page: 1, limit: 10, search: "john" },
    });
  });
});

describe("useEmployeeMutations", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should create employee and show toast", async () => {
    (api.post as ReturnType<typeof vi.fn>).mockResolvedValue({
      data: { data: { username: "john123" } },
    });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useEmployeeMutations(), { wrapper });

    result.current.createMutation.mutate({
      full_name: "John",
      email: "john@test.com",
      phone: "081234567890",
      department_id: 1,
      shift_id: 1,
      role_ids: [1],
      join_date: "2024-01-01",
      salary: 5000000,
    });

    await waitFor(() => expect(result.current.createMutation.isSuccess).toBe(true));
    expect(toast.success).toHaveBeenCalledWith("Employee created", {
      duration: 1000,
    });
  });

  it("should create employee without username in response", async () => {
    (api.post as ReturnType<typeof vi.fn>).mockResolvedValue({
      data: { message: "ok" },
    });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useEmployeeMutations(), { wrapper });

    result.current.createMutation.mutate({
      full_name: "John",
      email: "john@test.com",
      department_id: 1,
      shift_id: 1,
      base_salary: 5000000,
    });

    await waitFor(() =>
      expect(result.current.createMutation.isSuccess).toBe(true),
    );
    expect(toast.success).toHaveBeenCalledWith("Employee created successfully");
  });

  it("should update employee", async () => {
    (api.put as ReturnType<typeof vi.fn>).mockResolvedValue({
      data: { message: "ok" },
    });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useEmployeeMutations(), { wrapper });

    result.current.updateMutation.mutate({
      id: 1,
      data: { full_name: "Jane" },
    });

    await waitFor(() =>
      expect(result.current.updateMutation.isSuccess).toBe(true),
    );
    expect(toast.success).toHaveBeenCalledWith("Employee updated successfully");
  });

  it("should delete employee", async () => {
    (api.delete as ReturnType<typeof vi.fn>).mockResolvedValue({
      data: { message: "ok" },
    });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useEmployeeMutations(), { wrapper });

    result.current.deleteMutation.mutate(1);

    await waitFor(() =>
      expect(result.current.deleteMutation.isSuccess).toBe(true),
    );
    expect(toast.success).toHaveBeenCalledWith("Employee deleted");
  });

  it("should show error on create failure", async () => {
    (api.post as ReturnType<typeof vi.fn>).mockRejectedValue({
      response: { data: { message: "Email exists" } },
    });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useEmployeeMutations(), { wrapper });

    result.current.createMutation.mutate({
      full_name: "John",
      email: "john@test.com",
      phone: "081234567890",
      department_id: 1,
      shift_id: 1,
      role_ids: [1],
      join_date: "2024-01-01",
      salary: 5000000,
    });

    await waitFor(() =>
      expect(result.current.createMutation.isError).toBe(true),
    );
    expect(toast.error).toHaveBeenCalled();
  });
});

describe("useDashboardStats (admin)", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should fetch dashboard stats", async () => {
    const mockStats = { total_employees: 50, present_today: 45 };
    (api.get as ReturnType<typeof vi.fn>).mockResolvedValue({
      data: { data: mockStats },
    });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useDashboardStats(), { wrapper });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data).toEqual(mockStats);
  });
});

describe("useDepartments", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should fetch departments", async () => {
    const mockData = [{ id: 1, name: "Engineering" }];
    (api.get as ReturnType<typeof vi.fn>).mockResolvedValue({
      data: { data: mockData },
    });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useDepartments(), { wrapper });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data).toEqual(mockData);
  });
});

describe("useShifts", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should fetch shifts", async () => {
    const mockData = [{ id: 1, name: "Morning" }];
    (api.get as ReturnType<typeof vi.fn>).mockResolvedValue({
      data: { data: mockData },
    });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useShifts(), { wrapper });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data).toEqual(mockData);
  });
});
