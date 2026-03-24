import { Link, useLocation } from "react-router-dom";
import { cn } from "@/lib/utils";

import { menuItems } from "@/config/menu";
import { Button } from "@/components/ui/button";
import { useProfile } from "@/features/user/hooks/useProfile";
import { usePermissions } from "@/hooks/usePermissions";
import { Skeleton } from "@/components/ui/skeleton";
import type { MenuItem } from "@/config/types";

export function Sidebar({ className }: { className?: string }) {
  const location = useLocation();
  const { data: user, isLoading } = useProfile();
  const { hasPermission, hasAnyPermission } = usePermissions();

  const renderMenuItems = (items: MenuItem[]) => (
    <div className="space-y-1">
      {items.map((item) => (
        <Button
          key={item.href}
          variant={location.pathname === item.href ? "secondary" : "ghost"}
          className={cn(
            "w-full justify-start",
            location.pathname === item.href && "bg-slate-200"
          )}
          asChild
        >
          <Link to={item.href}>
            <item.icon className="mr-2 h-4 w-4" />
            {item.title}
          </Link>
        </Button>
      ))}
    </div>
  );

  return (
    <div className={cn("pb-12 min-h-screen border-r bg-background", className)}>
      <div className="space-y-4 py-4">
        {/* HEADER */}
        <div className="px-3 py-2">
          <h2 className="mb-2 px-4 text-lg font-semibold tracking-tight">
            BaseKarya
          </h2>

          {/* LOADING STATE */}
          {isLoading && (
            <div className="space-y-2 px-2 mt-4">
              <Skeleton className="h-9 w-full" />
              <Skeleton className="h-9 w-full" />
            </div>
          )}

          {/* MENU */}
          {!isLoading && user && (
            <div className="mt-4">
              <h3 className="mb-2 px-4 text-xs font-semibold uppercase text-muted-foreground tracking-wider">
                My Workspace
              </h3>
              {renderMenuItems(
                menuItems.filter((item) => {
                  if (!item.permission) return true;
                  if (Array.isArray(item.permission)) return hasAnyPermission(item.permission);
                  return hasPermission(item.permission);
                })
              )}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
