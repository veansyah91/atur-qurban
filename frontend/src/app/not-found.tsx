import Link from 'next/link';

export default function NotFound() {
  return (
    <div className="flex flex-col items-center justify-center min-h-screen bg-white p-4 text-center">
      <h2 className="text-4xl font-bold text-primary mb-2">404</h2>
      <h3 className="text-xl font-semibold text-gray-800 mb-4">Halaman Tidak Ditemukan</h3>
      <p className="text-gray-500 mb-8 max-w-md">
        Maaf, kami tidak dapat menemukan halaman yang Anda cari. Pastikan alamat yang Anda masukkan sudah benar.
      </p>
      <Link 
        href="/" 
        className="bg-primary hover:bg-primary-dark text-white font-medium px-6 py-2 rounded-lg transition-colors"
      >
        Kembali ke Beranda
      </Link>
    </div>
  );
}
