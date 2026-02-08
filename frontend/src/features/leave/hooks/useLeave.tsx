import type { Meta } from "@/features/admin/types";
import type {
  ApplyLeavePayload,
  LeaveActionPayload,
  LeaveRequest,
  LeaveType,
  UseLeavesParams,
} from "@/features/leave/types";
import { api } from "@/lib/axios";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";

export const useLeaves = (params: UseLeavesParams) => {
  return useQuery({
    queryKey: ["leaves", params],
    queryFn: async () => {
      const { data } = await api.get<{ data: LeaveRequest[]; meta: Meta }>(
        "/leaves",
        { params },
      );

      return data;
    },
    placeholderData: (previousData) => previousData,
  });
};

export const useApplyLeave = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (payload: ApplyLeavePayload) => {
      const { data } = await api.post("/leaves/apply", payload);
      return data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["leaves"] });
      toast.success("Pengajuan cuti berhasil dikirim!");
    },
    onError: (error: any) => {
      toast.error(error.response?.data.message || "Gagal Mengajukan cuti");
    },
  });
};

export const useLeaveAction = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({
      id,
      action,
      rejection_reason,
    }: LeaveActionPayload) => {
      const { data } = await api.put(`/leaves/${id}/action`, {
        action,
        rejection_reason,
      });

      return data;
    },
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: ["leave", variables.id.toString()],
      });
      queryClient.invalidateQueries({ queryKey: ["leaves"] });

      toast.success(
        `Pengajuan cuti berhasil di-${variables.action.toLowerCase()}`,
      );
    },
    onError: (error: any) => {
      const errMsg =
        error.response?.data?.error || "Terjadi kesalahan saat memproses aksi";
      toast.error(errMsg);
    },
  });
};

export const useLeaveDetail = (id: string) => {
  return useQuery({
    queryKey: ["leave", id],
    queryFn: async () => {
      const { data } = await api.get<{ data: LeaveRequest }>(`/leaves/${id}`);
      return data.data;
    },
    enabled: !!id,
  });
};

export const useLeaveTypes = () => {
  return useQuery({
    queryKey: ["leave-types"],
    queryFn: async () => {
      const { data } = await api.get<{ data: LeaveType[] }>("/leaves/types");
      return data.data;
    },
    staleTime: 1000 * 60 * 60,
  });
};
