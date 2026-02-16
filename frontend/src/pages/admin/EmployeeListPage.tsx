import { useState } from "react";
import {
  useAllEmployees,
  useEmployeeMutations,
} from "@/features/admin/hooks/useAdmin";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Loader2, Search, Users, Plus, Pencil, Trash2 } from "lucide-react";
import type { Employee } from "@/features/admin/types";
import { EmployeeFormDialog } from "@/features/admin/components/EmployeeFormDialog";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import { useDebounce } from "@/hooks/useDebounce";
import { PaginationControls } from "@/components/shared/PaginationControls";

export default function EmployeeListPage() {
  const [page, setPage] = useState(1);
  const [search, setSearch] = useState("");

  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [selectedEmployee, setSelectedEmployee] = useState<Employee | null>(
    null,
  );

  const [employeeToDelete, setEmployeeToDelete] = useState<number | null>(null);

  const debouncedSearch = useDebounce(search, 500);

  const { data, isLoading } = useAllEmployees(page, debouncedSearch);
  const { createMutation, updateMutation, deleteMutation } =
    useEmployeeMutations();

  const handleAdd = () => {
    setSelectedEmployee(null);
    setIsDialogOpen(true);
  };

  const handleEdit = (emp: Employee) => {
    setSelectedEmployee(emp);
    setIsDialogOpen(true);
  };

  const handleDeleteClick = (id: number) => {
    setEmployeeToDelete(id);
  };

  const confirmDelete = async () => {
    if (employeeToDelete) {
      await deleteMutation.mutateAsync(employeeToDelete);
      setEmployeeToDelete(null);
    }
  };

  const handleFormSubmit = async (values: any) => {
    if (selectedEmployee) {
      // Update Logic
      await updateMutation.mutateAsync({
        id: selectedEmployee.id,
        data: values,
      });
    } else {
      // Create Logic
      await createMutation.mutateAsync(values);
    }
    setIsDialogOpen(false);
  };

  const isFormLoading = createMutation.isPending || updateMutation.isPending;

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <div>
          <h2 className="text-3xl font-bold tracking-tight">Employees</h2>
          <p className="text-slate-500">Manage all registered employees.</p>
        </div>

        <Button onClick={handleAdd} className="bg-blue-600 hover:bg-blue-700">
          <Plus className="mr-2 h-4 w-4" /> Add Employee
        </Button>
      </div>

      <Card>
        <CardHeader>
          <div className="flex justify-between items-center">
            <CardTitle className="flex items-center gap-2">
              <Users className="h-5 w-5" /> Employee List
            </CardTitle>
            <div className="relative w-64">
              <Search className="absolute left-2 top-2.5 h-4 w-4 text-slate-500" />
              <Input
                placeholder="Search name or NIK..."
                className="pl-8"
                value={search}
                onChange={(e) => {
                  setSearch(e.target.value);
                  setPage(1);
                }}
              />
            </div>
          </div>
        </CardHeader>
        <CardContent>
          {isLoading ? (
            <div className="flex justify-center py-10">
              <Loader2 className="animate-spin" />
            </div>
          ) : (
            <>
              <div className="rounded-md border">
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>ID</TableHead>
                      <TableHead>Name</TableHead>
                      <TableHead>Department</TableHead>
                      <TableHead>Shift</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {data?.data.map((emp) => (
                      <TableRow key={emp.id}>
                        <TableCell className="font-mono">
                          {emp.username}
                        </TableCell>
                        <TableCell className="font-medium">
                          {emp.full_name}
                        </TableCell>
                        <TableCell>{emp.department_name}</TableCell>
                        <TableCell>{emp.shift_name}</TableCell>
                        <TableCell className="text-right">
                          <div className="flex justify-end gap-2">
                            <Button
                              variant="ghost"
                              size="icon"
                              onClick={() => handleEdit(emp)}
                            >
                              <Pencil className="h-4 w-4 text-slate-500" />
                            </Button>
                            <Button
                              variant="ghost"
                              size="icon"
                              onClick={() => handleDeleteClick(emp.id)}
                            >
                              <Trash2 className="h-4 w-4 text-red-500" />
                            </Button>
                          </div>
                        </TableCell>
                      </TableRow>
                    ))}
                    {data?.data.length === 0 && (
                      <TableRow>
                        <TableCell colSpan={5} className="text-center py-8">
                          No employees found.
                        </TableCell>
                      </TableRow>
                    )}
                  </TableBody>
                </Table>
              </div>

              {data?.meta && (
                <PaginationControls
                  meta={{
                    limit: 10,
                    page: data.meta.page,
                    total_page: data.meta.total_page,
                    total_data: data.meta.total_data,
                  }}
                  onPageChange={setPage}
                  isLoading={isLoading}
                />
              )}
            </>
          )}
        </CardContent>
      </Card>

      <EmployeeFormDialog
        open={isDialogOpen}
        onOpenChange={setIsDialogOpen}
        employeeToEdit={selectedEmployee}
        onSubmit={handleFormSubmit}
        isLoading={isFormLoading}
      />

      <AlertDialog
        open={!!employeeToDelete}
        onOpenChange={(open) => !open && setEmployeeToDelete(null)}
      >
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Are you absolutely sure?</AlertDialogTitle>
            <AlertDialogDescription>
              This action cannot be undone. This will permanently delete the
              employee account and their data.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction
              onClick={confirmDelete}
              className="bg-red-600 hover:bg-red-700"
            >
              {deleteMutation.isPending ? "Deleting..." : "Delete"}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  );
}
