import { useNavigate } from "react-router-dom";
import { useMutation } from "@tanstack/react-query";
import { api } from "@/lib/axios";
import { toast } from "sonner";
import type { LoginPayload, LoginResponse } from "../types";

export const useLogin = () => {
  const navigate = useNavigate();

  return useMutation({
    mutationFn: async (payload: LoginPayload) => {
      // axios do POST request
      const response = await api.post<LoginResponse>("/auth/login", payload);
      return response.data;
    },

    onSuccess: (data) => {
      // save token - the token is nested in data.data.token
      localStorage.setItem("token", data.data.token);

      // navigate to dashboard
      navigate("/dashboard");
    },

    onError: (error: any) => {
      console.error("Login error:", error);
    },
  });
};

export const useLogout = () => {
  const navigate = useNavigate();

  const logout = () => {
    // remove token from localStorage
    localStorage.removeItem("token");

    // show success toast
    toast.success("Logout successful", {
      description: "You have been logged out successfully",
    });

    // redirect to login page
    navigate("/login");
  };

  return { logout };
};
