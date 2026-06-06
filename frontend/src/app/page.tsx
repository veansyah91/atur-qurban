import { Navbar } from "@/components/navbar";
import { Hero } from "@/components/hero";
import { Features } from "@/components/features";
import { Footer } from "@/components/footer";

export default function Home() {
  return (
    <div className="min-h-screen bg-white text-gray-800 font-sans">
      <Navbar />
      <Hero />
      <Features />

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

      <Footer />
    </div>
  );
}

