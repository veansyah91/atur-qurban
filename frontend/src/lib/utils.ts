import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";

// Helper untuk menggabungkan class Tailwind secara aman
export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

// Format nomor telepon ke E.164 (+62...)
export function toE164(phone: string): string {
  const cleaned = phone.replace(/\D/g, "");
  if (cleaned.startsWith("0")) {
    return "+62" + cleaned.slice(1);
  }
  if (cleaned.startsWith("62")) {
    return "+" + cleaned;
  }
  return "+" + cleaned;
}
