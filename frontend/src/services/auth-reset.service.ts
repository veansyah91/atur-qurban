import api from "@/lib/api";

export async function forgotPassword(phone: string): Promise<void> {
  await api.post("/auth/forgot-password", { phone });
}

export async function resetPassword(payload: { phone: string; otp: string; new_password: string }): Promise<void> {
  await api.post("/auth/reset-password", payload);
}
