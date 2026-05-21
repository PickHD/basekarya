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
import { toast } from "sonner";
import { useCreateFinanceTransaction } from "@/features/finance/hooks/useFinance";
import { useFinanceCategories } from "@/features/finance/hooks/useFinanceCategory";
import type { FinanceFormDialogProps } from "@/features/finance/types";

const transactionSchema = z.object({
  finance_category_id: z.string().min(1, "Kategori wajib dipilih"),
  type: z.string().min(1, "Tipe wajib dipilih"),
  amount: z
    .string()
    .transform((val) => Number(val))
    .refine((val) => val > 0, "Jumlah harus lebih dari 0"),
  description: z.string().optional(),
  transaction_date: z.string().min(1, "Tanggal transaksi wajib diisi"),
  reference_number: z.string().optional(),
});

export function FinanceTransactionCreateDialog({ open, onOpenChange }: FinanceFormDialogProps) {
  const { mutate, isPending } = useCreateFinanceTransaction();
  const { data: categories } = useFinanceCategories();

  const form = useForm<any>({
    resolver: zodResolver(transactionSchema),
    defaultValues: {
      finance_category_id: "",
      type: "",
      amount: "",
      description: "",
      transaction_date: "",
      reference_number: "",
    },
  });

  const watchType = form.watch("type");

  useEffect(() => {
    if (open) {
      form.reset({
        finance_category_id: "",
        type: "",
        amount: "",
        description: "",
        transaction_date: "",
        reference_number: "",
      });
    }
  }, [open, form]);

  const filteredCategories = (categories || []).filter((cat) =>
    watchType ? cat.type === watchType : true
  );

  const onSubmit = (formData: any) => {
    const payload = {
      finance_category_id: Number(formData.finance_category_id),
      type: formData.type,
      amount: formData.amount,
      description: formData.description || "",
      transaction_date: formData.transaction_date,
      reference_number: formData.reference_number || "",
    };

    mutate(payload, {
      onSuccess: () => {
        toast.success("Transaksi keuangan berhasil dibuat!");
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
          <DialogTitle>Buat Transaksi Keuangan</DialogTitle>
        </DialogHeader>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
            <FormField
              control={form.control}
              name="type"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Tipe Transaksi</FormLabel>
                  <FormControl>
                    <select
                      className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2"
                      {...field}
                    >
                      <option value="">Pilih Tipe</option>
                      <option value="INCOME">Pemasukan</option>
                      <option value="EXPENSE">Pengeluaran</option>
                    </select>
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="finance_category_id"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Kategori</FormLabel>
                  <FormControl>
                    <select
                      className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2"
                      {...field}
                    >
                      <option value="">Pilih Kategori</option>
                      {filteredCategories.map((cat) => (
                        <option key={cat.id} value={cat.id}>
                          {cat.name}
                        </option>
                      ))}
                    </select>
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="amount"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Jumlah (Rp)</FormLabel>
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
              name="transaction_date"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Tanggal Transaksi</FormLabel>
                  <FormControl>
                    <Input type="date" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="reference_number"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>No. Referensi (Opsional)</FormLabel>
                  <FormControl>
                    <Input placeholder="INV-001" {...field} />
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
                  <FormLabel>Keterangan</FormLabel>
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
                {isPending ? "Menyimpan..." : "Buat Transaksi"}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}
