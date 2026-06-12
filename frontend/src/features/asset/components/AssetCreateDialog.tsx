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
import { Loader2 } from "lucide-react";
import { toast } from "sonner";
import { useCreateAsset, useAssetCategories } from "@/features/asset/hooks/useAsset";
import type { AssetFormDialogProps, CreateAssetPayload } from "@/features/asset/types";

const assetSchema = z.object({
  asset_category_id: z.string().min(1, "Kategori harus dipilih"),
  name: z.string().min(1, "Nama aset harus diisi"),
  description: z.string().optional(),
  serial_number: z.string().optional(),
  condition: z.string().optional(),
});

export function AssetCreateDialog({ open, onOpenChange }: AssetFormDialogProps) {
  const { mutate, isPending } = useCreateAsset();
  const { data: categories } = useAssetCategories();

  const form = useForm<any>({
    resolver: zodResolver(assetSchema),
    defaultValues: {
      asset_category_id: "",
      name: "",
      description: "",
      serial_number: "",
      condition: "GOOD",
    },
  });

  useEffect(() => {
    if (open) {
      form.reset({
        asset_category_id: "",
        name: "",
        description: "",
        serial_number: "",
        condition: "GOOD",
      });
    }
  }, [open, form]);

  const onSubmit = (data: any) => {
    const payload: CreateAssetPayload = {
      asset_category_id: Number(data.asset_category_id),
      name: data.name,
      description: data.description || "",
      serial_number: data.serial_number || "",
      condition: data.condition || "GOOD",
    };

    mutate(payload, {
      onSuccess: () => {
        toast.success("Berhasil menambahkan aset!");
        onOpenChange(false);
      },
    });
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[500px]">
        <DialogHeader>
          <DialogTitle>Tambah Aset Baru</DialogTitle>
        </DialogHeader>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
            <FormField
              control={form.control}
              name="asset_category_id"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Kategori</FormLabel>
                  <Select onValueChange={field.onChange} value={field.value}>
                    <FormControl>
                      <SelectTrigger>
                        <SelectValue placeholder="Pilih kategori" />
                      </SelectTrigger>
                    </FormControl>
                    <SelectContent>
                      {categories?.map((cat) => (
                        <SelectItem key={cat.id} value={cat.id.toString()}>
                          {cat.name}
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
              name="name"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Nama Aset</FormLabel>
                  <FormControl>
                    <Input placeholder="MacBook Pro 14" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="serial_number"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Serial Number</FormLabel>
                  <FormControl>
                    <Input placeholder="SN001234" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="condition"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Kondisi</FormLabel>
                  <Select onValueChange={field.onChange} value={field.value}>
                    <FormControl>
                      <SelectTrigger>
                        <SelectValue />
                      </SelectTrigger>
                    </FormControl>
                    <SelectContent>
                      <SelectItem value="GOOD">Baik</SelectItem>
                      <SelectItem value="FAIR">Cukup</SelectItem>
                      <SelectItem value="DAMAGED">Rusak</SelectItem>
                      <SelectItem value="LOST">Hilang</SelectItem>
                    </SelectContent>
                  </Select>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="description"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Deskripsi</FormLabel>
                  <FormControl>
                    <Textarea className="resize-none" rows={3} placeholder="Deskripsi aset..." {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <DialogFooter>
              <Button type="submit" disabled={isPending} className="w-full sm:w-auto">
                {isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                {isPending ? "Menyimpan..." : "Simpan Aset"}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}
