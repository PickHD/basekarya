import { describe, it, expect, beforeEach } from "vitest";
import { renderHook } from "@testing-library/react";
import { usePermissions } from "@/hooks/usePermissions";

function createMockToken(payload: Record<string, unknown>): string {
  const header = btoa(JSON.stringify({ alg: "HS256", typ: "JWT" }));
  const body = btoa(JSON.stringify(payload));
  const signature = btoa("mock-signature");
  return `${header}.${body}.${signature}`;
}

describe("usePermissions", () => {
  beforeEach(() => {
    localStorage.clear();
  });

  it("should return empty permissions when no token", () => {
    const { result } = renderHook(() => usePermissions());
    expect(result.current.permissions).toEqual([]);
    expect(result.current.isPlatformAdmin).toBe(false);
    expect(result.current.companyId).toBe(0);
  });

  it("should decode permissions from valid token", () => {
    const token = createMockToken({
      user_id: 1,
      role: "admin",
      permissions: ["VIEW_EMPLOYEE", "CREATE_EMPLOYEE"],
      company_id: 5,
      is_platform_admin: false,
      exp: Date.now() / 1000 + 3600,
    });
    localStorage.setItem("token", token);

    const { result } = renderHook(() => usePermissions());
    expect(result.current.permissions).toEqual(["VIEW_EMPLOYEE", "CREATE_EMPLOYEE"]);
    expect(result.current.isPlatformAdmin).toBe(false);
    expect(result.current.companyId).toBe(5);
  });

  it("should detect platform admin", () => {
    const token = createMockToken({
      user_id: 1,
      role: "superadmin",
      permissions: [],
      company_id: 0,
      is_platform_admin: true,
      exp: Date.now() / 1000 + 3600,
    });
    localStorage.setItem("token", token);

    const { result } = renderHook(() => usePermissions());
    expect(result.current.isPlatformAdmin).toBe(true);
  });

  describe("hasPermission", () => {
    it("should return true for platform admin regardless of permission", () => {
      const token = createMockToken({
        user_id: 1,
        role: "superadmin",
        permissions: [],
        company_id: 0,
        is_platform_admin: true,
        exp: Date.now() / 1000 + 3600,
      });
      localStorage.setItem("token", token);

      const { result } = renderHook(() => usePermissions());
      expect(result.current.hasPermission("VIEW_EMPLOYEE")).toBe(true);
    });

    it("should return true when user has the permission", () => {
      const token = createMockToken({
        user_id: 1,
        role: "admin",
        permissions: ["VIEW_EMPLOYEE"],
        company_id: 1,
        is_platform_admin: false,
        exp: Date.now() / 1000 + 3600,
      });
      localStorage.setItem("token", token);

      const { result } = renderHook(() => usePermissions());
      expect(result.current.hasPermission("VIEW_EMPLOYEE")).toBe(true);
    });

    it("should return false when user lacks the permission", () => {
      const token = createMockToken({
        user_id: 1,
        role: "employee",
        permissions: ["VIEW_SELF_ATTENDANCE"],
        company_id: 1,
        is_platform_admin: false,
        exp: Date.now() / 1000 + 3600,
      });
      localStorage.setItem("token", token);

      const { result } = renderHook(() => usePermissions());
      expect(result.current.hasPermission("VIEW_EMPLOYEE")).toBe(false);
    });
  });

  describe("hasAnyPermission", () => {
    it("should return true when empty array", () => {
      const { result } = renderHook(() => usePermissions());
      expect(result.current.hasAnyPermission([])).toBe(true);
    });

    it("should return true when platform admin", () => {
      const token = createMockToken({
        user_id: 1,
        role: "superadmin",
        permissions: [],
        company_id: 0,
        is_platform_admin: true,
        exp: Date.now() / 1000 + 3600,
      });
      localStorage.setItem("token", token);

      const { result } = renderHook(() => usePermissions());
      expect(result.current.hasAnyPermission(["VIEW_EMPLOYEE"])).toBe(true);
    });

    it("should return true when user has at least one permission", () => {
      const token = createMockToken({
        user_id: 1,
        role: "admin",
        permissions: ["CREATE_EMPLOYEE"],
        company_id: 1,
        is_platform_admin: false,
        exp: Date.now() / 1000 + 3600,
      });
      localStorage.setItem("token", token);

      const { result } = renderHook(() => usePermissions());
      expect(result.current.hasAnyPermission(["VIEW_EMPLOYEE", "CREATE_EMPLOYEE"])).toBe(true);
    });

    it("should return false when user has none of the permissions", () => {
      const token = createMockToken({
        user_id: 1,
        role: "employee",
        permissions: ["VIEW_SELF_ATTENDANCE"],
        company_id: 1,
        is_platform_admin: false,
        exp: Date.now() / 1000 + 3600,
      });
      localStorage.setItem("token", token);

      const { result } = renderHook(() => usePermissions());
      expect(result.current.hasAnyPermission(["VIEW_EMPLOYEE", "CREATE_EMPLOYEE"])).toBe(false);
    });
  });

  describe("hasAllPermissions", () => {
    it("should return true when empty array", () => {
      const { result } = renderHook(() => usePermissions());
      expect(result.current.hasAllPermissions([])).toBe(true);
    });

    it("should return true when user has all permissions", () => {
      const token = createMockToken({
        user_id: 1,
        role: "admin",
        permissions: ["VIEW_EMPLOYEE", "CREATE_EMPLOYEE"],
        company_id: 1,
        is_platform_admin: false,
        exp: Date.now() / 1000 + 3600,
      });
      localStorage.setItem("token", token);

      const { result } = renderHook(() => usePermissions());
      expect(result.current.hasAllPermissions(["VIEW_EMPLOYEE", "CREATE_EMPLOYEE"])).toBe(true);
    });

    it("should return false when user is missing a permission", () => {
      const token = createMockToken({
        user_id: 1,
        role: "admin",
        permissions: ["VIEW_EMPLOYEE"],
        company_id: 1,
        is_platform_admin: false,
        exp: Date.now() / 1000 + 3600,
      });
      localStorage.setItem("token", token);

      const { result } = renderHook(() => usePermissions());
      expect(result.current.hasAllPermissions(["VIEW_EMPLOYEE", "CREATE_EMPLOYEE"])).toBe(false);
    });
  });

  it("should handle invalid token gracefully", () => {
    const errorSpy = vi.spyOn(console, "error").mockImplementation(() => {});
    localStorage.setItem("token", "not-a-valid-jwt");

    const { result } = renderHook(() => usePermissions());
    expect(result.current.permissions).toEqual([]);
    expect(result.current.isPlatformAdmin).toBe(false);
    errorSpy.mockRestore();
  });
});
