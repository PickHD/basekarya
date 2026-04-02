import { useMutation } from "@tanstack/react-query";
import { api } from "@/lib/axios";
import { toast } from "sonner";
import type { CreateAnnouncementPayload, CreateAnnouncementResponse } from "@/features/announcement/types";

export const usePublishAnnouncement = () => {
  return useMutation({
    mutationFn: async (payload: CreateAnnouncementPayload) => {
      const { data } = await api.post<CreateAnnouncementResponse>(
        "/announcements/publish",
        payload
      );
      return data;
    },
    onError: (error: any) => {
      const responseData = error.response?.data;
      let title = "Gagal";
      let description = responseData?.message || "Gagal mempublish pengumuman";

      if (responseData?.error) {
        if (responseData.error.message) {
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
