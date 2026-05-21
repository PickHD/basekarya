import {
  Building2,
  Check,
  ArrowRight,
  Users,
  Clock,
  Shield,
  BarChart3,
  Crown,
} from "lucide-react";
import { Link } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { useScrollAnimation } from "@/hooks/useScrollAnimation";
import { cn } from "@/lib/utils";

const features = [
  {
    icon: Clock,
    title: "Absensi & Kehadiran",
    desc: "Tracking kehadiran real-time dengan geolocation dan rekap otomatis.",
  },
  {
    icon: Users,
    title: "Manajemen Karyawan",
    desc: "Kelola data karyawan, kontrak, dan dokumen HR dalam satu tempat.",
  },
  {
    icon: BarChart3,
    title: "Penggajian",
    desc: "Generate payroll otomatis, slip gaji PDF, dan kirim via email.",
  },
  {
    icon: Shield,
    title: "Izin & Cuti",
    desc: "Pengajuan cuti, lembur, pinjaman, dan reimbursement dengan approval flow.",
  },
];

const tiers = [
  {
    name: "Free",
    price: "Gratis",
    priceNote: "Selamanya",
    desc: "Untuk UMKM yang baru mulai",
    max: "5 Karyawan",
    features: ["Absensi & Kehadiran", "Cuti & Izin"],
    cta: "Mulai Gratis",
    href: "/register?plan=free",
    highlight: false,
  },
  {
    name: "Basic",
    price: "Rp99.000",
    priceNote: "/bulan",
    desc: "Untuk bisnis yang berkembang",
    max: "50 Karyawan",
    features: [
      "Semua fitur Free",
      "Lembur",
      "Pinjaman",
      "Reimbursement",
      "Penggajian & Payroll",
      "Manajemen Kontrak",
      "Keuangan & Finance",
    ],
    cta: "Pilih Basic",
    href: "/register?plan=basic",
    highlight: true,
  },
  {
    name: "Pro",
    price: "Rp249.000",
    priceNote: "/bulan",
    desc: "Untuk tim yang lengkap",
    max: "Unlimited Karyawan",
    features: [
      "Semua fitur Basic",
      "Rekrutmen & Seleksi",
      "Onboarding Karyawan",
      "Prioritas Support",
    ],
    cta: "Pilih Pro",
    href: "/register?plan=pro",
    highlight: false,
  },
];

function AnimatedSection({
  children,
  className,
  delay = 0,
}: {
  children: React.ReactNode;
  className?: string;
  delay?: number;
}) {
  const { ref, isVisible } = useScrollAnimation();
  return (
    <div
      ref={ref}
      className={cn("opacity-0", isVisible && "animate-fade-up", className)}
      style={{ animationDelay: `${delay}ms` }}
    >
      {children}
    </div>
  );
}

