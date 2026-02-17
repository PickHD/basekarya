import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import type { NotificationPayload, NotificationResponse } from "../types";
import { api } from "@/lib/axios";

export const useNotifications = () => {
  return useQuery({
    queryKey: ["notifications"],
    queryFn: async () => {
      const { data } = await api.get<NotificationResponse>("/notifications");

      return data.data;
    },
    placeholderData: (prev) => prev,
    refetchInterval: 60000,
  });
};

export const useMarkAsRead = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (id: number) => {
      await api.put(`/notifications/${id}/read`);
    },
    onSuccess: (_, variables) => {
      queryClient.setQueryData(
        ["notifications"],
        (oldData: NotificationPayload[] | undefined) => {
          if (!oldData) return [];

          return oldData.map((notif) =>
            notif.id === variables ? { ...notif, is_read: true } : notif,
          );
        },
      );
    },
  });
};
