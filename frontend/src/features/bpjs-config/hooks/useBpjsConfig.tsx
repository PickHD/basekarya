import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import { api } from "@/lib/axios";
import type { BPJSRateConfig, BPJSRateConfigPayload } from "../types";
import type { ApiResponse } from "@/types/api";

export function useBpjsConfigs() {
  return useQuery({
    queryKey: ["bpjs-configs"],
    queryFn: async () => {
      const { data } = await api.get<ApiResponse<BPJSRateConfig[]>>("/admin/bpjs/configs");
      return data.data;
    },
  });
}

export function useBpjsConfigMutations() {
  const queryClient = useQueryClient();

  const createMutation = useMutation({
    mutationFn: async (payload: BPJSRateConfigPayload) => {
      return await api.post("/admin/bpjs/configs", payload);
    },
    onSuccess: () => {
      toast.success("BPJS config created");
      queryClient.invalidateQueries({ queryKey: ["bpjs-configs"], type: "active" });
    },
    onError: (error: any) => {
      const msg = error?.response?.data?.error?.message || error?.response?.data?.message || "Failed to create BPJS config";
      toast.error("Create Failed", { description: msg });
    },
  });

  const updateMutation = useMutation({
    mutationFn: async ({ id, ...payload }: BPJSRateConfigPayload & { id: number }) => {
      return await api.put(`/admin/bpjs/configs/${id}`, payload);
    },
    onSuccess: () => {
      toast.success("BPJS config updated");
      queryClient.invalidateQueries({ queryKey: ["bpjs-configs"], type: "active" });
    },
    onError: (error: any) => {
      const msg = error?.response?.data?.error?.message || error?.response?.data?.message || "Failed to update BPJS config";
      toast.error("Update Failed", { description: msg });
    },
  });

  const deleteMutation = useMutation({
    mutationFn: async (id: number) => {
      return await api.delete(`/admin/bpjs/configs/${id}`);
    },
    onSuccess: () => {
      toast.success("BPJS config deleted");
      queryClient.invalidateQueries({ queryKey: ["bpjs-configs"], type: "active" });
    },
    onError: (error: any) => {
      const msg = error?.response?.data?.error?.message || error?.response?.data?.message || "Failed to delete BPJS config";
      toast.error("Delete Failed", { description: msg });
    },
  });

  return { createMutation, updateMutation, deleteMutation };
}
