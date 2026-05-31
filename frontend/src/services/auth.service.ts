import api from "@/lib/api";
import { useAuthStore } from "@/stores/auth.store";
import type {
  LoginAdminPayload,
  LoginOtpRequestPayload,
  LoginOtpVerifyPayload,
  AuthResponse,
  User,
} from "@/types/auth";

// Login admin dengan nomor telepon dan password
export async function loginAdmin(
  payload: LoginAdminPayload
): Promise<AuthResponse> {
  const { data } = await api.post<{ data: AuthResponse }>(
    "/auth/login",
    payload
  );
  return data.data;
}

// Kirim OTP ke nomor telepon peserta
export async function requestOtp(
  payload: LoginOtpRequestPayload
): Promise<void> {
  await api.post("/auth/otp/request", payload);
}

// Verifikasi OTP dan dapatkan token
export async function verifyOtp(
  payload: LoginOtpVerifyPayload
): Promise<AuthResponse> {
  const { data } = await api.post<{ data: AuthResponse }>(
    "/auth/otp/verify",
    payload
  );
  return data.data;
}

// Me - Ambil user yang login
export async function getMe(): Promise<User> {
  const { data } = await api.get<{ data: User }>("/auth/me");
  return data.data;
}

// Refresh - Perbarui Token 
export async function refreshAccessToken(refreshToken: string): Promise<AuthResponse> {
   const { data } = await api.post<{ data: AuthResponse }>("/auth/refresh", { refresh_token: refreshToken });
   return data.data;
}

// Logout — hapus refresh token di server
export async function logout(): Promise<void> {
  try {
    await api.post("/auth/logout");
  } catch (e) {
    // ignore
  } finally {
    useAuthStore.getState().clearAuth();
    if (typeof window !== "undefined") {
       window.location.href = "/login";
    }
  }
}
