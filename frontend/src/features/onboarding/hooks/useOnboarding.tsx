import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/axios";
import { toast } from "sonner";
import type {
  CreateWorkflowPayload,
  UseWorkflowsParams,
} from "@/features/onboarding/types";

// ── Workflow Hooks ────────────────────────────────────────────────────────────

export const useOnboardingWorkflows = (params: UseWorkflowsParams = {}) => {
  return useQuery({
    queryKey: ["onboarding-workflows", params],
    queryFn: async () => {
      const { data } = await api.get("/onboarding/workflows", { params });
      return data;
    },
  });
};

export const useOnboardingWorkflowDetail = (id: number | null) => {
  return useQuery({
    queryKey: ["onboarding-workflows", "detail", id],
    queryFn: async () => {
      const { data } = await api.get(`/onboarding/workflows/${id}`);
      return data?.data ?? data;
    },
    enabled: !!id,
  });
};

export const useCreateWorkflow = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (payload: CreateWorkflowPayload) => {
      const { data } = await api.post("/onboarding/workflows", payload);
      return data;
    },
    onSuccess: () => {
      toast.success("Onboarding workflow created");
      queryClient.invalidateQueries({ queryKey: ["onboarding-workflows"] });
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || "Failed to create workflow");
    },
  });
};

// ── Task Hooks ────────────────────────────────────────────────────────────────

export const useCompleteTask = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async ({ id, notes }: { id: number; notes?: string }) => {
      const { data } = await api.put(`/onboarding/tasks/${id}/complete`, { notes: notes ?? "" });
      return data;
    },
    onSuccess: (_, { id }) => {
      toast.success("Task completed!");
      queryClient.invalidateQueries({ queryKey: ["onboarding-workflows"] });
      void id;
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || "Failed to complete task");
    },
  });
};
