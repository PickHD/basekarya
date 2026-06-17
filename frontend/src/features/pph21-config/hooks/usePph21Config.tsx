import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import { api } from "@/lib/axios";
import type { TERBracket, TERBracketPayload, PTKPConfig, PTKPConfigPayload } from "../types";
import type { ApiResponse } from "@/types/api";

export function useTerBrackets(category: string) {
  return useQuery({
    queryKey: ["ter-brackets", category],
    queryFn: async () => {
      const { data } = await api.get<ApiResponse<TERBracket[]>>("/admin/tax/ter-brackets", {
        params: { category },
      });
      return data.data ?? [];
    },
  });
}

export function useTerBracketMutations() {
  const queryClient = useQueryClient();

  const createMutation = useMutation({
    mutationFn: async (payload: TERBracketPayload) => {
      return await api.post("/admin/tax/ter-brackets", payload);
    },
    onSuccess: () => {
      toast.success("TER bracket created");
      queryClient.invalidateQueries({ queryKey: ["ter-brackets"], type: "active" });
    },
    onError: (error: any) => {
      const msg = error?.response?.data?.error?.message || error?.response?.data?.message || "Failed to create TER bracket";
      toast.error("Create Failed", { description: msg });
    },
  });

  const updateMutation = useMutation({
    mutationFn: async ({ id, ...payload }: TERBracketPayload & { id: number }) => {
      return await api.put(`/admin/tax/ter-brackets/${id}`, payload);
    },
    onSuccess: () => {
      toast.success("TER bracket updated");
      queryClient.invalidateQueries({ queryKey: ["ter-brackets"], type: "active" });
    },
    onError: (error: any) => {
      const msg = error?.response?.data?.error?.message || error?.response?.data?.message || "Failed to update TER bracket";
      toast.error("Update Failed", { description: msg });
    },
  });

  const deleteMutation = useMutation({
    mutationFn: async (id: number) => {
      return await api.delete(`/admin/tax/ter-brackets/${id}`);
    },
    onSuccess: () => {
      toast.success("TER bracket deleted");
      queryClient.invalidateQueries({ queryKey: ["ter-brackets"], type: "active" });
    },
    onError: (error: any) => {
      const msg = error?.response?.data?.error?.message || error?.response?.data?.message || "Failed to delete TER bracket";
      toast.error("Delete Failed", { description: msg });
    },
  });

  return { createMutation, updateMutation, deleteMutation };
}

export function usePtkpConfigs(year?: number) {
  return useQuery({
    queryKey: ["ptkp-configs", year],
    queryFn: async () => {
      const { data } = await api.get<ApiResponse<PTKPConfig[]>>("/admin/tax/ptkp-configs", {
        params: year ? { effective_year: year } : {},
      });
      return data.data ?? [];
    },
  });
}

export function usePtkpConfigMutations() {
  const queryClient = useQueryClient();

  const createMutation = useMutation({
    mutationFn: async (payload: PTKPConfigPayload) => {
      return await api.post("/admin/tax/ptkp-configs", payload);
    },
    onSuccess: () => {
      toast.success("PTKP config created");
      queryClient.invalidateQueries({ queryKey: ["ptkp-configs"], type: "active" });
    },
    onError: (error: any) => {
      const msg = error?.response?.data?.error?.message || error?.response?.data?.message || "Failed to create PTKP config";
      toast.error("Create Failed", { description: msg });
    },
  });

  const updateMutation = useMutation({
    mutationFn: async ({ id, ...payload }: PTKPConfigPayload & { id: number }) => {
      return await api.put(`/admin/tax/ptkp-configs/${id}`, payload);
    },
    onSuccess: () => {
      toast.success("PTKP config updated");
      queryClient.invalidateQueries({ queryKey: ["ptkp-configs"], type: "active" });
    },
    onError: (error: any) => {
      const msg = error?.response?.data?.error?.message || error?.response?.data?.message || "Failed to update PTKP config";
      toast.error("Update Failed", { description: msg });
    },
  });

  const deleteMutation = useMutation({
    mutationFn: async (id: number) => {
      return await api.delete(`/admin/tax/ptkp-configs/${id}`);
    },
    onSuccess: () => {
      toast.success("PTKP config deleted");
      queryClient.invalidateQueries({ queryKey: ["ptkp-configs"], type: "active" });
    },
    onError: (error: any) => {
      const msg = error?.response?.data?.error?.message || error?.response?.data?.message || "Failed to delete PTKP config";
      toast.error("Delete Failed", { description: msg });
    },
  });

  return { createMutation, updateMutation, deleteMutation };
}
