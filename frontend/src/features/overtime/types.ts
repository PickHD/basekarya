export type OvertimeStatus = "PENDING" | "APPROVED" | "REJECTED" | "PAID";

export interface Overtime {
  id: number;
  user_id: number;
  employee_id: number;
  employee_name?: string;
  employee_nik?: string;

  date: string;
  start_time: string;
  end_time: string;
  duration_minutes: number;

  reason?: string;
  status: OvertimeStatus;
  rejection_reason?: string;
  created_at: string;
}

export interface CreateOvertimePayload {
  date: string;
  start_time: string;
  end_time: string;
  reason: string;
}

export interface OvertimeFilter {
  status?: string;
  page?: number;
  limit?: number;
}

export interface OvertimeActionPayload {
  id: number;
  action: "APPROVE" | "REJECT";
  rejection_reason?: string;
}

export interface OvertimeFormDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}
