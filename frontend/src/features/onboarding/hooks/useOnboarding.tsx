import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/axios";
import { toast } from "sonner";
import type {
  CreateTemplatePayload,
  UpdateTemplatePayload,
  CreateWorkflowPayload,
  UseWorkflowsParams,
} from "@/features/onboarding/types";

// ── Template Hooks ────────────────────────────────────────────────────────────

export const useOnboardingTemplates = () => {
  return useQuery({
    queryKey: ["onboarding-templates"],
    queryFn: async () => {
      const { data } = await api.get("/onboarding/templates");
      return (data?.data ?? data) as any[];
    },
  });
};

export const useCreateTemplate = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (payload: CreateTemplatePayload) => {
      const { data } = await api.post("/onboarding/templates", payload);
      return data;
    },
    onSuccess: () => {
      toast.success("Template created successfully");
      queryClient.invalidateQueries({ queryKey: ["onboarding-templates"] });
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || "Failed to create template");
    },
  });
};

export const useUpdateTemplate = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async ({ id, payload }: { id: number; payload: UpdateTemplatePayload }) => {
      const { data } = await api.put(`/onboarding/templates/${id}`, payload);
      return data;
    },
    onSuccess: () => {
      toast.success("Template updated successfully");
      queryClient.invalidateQueries({ queryKey: ["onboarding-templates"] });
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || "Failed to update template");
    },
  });
};

export const useDeleteTemplate = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (id: number) => {
      const { data } = await api.delete(`/onboarding/templates/${id}`);
      return data;
    },
    onSuccess: () => {
      toast.success("Template deleted");
      queryClient.invalidateQueries({ queryKey: ["onboarding-templates"] });
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || "Failed to delete template");
    },
  });
};

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
