// ── Template Types ────────────────────────────────────────────────────────────

export interface OnboardingTemplateItem {
  id: number;
  task_name: string;
  description: string;
  sort_order: number;
}

export interface OnboardingTemplate {
  id: number;
  name: string;
  department: string;
  items: OnboardingTemplateItem[];
  created_at: string;
}

export interface CreateTemplatePayload {
  name: string;
  department: string;
  items: { task_name: string; description: string; sort_order: number }[];
}

export interface UpdateTemplatePayload {
  name: string;
  department: string;
  items: { task_name: string; description: string; sort_order: number }[];
}

// ── Workflow Types ────────────────────────────────────────────────────────────

export type OnboardingStatus = "IN_PROGRESS" | "COMPLETED";

export interface OnboardingWorkflowList {
  id: number;
  new_hire_name: string;
  new_hire_email: string;
  position: string;
  department: string;
  start_date: string | null;
  status: OnboardingStatus;
  progress: number; // 0–100
  created_at: string;
}

export interface OnboardingTask {
  id: number;
  task_name: string;
  description: string;
  department: string;
  is_completed: boolean;
  completed_by: string;
  completed_at: string | null;
  notes: string;
  sort_order: number;
}

export interface OnboardingWorkflowDetail {
  id: number;
  new_hire_name: string;
  new_hire_email: string;
  position: string;
  department: string;
  start_date: string | null;
  status: OnboardingStatus;
  progress: number;
  welcome_email_sent: boolean;
  created_at: string;
  it_tasks: OnboardingTask[];
  hr_tasks: OnboardingTask[];
  other_tasks: OnboardingTask[];
}

export interface CreateWorkflowPayload {
  applicant_id?: number;
  employee_id?: number;
  new_hire_name: string;
  new_hire_email: string;
  position?: string;
  department?: string;
  start_date?: string;
}

export interface UseWorkflowsParams {
  page?: number;
  limit?: number;
  status?: string;
  search?: string;
}
