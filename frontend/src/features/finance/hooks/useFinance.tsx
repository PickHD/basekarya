import { api } from "@/lib/axios";
import type { Meta } from "@/types/api";
import type {
  FinanceTransaction,
  FinanceTransactionDetail,
  TransactionFilter,
  CreateTransactionPayload,
  TransactionActionPayload,
} from "@/features/finance/types";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";

export const useFinanceTransactions = (filter: TransactionFilter) => {
  return useQuery({
    queryKey: ["finance-transactions", filter],
    queryFn: async () => {
      const { data } = await api.get<{ data: FinanceTransaction[]; meta: Meta }>("/finances/transactions", {
        params: filter,
      });
      return data;
    },
    placeholderData: (prev) => prev,
  });
};

export const useFinanceTransaction = (id: string) => {
  return useQuery({
    queryKey: ["finance-transaction", id],
    queryFn: async () => {
      const { data } = await api.get<{ data: FinanceTransactionDetail }>(`/finances/transactions/${id}`);
      return data.data;
    },
    enabled: !!id,
  });
};

export const useCreateFinanceTransaction = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (payload: CreateTransactionPayload) => {
      const { data } = await api.post("/finances/transactions", payload);
      return data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["finance-transactions"] });
      toast.success("Transaksi keuangan berhasil dibuat!");
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || "Gagal membuat transaksi keuangan");
    },
  });
};

export const useFinanceAction = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({ id, action, rejection_reason }: TransactionActionPayload) => {
      const { data } = await api.put(`/finances/transactions/${id}/action`, {
        action,
        rejection_reason,
      });
      return data;
    },
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ["finance-transaction", variables.id.toString()] });
      queryClient.invalidateQueries({ queryKey: ["finance-transactions"] });
      toast.success(`Transaksi berhasil di-${variables.action.toLowerCase()}`);
    },
    onError: (error: any) => {
      const errMsg = error.response?.data?.error || "Terjadi kesalahan saat memproses aksi";
      toast.error(errMsg);
    },
  });
};

export const useExportFinanceTransactions = () => {
  return useMutation({
    mutationFn: async (params: TransactionFilter) => {
      const response = await api.get("/finances/transactions/export", {
        params,
        responseType: "blob",
      });
      return response.data;
    },
    onSuccess: (data) => {
      const url = window.URL.createObjectURL(new Blob([data]));
      const link = document.createElement("a");
      link.href = url;
      link.setAttribute("download", "finance_transactions.xlsx");
      document.body.appendChild(link);
      link.click();
      link.parentNode?.removeChild(link);
    },
    onError: () => {
      toast.error("Gagal mengunduh data");
    },
  });
};
