# Copilot Instructions — Atur Qurban Frontend

## Konteks Project

Aplikasi web manajemen qurban (hewan kurban). Terdapat dua peran pengguna:
- **Admin** — login dengan nomor telepon + password
- **Peserta** — login dengan nomor telepon + OTP

Backend: Go/Fiber REST API, JWT auth, semua endpoint di bawah `/api/v1`.

---

## Tech Stack

| Komponen | Library |
|---|---|
| Framework | Next.js 16+ (App Router) |
| Language | TypeScript |
| Styling | Tailwind CSS v4 |
| UI Components | shadcn/ui (base), `components/shared/` (custom) |
| Data Fetching | TanStack Query v5 |
| HTTP Client | Axios (`src/lib/api.ts`) |
| State (auth) | Zustand + persist middleware |
| Form | React Hook Form + Zod resolver |
| Icons | Lucide React |
| Theme | next-themes |

---

## Struktur Folder

```
src/
├── app/
│   ├── (auth)/           # Route group: /login, /otp
│   ├── (dashboard)/      # Route group: halaman setelah login
│   ├── providers.tsx     # QueryClient + ThemeProvider (client component)
│   └── layout.tsx        # Root layout — import Providers di sini
├── components/
│   ├── ui/               # Komponen shadcn/ui (auto-generate, jangan edit manual)
│   └── shared/           # Komponen reusable custom (Navbar, Sidebar, dll)
├── hooks/                # Custom hooks (useAuth, usePagination, dll)
├── lib/
│   ├── api.ts            # Axios instance + interceptors
│   └── utils.ts          # cn() helper dan utility functions
├── services/             # Fungsi pemanggil API, dipakai di queryFn/mutationFn
│   └── auth.service.ts
├── stores/               # Zustand stores
│   └── auth.store.ts     # Auth state + setAuth/clearAuth
├── types/                # TypeScript interfaces/types global
│   └── auth.ts
└── proxy.ts              # Next.js proxy (middleware) — proteksi route
```

---

## Pattern TanStack Query

### Query key convention
```ts
// Format: [domain, identifier?]
['sacrifices']           // daftar semua
['sacrifices', id]       // satu item
['participants', userId] // relasi
```

### Contoh useQuery
```ts
import { useQuery } from "@tanstack/react-query";
import { getSacrifices } from "@/services/sacrifice.service";

export function useSacrifices() {
  return useQuery({
    queryKey: ["sacrifices"],
    queryFn: getSacrifices,
  });
}
```

### Contoh useMutation + invalidasi cache
```ts
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { createSacrifice } from "@/services/sacrifice.service";

export function useCreateSacrifice() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: createSacrifice,
    onSuccess: () => {
      // Invalidasi agar list ter-refresh otomatis
      queryClient.invalidateQueries({ queryKey: ["sacrifices"] });
    },
  });
}
```

---

## Pattern Form (React Hook Form + Zod)

```ts
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";

const schema = z.object({
  phone: z.string().min(10, "Nomor telepon tidak valid"),
  password: z.string().min(6, "Password minimal 6 karakter"),
});

type FormValues = z.infer<typeof schema>;

export function LoginForm() {
  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<FormValues>({ resolver: zodResolver(schema) });

  const onSubmit = async (values: FormValues) => { /* ... */ };

  return (
    <form onSubmit={handleSubmit(onSubmit)}>
      <input {...register("phone")} />
      {errors.phone && <p className="text-red-500 text-sm">{errors.phone.message}</p>}
    </form>
  );
}
```

---

## Pattern Auth

### Akses token & user di client component
```ts
import { useAuthStore } from "@/stores/auth.store";

const { user, isAuthenticated, setAuth, clearAuth } = useAuthStore();
```

### Setelah login sukses
```ts
const data = await loginAdmin({ phone, password });
setAuth(data.user, data.access_token);
router.push("/dashboard");
```

### Logout
```ts
await logout(); // hapus refresh token di server
clearAuth();    // hapus state lokal + cookie
router.push("/login");
```

### Proteksi route
Proxy di `src/proxy.ts` membaca cookie `token`. Halaman di bawah `/dashboard` otomatis terlindungi.

---

## Konvensi Komponen

- Gunakan komponen shadcn/ui sebagai **base** (hasil `npx shadcn add ...`)
- Komponen yang dimodifikasi atau digabungkan → buat di `components/shared/`
- Semua komponen interaktif harus punya direktif `"use client"` jika menggunakan hooks
- Komponen murni/presentasi tidak perlu direktif (Server Component by default)

---

## Coding Conventions

### Bahasa
- **Komentar kode dalam Bahasa Indonesia**
- Nama variabel, fungsi, tipe, dan komponen tetap dalam Bahasa Inggris

### Error handling API
```ts
try {
  const data = await someApiCall();
} catch (error) {
  if (axios.isAxiosError(error)) {
    const message = error.response?.data?.message ?? "Terjadi kesalahan";
    toast.error(message);
  }
}
```

### Format response API (dari backend)
```ts
// Sukses
{ "message": "...", "data": { ... } }

// Error
{ "message": "...", "errors": { ... } }

// Paginasi
{ "message": "...", "data": [...], "meta": { "page": 1, "total": 100 } }
```

---

## Planning

Sebelum mengimplementasikan fitur baru, buat dokumen perencanaan di session folder (`plan.md`). Dokumen ini akan digunakan oleh AI model lain untuk implementasi.

### Format Plan (minimal)
1. Deskripsi singkat fitur
2. File baru yang perlu dibuat (service, hooks, page, components)
3. File yang perlu dimodifikasi
4. Route/URL yang digunakan
5. Query keys yang dipakai
6. Catatan khusus (validasi, permission, dsb)

---

## Code Review Checkpoint

Jalankan review setelah setiap fitur selesai:

- [ ] Semua constraint di `copilot-instructions.md` diikuti
- [ ] Tidak ada logic bisnis langsung di page component (pindah ke hooks/services)
- [ ] Form menggunakan Zod schema + RHF resolver
- [ ] Query key konsisten dengan konvensi `[domain, id?]`
- [ ] Komponen client sudah ada `"use client"` jika memakai hooks
- [ ] Tidak ada hard-coded URL — selalu lewat `src/lib/api.ts`
- [ ] Error API ditangani dan ditampilkan ke user
