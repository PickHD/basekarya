export type ReimbursementStatus = "PENDING" | "APPROVED" | "REJECTED";

export interface Reimbursement {
  id: number;
  user_id: number;
  title: string;
  description: string;
  amount: number;
  date_of_expense: string;
  proof_file_url: string;
  status: ReimbursementStatus;
  rejection_reason?: string;
  created_at: string;
  requester_name?: string;
  approved_by?: number;
}

export interface CreateReimbursementPayload {
  title: string;
  description: string;
  amount: number;
  date: string;
  proof_file: FileList;
}

export interface ReimbursementFilter {
  status?: string;
  page?: number;
  limit?: number;
}

export interface ReimbursementActionPayload {
  id: number;
  action: "APPROVE" | "REJECT";
  rejection_reason?: string;
}

export interface ReimbursementFormDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}
