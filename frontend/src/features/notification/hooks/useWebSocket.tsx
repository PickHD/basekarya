import { useEffect, useRef, useState } from "react";
import type { NotificationPayload } from "@/features/notification/types";
import { toast } from "sonner";
import { useProfile } from "@/features/user/hooks/useProfile";
import {
  useNotifications,
  useMarkAsRead,
} from "@/features/notification/hooks/useNotification";
import { useQueryClient } from "@tanstack/react-query";
import { sanitizeHtml } from "@/lib/sanitize";

const RECONNECT_INTERVAL = 3000;

export const useWebSocket = () => {
  const { data: user } = useProfile();
  const [isConnected, setIsConnected] = useState(false);
  const { data: notifications = [] } = useNotifications();
  const { mutate: markRead } = useMarkAsRead();

  const queryClient = useQueryClient();
  const unreadCount = notifications.filter((n) => !n.is_read).length;

  const socketRef = useRef<WebSocket | null>(null);
  const reconnectTimeoutRef = useRef<NodeJS.Timeout>(null);

  useEffect(() => {
    if (!user?.id) return;

    const getWebSocketUrl = () => {
      let token = localStorage.getItem("token") || "";
      token = token.replace(/^"|"$/g, "");

      const baseUrl = import.meta.env.VITE_API_URL || "http://localhost:8081";
      const wsProtocol = window.location.protocol === "https:" ? "wss:" : "ws:";

      try {
        const url = new URL(baseUrl);
        url.protocol = wsProtocol;
        url.pathname = "/api/v1/ws";
        url.searchParams.append("token", token);

        return url.toString();
      } catch {
        return "";
      }
    };

    function connect() {
      const wsUrl = getWebSocketUrl();
      if (!wsUrl) return;

      if (socketRef.current?.readyState === WebSocket.OPEN) return;

      const socket = new WebSocket(wsUrl);
      socketRef.current = socket;

      socket.onopen = () => {
        setIsConnected(true);
        if (reconnectTimeoutRef.current) {
          clearTimeout(reconnectTimeoutRef.current);
          reconnectTimeoutRef.current = null;
        }
      };

      socket.onmessage = (event) => {
        try {
          const payload = JSON.parse(event.data) as NotificationPayload;

          if (!payload.type) return;

          queryClient.setQueryData(
            ["notifications"],
            (oldData: NotificationPayload[] | undefined) => {
              const newNotif: NotificationPayload = {
                ...payload,
                type: payload.type,
                title: payload.title || payload.title,
                message: payload.message || payload.message,
                is_read: false,
                created_at: new Date().toISOString(),
                id: payload.id || payload.id || Math.random(),
              };

              return [newNotif, ...(oldData || [])];
            },
          );

          const title = payload.title || "Notification";
          const message = sanitizeHtml(payload.message || "");

          const toastMessage = (
            <div className="flex flex-col w-full">
              <span className="font-semibold">{title}</span>
              <div 
                className="text-muted-foreground mt-1 text-sm prose prose-sm dark:prose-invert line-clamp-3 prose-p:my-0 md:prose-p:my-0" 
                dangerouslySetInnerHTML={{ __html: message }} 
              />
            </div>
          );

          switch (payload.type) {
            case "PAYROLL_PAID":
            case "APPROVED":
              toast.success(toastMessage, { duration: 3000 });
              break;
            case "REJECTED":
              toast.error(toastMessage, { duration: 3000 });
              break;
            case "LEAVE_APPROVAL_REQ":
            case "REIMBURSE_APPROVAL_REQ":
              toast.info(toastMessage, { duration: 3000 });
              break;
            case "ANNOUNCEMENT":
              toast.info(toastMessage, { duration: 3000 });
              break;
            case "CONTRACT_EXPIRING":
              toast.info(toastMessage, { duration: 3000 });
              break;
            default:
              toast(toastMessage);
              break;
          }
        } catch {
          // ignore non-JSON messages
        }
      };

      socket.onclose = () => {
        setIsConnected(false);
        socketRef.current = null;

        reconnectTimeoutRef.current = setTimeout(() => {
          connect();
        }, RECONNECT_INTERVAL);
      };

      socket.onerror = () => {
        socket.close();
      };
    }

    connect();

    return () => {
      if (socketRef.current) {
        socketRef.current.onclose = null;
        socketRef.current.close();
      }
      if (reconnectTimeoutRef.current) {
        clearTimeout(reconnectTimeoutRef.current);
      }
    };
  }, [user?.id, queryClient]);

  const markAsRead = (id: number) => {
    markRead(id);
  };

  return {
    isConnected,
    notifications,
    unreadCount,
    markAsRead,
  };
};
