"use client";

import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { Loader2, ArrowLeft, RefreshCw } from "lucide-react";
import { Link, useNavigate, useLocation } from "react-router-dom";
import { useState, useRef, useEffect } from "react";

import { Button } from "@/components/ui/button";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { useVerifyOTP, useForgotPassword } from "@/features/auth/hooks/useAuth";
import { toast } from "sonner";

const formSchema = z.object({
  code: z.string().length(6, {
    message: "OTP harus terdiri dari 6 digit.",
  }),
});

type FormValues = z.infer<typeof formSchema>;

export function VerifyOTPForm() {
  const { mutate: verifyOTP, isPending } = useVerifyOTP();
  const { mutate: resendOTP, isPending: isResending } = useForgotPassword();
  const navigate = useNavigate();
  const location = useLocation();
  const email = location.state?.email;

  const [countdown, setCountdown] = useState(0);
  const inputsRef = useRef<(HTMLInputElement | null)[]>([]);

  useEffect(() => {
    let timer: NodeJS.Timeout;
    if (countdown > 0) {
      timer = setTimeout(() => setCountdown(countdown - 1), 1000);
    }
    return () => clearTimeout(timer);
  }, [countdown]);

  const form = useForm<FormValues>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      code: "",
    },
  });

  const onSubmit = (data: FormValues) => {
    verifyOTP(data, {
      onSuccess: (response) => {
        if (response.data.is_valid) {
          toast.success("Berhasil", {
            description: "Kode OTP valid.",
          });
          navigate("/reset-password", { state: { email, code: data.code } });
        } else {
          toast.error("Gagal", {
            description: "Kode OTP tidak valid atau sudah kadaluarsa.",
          });
        }
      },
    });
  };

  const handleResend = () => {
    if (!email) {
      toast.error("Error", { description: "Email tidak ditemukan." });
      return;
    }

    resendOTP({ email }, {
      onSuccess: () => {
        toast.success("Berhasil", {
          description: "Kode OTP baru telah dikirim ke email Anda.",
        });
        setCountdown(60);
      }
    });
  };

  const onInputChange = (index: number, value: string) => {
    const currentCode = form.getValues("code");
    const codeArray = currentCode.split("");
    codeArray[index] = value;
    const newCode = codeArray.join("").slice(0, 6);
    
    // update form value
    form.setValue("code", newCode.padEnd(6, " ").replace(/ /g, ""), { shouldValidate: true });

    // Focus next input or blur
    if (value && index < 5) {
      inputsRef.current[index + 1]?.focus();
    }
  };

  const onInputKeyDown = (index: number, e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === "Backspace" && !e.currentTarget.value && index > 0) {
      inputsRef.current[index - 1]?.focus();
    }
  };
  
  const handlePaste = (e: React.ClipboardEvent<HTMLDivElement>) => {
      e.preventDefault();
      const pastedData = e.clipboardData.getData('text/plain').trim().slice(0, 6);
      
      if (/^\d+$/.test(pastedData)) {
          form.setValue("code", pastedData, { shouldValidate: true });
          
          const nextFocusIndex = Math.min(pastedData.length, 5);
          inputsRef.current[nextFocusIndex]?.focus();
      }
  };


  if (!email) {
    return (
      <div className="text-center">
        <p className="text-red-500 mb-4">Sesi tidak valid atau telah kadaluarsa.</p>
        <Button asChild>
          <Link to="/forgot-password">Kembali</Link>
        </Button>
      </div>
    );
  }

  // Helper to mask email (e.g., test@example.com -> t**t@example.com) // simple mask
  const maskedEmail = email.replace(/(.{2})(.*)(?=@)/,
    (_gp1: string, gp2: string, gp3: string) => {
      return gp2 + gp3.replace(/./g, '*');
    });

  const otpValue = form.watch("code");

  return (
    <div className="space-y-6 w-full">
      <div className="flex flex-col space-y-2 text-left">
        <h2 className="text-3xl font-bold tracking-tight text-slate-950">
          Verifikasi OTP
        </h2>
        <p className="text-sm text-slate-500">
          Masukkan 6 digit kode yang telah kami kirimkan ke email <span className="font-semibold text-slate-900">{maskedEmail}</span>
        </p>
      </div>

      <Form {...form}>
        <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
          <div className="space-y-4">
            <FormField
              control={form.control}
              name="code"
              render={() => (
                <FormItem>
                  <FormControl>
                    <div className="flex justify-between max-w-sm mx-auto sm:max-w-none" onPaste={handlePaste}>
                      {[0, 1, 2, 3, 4, 5].map((index) => (
                        <Input
                          key={index}
                          ref={(el: HTMLInputElement | null) => { inputsRef.current[index] = el; }}
                          type="text"
                          inputMode="numeric"
                          pattern="\d*"
                          maxLength={1}
                          className="w-12 h-14 text-center text-xl font-bold rounded-lg border-slate-300 focus-visible:ring-blue-600 sm:w-14"
                          value={otpValue[index] || ""}
                          onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
                            const val = e.target.value.replace(/\D/g, "");
                            onInputChange(index, val);
                          }}
                          onKeyDown={(e: React.KeyboardEvent<HTMLInputElement>) => onInputKeyDown(index, e)}
                        />
                      ))}
                    </div>
                  </FormControl>
                  <div className="text-center">
                     <FormMessage />
                  </div>
                </FormItem>
              )}
            />
          </div>

          <div className="flex justify-center text-sm">
            <span className="text-slate-500 mr-2">Tidak menerima kode?</span>
            <button
              type="button"
              onClick={handleResend}
              disabled={countdown > 0 || isResending}
              className="text-blue-700 font-semibold hover:text-blue-800 disabled:text-slate-400 disabled:cursor-not-allowed flex items-center"
            >
              {isResending ? (
                 <Loader2 className="mr-2 h-3 w-3 animate-spin" />
              ) : countdown > 0 ? (
                `Kirim ulang dalam ${countdown}s`
              ) : (
                <>
                   <RefreshCw className="mr-1 h-3 w-3" />
                   Kirim Ulang
                </>
              )}
            </button>
          </div>

          <Button
            type="submit"
            className="w-full bg-blue-700 hover:bg-blue-800 text-white font-bold py-6 transition-all duration-200"
            disabled={isPending || otpValue.length !== 6}
          >
            {isPending ? (
              <>
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                Memverifikasi...
              </>
            ) : (
              "Verifikasi"
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
