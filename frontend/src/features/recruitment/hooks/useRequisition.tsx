import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/axios";
import { toast } from "sonner";
import type {
  CreateRequisitionPayload,
  RequisitionActionPayload,
  UseRequisitionsParams,
} from "../types";

export const useRequisitions = (params: UseRequisitionsParams = {}) => {
  return useQuery({
    queryKey: ["requisitions", params],
    queryFn: async () => {
      const { data } = await api.get("/recruitments/requisitions", { params });
      return data;
    },
  });
};

export const useRequisitionDetail = (id: number | null) => {
  return useQuery({
    queryKey: ["requisitions", "detail", id],
    queryFn: async () => {
      const { data } = await api.get(`/recruitments/requisitions/${id}`);
      return data?.data ?? data;
    },
    enabled: !!id,
  });
};

export const useCreateRequisition = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (payload: CreateRequisitionPayload) => {
      const { data } = await api.post("/recruitments/requisitions", payload);
      return data;
    },
    onSuccess: () => {
      toast.success("Requisition created successfully");
      queryClient.invalidateQueries({ queryKey: ["requisitions"] });
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || "Failed to create requisition");
    },
  });
};

export const useSubmitRequisition = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (id: number) => {
      const { data } = await api.put(`/recruitments/requisitions/${id}/submit`);
      return data;
    },
    onSuccess: () => {
      toast.success("Requisition submitted for approval");
      queryClient.invalidateQueries({ queryKey: ["requisitions"] });
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || "Failed to submit requisition");
    },
  });
};

export const useRequisitionAction = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async ({ id, payload }: { id: number; payload: RequisitionActionPayload }) => {
      const { data } = await api.put(`/recruitments/requisitions/${id}/action`, payload);
      return data;
    },
    onSuccess: () => {
      toast.success("Action processed successfully");
      queryClient.invalidateQueries({ queryKey: ["requisitions"] });
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || "Failed to process action");
    },
  });
};

export const useCloseRequisition = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (id: number) => {
      const { data } = await api.put(`/recruitments/requisitions/${id}/close`);
      return data;
    },
    onSuccess: () => {
      toast.success("Requisition closed");
      queryClient.invalidateQueries({ queryKey: ["requisitions"] });
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || "Failed to close requisition");
    },
  });
};

export const useDeleteRequisition = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (id: number) => {
      const { data } = await api.delete(`/recruitments/requisitions/${id}`);
      return data;
    },
    onSuccess: () => {
      toast.success("Requisition deleted");
      queryClient.invalidateQueries({ queryKey: ["requisitions"] });
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || "Failed to delete requisition");
    },
  });
};
