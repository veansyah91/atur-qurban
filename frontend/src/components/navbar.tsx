import Image from "next/image";
import Link from "next/link";
import { Button } from "@/components/ui/button";

export function Navbar() {
  return (
    <nav className="sticky top-0 z-50 bg-white border-b border-gray-100 shadow-sm">
      <div className="max-w-6xl mx-auto px-4 sm:px-6 flex items-center justify-between h-16">
        <Link href="/">
          <Image
            src="/assets/logo-horizontal.png"
            alt="Atur Qurban"
            width={140}
            height={35}
            priority
          />
        </Link>
        <div className="flex items-center gap-2 sm:gap-4">
          <Link
            href="/donate"
            className="hidden sm:block"
          >
            <Button variant="secondary" size="sm">Donasi</Button>
          </Link>
          <div className="h-6 w-[1px] bg-gray-200 hidden sm:block mx-1"></div>
          <Link
            href="/login"
          >
            <Button variant="ghost" size="sm">Masuk</Button>
          </Link>
          <Link
            href="/register"
          >
            <Button size="sm">Daftar</Button>
          </Link>
        </div>
      </div>
    </nav>
  );
}
