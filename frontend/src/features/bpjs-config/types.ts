export type BPJSComponentType = "KESEHATAN" | "JHT" | "JKK" | "JKM" | "JP";

export type IndustryRiskLevel = "VERY_LOW" | "LOW" | "MEDIUM" | "HIGH" | "VERY_HIGH";

export interface BPJSRateConfig {
  id: number;
  company_id: number | null;
  type: BPJSComponentType;
  employee_rate: number;
  employer_rate: number;
  max_salary_cap: number | null;
  industry_risk_level: IndustryRiskLevel | null;
  is_active: boolean;
  effective_from: string;
  effective_until: string | null;
  created_at: string;
  updated_at: string;
}

export interface BPJSRateConfigPayload {
  type: BPJSComponentType;
  employee_rate: number;
  employer_rate: number;
  max_salary_cap?: number | null;
  industry_risk_level?: IndustryRiskLevel | null;
  is_active?: boolean;
  effective_from: string;
  effective_until?: string | null;
}

export const BPJS_COMPONENT_LABELS: Record<BPJSComponentType, string> = {
  KESEHATAN: "BPJS Kesehatan",
  JHT: "Jaminan Hari Tua (JHT)",
  JKK: "Jaminan Kecelakaan Kerja (JKK)",
  JKM: "Jaminan Kematian (JKM)",
  JP: "Jaminan Pensiun (JP)",
};

export const RISK_LEVEL_LABELS: Record<IndustryRiskLevel, string> = {
  VERY_LOW: "Sangat Rendah",
  LOW: "Rendah",
  MEDIUM: "Menengah",
  HIGH: "Tinggi",
  VERY_HIGH: "Sangat Tinggi",
};
