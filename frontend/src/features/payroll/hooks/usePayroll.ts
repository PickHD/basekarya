import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import type { GeneratePayrollPayload, Payroll, PayrollFilter } from "../types";
import { api } from "@/lib/axios";
import { toast } from "sonner";
import type { Meta } from "@/types/api";

export const usePayrolls = (filter: PayrollFilter) => {
  return useQuery({
    queryKey: ["payrolls", filter],
    queryFn: async () => {
      const { data } = await api.get<{ data: Payroll[]; meta: Meta }>(
        "/admin/payrolls",
        { params: filter },
      );

      return data;
    },
    placeholderData: (prev) => prev,
  });
};

export const usePayroll = (id: string | number | null) => {
  return useQuery({
    queryKey: ["payroll", id],
    queryFn: async () => {
      const { data } = await api.get<{ data: Payroll }>(
        `/admin/payrolls/${id}`,
      );

      return data.data;
    },
    enabled: !!id,
  });
};

export const useGeneratePayroll = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (payload: GeneratePayrollPayload) => {
      const { data } = await api.post("/admin/payrolls/generate", payload);
      return data;
    },
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: ["payrolls"] });
      toast.success(
        `Berhasil generate ${data.data?.success_count || " "} slip gaji!`,
      );
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || "Gagal generate payroll");
    },
  });
};

export const downloadPayslip = async (id: number, filename: string) => {
  try {
    const response = await api.get(`/admin/payrolls/${id}/download`, {
      responseType: "blob",
    });

    const url = window.URL.createObjectURL(new Blob([response.data]));

    const link = document.createElement("a");
    link.href = url;
    link.setAttribute("download", filename);
    document.body.appendChild(link);
    link.click();

    link.remove();
    window.URL.revokeObjectURL(url);
  } catch (error) {
    console.error("Download failed", error);
    toast.error("Gagal mendownload slip gaji");
  }
};

export const useMarkAsPaid = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (id: number) => {
      const { data } = await api.put(`/admin/payrolls/${id}/status`, id);
      return data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["payrolls"] });
      queryClient.invalidateQueries({ queryKey: ["payroll"] });
      toast.success("Status berhasil diubah menjadi PAID");
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || "Gagal update status");
    },
  });
};
