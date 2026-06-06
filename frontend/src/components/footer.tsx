import Image from "next/image";

export function Footer() {
  return (
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
            width={80}
            height={80}
            className="rounded"
            loading="eager"
          />
        </a>
        <p className="text-xs text-gray-400">
          © {new Date().getFullYear()} Atur Qurban. All rights reserved.
        </p>
      </div>
    </footer>
  );
}
