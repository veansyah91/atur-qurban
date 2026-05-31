"use client";

import { create } from "zustand";
import { persist } from "zustand/middleware";
import type { User } from "@/types/auth";

interface AuthStore {
  user: User | null;
  accessToken: string | null;
  refreshToken: string | null;
  isAuthenticated: boolean;
  // Simpan data auth setelah login sukses
  setAuth: (user: User, accessToken: string, refreshToken: string) => void;
  // Update data user saja
  setUser: (user: User) => void;
  // Hapus data auth saat logout
  clearAuth: () => void;
}

export const useAuthStore = create<AuthStore>()(
  persist(
    (set) => ({
      user: null,
      accessToken: null,
      refreshToken: null,
      isAuthenticated: false,

      setAuth: (user, accessToken, refreshToken) => {
        // Simpan token ke cookie agar bisa dibaca Next.js middleware
        document.cookie = `token=${accessToken}; path=/; max-age=${60 * 60 * 24 * 7}`;
        set({ user, accessToken, refreshToken, isAuthenticated: true });
      },

      setUser: (user) => {
         set({ user });
      },

      clearAuth: () => {
        // Hapus cookie token
        document.cookie = "token=; path=/; max-age=0";
        set({ user: null, accessToken: null, refreshToken: null, isAuthenticated: false });
      },
    }),
    {
      name: "auth-storage",
      // Hanya persist field yang diperlukan
      partialize: (state) => ({
        user: state.user,
        accessToken: state.accessToken,
        refreshToken: state.refreshToken,
        isAuthenticated: state.isAuthenticated,
      }),
    }
  )
);
