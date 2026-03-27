import { useState } from "react";
import { useRoles, useRoleDetails } from "../hooks/useRole";
import { Button } from "@/components/ui/button";
import { RoleCreateDialog } from "@/features/role/components/RoleCreateDialog";
import { RolePermissionsDialog } from "@/features/role/components/RolePermissionsDialog";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Shield, Settings, Plus, Loader2 } from "lucide-react";
import type { Role } from "@/features/role/types";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";

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
          <h2 className="text-2xl sm:text-3xl font-bold tracking-tight">
            Roles Management
          </h2>
          <p className="text-sm sm:text-base text-slate-500">
            Manage system roles and their permissions.
          </p>
        </div>
        <div className="flex gap-2 w-full sm:w-auto">
          <Button
            onClick={() => setIsCreateOpen(true)}
            className="bg-blue-600 hover:bg-blue-700 w-full sm:w-auto"
          >
            <Plus className="mr-2 h-4 w-4" /> Add Role
          </Button>
        </div>
      </div>

      <Card>
        <CardHeader>
          <div className="flex flex-col md:flex-row justify-between md:items-center gap-4">
            <CardTitle className="flex items-center gap-2 text-lg">
              <Shield className="h-5 w-5" /> Existing Roles
            </CardTitle>
          </div>
        </CardHeader>
        <CardContent>
          {rolesLoading ? (
            <div className="flex justify-center py-10">
              <Loader2 className="animate-spin h-8 w-8 text-blue-600" />
            </div>
          ) : roles?.data?.length === 0 ? (
            <div className="col-span-full text-center py-20 border rounded-lg bg-slate-50 text-slate-500">
              <Shield className="h-12 w-12 mx-auto mb-4 text-slate-300" />
              <p className="text-lg font-medium text-slate-900">No roles found</p>
              <p>Get started by creating a new role for your system.</p>
            </div>
          ) : (
            <>
              <div className="grid grid-cols-1 gap-4 md:hidden">
                {roles?.data?.map((role: Role) => (
                  <div
                    key={role.id}
                    className="flex flex-col rounded-lg border bg-card p-4 shadow-sm space-y-3"
                  >
                    <div className="flex justify-between items-start gap-2">
                      <div>
                        <h4 className="font-semibold text-lg flex items-center gap-2">
                          <Shield className="h-4 w-4 text-blue-600" />
                          {role.name}
                        </h4>
                      </div>
                    </div>

                    <div className="pt-2 border-t mt-2">
                      <Button
                        variant="outline"
                        size="sm"
                        className="w-full"
                        onClick={() => handleManagePermissions(role.id)}
                        disabled={selectedRoleId === role.id && roleDetailsLoading}
                      >
                        {selectedRoleId === role.id && roleDetailsLoading ? (
                          <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                        ) : (
                          <Settings className="mr-2 h-4 w-4" />
                        )}
                        Manage Permissions
                      </Button>
                    </div>
                  </div>
                ))}
              </div>

              {/* Desktop View: Data Table */}
              <div className="hidden md:block rounded-md border">
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>Role Name</TableHead>
                      <TableHead className="text-right">Action</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {roles?.data?.map((role: Role) => (
                      <TableRow key={role.id}>
                        <TableCell className="font-medium">
                          <div className="flex items-center gap-2">
                            <Shield className="h-4 w-4 text-blue-600" />
                            {role.name}
                          </div>
                        </TableCell>
                        <TableCell className="text-right">
                          <Button
                            variant="ghost"
                            size="icon"
                            className="hover:bg-slate-100"
                            onClick={() => handleManagePermissions(role.id)}
                            disabled={selectedRoleId === role.id && roleDetailsLoading}
                            title="Manage Permissions"
                          >
                            {selectedRoleId === role.id && roleDetailsLoading ? (
                              <Loader2 className="h-4 w-4 animate-spin" />
                            ) : (
                              <Settings className="h-4 w-4 text-slate-500" />
                            )}
                          </Button>
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </div>
            </>
          )}
        </CardContent>
      </Card>

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
