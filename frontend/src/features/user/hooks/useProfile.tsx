import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/axios";
import { toast } from "sonner";
import type { UserProfile, PasswordPayload } from "@/features/user/types";

export const useProfile = () => {
  return useQuery({
    queryKey: ["profile"],
    queryFn: async () => {
      const { data } = await api.get<{ data: UserProfile }>("/users/me");

      return data.data;
    },
  });
};

export const useUpdateProfile = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (formData: FormData) => {
      return await api.put("/users/profile", formData);
    },
    onSuccess: () => {
      toast.success("Profile updated successfully");
      queryClient.invalidateQueries({ queryKey: ["profile"] });
    },
    onError: (error: any) => {
      const responseData = error.response?.data;

      let title = "Update Profile Failed";
      let description = responseData?.message || "Failed to update profile";

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

export const useChangePassword = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (payload: PasswordPayload) => {
      return await api.put("/users/change-password", payload);
    },
    onSuccess: async () => {
      toast.success("Password changed successfully");

      await queryClient.invalidateQueries({ queryKey: ["profile"] });
    },
    onError: (error: any) => {
      const responseData = error.response?.data;

      let title = "Change Password Failed";
      let description = responseData?.message || "Failed to change password";

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
