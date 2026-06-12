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
import { Textarea } from "@/components/ui/textarea";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Loader2 } from "lucide-react";
import { toast } from "sonner";
import { useCreateAssetAssignment, useAssets } from "@/features/asset/hooks/useAsset";
import type { AssetAssignmentFormDialogProps, CreateAssetAssignmentPayload } from "@/features/asset/types";

const assignmentSchema = z.object({
  asset_id: z.string().min(1, "Aset harus dipilih"),
  purpose: z.string().min(3, "Tujuan harus diisi"),
  expected_return_date: z.string().optional(),
});

export function AssetAssignmentCreateDialog({ open, onOpenChange }: AssetAssignmentFormDialogProps) {
  const { mutate, isPending } = useCreateAssetAssignment();
  const { data: assetsData } = useAssets({
    status: "AVAILABLE",
    page: 1,
    limit: 100,
  });

  const availableAssets = assetsData?.data || [];

  const form = useForm<any>({
    resolver: zodResolver(assignmentSchema),
    defaultValues: {
      asset_id: "",
      purpose: "",
      expected_return_date: "",
    },
  });

  useEffect(() => {
    if (open) {
      form.reset({
        asset_id: "",
        purpose: "",
        expected_return_date: "",
      });
    }
  }, [open, form]);

  const onSubmit = (data: any) => {
    const payload: CreateAssetAssignmentPayload = {
      asset_id: Number(data.asset_id),
      purpose: data.purpose,
      expected_return_date: data.expected_return_date || "",
    };

    mutate(payload, {
      onSuccess: () => {
        toast.success("Berhasil mengajukan permintaan aset!");
        onOpenChange(false);
      },
    });
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[500px]">
        <DialogHeader>
          <DialogTitle>Ajukan Permintaan Aset</DialogTitle>
        </DialogHeader>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
            <FormField
              control={form.control}
              name="asset_id"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Aset</FormLabel>
                  <Select onValueChange={field.onChange} value={field.value}>
                    <FormControl>
                      <SelectTrigger>
                        <SelectValue placeholder="Pilih aset yang tersedia" />
                      </SelectTrigger>
                    </FormControl>
                    <SelectContent>
                      {availableAssets.map((asset) => (
                        <SelectItem key={asset.id} value={asset.id.toString()}>
                          {asset.name} {asset.serial_number ? `(SN: ${asset.serial_number})` : ""}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                  {availableAssets.length === 0 && (
                    <p className="text-xs text-amber-600 mt-1">Tidak ada aset tersedia untuk dipinjam.</p>
                  )}
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="purpose"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Tujuan Penggunaan</FormLabel>
                  <FormControl>
                    <Textarea className="resize-none" rows={3} placeholder="Alasan membutuhkan aset ini..." {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="expected_return_date"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Estimasi Tanggal Kembali</FormLabel>
                  <FormControl>
                    <Input type="date" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <DialogFooter>
              <Button type="submit" disabled={isPending} className="w-full sm:w-auto">
                {isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                {isPending ? "Mengirim..." : "Kirim Permintaan"}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}
