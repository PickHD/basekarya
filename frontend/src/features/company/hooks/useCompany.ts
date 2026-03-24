import { api } from "@/lib/axios";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import type { CompanyProfile } from "@/features/company/types";
import { toast } from "sonner";

export const useCompanyProfile = () => {
  return useQuery({
    queryKey: ["company-profile"],
    queryFn: async () => {
      const { data } = await api.get<{ data: CompanyProfile }>(
        "/companies/profile",
      );

      return data.data;
    },
  });
};

export const useUpdateCompanyProfile = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (formData: FormData) => {
      return await api.put("/companies/profile", formData);
    },
    onSuccess: () => {
      toast.success("Company profile updated successfully");
      queryClient.invalidateQueries({ queryKey: ["company-profile"] });
    },
    onError: (error: any) => {
      const responseData = error.response?.data;

      let title = "Update company profile failed";
      let description =
        responseData?.message || "Failed to update company profile";

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
