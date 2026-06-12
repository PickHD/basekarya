import type {
  CreateAssetPayload,
  Asset,
  AssetCategory,
  AssetAssignment,
  AssetFilter,
  AssetAssignmentFilter,
  CreateAssetCategoryPayload,
  UpdateAssetCategoryPayload,
  UpdateAssetPayload,
  CreateAssetAssignmentPayload,
  AssetActionPayload,
} from "@/features/asset/types";
import { api } from "@/lib/axios";
import type { Meta } from "@/types/api";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";

export const useAssetCategories = () => {
  return useQuery({
    queryKey: ["assetCategories"],
    queryFn: async () => {
      const { data } = await api.get<{ data: AssetCategory[] }>("/assets/categories", {
        params: { page: 1, limit: 100 },
      });
      return data.data;
    },
  });
};

export const useCreateAssetCategory = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (payload: CreateAssetCategoryPayload) => {
      const { data } = await api.post("/assets/categories", payload);
      return data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["assetCategories"] });
      toast.success("Kategori aset berhasil dibuat!");
    },
    onError: (error: any) => {
      toast.error(error.response?.data.message || "Gagal membuat kategori aset");
    },
  });
};

export const useUpdateAssetCategory = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({ id, ...payload }: UpdateAssetCategoryPayload & { id: number }) => {
      const { data } = await api.put(`/assets/categories/${id}`, payload);
      return data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["assetCategories"] });
      toast.success("Kategori aset berhasil diperbarui!");
    },
    onError: (error: any) => {
      toast.error(error.response?.data.message || "Gagal memperbarui kategori aset");
    },
  });
};

export const useDeleteAssetCategory = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (id: number) => {
      const { data } = await api.delete(`/assets/categories/${id}`);
      return data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["assetCategories"] });
      toast.success("Kategori aset berhasil dihapus!");
    },
    onError: (error: any) => {
      toast.error(error.response?.data.message || "Gagal menghapus kategori aset");
    },
  });
};

export const useAssets = (filter: AssetFilter) => {
  return useQuery({
    queryKey: ["assets", filter],
    queryFn: async () => {
      const { data } = await api.get<{ data: Asset[]; meta: Meta }>("/assets", {
        params: filter,
      });
      return data;
    },
    placeholderData: (prev) => prev,
  });
};

export const useAsset = (id: string) => {
  return useQuery({
    queryKey: ["asset", id],
    queryFn: async () => {
      const { data } = await api.get<{ data: Asset }>(`/assets/${id}`);
      return data.data;
    },
    enabled: !!id,
  });
};

export const useCreateAsset = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (payload: CreateAssetPayload) => {
      const { data } = await api.post("/assets", payload);
      return data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["assets"] });
      toast.success("Aset berhasil dibuat!");
    },
    onError: (error: any) => {
      toast.error(error.response?.data.message || "Gagal membuat aset");
    },
  });
};

export const useUpdateAsset = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({ id, ...payload }: UpdateAssetPayload & { id: number }) => {
      const { data } = await api.put(`/assets/${id}`, payload);
      return data;
    },
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ["assets"] });
      queryClient.invalidateQueries({ queryKey: ["asset", variables.id.toString()] });
      toast.success("Aset berhasil diperbarui!");
    },
    onError: (error: any) => {
      toast.error(error.response?.data.message || "Gagal memperbarui aset");
    },
  });
};

export const useDeleteAsset = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (id: number) => {
      const { data } = await api.delete(`/assets/${id}`);
      return data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["assets"] });
      toast.success("Aset berhasil dihapus!");
    },
    onError: (error: any) => {
      toast.error(error.response?.data.message || "Gagal menghapus aset");
    },
  });
};

export const useAssetAssignments = (filter: AssetAssignmentFilter) => {
  return useQuery({
    queryKey: ["assetAssignments", filter],
    queryFn: async () => {
      const { data } = await api.get<{ data: AssetAssignment[]; meta: Meta }>("/assets/assignments", {
        params: filter,
      });
      return data;
    },
    placeholderData: (prev) => prev,
  });
};

export const useAssetAssignment = (id: string) => {
  return useQuery({
    queryKey: ["assetAssignment", id],
    queryFn: async () => {
      const { data } = await api.get<{ data: AssetAssignment }>(`/assets/assignments/${id}`);
      return data.data;
    },
    enabled: !!id,
  });
};

export const useCreateAssetAssignment = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (payload: CreateAssetAssignmentPayload) => {
      const { data } = await api.post("/assets/assignments", payload);
      return data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["assetAssignments"] });
      queryClient.invalidateQueries({ queryKey: ["assets"] });
      toast.success("Permintaan aset berhasil dikirim!");
    },
    onError: (error: any) => {
      toast.error(error.response?.data.message || "Gagal mengajukan permintaan aset");
    },
  });
};

export const useAssetAssignmentAction = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({ id, action, rejection_reason }: AssetActionPayload) => {
      const { data } = await api.put(`/assets/assignments/${id}/action`, {
        action,
        rejection_reason,
      });
      return data;
    },
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ["assetAssignment", variables.id.toString()] });
      queryClient.invalidateQueries({ queryKey: ["assetAssignments"] });
      queryClient.invalidateQueries({ queryKey: ["assets"] });
      toast.success(`Permintaan aset berhasil di-${variables.action.toLowerCase()}`);
    },
    onError: (error: any) => {
      const errMsg = error.response?.data?.error || "Terjadi kesalahan saat memproses aksi";
      toast.error(errMsg);
    },
  });
};

export const useReturnAssetAssignment = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (id: number) => {
      const { data } = await api.put(`/assets/assignments/${id}/return`);
      return data;
    },
    onSuccess: (_, id) => {
      queryClient.invalidateQueries({ queryKey: ["assetAssignment", id.toString()] });
      queryClient.invalidateQueries({ queryKey: ["assetAssignments"] });
      queryClient.invalidateQueries({ queryKey: ["assets"] });
      toast.success("Aset berhasil dikembalikan!");
    },
    onError: (error: any) => {
      toast.error(error.response?.data.message || "Gagal mengembalikan aset");
    },
  });
};

export const useExportAssets = () => {
  return useMutation({
    mutationFn: async (params: { status?: string; condition?: string; category_id?: number }) => {
      const response = await api.get("/assets/export", {
        params,
        responseType: "blob",
      });
      return response.data;
    },
    onSuccess: (data) => {
      const url = window.URL.createObjectURL(new Blob([data]));
      const link = document.createElement("a");
      link.href = url;
      link.setAttribute("download", "assets.xlsx");
      document.body.appendChild(link);
      link.click();
      link.parentNode?.removeChild(link);
    },
    onError: () => {
      toast.error("Gagal mengunduh data");
    },
  });
};
