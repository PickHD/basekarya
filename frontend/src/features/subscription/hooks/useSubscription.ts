import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/axios";
import { toast } from "sonner";
import type { ReviewPayload, UpgradePayload } from "../types";

export const useRequestUpgrade = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (payload: UpgradePayload) => {
      const { data } = await api.post("/subscriptions/upgrade", payload);
      return data;
    },
    onSuccess: () => {
      toast.success("Permintaan upgrade terkirim", {
        description: "Tim kami akan menghubungi Anda untuk konfirmasi pembayaran.",
      });
      queryClient.invalidateQueries({ queryKey: ["company-profile"] });
    },
    onError: (error: any) => {
      const msg = error.response?.data?.message || "Gagal mengirim permintaan upgrade";
      toast.error("Gagal", { description: msg });
    },
  });
};

export const usePendingRequests = () => {
  return useQuery({
    queryKey: ["pending-subscription-requests"],
    queryFn: async () => {
      const { data } = await api.get("/admin/subscriptions/pending");
      return data;
    },
  });
};

export const useAllRequests = () => {
  return useQuery({
    queryKey: ["all-subscription-requests"],
    queryFn: async () => {
      const { data } = await api.get("/admin/subscriptions/requests");
      return data;
    },
  });
};

export const useReviewRequest = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({ id, payload }: { id: number; payload: ReviewPayload }) => {
      const { data } = await api.put(`/admin/subscriptions/${id}/review`, payload);
      return data;
    },
    onSuccess: () => {
      toast.success("Berhasil", { description: "Permintaan berhasil diproses." });
      queryClient.invalidateQueries({ queryKey: ["pending-subscription-requests"] });
      queryClient.invalidateQueries({ queryKey: ["all-subscription-requests"] });
      queryClient.invalidateQueries({ queryKey: ["admin-companies"] });
      queryClient.invalidateQueries({ queryKey: ["admin-dashboard-stats"] });
    },
    onError: (error: any) => {
      const msg = error.response?.data?.message || "Gagal memproses permintaan";
      toast.error("Gagal", { description: msg });
    },
  });
};

export const useCompanies = (search: string = "") => {
  return useQuery({
    queryKey: ["admin-companies", search],
    queryFn: async () => {
      const { data } = await api.get("/admin/subscriptions/companies", {
        params: search ? { search } : undefined,
      });
      return data;
    },
  });
};

export const useCompanyDetail = (id: number) => {
  return useQuery({
    queryKey: ["admin-company-detail", id],
    queryFn: async () => {
      const { data } = await api.get(`/admin/subscriptions/companies/${id}`);
      return data;
    },
    enabled: !!id,
  });
};

export const useUpdateCompanyStatus = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({ id, status }: { id: number; status: string }) => {
      const { data } = await api.put(`/admin/subscriptions/companies/${id}/status`, {
        subscription_status: status,
      });
      return data;
    },
    onSuccess: () => {
      toast.success("Berhasil", { description: "Status perusahaan diperbarui." });
      queryClient.invalidateQueries({ queryKey: ["admin-companies"] });
      queryClient.invalidateQueries({ queryKey: ["admin-company-detail"] });
      queryClient.invalidateQueries({ queryKey: ["admin-dashboard-stats"] });
    },
    onError: (error: any) => {
      const msg = error.response?.data?.message || "Gagal memperbarui status";
      toast.error("Gagal", { description: msg });
    },
  });
};

export const useDashboardStats = () => {
  return useQuery({
    queryKey: ["admin-dashboard-stats"],
    queryFn: async () => {
      const { data } = await api.get("/admin/subscriptions/dashboard");
      return data;
    },
  });
};
