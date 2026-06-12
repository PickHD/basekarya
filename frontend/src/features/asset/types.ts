export type AssetStatus = "AVAILABLE" | "ASSIGNED" | "MAINTENANCE" | "RETIRED";
export type AssetCondition = "GOOD" | "FAIR" | "DAMAGED" | "LOST";
export type AssetAssignmentStatus = "PENDING" | "ACTIVE" | "RETURNED" | "REJECTED";

export interface AssetCategory {
  id: number;
  name: string;
  description: string;
  created_at: string;
  updated_at: string;
}

export interface Asset {
  id: number;
  name: string;
  description: string;
  serial_number: string;
  asset_category_id: number;
  category_name: string;
  status: AssetStatus;
  condition: AssetCondition;
  current_employee?: string;
  created_at: string;
  updated_at: string;
}

export interface AssetAssignment {
  id: number;
  asset_id: number;
  asset_name: string;
  employee_id: number;
  employee_name: string;
  employee_nik: string;
  purpose: string;
  expected_return_date: string | null;
  actual_return_date: string | null;
  notes: string;
  status: AssetAssignmentStatus;
  rejection_reason?: string;
  created_at: string;
}

export interface CreateAssetCategoryPayload {
  name: string;
  description: string;
}

export interface UpdateAssetCategoryPayload {
  name: string;
  description: string;
}

export interface CreateAssetPayload {
  asset_category_id: number;
  name: string;
  description: string;
  serial_number: string;
  condition: AssetCondition;
}

export interface UpdateAssetPayload {
  asset_category_id?: number;
  name?: string;
  description?: string;
  serial_number?: string;
  status?: AssetStatus;
  condition?: AssetCondition;
}

export interface CreateAssetAssignmentPayload {
  asset_id: number;
  purpose: string;
  expected_return_date?: string;
}

export interface AssetActionPayload {
  id: number;
  action: "APPROVE" | "REJECT";
  rejection_reason?: string;
}

export interface AssetFilter {
  status?: string;
  condition?: string;
  category_id?: number;
  page?: number;
  limit?: number;
}

export interface AssetAssignmentFilter {
  status?: string;
  page?: number;
  limit?: number;
}

export interface AssetFormDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

export interface AssetCategoryFormDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  category?: AssetCategory | null;
}

export interface AssetAssignmentFormDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

export interface AssetDetailDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  assetId: number | null;
}

export interface AssetAssignmentDetailDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  assignmentId: number | null;
}
