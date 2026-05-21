export interface SubscriptionPlan {
  id: number;
  name: string;
  slug: string;
  max_employees: number;
  price_monthly: number;
  features: string;
}

export interface UpgradePayload {
  plan_slug: string;
}

export interface UpgradeResponse {
  id: number;
  requested_plan_id: number;
  status: string;
}

export interface ReviewPayload {
  status: "APPROVED" | "REJECTED";
  notes?: string;
}

export interface SubscriptionRequestItem {
  id: number;
  company_name: string;
  current_plan_name: string;
  requested_plan_name: string;
  current_plan_price: number;
  requested_plan_price: number;
  price_difference: number;
  status: string;
  requested_by_name: string;
  requested_by_email: string;
  notes: string;
  created_at: string;
}

export interface CompanyListItem {
  id: number;
  name: string;
  email: string;
  phone_number: string;
  plan_name: string;
  plan_slug: string;
  subscription_status: string;
  subscription_expires_at: string;
  employee_count: number;
  created_at: string;
}

export interface CompanyDetail {
  id: number;
  name: string;
  email: string;
  phone_number: string;
  address: string;
  plan_name: string;
  plan_slug: string;
  max_employees: number;
  price_monthly: number;
  subscription_status: string;
  subscription_expires_at: string;
  employee_count: number;
  created_at: string;
}

export interface DashboardStats {
  total_companies: number;
  active_subscriptions: number;
  pending_payments: number;
  total_revenue: number;
  plan_distribution: PlanDistribution[];
}

export interface PlanDistribution {
  plan_name: string;
  plan_slug: string;
  count: number;
  revenue: number;
}
