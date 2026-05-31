import api from "@/lib/api";
import type {
  AuthResponse,
  LoginAdminPayload,
  LoginOtpRequestPayload,
  LoginOtpVerifyPayload,
} from "@/types/auth";

// Login admin dengan nomor telepon dan password
export async function loginAdmin(
  payload: LoginAdminPayload
): Promise<AuthResponse> {
  const { data } = await api.post<{ data: AuthResponse }>(
    "/auth/admin/login",
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

// Logout — hapus refresh token di server
export async function logout(): Promise<void> {
  await api.post("/auth/logout");
}
