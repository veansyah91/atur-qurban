"use client";

import Image from "next/image";
import Link from "next/link";
import React from "react";

interface AuthLayoutProps {
  children: React.ReactNode;
}

export default function AuthLayout({ children }: AuthLayoutProps) {
  return (
    <div className="flex flex-col lg:flex-row min-h-screen w-full">
      {/* Sisi Kiri: Logo (Desktop Only) */}
      <div className="hidden lg:flex lg:w-1/2 flex-col items-center justify-center bg-primary/5 p-12">
        <div className="max-w-md text-center">
          <Link href="/" className="cursor-pointer">
            <Image
              src="/assets/logo-vertical.png"
              alt="Atur Qurban"
              width={400}
              height={120}
              className="w-80 h-auto object-contain mb-8 mx-auto"
              priority
            />
          </Link>
          <h1 className="text-3xl font-bold text-primary mb-4">Atur Qurban</h1>
          <p className="text-muted-foreground text-lg">
            Solusi manajemen qurban yang efisien dan transparan.
          </p>
        </div>
      </div>

      {/* Sisi Kanan: Form (Full width on mobile, half on desktop) */}
      <div className="w-full lg:w-1/2 flex items-center justify-center bg-surface px-4 py-12">
        <div className="w-full max-w-md">
          {/* Logo mobile only */}
          <div className="lg:hidden flex flex-col items-center mb-8">
            <Link href="/" className="cursor-pointer">
              <Image
                src="/assets/logo-vertical.png"
                alt="Atur Qurban"
                width={200}
                height={60}
                className="w-40 h-auto object-contain"
                priority
              />
            </Link>
          </div>
          {children}
        </div>
      </div>
    </div>
  );
}
