"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { Eye, EyeOff } from "lucide-react";

import { resetPassword } from "@/services/auth-reset.service";
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
  new_password: z.string().min(6, { message: "Password baru minimal 6 karakter" }),
});

export default function ResetPasswordPage() {
  const router = useRouter();
  
  const [phone, setPhone] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [showPassword, setShowPassword] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const storedPhone = sessionStorage.getItem("reset_phone");
    if (!storedPhone) {
      router.push("/forgot-password");
    } else {
      setPhone(storedPhone);
    }
  }, [router]);

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      otp: "",
      new_password: "",
    },
  });

  async function onSubmit(values: z.infer<typeof formSchema>) {
    if (!phone) return;
    
    try {
      setIsLoading(true);
      setError(null);
      
      await resetPassword({
        phone: phone,
        otp: values.otp,
        new_password: values.new_password,
      });

      sessionStorage.removeItem("reset_phone");
      
      // Redirect back to login after successful reset
      router.push("/login?reset=success"); 
    } catch (err: any) {
      setError(
        err.response?.data?.message || "Gagal mereset password. Pastikan OTP benar."
      );
    } finally {
      setIsLoading(false);
    }
  }

  if (!phone) return null; // loading state

  return (
    <Card className="border-none shadow-lg">
      <CardHeader className="space-y-2 text-center">
        <CardTitle className="text-2xl font-bold tracking-tight">
          Reset Password
        </CardTitle>
        <CardDescription className="text-sm text-muted-foreground">
          Masukkan OTP yang dikirim ke <strong>{phone}</strong> beserta password baru Anda.
        </CardDescription>
      </CardHeader>
      <CardContent>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
            {error && (
              <div className="p-3 rounded-md bg-red-50 text-red-500 text-sm font-medium">
                {error}
              </div>
            )}
            
            <FormField
              control={form.control}
              name="otp"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Kode OTP</FormLabel>
                  <FormControl>
                    <Input
                      placeholder="000000"
                      maxLength={6}
                      disabled={isLoading}
                      {...field}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="new_password"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Password Baru</FormLabel>
                  <FormControl>
                    <div className="relative">
                      <Input
                        type={showPassword ? "text" : "password"}
                        placeholder="••••••••"
                        disabled={isLoading}
                        {...field}
                      />
                      <Button
                        type="button"
                        variant="ghost"
                        size="icon"
                        className="absolute right-0 top-0 h-full px-3 py-2 hover:bg-transparent"
                        onClick={() => setShowPassword(!showPassword)}
                      >
                        {showPassword ? (
                          <EyeOff className="h-4 w-4 text-slate-400" />
                        ) : (
                          <Eye className="h-4 w-4 text-slate-400" />
                        )}
                      </Button>
                    </div>
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            
            <Button type="submit" className="w-full" disabled={isLoading}>
              {isLoading ? "Menyimpan..." : "Simpan Password Baru"}
            </Button>
          </form>
        </Form>
      </CardContent>
    </Card>
  );
}
