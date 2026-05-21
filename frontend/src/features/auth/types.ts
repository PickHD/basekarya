export interface LoginPayload {
  username: string;
  password: string;
}

export interface LoginResponse {
  message: string;
  data: {
    token: string;
    must_change_password: boolean;
  };
}

export interface DecodedToken {
  user_id: number;
  role: string;
  permissions: string[];
  company_id: number;
  is_platform_admin: boolean;
  exp: number;
}

export interface ForgotPasswordPayload {
  email: string;
}

export interface VerifyOTPPayload {
  code: string;
}

export interface VerifyOTPResponse {
  message: string;
  data: {
    is_valid: boolean;
  };
}

export interface ResetPasswordPayload {
  code: string;
  password: string;
}

export interface RegisterPayload {
  company_name: string;
  admin_name: string;
  admin_email: string;
  password: string;
  phone_number: string;
  plan_slug: string;
}

export interface RegisterResponse {
  message: string;
  data: {
    username: string;
  };
}
