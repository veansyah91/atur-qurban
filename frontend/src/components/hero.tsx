import Link from "next/link";
import { Button } from "@/components/ui/button";

export function Hero() {
  return (
    <section className="bg-primary text-white overflow-hidden relative">
      {/* Simple decorative element */}
      <div className="absolute top-0 right-0 w-64 h-64 bg-primary-light rounded-full -mr-32 -mt-32 opacity-50 blur-3xl"></div>
      <div className="max-w-6xl mx-auto px-4 sm:px-6 py-20 sm:py-32 text-center relative z-10">
        <span className="inline-block bg-accent text-white text-xs font-semibold px-3 py-1 rounded-full mb-6 uppercase tracking-wide">
          Gratis Selamanya
        </span>
        <h1 className="text-3xl sm:text-5xl font-bold leading-tight mb-6">
          Kelola Qurban Anda
          <br className="hidden sm:block" />
          Lebih Mudah &amp; Digital
        </h1>
        <p className="text-base sm:text-lg text-green-50 max-w-xl mx-auto mb-10 opacity-90">
          Platform manajemen qurban digital — dari tabungan, kupon, hingga
          distribusi daging — semua dalam satu aplikasi, tanpa biaya.
        </p>
        <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
          <Link href="/register">
            <Button size="lg" variant="secondary" className="px-10 py-6 text-lg font-bold">
              Mulai Sekarang
            </Button>
          </Link>
          <Link href="/donate" className="sm:hidden">
            <Button size="lg" variant="outline" className="text-white border-white hover:bg-white/10 px-10 py-6 text-lg">
              Donasi
            </Button>
          </Link>
        </div>
      </div>
    </section>
  );
}
