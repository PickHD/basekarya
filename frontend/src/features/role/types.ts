export interface PermissionGroup {
  id: number;
  name: string;
}

export interface Permission {
  id: number;
  name: string;
  display_name?: string;
  description: string | null;
}

export interface PermissionsGroupedResponse {
  group: PermissionGroup;
  permissions: Permission[];
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
  data: PermissionsGroupedResponse[];
}