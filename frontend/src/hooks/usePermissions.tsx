import { useMemo } from "react";
import { jwtDecode } from "jwt-decode";
import type { DecodedToken } from "@/features/auth/types";

export const usePermissions = () => {
  const permissions = useMemo(() => {
    try {
      const token = localStorage.getItem("token");
      if (!token) return [];

      const decoded = jwtDecode<DecodedToken>(token);
      return decoded.permissions || [];
    } catch (error) {
      console.error("Failed to decode token for permissions:", error);
      return [];
    }
  }, []);

  const hasPermission = (permission: string) => {
    return permissions.includes(permission);
  };

  const hasAnyPermission = (requiredPermissions: string[]) => {
    if (!requiredPermissions || requiredPermissions.length === 0) return true;
    return requiredPermissions.some((p) => permissions.includes(p));
  };

  const hasAllPermissions = (requiredPermissions: string[]) => {
    if (!requiredPermissions || requiredPermissions.length === 0) return true;
    return requiredPermissions.every((p) => permissions.includes(p));
  };

  return {
    permissions,
    hasPermission,
    hasAnyPermission,
    hasAllPermissions,
  };
};
