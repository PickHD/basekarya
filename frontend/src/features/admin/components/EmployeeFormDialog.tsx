import { useEffect } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Loader2 } from "lucide-react";
import type { Employee, CreateEmployeePayload } from "@/features/admin/types";
import {
  useDepartments,
  useShifts,
} from "@/features/admin/hooks/useMasterData";

const formSchema = z.object({
  username: z.string().min(3, "Username minimal atleast 3 characters"),

  full_name: z.string().min(1, "full_name required").trim(),

  nik: z
    .string()
    .length(16, "NIK must be 16 digit")
    .regex(/^\d+$/, "NIK must be a numbers"),

  department_id: z.string().min(1, "Select Dept."),

  shift_id: z.string().min(1, "Select shift."),

  base_salary: z.preprocess(
    (val) => {
      if (!val) return 0;
      if (typeof val === "string") {
        return parseInt(val.replace(/\D/g, ""), 10);
      }
      return val;
    },
    z
      .number()
      .min(0, "Base salary cannot be a negative number")
      .max(1000000000, "Base salary exceed the limit"),
  ),
});

interface EmployeeFormDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  employeeToEdit?: Employee | null;
  onSubmit: (values: CreateEmployeePayload) => void;
  isLoading: boolean;
}

export function EmployeeFormDialog({
  open,
  onOpenChange,
  employeeToEdit,
  onSubmit,
  isLoading,
}: EmployeeFormDialogProps) {
  const isEdit = !!employeeToEdit;

  const { data: departments, isLoading: deptLoading } = useDepartments();
  const { data: shifts, isLoading: shiftLoading } = useShifts();

  const form = useForm<any>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      username: "",
      full_name: "",
      nik: "",
      department_id: "",
      shift_id: "",
      base_salary: 0,
    },
  });

  useEffect(() => {
    if (open) {
      if (employeeToEdit) {
        form.reset({
          username: employeeToEdit.username,
          full_name: employeeToEdit.full_name,
          nik: employeeToEdit.nik,
          department_id: employeeToEdit.department_name === "Umum" ? "1" : "2",
          shift_id: "1",
          base_salary: employeeToEdit.base_salary,
        });
      } else {
        form.reset({
          username: "",
          full_name: "",
          nik: "",
          department_id: "",
          shift_id: "",
          base_salary: 0,
        });
      }
    }
  }, [open, employeeToEdit, form]);

  const formatCurrency = (value: string) => {
    if (!value) return "";
    const number = value.replace(/\D/g, "");
    return new Intl.NumberFormat("id-ID").format(Number(number));
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>
            {isEdit ? "Edit Employee" : "Add New Employee"}
          </DialogTitle>
        </DialogHeader>

        <Form {...form}>
          <form
            onSubmit={form.handleSubmit((values) => {
              const payload: CreateEmployeePayload = {
                username: values.username,
                full_name: values.full_name,
                nik: values.nik,
                department_id: Number(values.department_id),
                shift_id: Number(values.shift_id),
                base_salary: Number(values.base_salary),
              };
              onSubmit(payload);
            })}
            className="space-y-4"
          >
            {!isEdit && (
              <FormField
                control={form.control}
                name="username"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Username</FormLabel>
                    <FormControl>
                      <Input {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
            )}

            <FormField
              control={form.control}
              name="full_name"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Full Name</FormLabel>
                  <FormControl>
                    <Input {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="nik"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>NIK</FormLabel>
                  <FormControl>
                    <Input {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <div className="grid grid-cols-2 gap-4">
              <FormField
                control={form.control}
                name="department_id"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Department</FormLabel>
                    <Select
                      onValueChange={field.onChange}
                      defaultValue={
                        field.value ? String(field.value) : undefined
                      }
                      value={field.value ? String(field.value) : undefined}
                    >
                      <FormControl>
                        <SelectTrigger>
                          <SelectValue
                            placeholder={
                              deptLoading ? "Loading..." : "Select Dept"
                            }
                          />
                        </SelectTrigger>
                      </FormControl>
                      <SelectContent>
                        {departments?.map((dept) => (
                          <SelectItem key={dept.id} value={String(dept.id)}>
                            {dept.name}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="shift_id"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Shift</FormLabel>
                    <Select
                      onValueChange={field.onChange}
                      defaultValue={
                        field.value ? String(field.value) : undefined
                      }
                      value={field.value ? String(field.value) : undefined}
                    >
                      <FormControl>
                        <SelectTrigger>
                          <SelectValue
                            placeholder={
                              shiftLoading ? "Loading..." : "Select Shift"
                            }
                          />
                        </SelectTrigger>
                      </FormControl>
                      <SelectContent>
                        {shifts?.map((shift) => (
                          <SelectItem key={shift.id} value={String(shift.id)}>
                            {shift.name}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>

            <FormField
              control={form.control}
              name="base_salary"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Base Salary (Rp)</FormLabel>
                  <FormControl>
                    <Input
                      placeholder="0"
                      {...field}
                      type="text"
                      value={
                        field.value
                          ? `Rp ${formatCurrency(String(field.value))}`
                          : ""
                      }
                      onChange={(e) => {
                        const rawValue = e.target.value.replace(/\D/g, "");
                        field.onChange(rawValue);
                      }}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <DialogFooter>
              <Button type="submit" disabled={isLoading}>
                {isLoading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                {isEdit ? "Save Changes" : "Create Employee"}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}
