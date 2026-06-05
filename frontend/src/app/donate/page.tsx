import Image from "next/image";
import Link from "next/link";
import { Button } from "@/components/ui/button";

export default function DonatePage() {
  return (
    <div className="min-h-screen bg-surface">
      <nav className="sticky top-0 z-50 bg-white border-b border-gray-100 shadow-sm">
        <div className="max-w-6xl mx-auto px-4 sm:px-6 flex items-center justify-between h-16">
          <Link href="/">
            <Image
              src="/assets/logo-horizontal.png"
              alt="Atur Qurban"
              width={140}
              height={40}
              className="w-[140px] h-auto object-contain"
            />
          </Link>
        </div>
      </nav>

      <main className="max-w-2xl mx-auto px-4 py-16 text-center">
        <h1 className="text-3xl font-bold text-gray-900 mb-4">Dukung Atur Qurban</h1>
        <p className="text-gray-600 mb-8 leading-relaxed">
          Atur Qurban adalah platform gratis yang dibangun untuk memudahkan pengelolaan qurban di komunitas. 
          Dukungan Anda membantu kami terus mengembangkan fitur-fitur baru dan menjaga platform tetap gratis untuk semua orang.
        </p>

        <div className="bg-white p-8 rounded-2xl shadow-sm border border-gray-100">
          <h2 className="text-xl font-semibold mb-6">Pilih Nominal Donasi</h2>
          <div className="grid grid-cols-2 gap-4 mb-8">
            <Button variant="outline">Rp 50.000</Button>
            <Button variant="outline">Rp 100.000</Button>
            <Button variant="outline">Rp 250.000</Button>
            <Button variant="outline">Nominal Lain</Button>
          </div>
          <Button className="w-full py-6 text-lg">Lanjutkan Pembayaran</Button>
        </div>
        
        <Link href="/" className="inline-block mt-8 text-sm text-gray-500 hover:text-primary transition-colors">
          &larr; Kembali ke Beranda
        </Link>
      </main>
    </div>
  );
}
