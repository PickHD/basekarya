export type FinanceType = "INCOME" | "EXPENSE";
export type FinanceStatus = "PENDING" | "APPROVED" | "REJECTED";

export interface FinanceCategory {
  id: number;
  name: string;
  type: FinanceType;
  description: string;
  created_at: string;
  updated_at: string;
}

export interface FinanceTransaction {
  id: number;
  creator_name: string;
  category_name: string;
  type: FinanceType;
  amount: number;
  transaction_date: string;
  reference_number: string;
  status: FinanceStatus;
  created_at: string;
}

export interface FinanceTransactionDetail {
  id: number;
  creator_name: string;
  category_name: string;
  category_type: FinanceType;
  type: FinanceType;
  amount: number;
  description: string;
  transaction_date: string;
  reference_number: string;
  status: FinanceStatus;
  rejection_reason: string;
  approved_by: number | null;
  approver_name: string;
  created_at: string;
}

export interface CreateTransactionPayload {
  finance_category_id: number;
  type: string;
  amount: number;
  description?: string;
  transaction_date: string;
  reference_number?: string;
}

export interface TransactionFilter {
  type?: string;
  status?: string;
  start_date?: string;
  end_date?: string;
  cursor?: string;
  limit?: number;
}

export interface TransactionActionPayload {
  id: number;
  action: "APPROVE" | "REJECT";
  rejection_reason?: string;
}

export interface CategoryPayload {
  name: string;
  type: string;
  description?: string;
}

export interface MonthlySummaryItem {
  month: string;
  income: number;
  expense: number;
}

export interface CategoryBreakdownItem {
  category_name: string;
  type: string;
  total: number;
}

export interface FinanceDashboard {
  total_income: number;
  total_expense: number;
  net_balance: number;
  transaction_count: number;
  monthly_summary: MonthlySummaryItem[];
  category_breakdown: CategoryBreakdownItem[];
  recent_transactions: FinanceTransaction[];
}

export interface FinanceFormDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}
