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

export function Features() {
  return (
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
  );
}
