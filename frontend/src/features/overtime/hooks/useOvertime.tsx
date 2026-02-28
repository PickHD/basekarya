import type {
  CreateOvertimePayload,
  Overtime,
  OvertimeFilter,
} from "@/features/overtime/types";
import { api } from "@/lib/axios";
import type { Meta } from "@/types/api";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import type { OvertimeActionPayload } from "../types";
import { toast } from "sonner";

export const useOvertimes = (filter: OvertimeFilter) => {
  return useQuery({
    queryKey: ["overtimes", filter],
    queryFn: async () => {
      const { data } = await api.get<{ data: Overtime[]; meta: Meta }>("/overtimes", {
        params: filter,
      });

      return data;
    },

    placeholderData: (prev) => prev,
  });
};

export const useOvertime = (id: string) => {
  return useQuery({
    queryKey: ["overtime", id],
    queryFn: async () => {
      const { data } = await api.get<{ data: Overtime }>(`/overtimes/${id}`);

      return data.data;
    },
    enabled: !!id,
  });
};

export const useCreateOvertime = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (payload: CreateOvertimePayload) => {
      const { data } = await api.post("/overtimes", payload);
      return data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["overtimes"] });
      toast.success("Pengajuan lembur berhasil dikirim!");
    },
    onError: (error: any) => {
      toast.error(error.response?.data.message || "Gagal Mengajukan lembur");
    },
  });
};

export const useOvertimeAction = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({ id, action, rejection_reason }: OvertimeActionPayload) => {
      const { data } = await api.put(`/overtimes/${id}/action`, {
        action,
        rejection_reason,
      });

      return data;
    },
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: ["overtime", variables.id.toString()],
      });
      queryClient.invalidateQueries({ queryKey: ["overtimes"] });

      let actionMsg = "";
      if (variables.action === "APPROVE") actionMsg = "disetujui";
      else if (variables.action === "REJECT") actionMsg = "ditolak";
      
      toast.success(`Lembur berhasil ${actionMsg}`);
    },
    onError: (error: any) => {
      const errMsg =
        error.response?.data?.error || "Terjadi kesalahan saat memproses aksi";
      toast.error(errMsg);
    },
  });
};

export const useExportOvertimes = () => {
  return useMutation({
    mutationFn: async (params: { status?: string; search?: string }) => {
      const response = await api.get("/overtimes/export", {
        params,
        responseType: "blob",
      });
      return response.data;
    },
    onSuccess: (data) => {
      const url = window.URL.createObjectURL(new Blob([data]));
      const link = document.createElement("a");
      link.href = url;
      link.setAttribute("download", "overtimes.xlsx");
      document.body.appendChild(link);
      link.click();
      link.parentNode?.removeChild(link);
    },
    onError: () => {
      toast.error("Gagal mengunduh data");
    },
  });
};
