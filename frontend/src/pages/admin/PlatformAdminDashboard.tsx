import { useDashboardStats } from "@/features/subscription/hooks/useSubscription";
import { Building2, CheckCircle, Clock, DollarSign, Users } from "lucide-react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { cn } from "@/lib/utils";

export default function PlatformAdminDashboard() {
  const { data, isLoading } = useDashboardStats();
  const stats = data?.data;

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary" />
      </div>
    );
  }

  const statCards = [
    {
      title: "Total Perusahaan",
      value: stats?.total_companies ?? 0,
      icon: Building2,
      color: "text-blue-600",
      bg: "bg-blue-50",
    },
    {
      title: "Langganan Aktif",
      value: stats?.active_subscriptions ?? 0,
      icon: CheckCircle,
      color: "text-emerald-600",
      bg: "bg-emerald-50",
    },
    {
      title: "Menunggu Pembayaran",
      value: stats?.pending_payments ?? 0,
      icon: Clock,
      color: "text-amber-600",
      bg: "bg-amber-50",
    },
    {
      title: "Estimasi Pendapatan/Bulan",
      value: new Intl.NumberFormat("id-ID", {
        style: "currency",
        currency: "IDR",
        minimumFractionDigits: 0,
      }).format(stats?.total_revenue ?? 0),
      icon: DollarSign,
      color: "text-violet-600",
      bg: "bg-violet-50",
    },
  ];

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-foreground">Platform Dashboard</h1>
        <p className="text-muted-foreground mt-1">
          Ringkasan data platform BaseKarya
        </p>
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        {statCards.map((card) => (
          <Card key={card.title}>
            <CardHeader className="flex flex-row items-center justify-between pb-2">
              <CardTitle className="text-sm font-medium text-muted-foreground">
                {card.title}
              </CardTitle>
              <div className={cn("p-2 rounded-lg", card.bg)}>
                <card.icon className={cn("h-4 w-4", card.color)} />
              </div>
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{card.value}</div>
            </CardContent>
          </Card>
        ))}
      </div>

      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Users className="h-5 w-5" />
            Distribusi Paket
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            {stats?.plan_distribution?.map((plan: any) => (
              <div key={plan.plan_slug} className="flex items-center gap-4">
                <div className="w-32 font-medium text-sm">{plan.plan_name}</div>
                <div className="flex-1">
                  <div className="h-6 bg-muted rounded-full overflow-hidden">
                    <div
                      className={cn(
                        "h-full rounded-full transition-all",
                        plan.plan_slug === "free"
                          ? "bg-slate-400"
                          : plan.plan_slug === "basic"
                          ? "bg-blue-500"
                          : "bg-violet-500"
                      )}
                      style={{
                        width: `${Math.max(
                          (plan.count / (stats?.total_companies || 1)) * 100,
                          5
                        )}%`,
                      }}
                    />
                  </div>
                </div>
                <div className="w-16 text-right text-sm font-medium">
                  {plan.count} co
                </div>
                <div className="w-32 text-right text-sm text-muted-foreground">
                  {plan.revenue > 0
                    ? new Intl.NumberFormat("id-ID", {
                        style: "currency",
                        currency: "IDR",
                        minimumFractionDigits: 0,
                      }).format(plan.revenue)
                    : "Gratis"}
                </div>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
