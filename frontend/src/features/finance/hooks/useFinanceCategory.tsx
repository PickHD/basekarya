import { api } from "@/lib/axios";
import type { FinanceCategory, CategoryPayload } from "@/features/finance/types";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";

export const useFinanceCategories = (type?: string) => {
  return useQuery({
    queryKey: ["finance-categories", type],
    queryFn: async () => {
      const { data } = await api.get<{ data: FinanceCategory[] }>("/finances/categories", {
        params: type ? { type } : {},
      });
      return data.data;
    },
    staleTime: 1000 * 60 * 60,
  });
};

export const useCreateFinanceCategory = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (payload: CategoryPayload) => {
      const { data } = await api.post("/finances/categories", payload);
      return data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["finance-categories"] });
      toast.success("Kategori berhasil ditambahkan!");
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || "Gagal menambahkan kategori");
    },
  });
};

export const useUpdateFinanceCategory = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({ id, ...payload }: CategoryPayload & { id: number }) => {
      const { data } = await api.put(`/finances/categories/${id}`, payload);
      return data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["finance-categories"] });
      toast.success("Kategori berhasil diperbarui!");
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || "Gagal memperbarui kategori");
    },
  });
};

export const useDeleteFinanceCategory = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (id: number) => {
      const { data } = await api.delete(`/finances/categories/${id}`);
      return data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["finance-categories"] });
      toast.success("Kategori berhasil dihapus!");
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || "Gagal menghapus kategori");
    },
  });
};
