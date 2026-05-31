"use client";

import { useEffect } from "react";
import { getMe } from "@/services/auth.service";
import { useAuthStore } from "@/stores/auth.store";

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const { accessToken, setUser, clearAuth } = useAuthStore();

  useEffect(() => {
    // Memanggil API profil (Me) saat aplikasi di-refresh untuk validasi sesi
    const validateSession = async () => {
      if (accessToken) {
        try {
          const user = await getMe();
          setUser(user);
        } catch (error) {
          // Error 401 should be handled by interceptor, but if we get here we can clear
          console.error("Failed to validate session", error);
        }
      }
    };

    validateSession();
  }, [accessToken, setUser]);

  return <>{children}</>;
}