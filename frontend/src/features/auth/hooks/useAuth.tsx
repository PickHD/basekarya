import { useNavigate } from "react-router-dom";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/axios";
import { toast } from "sonner";
import type {
  LoginPayload,
  LoginResponse,
  ForgotPasswordPayload,
  VerifyOTPPayload,
  VerifyOTPResponse,
  ResetPasswordPayload,
  RegisterPayload,
  RegisterResponse,
} from "@/features/auth/types";

export const useLogin = () => {
  return useMutation({
    mutationFn: async (payload: LoginPayload) => {
      // axios do POST request
      const { data } = await api.post<LoginResponse>("/auth/login", payload);
      return data;
    },

    onSuccess: (data) => {
      // save token - the token is nested in data.data.token
      localStorage.setItem("token", data.data.token);
    },

    onError: (error: any) => {
      const responseData = error.response?.data;

      let title = "Login Failed";
      let description = responseData?.message || "Failed to login";

      if (responseData?.error) {
        if (
          responseData.error.errors &&
          Array.isArray(responseData.error.errors)
        ) {
          title = "Validation Failed";
          description = responseData.error.errors
            .map((err: any) => err.message)
            .join(", ");
        } else if (responseData.error.message) {
          description = responseData.error.message;
        } else if (typeof responseData.error === "string") {
          description = responseData.error;
        }
      }

      toast.error(title, {
        description: description,
      });
    },
  });
};

export const useLogout = () => {
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  const logout = () => {
    // remove token from localStorage
    localStorage.removeItem("token");

    // clear cache
    queryClient.removeQueries();
    queryClient.clear();

    // show success toast
    toast.success("Logout successful", {
      description: "You have been logged out successfully",
    });

    // redirect to login page
    navigate("/login");
  };

  return { logout };
};

export const useForgotPassword = () => {
  return useMutation({
    mutationFn: async (payload: ForgotPasswordPayload) => {
      const { data } = await api.post("/auth/forgot-password", payload);
      return data;
    },
    onError: (error: any) => {
      const responseData = error.response?.data;
      let description = responseData?.message || "Gagal mengirim OTP";

      if (responseData?.error) {
         if (responseData.error.message) {
          description = responseData.error.message;
        } else if (typeof responseData.error === "string") {
          description = responseData.error;
        }
      }

      toast.error("Gagal", {
        description: description,
      });
    },
  });
};

export const useVerifyOTP = () => {
  return useMutation({
    mutationFn: async (payload: VerifyOTPPayload) => {
      const { data } = await api.post<VerifyOTPResponse>("/auth/verify-otp", payload);
      return data;
    },
    onError: (error: any) => {
      const responseData = error.response?.data;
      let description = responseData?.message || "Gagal verifikasi OTP";

      if (responseData?.error) {
         if (responseData.error.message) {
          description = responseData.error.message;
        } else if (typeof responseData.error === "string") {
          description = responseData.error;
        }
      }

      toast.error("Gagal", {
        description: description,
      });
    },
  });
};

export const useResetPassword = () => {
  return useMutation({
    mutationFn: async (payload: ResetPasswordPayload) => {
      const { data } = await api.post("/auth/reset-password", payload);
      return data;
    },
    onError: (error: any) => {
      const responseData = error.response?.data;
      let description = responseData?.message || "Gagal mereset kata sandi";

      if (responseData?.error) {
         if (responseData.error.message) {
          description = responseData.error.message;
        } else if (typeof responseData.error === "string") {
          description = responseData.error;
        }
      }

      toast.error("Gagal", {
        description: description,
      });
    },
  });
};

export const useRegister = () => {
  return useMutation({
    mutationFn: async (payload: RegisterPayload) => {
      const { data } = await api.post<RegisterResponse>("/auth/register", payload);
      return data;
    },
    onError: (error: any) => {
      const responseData = error.response?.data;
      let title = "Registrasi Gagal";
      let description = responseData?.message || "Gagal mendaftar";

      if (responseData?.error) {
        if (responseData.error.errors && Array.isArray(responseData.error.errors)) {
          title = "Validation Failed";
          description = responseData.error.errors
            .map((err: any) => err.message)
            .join(", ");
        } else if (responseData.error.message) {
          description = responseData.error.message;
        } else if (typeof responseData.error === "string") {
          description = responseData.error;
        }
      }

      toast.error(title, { description });
    },
  });
};

export const useSubscriptionPlans = () => {
  return useQuery({
    queryKey: ["subscription-plans"],
    queryFn: async () => {
      const { data } = await api.get("/subscriptions/plans");
      return data;
    },
    staleTime: 1000 * 60 * 30,
  });
};
