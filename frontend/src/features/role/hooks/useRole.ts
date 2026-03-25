import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/axios";
import { toast } from "sonner";
import type { AssignPermissionsPayload, RolesResponse, Role, PermissionsResponse } from "@/features/role/types";

export function useRoles() {
  return useQuery({
    queryKey: ["roles"],
    queryFn: async () => {
      const res = await api.get<RolesResponse>("/roles");
      return res.data;
    },
  });
}

export function useCreateRole() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (payload: { name: string; description?: string }) => {
      const res = await api.post("/roles", payload);
      return res.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["roles"] });
      toast.success("Role Created");
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || "Failed to create role");
    },
  });
}

export function useRoleDetails(id: string) {
  return useQuery({
    queryKey: ["roles", id],
    queryFn: async () => {
      if (!id) return null;
      const res = await api.get<any>(`/roles/${id}/permissions`);
      const data = res.data.data;
      return {
        id: data.role_id,
        name: data.role_name,
        permissions: data.permissions,
        description: "",
        is_active: true
      } as Role;
    },
    enabled: !!id,
  });
}

export function useAllPermissions() {
  return useQuery({
    queryKey: ["permissions"],
    queryFn: async () => {
      const res = await api.get<PermissionsResponse>("/permissions");
      return res.data.data;
    },
  });
}

export function useAssignPermissions() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (payload: AssignPermissionsPayload) => {
      const res = await api.put(
        `/roles/${payload.role_id}/permissions`,
        { permission_ids: payload.permission_ids }
      );
      return res.data;
    },
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ["roles"] });
      queryClient.invalidateQueries({ queryKey: ["roles", variables.role_id.toString()] });
      toast.success("Permissions Updated");
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || "Failed to update permissions");
    },
  });
}
