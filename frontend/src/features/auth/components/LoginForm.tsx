"use client";

import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { Loader2 } from "lucide-react";

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
import { useLogin } from "@/features/auth/hooks/useAuth";
import { toast } from "sonner";
import { useCallback, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { jwtDecode } from "jwt-decode";
import type { DecodedToken } from "../types";

const formSchema = z.object({
  username: z.string().min(1, {
    message: "Username is required",
  }),
  password: z.string().min(1, {
    message: "Password is required.",
  }),
});

type FormValues = z.infer<typeof formSchema>;

export function LoginForm() {
  const { mutate: login, isPending, error } = useLogin();
  const navigate = useNavigate();

  const handleRedirect = useCallback(
    (token: string) => {
      try {
        const decoded = jwtDecode<DecodedToken>(token);

        if (decoded.role === "SUPERADMIN") {
          navigate("/admin/recap", { replace: true });
        } else {
          navigate("/dashboard", { replace: true });
        }
      } catch (error) {
        console.error("Invalid token:", error);
        localStorage.removeItem("token");
      }
    },
    [navigate]
  );

  useEffect(() => {
    const token = localStorage.getItem("token");
    if (token) {
      handleRedirect(token);
    }
  }, [handleRedirect]);

  const form = useForm<FormValues>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      username: "",
      password: "",
    },
  });

  async function onSubmit(data: FormValues) {
    login(data, {
      onSuccess: (response: any) => {
        const token = response.token || response.data?.token;
        const mustChange =
          response.must_change_password || response.data?.must_change_password;

        if (token) {
          localStorage.setItem("token", token);

          if (mustChange) {
            toast.warning("Security Alert", {
              description:
                "Anda harus mengubah kata sandi sebelum melanjutkan, buka profil > buka tab keamanan",
            });
            navigate("/profile");
            return;
          }

          toast.success("Login berhasil", {
            description: "Redirecting to your workspace...",
          });

          handleRedirect(token);
        } else {
          toast.error("Login gagal", { description: "No token received" });
        }
      },
      onError: (err: any) => {
        const msg = err.response?.data?.message || "Login Gagal";
        toast.error("Login gagal", {
          description: msg,
        });
      },
    });
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
        <div className="space-y-4">
          <FormField
            control={form.control}
            name="username"
            render={({ field }) => (
              <FormItem>
                <FormLabel className="text-slate-900 font-semibold">
                  Username atau Employee ID
                </FormLabel>
                <FormControl>
                  <Input
                    placeholder="1293812391293"
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
                  Password
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

          {error && (
            <div className="text-red-500 text-sm font-medium">
              {(error as any).response?.data?.message || "Something when wrong"}
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
              Sedang memverifikasi...
            </>
          ) : (
            "Masuk"
          )}
        </Button>
      </form>
    </Form>
  );
}
