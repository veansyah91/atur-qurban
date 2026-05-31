"use client";

import { useState, useEffect } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";

import { useAuthStore } from "@/stores/auth.store";
import { verifyRegisterOtp, resendRegisterOtp } from "@/services/auth-register.service";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

const formSchema = z.object({
  otp: z.string().min(6, { message: "OTP harus 6 digit" }).max(6),
});

import Image from "next/image";

// Implement component logic
export default function VerifyPage() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const phone = searchParams.get("phone") ?? "";
  const setAuth = useAuthStore((s) => s.setAuth);

  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [successMsg, setSuccessMsg] = useState<string | null>(null);
  const [countdown, setCountdown] = useState<number>(0);

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: { otp: "" },
  });

  useEffect(() => {
    let timer: number | undefined;
    if (countdown > 0) {
      timer = window.setInterval(() => setCountdown((c) => c - 1), 1000);
    }
    return () => {
      if (timer) clearInterval(timer);
    };
  }, [countdown]);

  async function onSubmit(values: z.infer<typeof formSchema>) {
    try {
      setIsLoading(true);
      setError(null);
      const data = await verifyRegisterOtp({ phone, otp: values.otp });
      // Expecting data to include user and tokens
      setAuth(data.user, data.access_token, data.refresh_token);
      router.push("/dashboard");
    } catch (err: any) {
      setError(err?.response?.data?.message || "Kode OTP tidak valid.");
    } finally {
      setIsLoading(false);
    }
  }

  async function handleResend() {
    try {
      setIsLoading(true);
      setError(null);
      await resendRegisterOtp(phone);
      setSuccessMsg("OTP berhasil dikirim ulang.");
      setCountdown(60);
    } catch (err: any) {
      setError(err?.response?.data?.message || "Gagal mengirim ulang OTP.");
    } finally {
      setIsLoading(false);
    }
  }

  if (!phone) return null; // loading state

  return (
    <Card className="border-none shadow-lg">
      <CardHeader className="space-y-1 text-center">
        <CardTitle className="text-2xl font-bold tracking-tight">
          Verifikasi OTP
        </CardTitle>
        <CardDescription className="text-sm text-muted-foreground">
          Masukkan 6 digit kode OTP yang telah dikirim ke nomor <strong>{phone}</strong>
        </CardDescription>
      </CardHeader>
      <CardContent>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
            {error && (
              <div className="p-3 rounded-md bg-red-50 text-red-500 text-sm font-medium text-center">
                {error}
              </div>
            )}
            {successMsg && (
              <div className="p-3 rounded-md bg-green-50 text-green-600 text-sm font-medium text-center">
                {successMsg}
              </div>
            )}
            
            <FormField
              control={form.control}
              name="otp"
              render={({ field }) => (
                <FormItem className="text-center">
                  <FormLabel className="sr-only">Kode OTP</FormLabel>
                  <FormControl>
                    <Input
                      placeholder="000000"
                      className="text-center text-2xl tracking-widest focus-visible:ring-primary"
                      maxLength={6}
                      disabled={isLoading}
                      {...field}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            
            <Button type="submit" className="w-full" disabled={isLoading}>
              {isLoading ? "Memverifikasi..." : "Verifikasi"}
            </Button>
            
            <div className="text-center text-sm text-slate-500 mt-4">
              Belum menerima OTP?{" "}
              <button
                type="button"
                onClick={handleResend}
                disabled={countdown > 0 || isLoading}
                className="font-semibold text-primary hover:text-primary-dark disabled:text-slate-400 disabled:hover:text-slate-400"
              >
                {countdown > 0 ? `Kirim ulang dalam ${countdown}s` : "Kirim Ulang OTP"}
              </button>
            </div>
          </form>
        </Form>
      </CardContent>
    </Card>
  );
}

