import { describe, it, expect, beforeEach, vi } from "vitest";
import { renderHook, waitFor } from "@testing-library/react";
import { createQueryWrapper } from "@/test/utils";
import {
  useRoles,
  useCreateRole,
  useRoleDetails,
  useAllPermissions,
  useAssignPermissions,
} from "@/features/role/hooks/useRole";

vi.mock("@/lib/axios", () => ({
  api: { get: vi.fn(), post: vi.fn(), put: vi.fn() },
}));

vi.mock("sonner", () => ({
  toast: { success: vi.fn(), error: vi.fn() },
}));

import { api } from "@/lib/axios";
import { toast } from "sonner";

describe("useRoles", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should fetch roles", async () => {
    const mockRoles = { data: [{ id: 1, name: "admin" }, { id: 2, name: "employee" }] };
    (api.get as ReturnType<typeof vi.fn>).mockResolvedValue({ data: mockRoles });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useRoles(), { wrapper });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data).toEqual(mockRoles);
  });
});

describe("useCreateRole", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should create role and show toast", async () => {
    (api.post as ReturnType<typeof vi.fn>).mockResolvedValue({ data: { message: "ok" } });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useCreateRole(), { wrapper });

    result.current.mutate({ name: "manager" });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(toast.success).toHaveBeenCalledWith("Role Created");
  });

  it("should show error on failure", async () => {
    (api.post as ReturnType<typeof vi.fn>).mockRejectedValue({
      response: { data: { message: "Role exists" } },
    });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useCreateRole(), { wrapper });

    result.current.mutate({ name: "admin" });

    await waitFor(() => expect(result.current.isError).toBe(true));
    expect(toast.error).toHaveBeenCalledWith("Role exists");
  });
});

describe("useRoleDetails", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should fetch role details with permissions", async () => {
    (api.get as ReturnType<typeof vi.fn>).mockResolvedValue({
      data: { data: { role_id: 1, role_name: "admin", permissions: ["VIEW_EMPLOYEE"] } },
    });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useRoleDetails("1"), { wrapper });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data?.name).toBe("admin");
    expect(result.current.data?.permissions).toEqual(["VIEW_EMPLOYEE"]);
  });

  it("should not fetch when id is empty", () => {
    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useRoleDetails(""), { wrapper });
    expect(result.current.fetchStatus).toBe("idle");
  });
});

describe("useAllPermissions", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should fetch all permissions", async () => {
    const mockPerms = { data: [{ id: 1, name: "VIEW_EMPLOYEE" }] };
    (api.get as ReturnType<typeof vi.fn>).mockResolvedValue({ data: mockPerms });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useAllPermissions(), { wrapper });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data).toEqual(mockPerms.data);
  });
});

describe("useAssignPermissions", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should assign permissions and show toast", async () => {
    (api.put as ReturnType<typeof vi.fn>).mockResolvedValue({ data: { message: "ok" } });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useAssignPermissions(), { wrapper });

    result.current.mutate({ role_id: 1, permission_ids: [1, 2, 3] });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(toast.success).toHaveBeenCalledWith("Permissions Updated");
    expect(api.put).toHaveBeenCalledWith("/roles/1/permissions", {
      permission_ids: [1, 2, 3],
    });
  });
});
