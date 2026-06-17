export interface TERBracket {
  id: number;
  company_id: number | null;
  category: string;
  bracket_number: number;
  min_monthly_salary: number;
  rate: number;
  effective_from: string;
  effective_until: string | null;
  created_at: string;
  updated_at: string;
}

export interface TERBracketPayload {
  category: string;
  bracket_number: number;
  min_monthly_salary: number;
  rate: number;
  effective_from: string;
  effective_until?: string | null;
}

export interface PTKPConfig {
  id: number;
  code: string;
  annual_amount: number;
  effective_year: number;
  created_at: string;
  updated_at: string;
}

export interface PTKPConfigPayload {
  code: string;
  annual_amount: number;
  effective_year: number;
}

export const TER_CATEGORIES = ["A", "B", "C"] as const;

export const PTKP_CODES = ["TK/0", "TK/1", "TK/2", "TK/3", "K/0", "K/1", "K/2", "K/3"] as const;
