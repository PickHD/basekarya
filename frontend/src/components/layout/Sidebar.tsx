import { useState, useMemo } from "react";
import { Link, useLocation } from "react-router-dom";
import { cn } from "@/lib/utils";

import { menuItems } from "@/config/menu";
import { Button } from "@/components/ui/button";
import { useProfile } from "@/features/user/hooks/useProfile";
import { usePermissions } from "@/hooks/usePermissions";
import { usePlanModules } from "@/hooks/usePlanModules";
import { Skeleton } from "@/components/ui/skeleton";
import type { MenuItem } from "@/config/types";
import { ChevronDown, ChevronRight, PanelLeftClose, PanelLeft } from "lucide-react";

interface SidebarProps {
  className?: string;
  isCollapsed?: boolean;
  onToggleCollapse?: () => void;
}

export function Sidebar({ className, isCollapsed = false, onToggleCollapse }: SidebarProps) {
  const location = useLocation();
  const { data: user, isLoading } = useProfile();
  const { hasPermission, hasAnyPermission, isPlatformAdmin } = usePermissions();
  const { hasModule } = usePlanModules();
  const [collapsedGroups, setCollapsedGroups] = useState<Set<string>>(new Set());

  const filteredItems = useMemo(() => {
    return menuItems.filter((item) => {
      if (isPlatformAdmin && item.hideForPlatformAdmin) return false;
      if (!isPlatformAdmin && item.platformAdminOnly) return false;
      if (item.requiredModule && !hasModule(item.requiredModule)) return false;
      if (!item.permission) return true;
      if (Array.isArray(item.permission)) return hasAnyPermission(item.permission);
      return hasPermission(item.permission);
    });
  }, [hasPermission, hasAnyPermission, isPlatformAdmin, hasModule]);

  const groupedItems = useMemo(() => {
    const groups = new Map<string, MenuItem[]>();
    for (const item of filteredItems) {
      const group = item.group || "Lainnya";
      if (!groups.has(group)) {
        groups.set(group, []);
      }
      groups.get(group)!.push(item);
    }
    return groups;
  }, [filteredItems]);

  const isItemExactActive = (href: string) => {
    return location.pathname === href;
  };

  const activeGroup = useMemo(() => {
    let bestGroup: string | null = null;
    let bestLength = 0;
    for (const item of filteredItems) {
      if (location.pathname === item.href || location.pathname.startsWith(item.href + "/")) {
        if (item.href.length > bestLength) {
          bestGroup = item.group || "Lainnya";
          bestLength = item.href.length;
        }
      }
    }
    return bestGroup;
  }, [filteredItems, location.pathname]);

  const toggleGroup = (group: string) => {
    setCollapsedGroups((prev) => {
      const next = new Set(prev);
      if (next.has(group)) {
        next.delete(group);
      } else {
        next.add(group);
      }
      return next;
    });
  };

  const isGroupCollapsed = (group: string) => {
    if (group === activeGroup) return false;
    return collapsedGroups.has(group);
  };

  if (isCollapsed) {
    return (
      <div className={cn("flex flex-col h-screen bg-sidebar text-sidebar-foreground", className)}>
        <div className="flex-1 py-4 px-2 overflow-y-auto">
          <div className="flex items-center justify-center h-10 mb-4">
            <span className="text-xl font-bold text-white">B</span>
          </div>

          {isLoading ? (
            <div className="space-y-2">
              <Skeleton className="h-8 w-full bg-sidebar-muted" />
              <Skeleton className="h-8 w-full bg-sidebar-muted" />
            </div>
          ) : (
            <div className="space-y-1">
              {filteredItems.map((item) => {
                const isActive = isItemExactActive(item.href);

                return (
                  <Button
                    key={item.href}
                    variant="ghost"
                    size="icon"
                    className={cn(
                      "w-full text-sidebar-foreground/70 hover:text-sidebar-foreground hover:bg-sidebar-muted transition-all duration-200",
                      isActive && "bg-sidebar-accent text-sidebar-accent-foreground hover:bg-sidebar-accent hover:text-sidebar-accent-foreground"
                    )}
                    asChild
                  >
                    <Link to={item.href} title={item.title}>
                      <item.icon className="h-4 w-4" />
                    </Link>
                  </Button>
                );
              })}
            </div>
          )}
        </div>

        {onToggleCollapse && (
          <div className="border-t border-sidebar-border p-2 shrink-0">
            <Button
              variant="ghost"
              size="icon"
              onClick={onToggleCollapse}
              className="w-full text-sidebar-foreground/60 hover:text-sidebar-foreground hover:bg-sidebar-muted transition-all duration-200"
            >
              <PanelLeft className="h-4 w-4" />
            </Button>
          </div>
        )}
      </div>
    );
  }

  return (
    <div className={cn("flex flex-col h-screen bg-sidebar text-sidebar-foreground", className)}>
      <div className="flex-1 space-y-4 py-4 overflow-y-auto">
        <div className="px-3 py-2">
          <h2 className="mb-2 px-4 text-lg font-bold tracking-tight text-white transition-all duration-300">
            BaseKarya
          </h2>

          {isLoading && (
            <div className="space-y-2 px-2 mt-4">
              <Skeleton className="h-9 w-full bg-sidebar-muted" />
              <Skeleton className="h-9 w-full bg-sidebar-muted" />
            </div>
          )}

          {!isLoading && user && (
            <div className="mt-4 space-y-3">
              {Array.from(groupedItems.entries()).map(([group, items]) => (
                <div key={group}>
                  <button
                    onClick={() => toggleGroup(group)}
                    className="flex items-center justify-between w-full px-4 py-1.5 text-xs font-semibold uppercase text-sidebar-foreground/60 tracking-wider hover:text-sidebar-foreground transition-colors duration-200"
                  >
                    <span>{group}</span>
                    {isGroupCollapsed(group) ? (
                      <ChevronRight className="h-3 w-3 transition-transform duration-200" />
                    ) : (
                      <ChevronDown className="h-3 w-3 transition-transform duration-200" />
                    )}
                  </button>

                  <div
                    className={cn(
                      "overflow-hidden transition-all duration-300 ease-in-out",
                      isGroupCollapsed(group) ? "max-h-0 opacity-0" : "max-h-96 opacity-100"
                    )}
                  >
                    <div className="space-y-0.5 mt-1">
                      {items.map((item) => {
                        const isActive = isItemExactActive(item.href);

                        return (
                          <Button
                            key={item.href}
                            variant="ghost"
                            className={cn(
                              "w-full justify-start text-sidebar-foreground/70 hover:text-sidebar-foreground hover:bg-sidebar-muted transition-all duration-200",
                              isActive && "bg-sidebar-accent text-sidebar-accent-foreground hover:bg-sidebar-accent hover:text-sidebar-accent-foreground font-medium"
                            )}
                            asChild
                          >
                            <Link to={item.href}>
                              <item.icon className="mr-2 h-4 w-4 shrink-0" />
                              <span className="truncate">{item.title}</span>
                            </Link>
                          </Button>
                        );
                      })}
                    </div>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>

      {onToggleCollapse && (
        <div className="border-t border-sidebar-border p-2 shrink-0">
          <Button
            variant="ghost"
            onClick={onToggleCollapse}
            className="w-full justify-start text-sidebar-foreground/60 hover:text-sidebar-foreground hover:bg-sidebar-muted transition-all duration-200"
          >
            <PanelLeftClose className="h-4 w-4 mr-2 shrink-0" />
            <span className="text-xs">Collapse</span>
          </Button>
        </div>
      )}
    </div>
  );
}
