import { useState } from "react";
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
import { Textarea } from "@/components/ui/textarea";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Loader2, Plus, Pencil, Trash2 } from "lucide-react";
import {
  useFinanceCategories,
  useCreateFinanceCategory,
  useUpdateFinanceCategory,
  useDeleteFinanceCategory,
} from "@/features/finance/hooks/useFinanceCategory";

interface FinanceCategoryManagerProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

const categorySchema = z.object({
  name: z.string().min(2, "Nama minimal 2 karakter"),
  type: z.string().min(1, "Tipe wajib dipilih"),
  description: z.string().optional(),
});

export function FinanceCategoryManager({ open, onOpenChange }: FinanceCategoryManagerProps) {
  const { data: categories, isLoading } = useFinanceCategories();
  const { mutate: createCat, isPending: isCreating } = useCreateFinanceCategory();
  const { mutate: updateCat, isPending: isUpdating } = useUpdateFinanceCategory();
  const { mutate: deleteCat, isPending: isDeleting } = useDeleteFinanceCategory();

  const [editingId, setEditingId] = useState<number | null>(null);
  const [showForm, setShowForm] = useState(false);

  const form = useForm<any>({
    resolver: zodResolver(categorySchema),
    defaultValues: {
      name: "",
      type: "",
      description: "",
    },
  });

  const resetForm = () => {
    form.reset({ name: "", type: "", description: "" });
    setEditingId(null);
  };

  const handleOpenForm = () => {
    resetForm();
    setShowForm(true);
  };

  const handleCloseForm = () => {
    resetForm();
    setShowForm(false);
  };

  const handleEdit = (cat: any) => {
    setEditingId(cat.id);
    form.reset({
      name: cat.name,
      type: cat.type,
      description: cat.description || "",
    });
    setShowForm(true);
  };

  const handleDelete = (id: number) => {
    if (confirm("Yakin ingin menghapus kategori ini?")) {
      deleteCat(id);
    }
  };

  const onSubmit = (formData: any) => {
    const payload = {
      name: formData.name,
      type: formData.type,
      description: formData.description || "",
    };

    if (editingId) {
      updateCat({ id: editingId, ...payload }, { onSuccess: () => handleCloseForm() });
    } else {
      createCat(payload, { onSuccess: () => handleCloseForm() });
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-3xl max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>Kelola Kategori Keuangan</DialogTitle>
        </DialogHeader>

        <div className="space-y-4">
          {!showForm ? (
            <Button onClick={handleOpenForm} className="w-full sm:w-auto">
              <Plus className="mr-2 h-4 w-4" /> Tambah Kategori
            </Button>
          ) : (
            <div className="border rounded-lg p-4 space-y-4 bg-slate-50">
              <h4 className="font-medium">
                {editingId ? "Edit Kategori" : "Tambah Kategori Baru"}
              </h4>
              <Form {...form}>
                <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
                  <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                    <FormField
                      control={form.control}
                      name="name"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Nama</FormLabel>
                          <FormControl>
                            <Input placeholder="Nama kategori" {...field} />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />

                    <FormField
                      control={form.control}
                      name="type"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Tipe</FormLabel>
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
                  </div>

                  <FormField
                    control={form.control}
                    name="description"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Deskripsi (Opsional)</FormLabel>
                        <FormControl>
                          <Textarea className="resize-none" rows={2} {...field} />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <div className="flex gap-2">
                    <Button
                      type="submit"
                      disabled={isCreating || isUpdating}
                      className="bg-blue-600 hover:bg-blue-700"
                    >
                      {(isCreating || isUpdating) && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                      {editingId ? "Simpan Perubahan" : "Tambah Kategori"}
                    </Button>
                    <Button
                      type="button"
                      variant="outline"
                      onClick={handleCloseForm}
                    >
                      Batal
                    </Button>
                  </div>
                </form>
              </Form>
            </div>
          )}

          {isLoading ? (
            <div className="flex justify-center py-6">
              <Loader2 className="animate-spin h-6 w-6 text-blue-600" />
            </div>
          ) : (
            <>
              <div className="hidden md:block rounded-md border">
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>Nama</TableHead>
                      <TableHead>Tipe</TableHead>
                      <TableHead>Deskripsi</TableHead>
                      <TableHead className="text-right">Aksi</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {(categories || []).map((cat) => (
                      <TableRow key={cat.id}>
                        <TableCell className="font-medium">{cat.name}</TableCell>
                        <TableCell>
                          <span className={`font-medium ${cat.type === "INCOME" ? "text-green-600" : "text-red-600"}`}>
                            {cat.type === "INCOME" ? "Pemasukan" : "Pengeluaran"}
                          </span>
                        </TableCell>
                        <TableCell className="text-slate-500">{cat.description || "-"}</TableCell>
                        <TableCell className="text-right">
                          <div className="flex justify-end gap-1">
                            <Button
                              variant="ghost"
                              size="icon"
                              onClick={() => handleEdit(cat)}
                            >
                              <Pencil className="h-4 w-4 text-blue-500" />
                            </Button>
                            <Button
                              variant="ghost"
                              size="icon"
                              onClick={() => handleDelete(cat.id)}
                              disabled={isDeleting}
                            >
                              <Trash2 className="h-4 w-4 text-red-500" />
                            </Button>
                          </div>
                        </TableCell>
                      </TableRow>
                    ))}
                    {(categories || []).length === 0 && (
                      <TableRow>
                        <TableCell colSpan={4} className="text-center py-6 text-slate-500">
                          Belum ada kategori.
                        </TableCell>
                      </TableRow>
                    )}
                  </TableBody>
                </Table>
              </div>

              <div className="md:hidden space-y-3">
                {(categories || []).map((cat) => (
                  <div key={cat.id} className="flex items-center justify-between rounded-lg border p-3">
                    <div>
                      <p className="font-medium">{cat.name}</p>
                      <span className={`text-xs font-medium ${cat.type === "INCOME" ? "text-green-600" : "text-red-600"}`}>
                        {cat.type === "INCOME" ? "Pemasukan" : "Pengeluaran"}
                      </span>
                    </div>
                    <div className="flex gap-1">
                      <Button variant="ghost" size="icon" onClick={() => handleEdit(cat)}>
                        <Pencil className="h-4 w-4 text-blue-500" />
                      </Button>
                      <Button variant="ghost" size="icon" onClick={() => handleDelete(cat.id)}>
                        <Trash2 className="h-4 w-4 text-red-500" />
                      </Button>
                    </div>
                  </div>
                ))}
              </div>
            </>
          )}
        </div>
      </DialogContent>
    </Dialog>
  );
}
