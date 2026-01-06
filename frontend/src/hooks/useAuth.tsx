import { useNavigate } from "react-router-dom";
import { useMutation } from "@tanstack/react-query";
import { api } from "@/lib/axios";

interface LoginPayload {
  username: string;
  password: string;
}

interface LoginResponse {
  token: string;
  message: string;
}

export const useLogin = () => {
  const navigate = useNavigate();

  return useMutation({
    mutationFn: async (payload: LoginPayload) => {
      // axios do POST request
      const response = await api.post<LoginResponse>("/auth/login", payload);
      console.log(payload);

      return response.data;
    },

    onSuccess: (data) => {
      // save token
      localStorage.setItem("token", data.token);

      // navigate to dashboard
      navigate("/dashboard");
    },

    onError: (error: any) => {
      console.error("Login error:", error);
    },
  });
};
