export type RequisitionStatus = "DRAFT" | "PENDING" | "APPROVED" | "REJECTED" | "CLOSED";
export type RequisitionPriority = "LOW" | "MEDIUM" | "HIGH" | "URGENT";
export type ApplicantStage = "SCREENING" | "INTERVIEW" | "OFFERING" | "HIRED" | "REJECTED";
export type EmploymentType = "PKWT" | "PKWTT";

export interface JobRequisition {
  id: number;
  title: string;
  description?: string;
  department_id: number;
  department_name: string;
  employment_type: EmploymentType;
  quantity: number;
  priority: RequisitionPriority;
  status: RequisitionStatus;
  requester_id: number;
  requester_name: string;
  approved_by?: number;
  approver_name?: string;
  rejection_reason?: string;
  target_date?: string;
  created_at: string;
  updated_at?: string;
}

export interface CreateRequisitionPayload {
  department_id: number;
  title: string;
  description?: string;
  quantity: number;
  employment_type: EmploymentType;
  priority: RequisitionPriority;
  target_date?: string;
}

export interface RequisitionActionPayload {
  action: "APPROVE" | "REJECT";
  rejection_reason?: string;
}

export interface Applicant {
  id: number;
  job_requisition_id: number;
  full_name: string;
  email: string;
  phone_number?: string;
  resume_url?: string;
  stage: ApplicantStage;
  stage_order: number;
  notes?: string;
  rejection_reason?: string;
  created_at: string;
}

export interface ApplicantDetail extends Applicant {
  stage_histories: StageHistory[];
}

export interface StageHistory {
  id: number;
  from_stage: string;
  to_stage: string;
  changed_by_name: string;
  notes: string;
  created_at: string;
}

export interface CreateApplicantPayload {
  full_name: string;
  email: string;
  phone_number?: string;
  resume_base64?: string;
}

export interface UpdateApplicantStagePayload {
  stage: ApplicantStage;
  notes?: string;
  rejection_reason?: string;
}

export interface KanbanBoard {
  SCREENING: Applicant[];
  INTERVIEW: Applicant[];
  OFFERING: Applicant[];
  HIRED: Applicant[];
  REJECTED: Applicant[];
}

export interface UseRequisitionsParams {
  page?: number;
  limit?: number;
  status?: string;
  priority?: string;
  department_id?: number;
  search?: string;
}
