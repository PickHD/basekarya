import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/axios";
import { toast } from "sonner";
import type { CreateApplicantPayload, UpdateApplicantStagePayload } from "../types";

export const useApplicants = (requisitionId: number | null) => {
  return useQuery({
    queryKey: ["applicants", "kanban", requisitionId],
    queryFn: async () => {
      const { data } = await api.get(
        `/recruitments/requisitions/${requisitionId}/applicants`
      );
      return data?.data ?? data;
    },
    enabled: !!requisitionId,
  });
};

export const useApplicantDetail = (id: number | null) => {
  return useQuery({
    queryKey: ["applicants", "detail", id],
    queryFn: async () => {
      const { data } = await api.get(`/recruitments/applicants/${id}`);
      return data?.data ?? data;
    },
    enabled: !!id,
  });
};

export const useAddApplicant = (requisitionId: number) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (payload: CreateApplicantPayload) => {
      const { data } = await api.post(
        `/recruitments/requisitions/${requisitionId}/applicants`,
        payload
      );
      return data;
    },
    onSuccess: () => {
      toast.success("Applicant added successfully");
      queryClient.invalidateQueries({ queryKey: ["applicants", "kanban", requisitionId] });
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || "Failed to add applicant");
    },
  });
};

export const useUpdateApplicantStage = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async ({
      id,
      payload,
    }: {
      id: number;
      payload: UpdateApplicantStagePayload;
    }) => {
      const { data } = await api.put(`/recruitments/applicants/${id}/stage`, payload);
      return data;
    },
    onSuccess: () => {
      toast.success("Stage updated");
      queryClient.invalidateQueries({ queryKey: ["applicants"] });
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || "Failed to update stage");
    },
  });
};
