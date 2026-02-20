import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/axios";
import type { LookupItem } from "@/features/admin/types";

export const useDepartments = () => {
  return useQuery({
    queryKey: ["departments"],
    queryFn: async () => {
      const { data } = await api.get<{ data: LookupItem[] }>(
        "/admin/departments"
      );
      return data.data;
    },
    staleTime: 1000 * 60 * 60,
  });
};

export const useShifts = () => {
  return useQuery({
    queryKey: ["shifts"],
    queryFn: async () => {
      const { data } = await api.get<{ data: LookupItem[] }>("/admin/shifts");
      return data.data;
    },
    staleTime: 1000 * 60 * 60,
  });
};
