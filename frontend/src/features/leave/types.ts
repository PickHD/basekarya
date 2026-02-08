export type LeaveStatus = "PENDING" | "APPROVED" | "REJECTED";

export interface LeaveType {
  id: number;
  name: string;
  default_quota: number;
  is_deducted: boolean;
}

export interface LeaveRequest {
  id: number;
  employee_id: number;
  employee_name?: string;
  employee_nik?: string;
  leave_type_id: number;
  leave_type: LeaveType;

  start_date: string;
  end_date: string;
  total_days: number;
  reason: string;
  attachment_url?: string;

  status: LeaveStatus;
  rejection_reason?: string;
  created_at: string;
}

export interface ApplyLeavePayload {
  leave_type_id: number;
  start_date: string;
  end_date: string;
  reason: string;
  attachment_base64?: string;
}

export interface LeaveActionPayload {
  id: number;
  action: "APPROVE" | "REJECT";
  rejection_reason?: string;
}

export interface UseLeavesParams {
  page: number;
  limit: number;
  status?: string;
}
