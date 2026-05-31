"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";

import { forgotPassword } from "@/services/auth-reset.service";
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
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

const formSchema = z.object({
  phone: z.string().min(10, { message: "Nomor telepon tidak valid" }),
});

export default function ForgotPasswordPage() {
  const router = useRouter();
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      phone: "",
    },
  });

  async function onSubmit(values: z.infer<typeof formSchema>) {
    try {
      setIsLoading(true);
      setError(null);
      await forgotPassword(values.phone);
      
      // Save phone for reset page
      sessionStorage.setItem("reset_phone", values.phone);
      
      // Redirect to reset password page
      router.push("/reset-password");
    } catch (err: any) {
      setError(
        err.response?.data?.message || "Gagal meminta reset password. Pastikan nomor terdaftar."
      );
    } finally {
      setIsLoading(false);
    }
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-slate-50 px-4 py-12">
      <Card className="w-full max-w-md">
        <CardHeader className="space-y-2 text-center">
          <CardTitle className="text-2xl font-bold tracking-tight">
            Lupa Password
          </CardTitle>
          <CardDescription className="text-sm text-muted-foreground">
            Masukkan nomor telepon Anda. Kami akan mengirimkan kode OTP untuk mereset password.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
              {error && (
                <div className="p-3 rounded-md bg-red-50 text-red-500 text-sm font-medium">
                  {error}
                </div>
              )}
              <FormField
                control={form.control}
                name="phone"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Nomor Telepon</FormLabel>
                    <FormControl>
                      <Input
                        placeholder="Contoh: 081234567890"
                        type="tel"
                        disabled={isLoading}
                        {...field}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <Button type="submit" className="w-full" disabled={isLoading}>
                {isLoading ? "Mengirim Permintaan..." : "Kirim Kode OTP"}
              </Button>
            </form>
          </Form>
        </CardContent>
        <CardFooter className="flex flex-col space-y-4">
          <div className="text-center text-sm text-slate-500">
            Ingat password Anda?{" "}
            <Link
              href="/login"
              className="font-semibold text-blue-600 hover:text-blue-500"
            >
              Kembali ke Login
            </Link>
          </div>
        </CardFooter>
      </Card>
    </div>
  );
}
