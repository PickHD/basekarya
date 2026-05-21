import { api } from "@/lib/axios";
import type { FinanceDashboard } from "@/features/finance/types";
import { useQuery } from "@tanstack/react-query";

export const useFinanceDashboard = (startDate?: string, endDate?: string) => {
  return useQuery({
    queryKey: ["finance-dashboard", startDate, endDate],
    queryFn: async () => {
      const params: Record<string, string> = {};
      if (startDate) params.start_date = startDate;
      if (endDate) params.end_date = endDate;

      const { data } = await api.get<{ data: FinanceDashboard }>("/finances/dashboard", { params });
      return data.data;
    },
    refetchInterval: 1000 * 60 * 5,
  });
};
