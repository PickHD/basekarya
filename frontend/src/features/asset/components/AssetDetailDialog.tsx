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
import { useAsset, useUpdateAsset, useAssetCategories } from "@/features/asset/hooks/useAsset";
import type { AssetDetailDialogProps, Asset } from "@/features/asset/types";
import { AssetStatusBadge, AssetConditionBadge } from "./StatusBadge";
import type { AssetCategory } from "@/features/asset/types";

const updateAssetSchema = z.object({
  asset_category_id: z.string().optional(),
  name: z.string().optional(),
  serial_number: z.string().optional(),
  condition: z.string().optional(),
  status: z.string().optional(),
  description: z.string().optional(),
});

function AssetEditForm({
  data,
  categories,
  isUpdating,
  onSubmit,
}: {
  data: Asset;
  categories?: AssetCategory[];
  isUpdating: boolean;
  onSubmit: (formData: any) => void;
}) {
  const form = useForm<any>({
    resolver: zodResolver(updateAssetSchema),
    defaultValues: {
      asset_category_id: data.asset_category_id?.toString() || "",
      name: data.name || "",
      serial_number: data.serial_number || "",
      condition: data.condition || "",
      status: data.status || "",
      description: data.description || "",
    },
  });

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
        <div className="grid md:grid-cols-2 gap-4">
          <FormField
            control={form.control}
            name="name"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Nama Aset</FormLabel>
                <FormControl>
                  <Input {...field} />
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
                  <Input {...field} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />

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
            name="status"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Status</FormLabel>
                <Select onValueChange={field.onChange} value={field.value}>
                  <FormControl>
                    <SelectTrigger>
                      <SelectValue />
                    </SelectTrigger>
                  </FormControl>
                  <SelectContent>
                    <SelectItem value="AVAILABLE">Tersedia</SelectItem>
                    <SelectItem value="ASSIGNED">Digunakan</SelectItem>
                    <SelectItem value="MAINTENANCE">Perbaikan</SelectItem>
                    <SelectItem value="RETIRED">Dihapus</SelectItem>
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
                  <Textarea className="resize-none" rows={2} {...field} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
        </div>

        <DialogFooter>
          <Button type="submit" disabled={isUpdating} className="w-full sm:w-auto">
            {isUpdating && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
            {isUpdating ? "Menyimpan..." : "Simpan Perubahan"}
          </Button>
        </DialogFooter>
      </form>
    </Form>
  );
}

export function AssetDetailDialog({ open, onOpenChange, assetId }: AssetDetailDialogProps) {
  const { data, isLoading } = useAsset(assetId?.toString() || "");
  const { mutate: updateMutate, isPending: isUpdating } = useUpdateAsset();
  const { data: categories } = useAssetCategories();

  const handleSubmit = (formData: any) => {
    if (!data) return;
    const payload: any = {};
    if (formData.asset_category_id) payload.asset_category_id = Number(formData.asset_category_id);
    if (formData.name) payload.name = formData.name;
    if (formData.serial_number) payload.serial_number = formData.serial_number;
    if (formData.condition) payload.condition = formData.condition;
    if (formData.status) payload.status = formData.status;
    if (formData.description) payload.description = formData.description;

    updateMutate(
      { id: data.id, ...payload },
      {
        onSuccess: () => {
          toast.success("Aset berhasil diperbarui!");
          onOpenChange(false);
        },
      }
    );
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-2xl max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle className="flex justify-between items-center pr-8">
            <span>Detail Aset</span>
            {data && (
              <div className="flex gap-2">
                <AssetStatusBadge status={data.status} />
                <AssetConditionBadge condition={data.condition} />
              </div>
            )}
          </DialogTitle>
        </DialogHeader>

        {isLoading ? (
          <div className="flex justify-center py-10">
            <Loader2 className="h-8 w-8 animate-spin text-blue-600" />
          </div>
        ) : data ? (
          <div className="grid gap-6 py-4">
            <div className="bg-slate-50 p-4 rounded-lg border">
              <h3 className="font-bold text-lg text-slate-900">{data.name}</h3>
              {data.serial_number && (
                <p className="text-sm text-slate-500">SN: {data.serial_number}</p>
              )}
              <p className="text-sm text-slate-500">Kategori: {data.category_name || "-"}</p>
              {data.current_employee && (
                <p className="text-sm text-blue-600 mt-1">Saat ini digunakan: {data.current_employee}</p>
              )}
            </div>

            <AssetEditForm
              key={data.id}
              data={data}
              categories={categories}
              isUpdating={isUpdating}
              onSubmit={handleSubmit}
            />
          </div>
        ) : (
          <div className="py-10 text-center text-slate-500">Data tidak ditemukan.</div>
        )}
      </DialogContent>
    </Dialog>
  );
}
