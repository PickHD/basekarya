"use client";

import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { Loader2, ArrowLeft } from "lucide-react";
import { Link, useNavigate } from "react-router-dom";

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
import { useForgotPassword } from "@/features/auth/hooks/useAuth";
import { toast } from "sonner";

const formSchema = z.object({
  email: z.string().email({
    message: "Masukkan email yang valid.",
  }),
});

type FormValues = z.infer<typeof formSchema>;

export function ForgotPasswordForm() {
  const { mutate: forgotPassword, isPending } = useForgotPassword();
  const navigate = useNavigate();

  const form = useForm<FormValues>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      email: "",
    },
  });

  function onSubmit(data: FormValues) {
    forgotPassword(data, {
      onSuccess: () => {
        toast.success("Berhasil", {
          description: "Kode verifikasi telah dikirim ke email Anda.",
        });
        navigate("/verify-otp", { state: { email: data.email } });
      },
    });
  }

  return (
    <div className="space-y-6 w-full">
      <div className="flex flex-col space-y-2 text-left">
        <h2 className="text-3xl font-bold tracking-tight text-slate-950">
          Lupa Kata Sandi
        </h2>
        <p className="text-sm text-slate-500">
          Masukkan email yang terdaftar untuk menerima kode verifikasi OTP.
        </p>
      </div>

      <Form {...form}>
        <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
          <div className="space-y-4">
            <FormField
              control={form.control}
              name="email"
              render={({ field }) => (
                <FormItem>
                  <FormLabel className="text-slate-900 font-semibold">
                    Email
                  </FormLabel>
                  <FormControl>
                    <Input
                      placeholder="nama@email.com"
                      type="email"
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
                Mengirim...
              </>
            ) : (
              "Kirim Kode OTP"
            )}
          </Button>

          <Button
            type="button"
            variant="ghost"
            className="w-full text-slate-600 hover:text-slate-900"
            asChild
          >
            <Link to="/login">
              <ArrowLeft className="mr-2 h-4 w-4" />
              Kembali ke Login
            </Link>
          </Button>
        </form>
      </Form>
    </div>
  );
}
