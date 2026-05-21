export interface CompanyProfile {
  id: number;
  name: string;
  address: string;
  email: string;
  phone_number: string;
  website: string;
  tax_number: string;
  logo_url: string;
  subscription_plan_name: string;
  subscription_status: string;
  subscription_expires_at?: string;
  max_employees: number;
  plan_modules: string;
}

export interface CompanyProfilePayload {
  name: string;
  address: string;
  email: string;
  phone_number: string;
  website: string;
  tax_number: string;
  logo_url?: File;
}
