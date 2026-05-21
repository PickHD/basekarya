import { useMemo } from "react";
import { jwtDecode } from "jwt-decode";
import type { DecodedToken } from "@/features/auth/types";

export const usePermissions = () => {
  const { permissions, isPlatformAdmin, companyId } = useMemo(() => {
    try {
      const token = localStorage.getItem("token");
      if (!token) return { permissions: [] as string[], isPlatformAdmin: false, companyId: 0 };

      const decoded = jwtDecode<DecodedToken>(token);
      return {
        permissions: decoded.permissions || [],
        isPlatformAdmin: decoded.is_platform_admin || false,
        companyId: decoded.company_id || 0,
      };
    } catch (error) {
      console.error("Failed to decode token for permissions:", error);
      return { permissions: [] as string[], isPlatformAdmin: false, companyId: 0 };
    }
  }, []);

  const hasPermission = (permission: string) => {
    if (isPlatformAdmin) return true;
    return permissions.includes(permission);
  };

  const hasAnyPermission = (requiredPermissions: string[]) => {
    if (!requiredPermissions || requiredPermissions.length === 0) return true;
    if (isPlatformAdmin) return true;
    return requiredPermissions.some((p) => permissions.includes(p));
  };

  const hasAllPermissions = (requiredPermissions: string[]) => {
    if (!requiredPermissions || requiredPermissions.length === 0) return true;
    if (isPlatformAdmin) return true;
    return requiredPermissions.every((p) => permissions.includes(p));
  };

  return {
    permissions,
    hasPermission,
    hasAnyPermission,
    hasAllPermissions,
    isPlatformAdmin,
    companyId,
  };
};
