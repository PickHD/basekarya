import {
  Building2,
  MapPin,
  Mail,
  Phone,
  Loader2,
  FileText,
  CreditCard,
  Users,
  CheckCircle2,
  Crown,
  ArrowUpRight,
} from "lucide-react";

import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { useCompanyProfile } from "@/features/company/hooks/useCompany";
import { CompanyProfileForm } from "@/features/company/components/CompanyProfileForm";
import { useSubscriptionPlans } from "@/features/auth/hooks/useAuth";
import { useRequestUpgrade } from "@/features/subscription/hooks/useSubscription";
import { cn } from "@/lib/utils";

export default function CompanySettingsPage() {
  const { data: company, isLoading, isError } = useCompanyProfile();
  const { data: plansData } = useSubscriptionPlans();
  const { mutate: requestUpgrade, isPending: isUpgrading } = useRequestUpgrade();

  if (isLoading) {
    return (
      <div className="flex justify-center p-10">
        <Loader2 className="animate-spin h-8 w-8 text-blue-600" />
      </div>
    );
  }

  if (isError || !company) {
    return (
      <div className="p-10 text-center text-red-500">
        Failed to load company profile. Please try again later.
      </div>
    );
  }

  const plans = plansData?.data || [];
  const currentFeatures: string[] = company.plan_modules
    ? (() => {
        try {
          return JSON.parse(company.plan_modules).modules || [];
        } catch {
          return [];
        }
      })()
    : [];

  return (
    <div className="space-y-6 max-w-6xl mx-auto pb-10">
      <div className="flex flex-col gap-2">
        <h2 className="text-3xl font-bold tracking-tight text-slate-900">
          Company Settings
        </h2>
        <p className="text-slate-500">
          Manage organization identity, branding, and billing details.
        </p>
      </div>

      <div className="grid gap-6 md:grid-cols-12">
        <div className="md:col-span-4 space-y-6">
          <Card className="border-t-4 border-t-blue-600 shadow-sm">
            <CardHeader className="text-center">
              <div className="mx-auto w-32 h-32 mb-4 relative">
                <Avatar className="w-32 h-32 border-4 border-slate-50 shadow-md">
                  <AvatarImage
                    src={company.logo_url}
                    className="object-contain bg-white"
                  />
                  <AvatarFallback className="text-4xl bg-slate-100 text-slate-400">
                    <Building2 className="w-12 h-12" />
                  </AvatarFallback>
                </Avatar>
              </div>

              <CardTitle className="text-xl">{company.name}</CardTitle>
              <CardDescription className="font-mono text-blue-600 break-all">
                {company.website || "No website"}
              </CardDescription>
            </CardHeader>

            <CardContent>
              <div className="flex justify-center mb-6">
                <Badge
                  variant="outline"
                  className="px-3 py-1 uppercase bg-slate-50"
                >
                  Headquarters
                </Badge>
              </div>

              <div className="space-y-4 text-sm border-t pt-4">
                <div className="flex items-start gap-3">
                  <MapPin className="w-4 h-4 text-slate-400 mt-0.5 shrink-0" />
                  <span className="font-medium text-slate-700 leading-snug">
                    {company.address || "-"}
                  </span>
                </div>

                <div className="flex items-center gap-3">
                  <Mail className="w-4 h-4 text-slate-400 shrink-0" />
                  <span className="font-medium text-slate-700">
                    {company.email}
                  </span>
                </div>

                <div className="flex items-center gap-3">
                  <Phone className="w-4 h-4 text-slate-400 shrink-0" />
                  <span className="font-medium text-slate-700">
                    {company.phone_number || "-"}
                  </span>
                </div>

                <div className="flex items-center gap-3">
                  <FileText className="w-4 h-4 text-slate-400 shrink-0" />
                  <span className="font-medium text-slate-700">
                    Tax: {company.tax_number || "-"}
                  </span>
                </div>
              </div>
            </CardContent>
          </Card>

          <Card className="border-t-4 border-t-emerald-600 shadow-sm">
            <CardHeader>
              <div className="flex items-center gap-2">
                <CreditCard className="w-5 h-5 text-emerald-600" />
                <CardTitle className="text-base">Langganan</CardTitle>
              </div>
            </CardHeader>
            <CardContent className="space-y-4 text-sm">
              <div className="flex items-center justify-between">
                <span className="text-slate-500">Paket</span>
                <Badge className="bg-emerald-100 text-emerald-800 hover:bg-emerald-100">
                  <Crown className="w-3 h-3 mr-1" />
                  {company.subscription_plan_name || "Free"}
                </Badge>
              </div>
              <div className="flex items-center justify-between">
                <span className="text-slate-500">Status</span>
                <Badge
                  variant="outline"
                  className={cn(
                    "capitalize",
                    company.subscription_status === "ACTIVE"
                      ? "text-emerald-600"
                      : "text-amber-600"
                  )}
                >
                  <CheckCircle2 className="w-3 h-3 mr-1" />
                  {company.subscription_status === "PENDING_PAYMENT"
                    ? "Menunggu Pembayaran"
                    : company.subscription_status}
                </Badge>
              </div>
              <div className="flex items-center justify-between">
                <span className="text-slate-500">Karyawan</span>
                <div className="flex items-center gap-1">
                  <Users className="w-4 h-4 text-slate-400" />
                  <span className="font-medium">
                    {company.max_employees === 0
                      ? "Unlimited"
                      : `Max ${company.max_employees}`}
                  </span>
                </div>
              </div>
              {company.subscription_expires_at && (
                <div className="flex items-center justify-between">
                  <span className="text-slate-500">Berlaku hingga</span>
                  <span className="font-medium">
                    {new Date(
                      company.subscription_expires_at
                    ).toLocaleDateString("id-ID")}
                  </span>
                </div>
              )}
              {currentFeatures.length > 0 && (
                <div className="pt-2 border-t">
                  <p className="text-slate-500 mb-2">Fitur aktif:</p>
                  <div className="flex flex-wrap gap-1">
                    {currentFeatures.map((f) => (
                      <Badge
                        key={f}
                        variant="secondary"
                        className="text-xs capitalize"
                      >
                        {f}
                      </Badge>
                    ))}
                  </div>
                </div>
              )}
            </CardContent>
          </Card>
        </div>

        <div className="md:col-span-8 space-y-6">
          <Card>
            <CardHeader>
              <div className="flex items-center gap-2">
                <Building2 className="w-5 h-5 text-slate-500" />
                <CardTitle>General Information</CardTitle>
              </div>
              <CardDescription>
                Update your company logo, official address, and contact
                information.
              </CardDescription>
            </CardHeader>
            <CardContent>
              <CompanyProfileForm initialData={company} />
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <div className="flex items-center gap-2">
                <Crown className="w-5 h-5 text-slate-500" />
                <CardTitle>Upgrade Paket</CardTitle>
              </div>
              <CardDescription>
                Tingkatkan paket Anda untuk membuka lebih banyak fitur dan
                karyawan.
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                {plans.map((plan: any) => {
                  const isCurrentPlan =
                    plan.name.toLowerCase() ===
                    company.subscription_plan_name?.toLowerCase();
                  const isHigher =
                    plan.price_monthly >
                    (plans.find(
                      (p: any) =>
                        p.name.toLowerCase() ===
                        company.subscription_plan_name?.toLowerCase()
                    )?.price_monthly || 0);

                  return (
                    <div
                      key={plan.id}
                      className={cn(
                        "rounded-lg border-2 p-4 transition-all",
                        isCurrentPlan
                          ? "border-blue-600 bg-blue-50"
                          : "border-slate-200"
                      )}
                    >
                      <h3 className="font-bold text-slate-900">{plan.name}</h3>
                      <p className="text-2xl font-bold text-blue-700 mt-1">
                        {plan.price_monthly === 0
                          ? "Gratis"
                          : `Rp${plan.price_monthly.toLocaleString("id-ID")}`}
                      </p>
                      {plan.price_monthly > 0 && (
                        <p className="text-xs text-slate-500">/bulan</p>
                      )}
                      <p className="text-sm text-slate-600 mt-2">
                        {plan.max_employees === 0
                          ? "Unlimited"
                          : `Max ${plan.max_employees}`}{" "}
                        karyawan
                      </p>

                      {isCurrentPlan ? (
                        <Button
                          className="w-full mt-4"
                          variant="outline"
                          disabled
                        >
                          Paket Saat Ini
                        </Button>
                      ) : isHigher ? (
                        <Button
                          className="w-full mt-4"
                          onClick={() =>
                            requestUpgrade({ plan_slug: plan.slug })
                          }
                          disabled={isUpgrading}
                        >
                          <ArrowUpRight className="w-4 h-4 mr-1" />
                          Upgrade
                        </Button>
                      ) : (
                        <Button
                          className="w-full mt-4"
                          variant="outline"
                          disabled
                        >
                          Tidak Tersedia
                        </Button>
                      )}
                    </div>
                  );
                })}
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}
