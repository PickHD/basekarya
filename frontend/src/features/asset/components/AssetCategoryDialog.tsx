import { useEffect } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
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
import { Loader2 } from "lucide-react";
import {
  useCreateAssetCategory,
  useUpdateAssetCategory,
  useDeleteAssetCategory,
  useAssetCategories,
} from "@/features/asset/hooks/useAsset";
import type { AssetCategoryFormDialogProps, AssetCategory } from "@/features/asset/types";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Pencil, Trash2 } from "lucide-react";
import { useState } from "react";

const categorySchema = z.object({
  name: z.string().min(1, "Nama kategori harus diisi"),
  description: z.string().optional(),
});

export function AssetCategoryDialog({ open, onOpenChange }: AssetCategoryFormDialogProps) {
  const { data: categories, isLoading } = useAssetCategories();
  const { mutate: createMutate, isPending: isCreating } = useCreateAssetCategory();
  const { mutate: updateMutate, isPending: isUpdating } = useUpdateAssetCategory();
  const { mutate: deleteMutate } = useDeleteAssetCategory();

  const [editingCategory, setEditingCategory] = useState<AssetCategory | null>(null);
  const [isEditing, setIsEditing] = useState(false);

  const form = useForm<any>({
    resolver: zodResolver(categorySchema),
    defaultValues: {
      name: "",
      description: "",
    },
  });

  useEffect(() => {
    if (open) {
      form.reset({ name: "", description: "" });
      setEditingCategory(null);
      setIsEditing(false);
    }
  }, [open, form]);

  const onSubmit = (data: any) => {
    if (isEditing && editingCategory) {
      updateMutate(
        { id: editingCategory.id, name: data.name, description: data.description || "" },
        {
          onSuccess: () => {
            form.reset({ name: "", description: "" });
            setEditingCategory(null);
            setIsEditing(false);
          },
        }
      );
    } else {
      createMutate(
        { name: data.name, description: data.description || "" },
        {
          onSuccess: () => {
            form.reset({ name: "", description: "" });
          },
        }
      );
    }
  };

  const handleEdit = (category: AssetCategory) => {
    setEditingCategory(category);
    setIsEditing(true);
    form.reset({
      name: category.name,
      description: category.description || "",
    });
  };

  const handleDelete = (id: number) => {
    if (confirm("Apakah Anda yakin ingin menghapus kategori ini?")) {
      deleteMutate(id);
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[600px] max-h-[80vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>Kelola Kategori Aset</DialogTitle>
        </DialogHeader>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
            <div className="flex gap-2 items-end">
              <FormField
                control={form.control}
                name="name"
                render={({ field }) => (
                  <FormItem className="flex-1">
                    <FormLabel>Nama Kategori</FormLabel>
                    <FormControl>
                      <Input placeholder="Laptop, Monitor, dll" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <Button type="submit" disabled={isCreating || isUpdating}>
                {(isCreating || isUpdating) && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                {isEditing ? "Update" : "Tambah"}
              </Button>
              {isEditing && (
                <Button
                  type="button"
                  variant="outline"
                  onClick={() => {
                    setIsEditing(false);
                    setEditingCategory(null);
                    form.reset({ name: "", description: "" });
                  }}
                >
                  Batal
                </Button>
              )}
            </div>
          </form>
        </Form>

        <div className="border rounded-md">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Nama</TableHead>
                <TableHead className="text-right">Aksi</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {isLoading ? (
                <TableRow>
                  <TableCell colSpan={2} className="text-center py-4">
                    <Loader2 className="animate-spin h-5 w-5 text-blue-600 mx-auto" />
                  </TableCell>
                </TableRow>
              ) : categories?.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={2} className="text-center py-4 text-slate-500">
                    Belum ada kategori.
                  </TableCell>
                </TableRow>
              ) : (
                categories?.map((cat) => (
                  <TableRow key={cat.id}>
                    <TableCell className="font-medium">{cat.name}</TableCell>
                    <TableCell className="text-right">
                      <Button variant="ghost" size="icon" onClick={() => handleEdit(cat)}>
                        <Pencil className="h-4 w-4 text-slate-500" />
                      </Button>
                      <Button variant="ghost" size="icon" onClick={() => handleDelete(cat.id)}>
                        <Trash2 className="h-4 w-4 text-red-500" />
                      </Button>
                    </TableCell>
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>
        </div>
      </DialogContent>
    </Dialog>
  );
}
