import { Building2 } from "lucide-react";
import { Link } from "react-router-dom";
import { RegisterForm } from "@/features/auth/components/RegisterForm";

export default function RegisterPage() {
  return (
    <div className="min-h-screen w-full flex bg-white">
      <div className="hidden md:flex w-1/2 bg-slate-950 relative overflow-hidden items-center justify-center p-12">
        <div className="absolute inset-0 bg-grid-white/[0.05] bg-[size:40px_40px]" />
        <div className="absolute inset-0 bg-gradient-to-t from-blue-950/50 to-slate-950/20" />

        <div className="relative z-10 flex flex-col items-start text-white max-w-lg">
          <div className="flex items-center gap-3 mb-8">
            <div className="p-3 bg-blue-600 rounded-lg">
              <Building2 className="h-8 w-8 text-white" />
            </div>
            <h1 className="text-3xl font-bold tracking-tight">BaseKarya</h1>
          </div>
          <blockquote className="space-y-2 border-l-4 border-blue-600 pl-6">
            <p className="text-lg font-medium leading-relaxed">
              "Mulai kelola HR perusahaan Anda dengan mudah. Daftar gratis dan
              nikmati fitur manajemen karyawan, absensi, cuti, dan banyak lagi."
            </p>
            <footer className="text-sm text-slate-400">BaseKarya v2.12</footer>
          </blockquote>
        </div>
      </div>

      <div className="flex w-full md:w-1/2 flex-col items-center justify-center p-8 lg:p-24 bg-white">
        <div className="w-full max-w-md space-y-8 flex flex-col justify-center h-full">
          <div className="flex flex-col space-y-2 text-left">
            <h2 className="text-3xl font-bold tracking-tight text-slate-950">
              Daftar
            </h2>
            <p className="text-sm text-slate-500">
              Buat akun perusahaan baru untuk mulai menggunakan BaseKarya.
            </p>
          </div>

          <div className="mt-4">
            <RegisterForm />
            <div className="mt-6 text-center text-sm">
              Sudah punya akun?{" "}
              <Link
                to="/login"
                className="text-blue-600 hover:text-blue-800 font-semibold transition-colors duration-200"
              >
                Masuk di sini
              </Link>
            </div>
          </div>

          <p className="px-8 text-center text-sm text-slate-500 w-full mt-auto">
            <Link
              to="/"
              className="text-blue-600 hover:text-blue-800 font-semibold"
            >
              &larr; Kembali ke Beranda
            </Link>
            <span className="mx-2 text-slate-300">|</span>
            BaseKarya v2.12 &copy; 2026.
          </p>
        </div>
      </div>
    </div>
  );
}
