import { useProfile } from "@/features/user/hooks/useProfile";
import { Navigate, Outlet } from "react-router-dom";
import { Loader2 } from "lucide-react";
import type { ReactNode } from "react";
import { usePermissions } from "@/hooks/usePermissions";

interface ProtectedRouteProps {
  children?: ReactNode;
  requiredPermissions?: string[];
}

export const ProtectedRoute = ({
  children,
  requiredPermissions,
}: ProtectedRouteProps) => {
  const { data: user, isLoading } = useProfile();
  const { hasAnyPermission } = usePermissions();

  if (isLoading) {
    return (
      <div className="h-screen w-full flex items-center justify-center bg-slate-50">
        <div className="flex flex-col items-center gap-2">
          <Loader2 className="h-8 w-8 animate-spin text-blue-600" />
          <p className="text-sm text-slate-500">Verifying access...</p>
        </div>
      </div>
    );
  }

  if (!user) {
    return <Navigate to="/login" replace />;
  }

  if (requiredPermissions && !hasAnyPermission(requiredPermissions)) {
    return <Navigate to="/dashboard" replace />;
  }

  return children ? <>{children}</> : <Outlet />;
};
