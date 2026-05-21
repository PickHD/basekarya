"use client";

import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { Loader2, Check } from "lucide-react";
import { useNavigate } from "react-router-dom";
import { toast } from "sonner";

import { Button } from "@/components/ui/button";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { PasswordInput } from "@/components/ui/password-input";
import { useRegister, useSubscriptionPlans } from "@/features/auth/hooks/useAuth";
import { cn } from "@/lib/utils";
import { useSearchParams } from "react-router-dom";

const formSchema = z
  .object({
    company_name: z.string().min(2, { message: "Nama perusahaan minimal 2 karakter" }),
    admin_name: z.string().min(2, { message: "Nama admin minimal 2 karakter" }),
    admin_email: z.string().email({ message: "Email tidak valid" }),
    password: z.string().min(6, { message: "Kata sandi minimal 6 karakter" }),
    confirm_password: z.string().min(1, { message: "Konfirmasi kata sandi diperlukan" }),
    phone_number: z.string().min(8, { message: "Nomor telepon minimal 8 karakter" }),
    plan_slug: z.string().min(1, { message: "Pilih paket terlebih dahulu" }),
  })
  .refine((data) => data.password === data.confirm_password, {
    path: ["confirm_password"],
    message: "Kata sandi tidak cocok",
  });

type FormValues = z.infer<typeof formSchema>;

export function RegisterForm() {
  const { mutate: register, isPending } = useRegister();
  const { data: plansData } = useSubscriptionPlans();
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();

  const plans = plansData?.data || [];

  const form = useForm<FormValues>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      company_name: "",
      admin_name: "",
      admin_email: "",
      password: "",
      confirm_password: "",
      phone_number: "",
      plan_slug: searchParams.get("plan") || "free",
    },
  });

  const selectedPlan = form.watch("plan_slug");

  async function onSubmit(data: FormValues) {
    const payload = {
      company_name: data.company_name,
      admin_name: data.admin_name,
      admin_email: data.admin_email,
      password: data.password,
      phone_number: data.phone_number,
      plan_slug: data.plan_slug,
    };

    register(payload, {
      onSuccess: (response: any) => {
        const username = response.data?.username || response.username;
        if (data.plan_slug !== "free") {
          toast.success("Registrasi berhasil", {
            description:
              "Akun Anda menunggu konfirmasi pembayaran. Tim kami akan menghubungi Anda segera.",
            duration: 8000,
          });
        } else {
          toast.success("Registrasi berhasil", {
            description: `Username Anda: ${username} — Silakan login untuk melanjutkan.`,
            duration: 10000,
          });
        }
        navigate("/login", { replace: true });
      },
    });
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
        <div className="space-y-4">
          <FormField
            control={form.control}
            name="plan_slug"
            render={({ field }) => (
              <FormItem>
                <FormLabel className="text-slate-900 font-semibold">
                  Pilih Paket
                </FormLabel>
                <div className="grid grid-cols-3 gap-3 mt-2">
                  {plans.map((plan: any) => (
                    <button
                      key={plan.slug}
                      type="button"
                      onClick={() => field.onChange(plan.slug)}
                      className={cn(
                        "relative rounded-lg border-2 p-3 text-left transition-all",
                        selectedPlan === plan.slug
                          ? "border-blue-600 bg-blue-50"
                          : "border-slate-200 hover:border-slate-300"
                      )}
                    >
                      {selectedPlan === plan.slug && (
                        <div className="absolute -top-2 -right-2 h-5 w-5 rounded-full bg-blue-600 flex items-center justify-center">
                          <Check className="h-3 w-3 text-white" />
                        </div>
                      )}
                      <p className="font-bold text-sm text-slate-900">
                        {plan.name}
                      </p>
                      <p className="text-xs text-slate-500 mt-0.5">
                        {plan.max_employees === 0
                          ? "Unlimited"
                          : `Max ${plan.max_employees}`}{" "}
                        karyawan
                      </p>
                      <p className="font-bold text-blue-700 mt-1 text-sm">
                        {plan.price_monthly === 0
                          ? "Gratis"
                          : `Rp${plan.price_monthly.toLocaleString("id-ID")}/bln`}
                      </p>
                    </button>
                  ))}
                </div>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="company_name"
            render={({ field }) => (
              <FormItem>
                <FormLabel className="text-slate-900 font-semibold">
                  Nama Perusahaan
                </FormLabel>
                <FormControl>
                  <Input
                    placeholder="PT Contoh Sukses"
                    {...field}
                    className="border-slate-300 focus-visible:ring-blue-600"
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="admin_name"
            render={({ field }) => (
              <FormItem>
                <FormLabel className="text-slate-900 font-semibold">
                  Nama Lengkap Admin
                </FormLabel>
                <FormControl>
                  <Input
                    placeholder="John Doe"
                    {...field}
                    className="border-slate-300 focus-visible:ring-blue-600"
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="admin_email"
            render={({ field }) => (
              <FormItem>
                <FormLabel className="text-slate-900 font-semibold">
                  Email Admin
                </FormLabel>
                <FormControl>
                  <Input
                    type="email"
                    placeholder="admin@contoh.com"
                    {...field}
                    className="border-slate-300 focus-visible:ring-blue-600"
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="phone_number"
            render={({ field }) => (
              <FormItem>
                <FormLabel className="text-slate-900 font-semibold">
                  Nomor Telepon
                </FormLabel>
                <FormControl>
                  <Input
                    placeholder="081234567890"
                    {...field}
                    className="border-slate-300 focus-visible:ring-blue-600"
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="password"
            render={({ field }) => (
              <FormItem>
                <FormLabel className="text-slate-900 font-semibold">
                  Kata Sandi
                </FormLabel>
                <FormControl>
                  <PasswordInput
                    placeholder="••••••••"
                    {...field}
                    className="border-slate-300 focus-visible:ring-blue-600"
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="confirm_password"
            render={({ field }) => (
              <FormItem>
                <FormLabel className="text-slate-900 font-semibold">
                  Konfirmasi Kata Sandi
                </FormLabel>
                <FormControl>
                  <PasswordInput
                    placeholder="••••••••"
                    {...field}
                    className="border-slate-300 focus-visible:ring-blue-600"
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />

          {selectedPlan !== "free" && (
            <div className="rounded-lg border border-amber-200 bg-amber-50 p-3 text-sm text-amber-800">
              Paket berbayar memerlukan konfirmasi pembayaran. Tim kami akan
              menghubungi Anda setelah registrasi.
            </div>
          )}
        </div>

        <Button
          type="submit"
          className="w-full bg-blue-700 hover:bg-blue-800 text-white font-bold py-6 transition-all duration-200"
          disabled={isPending}
        >
          {isPending ? (
            <>
              <Loader2 className="mr-2 h-4 w-4 animate-spin" />
              Mendaftar...
            </>
          ) : (
            "Daftar Sekarang"
          )}
        </Button>
      </form>
    </Form>
  );
}
