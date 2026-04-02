"use client";

import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { Loader2 } from "lucide-react";
import { Link, useNavigate, useLocation } from "react-router-dom";

import { Button } from "@/components/ui/button";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { PasswordInput } from "@/components/ui/password-input";
import { useResetPassword } from "@/features/auth/hooks/useAuth";
import { toast } from "sonner";

const formSchema = z.object({
  password: z.string().min(8, {
    message: "Kata sandi minimal 8 karakter.",
  }),
  confirmPassword: z.string().min(1, {
    message: "Konfirmasi kata sandi diperlukan.",
  }),
}).refine((data) => data.password === data.confirmPassword, {
  path: ["confirmPassword"],
  message: "Kata sandi tidak cocok.",
});

type FormValues = z.infer<typeof formSchema>;

export function ResetPasswordForm() {
  const { mutate: resetPassword, isPending } = useResetPassword();
  const navigate = useNavigate();
  const location = useLocation();
  const code = location.state?.code;

  const form = useForm<FormValues>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      password: "",
      confirmPassword: "",
    },
  });

  const onSubmit = (data: FormValues) => {
    if (!code) {
      toast.error("Error", { description: "Sesi tidak valid." });
      return;
    }
    
    resetPassword({ code, password: data.password }, {
      onSuccess: () => {
        toast.success("Berhasil", {
          description: "Kata sandi Anda telah berhasil direset.",
        });
        navigate("/login");
      },
    });
  };

  if (!code) {
    return (
      <div className="text-center">
        <p className="text-red-500 mb-4">Sesi tidak valid atau telah kadaluarsa.</p>
        <Button asChild>
          <Link to="/forgot-password">Kembali</Link>
        </Button>
      </div>
    );
  }

  return (
    <div className="space-y-6 w-full">
      <div className="flex flex-col space-y-2 text-left">
        <h2 className="text-3xl font-bold tracking-tight text-slate-950">
          Reset Kata Sandi
        </h2>
        <p className="text-sm text-slate-500">
          Silakan buat kata sandi baru untuk akun Anda.
        </p>
      </div>

      <Form {...form}>
        <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
          <div className="space-y-4">
            <FormField
              control={form.control}
              name="password"
              render={({ field }) => (
                <FormItem>
                  <FormLabel className="text-slate-900 font-semibold">
                    Kata Sandi Baru
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
              name="confirmPassword"
              render={({ field }) => (
                <FormItem>
                  <FormLabel className="text-slate-900 font-semibold">
                    Konfirmasi Kata Sandi Baru
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
          </div>

          <Button
            type="submit"
            className="w-full bg-blue-700 hover:bg-blue-800 text-white font-bold py-6 transition-all duration-200"
            disabled={isPending}
          >
            {isPending ? (
              <>
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                Mereset...
              </>
            ) : (
              "Simpan Kata Sandi"
            )}
          </Button>
        </form>
      </Form>
    </div>
  );
}
