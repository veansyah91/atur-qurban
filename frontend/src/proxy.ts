import { NextRequest, NextResponse } from "next/server";

// Rute yang membutuhkan autentikasi
const PROTECTED_PATHS = ["/dashboard"];
// Rute yang hanya boleh diakses saat belum login
const AUTH_PATHS = ["/login", "/otp"];

export function proxy(request: NextRequest) {
  const { pathname } = request.nextUrl;
  const token = request.cookies.get("token")?.value;

  const isProtected = PROTECTED_PATHS.some((path) =>
    pathname.startsWith(path)
  );
  const isAuthPage = AUTH_PATHS.some((path) => pathname.startsWith(path));

  // Redirect ke login jika akses halaman protected tanpa token
  if (isProtected && !token) {
    const loginUrl = new URL("/login", request.url);
    loginUrl.searchParams.set("redirect", pathname);
    return NextResponse.redirect(loginUrl);
  }

  // Redirect ke dashboard jika sudah login tapi akses halaman auth
  if (isAuthPage && token) {
    return NextResponse.redirect(new URL("/dashboard", request.url));
  }

  return NextResponse.next();
}

export const config = {
  matcher: ["/dashboard/:path*", "/login", "/otp"],
};
