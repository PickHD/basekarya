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
import { Textarea } from "@/components/ui/textarea";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Loader2 } from "lucide-react";
import type {
  CreateReimbursementPayload,
  ReimbursementFormDialogProps,
} from "../types";
import { useCreateReimbursement } from "../hooks/useReimbursement";
import { toast } from "sonner";

const reimbursementSchema = z.object({
  title: z.string().min(3, "Judul minimal 3 karakter"),
  description: z.string().optional(),
  amount: z
    .string()
    .transform((val) => Number(val))
    .refine((val) => val > 0, "Nominal harus lebih dari 0"),
  date: z.string().min(1, "Tanggal wajib diisi"),
  proof_file: z
    .any()
    .refine((files) => files?.length === 1, "Bukti struk wajib diupload")
    .refine((files) => files?.[0]?.size <= 5000000, "Ukuran maksimal 5MB")
    .refine(
      (files) =>
        ["image/jpeg", "image/png", "application/pdf"].includes(
          files?.[0]?.type,
        ),
      "Format harus JPG, PNG, atau PDF",
    ),
});

export function ReimbursementFormDialog({
  open,
  onOpenChange,
}: ReimbursementFormDialogProps) {
  const { mutate, isPending } = useCreateReimbursement();

  const form = useForm<any>({
    resolver: zodResolver(reimbursementSchema),
    defaultValues: {
      title: "",
      amount: "",
      date: "",
      description: "",
      proof_file: undefined,
    },
  });

  useEffect(() => {
    if (open) {
      form.reset({
        title: "",
        amount: "",
        date: "",
        description: "",
        proof_file: undefined,
      });
    }
  }, [open, form]);

  const onSubmit = (data: any) => {
    const payload: CreateReimbursementPayload = {
      title: data.title,
      description: data.description || "",
      amount: data.amount,
      date: data.date,
      proof_file: data.proof_file,
    };

    mutate(payload, {
      onSuccess: () => {
        toast.success("Berhasil menambahkan pengajuan reimbursement!");
        onOpenChange(false);
      },
    });
  };

  const formatCurrency = (value: string) => {
    if (!value) return "";
    const number = value.replace(/\D/g, "");
    return new Intl.NumberFormat("id-ID").format(Number(number));
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[500px]">
        <DialogHeader>
          <DialogTitle>Ajukan Reimbursement</DialogTitle>
        </DialogHeader>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
            <FormField
              control={form.control}
              name="title"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Judul Pengeluaran</FormLabel>
                  <FormControl>
                    <Input
                      placeholder="Contoh: Makan Siang Client"
                      {...field}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <div className="grid grid-cols-2 gap-4">
              <FormField
                control={form.control}
                name="amount"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Nominal (Rp)</FormLabel>
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

              <FormField
                control={form.control}
                name="date"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Tanggal</FormLabel>
                    <FormControl>
                      <Input type="date" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>

            <FormField
              control={form.control}
              name="proof_file"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Bukti Struk</FormLabel>
                  <FormControl>
                    <Input
                      type="file"
                      accept="image/*,.pdf"
                      name={field.name}
                      onBlur={field.onBlur}
                      ref={field.ref}
                      onChange={(event) => {
                        field.onChange(event.target.files);
                      }}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="description"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Keterangan (Opsional)</FormLabel>
                  <FormControl>
                    <Textarea className="resize-none" rows={3} {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <DialogFooter>
              <Button
                type="submit"
                disabled={isPending}
                className="w-full sm:w-auto"
              >
                {isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                {isPending ? "Mengirim..." : "Kirim Pengajuan"}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}
