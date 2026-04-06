import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/axios";
import { toast } from "sonner";
import type {
  UpsertContractPayload,
  UseContractsParams,
} from "../types";

export const useContracts = (params: UseContractsParams) => {
  return useQuery({
    queryKey: ["contracts", params],
    queryFn: async () => {
      const { data } = await api.get("/contracts", { params });
      return data;
    },
  });
};

export const useContractByEmployee = (employeeId: number) => {
  return useQuery({
    queryKey: ["contracts", "employee", employeeId],
    queryFn: async () => {
      const { data } = await api.get(`/contracts/employee/${employeeId}`);
      return data?.data ?? data;
    },
    enabled: !!employeeId,
  });
};

export const useContractDetail = (id: number | null) => {
  return useQuery({
    queryKey: ["contracts", "detail", id],
    queryFn: async () => {
      const { data } = await api.get(`/contracts/${id}`);
      return data?.data ?? data;
    },
    enabled: !!id,
  });
};

export const useUpsertContract = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (payload: UpsertContractPayload) => {
      const { data } = await api.put("/contracts", payload);
      return data;
    },
    onSuccess: () => {
      toast.success("Contract saved successfully");
      queryClient.invalidateQueries({ queryKey: ["contracts"] });
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || "Failed to save contract");
    },
  });
};

export const useDeleteContract = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (id: number) => {
      const { data } = await api.delete(`/contracts/${id}`);
      return data;
    },
    onSuccess: () => {
      toast.success("Contract deleted successfully");
      queryClient.invalidateQueries({ queryKey: ["contracts"] });
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || "Failed to delete contract");
    },
  });
};
