import Image from "next/image";
import Link from "next/link";
import { Button } from "@/components/ui/button";

const features = [
  {
    icon: "🎫",
    title: "Buat Kupon Otomatis",
    description:
      "Generate kupon pengambilan daging secara digital, cepat dan tanpa kertas.",
  },
  {
    icon: "📱",
    title: "Scan Kupon",
    description:
      "Validasi kupon saat pengambilan daging langsung dari smartphone.",
  },
  {
    icon: "💰",
    title: "Catat Tabungan Qurban",
    description:
      "Pantau progres tabungan qurban setiap peserta secara real-time.",
  },
  {
    icon: "🐄",
    title: "Manajemen Hewan Qurban",
    description:
      "Catat data hewan qurban, jenis, dan rencana pembagian daging.",
  },
  {
    icon: "👥",
    title: "Data Peserta",
    description:
      "Kelola daftar peserta qurban dengan mudah dalam satu tempat.",
  },
  {
    icon: "📊",
    title: "Laporan & Rekap",
    description:
      "Lihat rekap distribusi daging dan tabungan secara ringkas.",
  },
];

export default function Home() {
  return (
    <div className="min-h-screen bg-white text-gray-800 font-sans">
      {/* Navbar */}
      <nav className="sticky top-0 z-50 bg-white border-b border-gray-100 shadow-sm">
        <div className="max-w-6xl mx-auto px-4 sm:px-6 flex items-center justify-between h-16">
          <Link href="/">
            <Image
              src="/assets/logo-horizontal.png"
              alt="Atur Qurban"
              width={140}
              height={40}
              // PERBAIKAN 1: Menggunakan inline style untuk width & height auto
              style={{ width: "140px", height: "auto" }}
              className="object-contain"
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

      {/* Hero Section */}
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

      {/* Features Section */}
      <section className="py-16 sm:py-24 bg-surface">
        <div className="max-w-6xl mx-auto px-4 sm:px-6">
          <h2 className="text-2xl sm:text-3xl font-bold text-center text-gray-900 mb-4">
            Fitur Unggulan
          </h2>
          <p className="text-center text-gray-500 mb-12 max-w-md mx-auto">
            Semua yang Anda butuhkan untuk mengelola kegiatan qurban dengan
            efisien.
          </p>
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
            {features.map((f) => (
              <div
                key={f.title}
                className="bg-white rounded-xl p-6 shadow-sm border border-gray-100 hover:shadow-md transition-shadow"
              >
                <div className="text-3xl mb-4">{f.icon}</div>
                <h3 className="font-semibold text-gray-900 mb-2">{f.title}</h3>
                <p className="text-sm text-gray-500 leading-relaxed">
                  {f.description}
                </p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* Kenapa Gratis Section */}
      <section className="py-16 sm:py-24 bg-white">
        <div className="max-w-2xl mx-auto px-4 sm:px-6 text-center">
          <h2 className="text-2xl sm:text-3xl font-bold text-gray-900 mb-4">
            Kenapa Gratis?
          </h2>
          <p className="text-gray-500 leading-relaxed">
            Atur Qurban dibangun dengan semangat gotong royong. Kami percaya
            setiap komunitas berhak mendapatkan alat yang baik tanpa hambatan
            biaya. Aplikasi ini sepenuhnya gratis, didukung oleh komunitas dan
            mitra yang peduli.
          </p>
        </div>
      </section>

      {/* Footer */}
      <footer className="border-t border-gray-100 py-10 bg-surface">
        <div className="max-w-6xl mx-auto px-4 sm:px-6 flex flex-col items-center gap-4">
          <p className="text-xs text-gray-400 uppercase tracking-widest">
            Supported by
          </p>
          <a
            href="https://keuanganumum.com"
            target="_blank"
            rel="noopener noreferrer"
            className="opacity-80 hover:opacity-100 transition-opacity"
          >
            <Image
              src="/assets/keuangan-umum.jpeg"
              alt="Keuangan Umum"
              width={120}
              height={40}
              // PERBAIKAN 2: Menggunakan inline style untuk width & height auto
              style={{ width: "120px", height: "auto" }}
              className="object-contain rounded"
              loading="eager"
            />
          </a>
          <p className="text-xs text-gray-400">
            © {new Date().getFullYear()} Atur Qurban. All rights reserved.
          </p>
        </div>
      </footer>
    </div>
  );
}
