export type ContractType = "PKWT" | "PKWTT";

export interface Contract {
  id: number;
  employee_id: number;
  employee_name?: string;
  employee_nik?: string;
  contract_type: ContractType;
  contract_number?: string;
  start_date: string;
  end_date?: string;       // null for PKWTT
  notes?: string;
  attachment_url?: string;
  created_at: string;
}

export interface UpsertContractPayload {
  employee_id: number;
  contract_type: ContractType;
  contract_number?: string;
  start_date: string;
  end_date?: string;
  notes?: string;
  attachment_base64?: string;
}

export interface UseContractsParams {
  page: number;
  limit: number;
  contract_type?: string;
  search?: string;
  expiring_within_days?: number;
}