export default function LandingPage() {
  return (
    <div className="min-h-screen bg-white text-foreground">
      {/* Navbar */}
      <nav className="sticky top-0 z-50 border-b bg-white/80 backdrop-blur-md animate-fade-in">
        <div className="max-w-7xl mx-auto px-6 h-16 flex items-center justify-between">
          <Link to="/" className="flex items-center gap-2 group">
            <div className="p-2 bg-primary rounded-lg transition-transform group-hover:scale-105">
              <Building2 className="h-5 w-5 text-primary-foreground" />
            </div>
            <span className="text-xl font-bold text-foreground">BaseKarya</span>
          </Link>
          <div className="flex items-center gap-3">
            <Link to="/login">
              <Button variant="ghost">Masuk</Button>
            </Link>
            <Link to="/register">
              <Button className="bg-primary hover:bg-primary/90">
                Daftar Gratis
              </Button>
            </Link>
          </div>
        </div>
      </nav>

      {/* Hero */}
      <section className="relative overflow-hidden">
        <div className="absolute inset-0 bg-gradient-to-br from-primary/5 via-white to-muted" />
        <div className="absolute top-20 left-10 w-72 h-72 bg-primary/10 rounded-full blur-3xl animate-float" />
        <div
          className="absolute bottom-10 right-10 w-96 h-96 bg-primary/5 rounded-full blur-3xl animate-float"
          style={{ animationDelay: "1.5s" }}
        />

        <div className="relative max-w-7xl mx-auto px-6 py-24 md:py-36 text-center">
          <AnimatedSection>
            <Badge className="mb-6 bg-primary/10 text-primary hover:bg-primary/10 px-4 py-1.5 text-sm border-primary/20">
              HRIS untuk UMKM Indonesia
            </Badge>
          </AnimatedSection>

          <AnimatedSection delay={100}>
            <h1 className="text-4xl md:text-6xl font-extrabold tracking-tight text-foreground max-w-4xl mx-auto leading-tight">
              Kelola Tim Anda dengan{" "}
              <span className="text-primary">Lebih Mudah</span>
            </h1>
          </AnimatedSection>

          <AnimatedSection delay={200}>
            <p className="mt-6 text-lg md:text-xl text-muted-foreground max-w-2xl mx-auto leading-relaxed">
              Sistem manajemen sumber daya manusia terintegrasi — absensi,
              penggajian, cuti, dan lainnya dalam satu platform.
            </p>
          </AnimatedSection>

          <AnimatedSection delay={300}>
            <div className="mt-10 flex flex-col sm:flex-row gap-4 justify-center">
              <Link to="/register">
                <Button
                  size="lg"
                  className="bg-primary hover:bg-primary/90 text-base px-8 py-6 shadow-lg shadow-primary/25 hover:shadow-primary/40 transition-shadow"
                >
                  Mulai Gratis Sekarang
                  <ArrowRight className="ml-2 h-5 w-5" />
                </Button>
              </Link>
              <Link to="/login">
                <Button
                  size="lg"
                  variant="outline"
                  className="text-base px-8 py-6 border-border hover:bg-muted"
                >
                  Sudah punya akun? Masuk
                </Button>
              </Link>
            </div>
            <p className="mt-4 text-sm text-muted-foreground">
              Tanpa kartu kredit. Gratis untuk 5 karyawan.
            </p>
          </AnimatedSection>
        </div>
      </section>

      {/* Features */}
      <section className="py-20 bg-muted/50">
        <div className="max-w-7xl mx-auto px-6">
          <AnimatedSection className="text-center mb-16">
            <h2 className="text-3xl md:text-4xl font-bold text-foreground">
              Semua yang Anda Butuhkan
            </h2>
            <p className="mt-4 text-lg text-muted-foreground max-w-2xl mx-auto">
              Fitur lengkap untuk mengelola SDM perusahaan Anda secara efisien.
            </p>
          </AnimatedSection>
          <div className="grid md:grid-cols-2 lg:grid-cols-4 gap-8">
            {features.map((f, i) => (
              <AnimatedSection key={f.title} delay={i * 100}>
                <div className="bg-card rounded-xl p-6 shadow-sm border border-border hover:shadow-md hover:border-primary/20 transition-all duration-300 h-full">
                  <div className="p-3 bg-primary/10 rounded-lg w-fit mb-4">
                    <f.icon className="h-6 w-6 text-primary" />
                  </div>
                  <h3 className="text-lg font-semibold text-foreground">
                    {f.title}
                  </h3>
                  <p className="mt-2 text-sm text-muted-foreground leading-relaxed">
                    {f.desc}
                  </p>
                </div>
              </AnimatedSection>
            ))}
          </div>
        </div>
      </section>

      {/* Pricing */}
      <section className="py-20">
        <div className="max-w-7xl mx-auto px-6">
          <AnimatedSection className="text-center mb-16">
            <h2 className="text-3xl md:text-4xl font-bold text-foreground">
              Pilih Paket yang Tepat
            </h2>
            <p className="mt-4 text-lg text-muted-foreground">
              Mulai gratis, upgrade kapan saja sesuai kebutuhan.
            </p>
          </AnimatedSection>
          <div className="grid md:grid-cols-3 gap-8 max-w-5xl mx-auto items-start">
            {tiers.map((tier, i) => (
              <AnimatedSection key={tier.name} delay={i * 150}>
                <div
                  className={cn(
                    "relative rounded-2xl border-2 p-8 flex flex-col transition-all duration-300 hover:shadow-lg",
                    tier.highlight
                      ? "border-primary shadow-xl md:scale-105 hover:shadow-primary/20"
                      : "border-border hover:border-primary/30",
                  )}
                >
                  {tier.highlight && (
                    <div className="absolute -top-4 left-1/2 -translate-x-1/2">
                      <Badge className="bg-primary text-primary-foreground hover:bg-primary px-3 py-1 shadow-md">
                        <Crown className="w-3 h-3 mr-1" />
                        Populer
                      </Badge>
                    </div>
                  )}
                  <div>
                    <h3 className="text-xl font-bold text-foreground">
                      {tier.name}
                    </h3>
                    <p className="text-sm text-muted-foreground mt-1">
                      {tier.desc}
                    </p>
                    <div className="mt-4 flex items-baseline gap-1">
                      <span className="text-4xl font-extrabold text-foreground">
                        {tier.price}
                      </span>
                      <span className="text-muted-foreground">
                        {tier.priceNote}
                      </span>
                    </div>
                    <p className="text-sm text-muted-foreground mt-2">
                      Maks. {tier.max}
                    </p>
                  </div>
                  <ul className="mt-8 space-y-3 flex-1">
                    {tier.features.map((f) => (
                      <li key={f} className="flex items-start gap-2 text-sm">
                        <Check className="h-4 w-4 text-primary mt-0.5 shrink-0" />
                        <span className="text-foreground">{f}</span>
                      </li>
                    ))}
                  </ul>
                  <Link to={tier.href} className="mt-8">
                    <Button
                      className={cn(
                        "w-full py-5 transition-all duration-200",
                        tier.highlight
                          ? "bg-primary hover:bg-primary/90 shadow-md shadow-primary/25 hover:shadow-primary/40"
                          : "border-border hover:bg-muted hover:text-foreground",
                      )}
                      variant={tier.highlight ? "default" : "outline"}
                    >
                      {tier.cta}
                    </Button>
                  </Link>
                </div>
              </AnimatedSection>
            ))}
          </div>
        </div>
      </section>

      {/* CTA */}
      <section className="py-20 bg-foreground">
        <AnimatedSection className="max-w-4xl mx-auto px-6 text-center">
          <h2 className="text-3xl md:text-4xl font-bold text-primary-foreground">
            Siap Memulai?
          </h2>
          <p className="mt-4 text-lg text-muted-foreground">
            Daftar sekarang dan kelola tim Anda dengan lebih efisien.
          </p>
          <Link to="/register">
            <Button
              size="lg"
              className="mt-8 bg-primary hover:bg-primary/90 text-base px-8 py-6 shadow-lg shadow-primary/25"
            >
              Daftar Gratis
              <ArrowRight className="ml-2 h-5 w-5" />
            </Button>
          </Link>
        </AnimatedSection>
      </section>

      {/* Footer */}
      <footer className="border-t border-border py-8">
        <div className="max-w-7xl mx-auto px-6 flex flex-col md:flex-row items-center justify-between gap-4">
          <div className="flex items-center gap-2">
            <Building2 className="h-4 w-4 text-muted-foreground" />
            <span className="text-sm text-muted-foreground">
              BaseKarya v2.10 &copy; 2026
            </span>
          </div>
          <div className="flex gap-6 text-sm text-muted-foreground">
            <Link
              to="/login"
              className="hover:text-foreground transition-colors"
            >
              Masuk
            </Link>
            <Link
              to="/register"
              className="hover:text-foreground transition-colors"
            >
              Daftar
            </Link>
          </div>
        </div>
      </footer>
    </div>
  );
}
