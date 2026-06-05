// Tipe data untuk auth domain

export interface User {
  id: string;
  phone: string;
  name: string;
  is_admin: boolean;
  last_login_at: string | null;
}

export interface AuthState {
  user: User | null;
  accessToken: string | null;
  isAuthenticated: boolean;
}

export interface LoginAdminPayload {
  phone: string;
  password: string;
}

export interface LoginOtpRequestPayload {
  phone: string;
}

export interface LoginOtpVerifyPayload {
  phone: string;
  otp: string;
}

export interface AuthResponse {
  access_token: string;
  refresh_token: string;
  user: User;
}
