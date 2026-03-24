export interface Permission {
  id: number;
  name: string;
  description: string | null;
  module: string | null;
}

export interface Role {
  id: number;
  name: string;
  description: string;
  is_active: boolean;
  permissions: Permission[];
}

export interface AssignPermissionsPayload {
  role_id: number;
  permission_ids: number[];
}

export interface RolesResponse {
  data: Role[];
}

export interface RoleResponse {
  data: Role;
}

export interface PermissionsResponse {
  data: Permission[];
}