export { useDepartments } from "./useAdmin";
import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/axios";
import type { LookupItem } from "@/features/admin/types";

export const useShifts = () => {
  return useQuery({
    queryKey: ["shifts"],
    queryFn: async () => {
      const { data } = await api.get<{ data: LookupItem[] }>("/masters/shifts");
      return data.data;
    },
    staleTime: 1000 * 60 * 60,
  });
};
