import { useEffect, useRef, useState } from "react";
import type { NotificationPayload } from "../types";
import { toast } from "sonner";
import { useProfile } from "@/features/user/hooks/useProfile";
import {
  useNotifications,
  useMarkAsRead,
} from "@/features/notification/hooks/useNotification";
import { useQueryClient } from "@tanstack/react-query";

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
      } catch (e) {
        console.error("Invalid URL:", e);
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
        console.log("[WS] Connected");
        setIsConnected(true);
        if (reconnectTimeoutRef.current) {
          clearTimeout(reconnectTimeoutRef.current);
          reconnectTimeoutRef.current = null;
        }
      };

      socket.onmessage = (event) => {
        try {
          console.log("Get messages: ", event.data);
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

          const title = payload.title || payload.title || "Notification";
          const message = payload.message || payload.message || "";

          switch (payload.type) {
            case "APPROVED":
              toast.success(title, { description: message, duration: 3000 });
              break;
            case "REJECTED":
              toast.error(title, { description: message, duration: 3000 });
              break;
            case "LEAVE_APPROVAL_REQ":
            case "REIMBURSE_APPROVAL_REQ":
              toast.info(title, { description: message, duration: 3000 });
              break;
            default:
              toast(title, { description: message });
              break;
          }
        } catch (err) {
          console.error("[WS] Parse Error:", err);
        }
      };

      socket.onclose = () => {
        console.log("[WS] Disconnected");
        setIsConnected(false);
        socketRef.current = null;

        reconnectTimeoutRef.current = setTimeout(() => {
          console.log("[WS] Attempting Reconnect...");
          connect();
        }, RECONNECT_INTERVAL);
      };

      socket.onerror = (error) => {
        console.error("[WS] Error:", error);
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
