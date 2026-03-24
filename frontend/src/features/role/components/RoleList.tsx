import { useState } from "react";
import { useRoles, useRoleDetails } from "../hooks/useRole";
import { Button } from "@/components/ui/button";
import { RoleCreateDialog } from "./RoleCreateDialog";
import { RolePermissionsDialog } from "./RolePermissionsDialog";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Shield, Settings, Plus, Loader2 } from "lucide-react";
import type { Role } from "../types";

export function RoleList() {
  const { data: roles, isLoading: rolesLoading } = useRoles();
  const [isCreateOpen, setIsCreateOpen] = useState(false);
  const [selectedRoleId, setSelectedRoleId] = useState<number | null>(null);

  const { data: selectedRoleData, isLoading: roleDetailsLoading } = useRoleDetails(
    selectedRoleId?.toString() || ""
  );

  const isPermissionsDialogOpen = !!selectedRoleId && !!selectedRoleData;

  const handleManagePermissions = (roleId: number) => {
    setSelectedRoleId(roleId);
  };

  const handleClosePermissionsDialog = () => {
    setSelectedRoleId(null);
  };

  return (
    <div className="space-y-6">
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <div>
          <h2 className="text-2xl sm:text-3xl font-bold tracking-tight">Roles Management</h2>
          <p className="text-sm sm:text-base text-slate-500">
            Configure system roles and define access permissions.
          </p>
        </div>
        <Button
          onClick={() => setIsCreateOpen(true)}
          className="bg-blue-600 hover:bg-blue-700 w-full sm:w-auto"
        >
          <Plus className="mr-2 h-4 w-4" /> Add Role
        </Button>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {rolesLoading ? (
          <div className="col-span-full flex justify-center py-20">
            <Loader2 className="h-8 w-8 animate-spin text-blue-600" />
          </div>
        ) : roles?.data?.length === 0 ? (
          <div className="col-span-full text-center py-20 border rounded-lg bg-slate-50 text-slate-500">
            <Shield className="h-12 w-12 mx-auto mb-4 text-slate-300" />
            <p className="text-lg font-medium text-slate-900">No roles found</p>
            <p>Get started by creating a new role for your system.</p>
          </div>
        ) : (
          roles?.data?.map((role: Role) => (
            <Card key={role.id} className="flex flex-col">
              <CardHeader className="pb-3 border-b">
                <CardTitle className="text-xl flex items-center justify-between">
                  <span className="flex items-center gap-2">
                    <Shield className="h-5 w-5 text-blue-600" />
                    {role.name}
                  </span>
                </CardTitle>
              </CardHeader>
              <CardContent className="pt-4 flex-1 flex flex-col">
                <p className="text-sm text-slate-600 mb-6 flex-1">
                  {role.description || "No description provided."}
                </p>

                <Button
                  className="w-full mt-auto bg-slate-100 hover:bg-slate-200 text-slate-900 border"
                  onClick={() => handleManagePermissions(role.id)}
                  disabled={selectedRoleId === role.id && roleDetailsLoading}
                >
                  {selectedRoleId === role.id && roleDetailsLoading ? (
                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  ) : (
                    <Settings className="mr-2 h-4 w-4 text-slate-600" />
                  )}
                  Manage Permissions
                </Button>
              </CardContent>
            </Card>
          ))
        )}
      </div>

      <RoleCreateDialog open={isCreateOpen} onOpenChange={setIsCreateOpen} />

      {selectedRoleData && (
        <RolePermissionsDialog
          key={selectedRoleData.id}
          open={isPermissionsDialogOpen}
          onOpenChange={handleClosePermissionsDialog}
          role={selectedRoleData}
        />
      )}
    </div>
  );
}
