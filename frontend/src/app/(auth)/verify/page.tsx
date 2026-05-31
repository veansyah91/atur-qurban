"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
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

export default function VerifyPage() {
  const router = useRouter();
  const setAuth = useAuthStore((state) => state.setAuth);
  
  const [phone, setPhone] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [successMsg, setSuccessMsg] = useState<string | null>(null);
  const [countdown, setCountdown] = useState(0);

  useEffect(() => {
    const storedPhone = sessionStorage.getItem("register_phone");
    if (!storedPhone) {
      router.push("/register");
    } else {
      setPhone(storedPhone);
      setCountdown(60); // Start 60s countdown
    }
  }, [router]);

  useEffect(() => {
    if (countdown > 0) {
      const timer = setTimeout(() => setCountdown(countdown - 1), 1000);
      return () => clearTimeout(timer);
    }
  }, [countdown]);

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      otp: "",
    },
  });

  async function onSubmit(values: z.infer<typeof formSchema>) {
    if (!phone) return;
    
    try {
      setIsLoading(true);
      setError(null);
      setSuccessMsg(null);
      
      const response = await verifyRegisterOtp({
        phone: phone,
        otp: values.otp,
      });

      setAuth(response.user, response.access_token, response.refresh_token);
      sessionStorage.removeItem("register_phone");
      router.push("/dashboard"); 
    } catch (err: any) {
      setError(
        err.response?.data?.message || "OTP tidak valid. Silakan coba lagi."
      );
    } finally {
      setIsLoading(false);
    }
  }

  async function handleResend() {
    if (!phone || countdown > 0) return;
    
    try {
      setIsLoading(true);
      setError(null);
      await resendRegisterOtp(phone);
      setSuccessMsg("OTP berhasil dikirim ulang");
      setCountdown(60);
    } catch (err: any) {
      setError(
        err.response?.data?.message || "Gagal mengirim ulang OTP."
      );
    } finally {
      setIsLoading(false);
    }
  }

  if (!phone) return null; // loading state

  return (
    <div className="flex min-h-screen items-center justify-center bg-slate-50 px-4 py-12">
      <Card className="w-full max-w-md">
        <CardHeader className="space-y-2 text-center">
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
                        className="text-center text-2xl tracking-widest"
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
                  className="font-semibold text-blue-600 hover:text-blue-500 disabled:text-slate-400 disabled:hover:text-slate-400"
                >
                  {countdown > 0 ? `Kirim ulang dalam ${countdown}s` : "Kirim Ulang OTP"}
                </button>
              </div>
            </form>
          </Form>
        </CardContent>
      </Card>
    </div>
  );
}
