import api from "@/lib/api";

export interface RegisterRequestPayload {
  name: string;
  phone: string;
  password?: string;
  confirm_password?: string;
}

export async function requestRegisterOtp(payload: RegisterRequestPayload): Promise<void> {
  await api.post("/auth/register/request-otp", payload);
}

export async function resendRegisterOtp(phone: string): Promise<void> {
  await api.post("/auth/register/resend-otp", { phone });
}

export async function verifyRegisterOtp(payload: { phone: string; otp: string }): Promise<any> {
  const { data } = await api.post("/auth/register/verify", payload);
  return data.data;
}

// Keep the others for compatibility if needed, but the ones above are specifically for the new user register flow
