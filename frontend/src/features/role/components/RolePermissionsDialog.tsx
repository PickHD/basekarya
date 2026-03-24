import { useState } from "react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import { Loader2, ShieldAlert } from "lucide-react";
import { useAllPermissions, useAssignPermissions } from "../hooks/useRole";
import type { Role, Permission } from "../types";

interface RolePermissionsDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  role: Role | null;
}

export function RolePermissionsDialog({
  open,
  onOpenChange,
  role,
}: RolePermissionsDialogProps) {
  const { data: allPermissions, isLoading: isLoadingPermissions } = useAllPermissions();
  const { mutate: assignPermissions, isPending } = useAssignPermissions();

  const [selectedIds, setSelectedIds] = useState<number[]>(
    role?.permissions?.map((p) => p.id) || []
  );

  const handleToggle = (id: number) => {
    setSelectedIds((prev) =>
      prev.includes(id) ? prev.filter((pId) => pId !== id) : [...prev, id]
    );
  };

  const handleSelectAll = (modulePermissions: Permission[]) => {
    const moduleIds = modulePermissions.map((p) => p.id);
    const allSelected = moduleIds.every((id) => selectedIds.includes(id));

    if (allSelected) {
      setSelectedIds((prev) => prev.filter((id) => !moduleIds.includes(id)));
    } else {
      setSelectedIds((prev) => Array.from(new Set([...prev, ...moduleIds])));
    }
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!role) return;

    assignPermissions(
      { role_id: role.id, permission_ids: selectedIds },
      {
        onSuccess: () => {
          onOpenChange(false);
        },
      }
    );
  };

  const groupedPermissions = allPermissions?.reduce((acc: Record<string, Permission[]>, curr: Permission) => {
    const moduleName = curr.module || "General";
    if (!acc[moduleName]) {
      acc[moduleName] = [];
    }
    acc[moduleName].push(curr);
    return acc;
  }, {} as Record<string, Permission[]>);

  if (!role) return null;

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-3xl max-h-[90vh] flex flex-col">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <ShieldAlert className="h-5 w-5 text-blue-600" />
            Manage Permissions: {role.name}
          </DialogTitle>
          <DialogDescription>
            Select the access levels and permissions to grant to this role across various modules.
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit} className="flex flex-col overflow-hidden h-full">
          <div className="flex-1 overflow-y-auto py-4 px-1 space-y-6">
            {isLoadingPermissions ? (
              <div className="flex justify-center py-10">
                <Loader2 className="h-8 w-8 animate-spin text-blue-600" />
              </div>
            ) : groupedPermissions ? (
              Object.entries(groupedPermissions).map(([module, permissions]) => {
                const moduleIds = permissions.map((p) => p.id);
                const allSelected = moduleIds.every((id) => selectedIds.includes(id));
                const someSelected = moduleIds.some((id) => selectedIds.includes(id)) && !allSelected;

                return (
                  <div key={module} className="border rounded-lg overflow-hidden">
                    <div className="bg-slate-50 px-4 py-3 flex items-center justify-between border-b">
                      <h4 className="font-semibold text-slate-800 capitalize">{module}</h4>
                      <div className="flex items-center gap-2">
                        <span className="text-xs text-slate-500">Select All</span>
                        <Checkbox
                          checked={allSelected ? true : someSelected ? "indeterminate" : false}
                          onCheckedChange={() => handleSelectAll(permissions)}
                        />
                      </div>
                    </div>
                    <div className="p-4 grid grid-cols-1 md:grid-cols-2 gap-4">
                      {permissions.map((perm) => (
                        <div key={perm.id} className="flex items-start space-x-3">
                          <Checkbox
                            id={`perm-${perm.id}`}
                            checked={selectedIds.includes(perm.id)}
                            onCheckedChange={() => handleToggle(perm.id)}
                            className="mt-1"
                          />
                          <div className="space-y-1 leading-none">
                            <label
                              htmlFor={`perm-${perm.id}`}
                              className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70 cursor-pointer"
                            >
                              {perm.name}
                            </label>
                            {perm.description && (
                              <p className="text-xs text-muted-foreground">
                                {perm.description}
                              </p>
                            )}
                          </div>
                        </div>
                      ))}
                    </div>
                  </div>
                );
              })
            ) : (
              <div className="text-center py-10 text-slate-500">Failed to load permissions.</div>
            )}
          </div>

          <DialogFooter className="pt-4 border-t mt-auto">
            <Button
              type="button"
              variant="outline"
              onClick={() => onOpenChange(false)}
              disabled={isPending}
            >
              Cancel
            </Button>
            <Button type="submit" className="bg-blue-600 hover:bg-blue-700" disabled={isPending || isLoadingPermissions}>
              {isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              Save Permissions
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
