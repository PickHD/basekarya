import { useState, useEffect } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Textarea } from "@/components/ui/textarea";
import { useUpsertContract, useContractByEmployee } from "../hooks/useContract";
import { useAllEmployees } from "@/features/admin/hooks/useAdmin";
import { format } from "date-fns";
import { Loader2 } from "lucide-react";

const schema = z.object({
  employee_id: z.string().min(1, "Employee is required"),
  contract_type: z.enum(["PKWT", "PKWTT"]),
  contract_number: z.string().optional(),
  start_date: z.string().min(1, "Start date is required"),
  end_date: z.string().optional(),
  notes: z.string().optional(),
});

type FormValues = z.infer<typeof schema>;

interface Props {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  initialEmployeeId?: number | null;
}

export function ContractFormDialog({ open, onOpenChange, initialEmployeeId }: Props) {
  const [selectedEmployeeId, setSelectedEmployeeId] = useState<number | null>(initialEmployeeId ?? null);
  const [attachment, setAttachment] = useState<File | null>(null);
  
  const { data: employeeData, isLoading: isLoadingEmployees } = useAllEmployees(1, "");
  const { data: existingContract, isLoading: isLoadingContract } = useContractByEmployee(selectedEmployeeId!);
  const { mutate, isPending } = useUpsertContract();

  const form = useForm<FormValues>({
    resolver: zodResolver(schema),
    defaultValues: {
      employee_id: "",
      contract_type: "PKWT",
      contract_number: "",
      start_date: format(new Date(), "yyyy-MM-dd"),
      end_date: "",
      notes: "",
    },
  });

  const contractType = form.watch("contract_type");

  useEffect(() => {
    if (existingContract) {
      form.reset({
        employee_id: existingContract.employee_id.toString(),
        contract_type: existingContract.contract_type as any,
        contract_number: existingContract.contract_number || "",
        start_date: existingContract.start_date ? format(new Date(existingContract.start_date), "yyyy-MM-dd") : "",
        end_date: existingContract.end_date ? format(new Date(existingContract.end_date), "yyyy-MM-dd") : "",
        notes: existingContract.notes || "",
      });
    } else if (selectedEmployeeId) {
      form.setValue("contract_type", "PKWT");
      form.setValue("contract_number", "");
      form.setValue("start_date", format(new Date(), "yyyy-MM-dd"));
      form.setValue("end_date", "");
      form.setValue("notes", "");
    }
  }, [existingContract, selectedEmployeeId, form]);

  // Sync selectedEmployeeId when the dialog opens (handles edit flow)
  useEffect(() => {
    if (open) {
      setSelectedEmployeeId(initialEmployeeId ?? null);
      setAttachment(null);
      if (!initialEmployeeId) {
        form.reset({
          employee_id: "",
          contract_type: "PKWT",
          contract_number: "",
          start_date: format(new Date(), "yyyy-MM-dd"),
          end_date: "",
          notes: "",
        });
      } else {
        // Pre-fill the employee_id field so the Select shows the right value
        form.setValue("employee_id", initialEmployeeId.toString());
      }
    }
  }, [open, initialEmployeeId]); // eslint-disable-line react-hooks/exhaustive-deps

  // Reset on close
  useEffect(() => {
    if (!open) {
      form.reset();
      setSelectedEmployeeId(null);
      setAttachment(null);
    }
  }, [open, form]);

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      setAttachment(e.target.files[0]);
    }
  };

  const onSubmit = async (values: FormValues) => {
    let base64 = "";
    if (attachment) {
      const reader = new FileReader();
      const readAsDataURL = new Promise<string>((resolve, reject) => {
        reader.onload = () => resolve(reader.result as string);
        reader.onerror = error => reject(error);
      });
      reader.readAsDataURL(attachment);
      base64 = await readAsDataURL;
    }

    let end_date = values.end_date;
    if (values.contract_type === "PKWTT") {
        end_date = "";
    }

    mutate({
      employee_id: Number(values.employee_id),
      contract_type: values.contract_type,
      contract_number: values.contract_number,
      start_date: values.start_date,
      end_date: end_date,
      notes: values.notes,
      attachment_base64: base64,
    }, {
      onSuccess: () => {
        onOpenChange(false);
      }
    });
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-xl max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>{initialEmployeeId ? "Edit Kontrak" : "Buat Kontrak"}</DialogTitle>
        </DialogHeader>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4 py-4">
            <FormField
              control={form.control}
              name="employee_id"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Employee</FormLabel>
                  <Select 
                    onValueChange={(val) => {
                      field.onChange(val);
                      setSelectedEmployeeId(Number(val));
                    }} 
                    value={field.value}
                    disabled={!!initialEmployeeId}
                  >
                    <FormControl>
                      <SelectTrigger>
                        <SelectValue placeholder="Pilih Karyawan" />
                      </SelectTrigger>
                    </FormControl>
                    <SelectContent>
                      {isLoadingEmployees ? (
                        <div className="p-2 text-sm text-muted-foreground flex items-center justify-center">
                           <Loader2 className="h-4 w-4 animate-spin mr-2"/> Loading...
                        </div>
                      ) : employeeData?.data?.map((emp: any) => (
                        <SelectItem key={emp.id} value={emp.id.toString()}>
                          {emp.full_name} ({emp.nik})
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                  <FormMessage />
                </FormItem>
              )}
            />

            {isLoadingContract && selectedEmployeeId ? (
                <div className="py-4 text-center text-sm text-slate-500">
                    <Loader2 className="h-5 w-5 animate-spin mx-auto mb-2 text-slate-400" />
                    Loading contract data...
                </div>
            ) : (
                <>
                    <div className="grid grid-cols-2 gap-4">
                        <FormField
                            control={form.control}
                            name="contract_type"
                            render={({ field }) => (
                            <FormItem>
                                <FormLabel>Contract Type</FormLabel>
                                <Select onValueChange={field.onChange} value={field.value}>
                                <FormControl>
                                    <SelectTrigger>
                                    <SelectValue placeholder="Select Type" />
                                    </SelectTrigger>
                                </FormControl>
                                <SelectContent>
                                    <SelectItem value="PKWT">PKWT (Kontrak)</SelectItem>
                                    <SelectItem value="PKWTT">PKWTT (Tetap)</SelectItem>
                                </SelectContent>
                                </Select>
                                <FormMessage />
                            </FormItem>
                            )}
                        />

                        <FormField
                            control={form.control}
                            name="contract_number"
                            render={({ field }) => (
                            <FormItem>
                                <FormLabel>Contract Number</FormLabel>
                                <FormControl>
                                <Input placeholder="Contract Number" {...field} />
                                </FormControl>
                                <FormMessage />
                            </FormItem>
                            )}
                        />
                    </div>

                    <div className="grid grid-cols-2 gap-4">
                        <FormField
                            control={form.control}
                            name="start_date"
                            render={({ field }) => (
                            <FormItem>
                                <FormLabel>Start Date</FormLabel>
                                <FormControl>
                                <Input type="date" {...field} />
                                </FormControl>
                                <FormMessage />
                            </FormItem>
                            )}
                        />

                        {contractType === "PKWT" && (
                            <FormField
                                control={form.control}
                                name="end_date"
                                render={({ field }) => (
                                <FormItem>
                                    <FormLabel>Tanggal Berakhir</FormLabel>
                                    <FormControl>
                                    <Input type="date" {...field} />
                                    </FormControl>
                                    <FormMessage />
                                </FormItem>
                                )}
                            />
                        )}
                    </div>

                    <FormField
                        control={form.control}
                        name="notes"
                        render={({ field }) => (
                        <FormItem>
                            <FormLabel>Notes</FormLabel>
                            <FormControl>
                            <Textarea placeholder="Notes..." {...field} />
                            </FormControl>
                            <FormMessage />
                        </FormItem>
                        )}
                    />

                    <div className="space-y-2">
                        <label className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">
                            Upload Document (Optional, Image)
                        </label>
                        <Input type="file" accept="image/*" onChange={handleFileChange} />
                        {existingContract?.attachment_url && !attachment && (
                            <p className="text-xs text-muted-foreground mt-1">
                            Document already exists. Uploading a new one will replace the old document.
                            </p>
                        )}
                    </div>
                </>
            )}

            <div className="flex justify-end gap-2 pt-4 border-t">
              <Button
                type="button"
                variant="outline"
                onClick={() => onOpenChange(false)}
                disabled={isPending}
              >
                Cancel
              </Button>
              <Button type="submit" disabled={isPending || (!selectedEmployeeId)}>
                {isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                Save Contract
              </Button>
            </div>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}
