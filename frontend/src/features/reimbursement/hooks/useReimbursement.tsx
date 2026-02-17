import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import type {
  CreateReimbursementPayload,
  ReimbursementActionPayload,
  ReimbursementFilter,
  Reimbursement,
} from "../types";
import { api } from "@/lib/axios";
import { useNavigate } from "react-router-dom";
import { toast } from "sonner";
import type { Meta } from "@/types/api";

export const useReimbursements = (filter: ReimbursementFilter) => {
  return useQuery({
    queryKey: ["reimbursements", filter],
    queryFn: async () => {
      const { data } = await api.get<{ data: Reimbursement[]; meta: Meta }>(
        "/reimbursements",
        { params: filter },
      );

      return data;
    },
    placeholderData: (prev) => prev,
  });
};

export const useReimbursement = (id: string) => {
  return useQuery({
    queryKey: ["reimbursement", id],
    queryFn: async () => {
      const { data } = await api.get<{ data: Reimbursement }>(
        `/reimbursements/${id}`,
      );

      return data.data;
    },
    enabled: !!id,
  });
};

export const useCreateReimbursement = () => {
  const queryClient = useQueryClient();
  const navigate = useNavigate();

  return useMutation({
    mutationFn: async (payload: CreateReimbursementPayload) => {
      const formData = new FormData();
      formData.append("title", payload.title);
      formData.append("description", payload.description);
      formData.append("amount", payload.amount.toString());
      formData.append("date", payload.date);
      if (payload.proof_file && payload.proof_file.length > 0) {
        formData.append("file", payload.proof_file[0]);
      }

      const { data } = await api.post("/reimbursements", formData);
      return data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["reimbursements"] });
      navigate("/reimbursement");
    },
  });
};

export const useReimbursementAction = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({
      id,
      action,
      rejection_reason,
    }: ReimbursementActionPayload) => {
      const { data } = await api.put(`/reimbursements/${id}/action`, {
        action,
        rejection_reason,
      });
      return data;
    },
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: ["reimbursement", variables.id.toString()],
      });
      queryClient.invalidateQueries({ queryKey: ["reimbursements"] });

      toast.success(
        `Reimbursement berhasil di-${variables.action.toLowerCase()}`,
      );
    },
    onError: (error: any) => {
      const errMsg =
        error.response?.data?.error || "Terjadi kesalahan saat memproses aksi";
      toast.error(errMsg);
    },
  });
};
