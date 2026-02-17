import * as Popover from "@radix-ui/react-popover";
import { BellIcon } from "@radix-ui/react-icons";
import { clsx } from "clsx";
import { useWebSocket } from "@/features/notification/hooks/useWebSocket"; // Sesuaikan path import

export const NotificationBell = () => {
  const { isConnected, notifications, unreadCount, markAsRead } =
    useWebSocket();

  const formatDate = (dateString?: string) => {
    if (!dateString) return "";
    return new Date(dateString).toLocaleString("id-ID", {
      day: "numeric",
      month: "short",
      hour: "2-digit",
      minute: "2-digit",
    });
  };

  return (
    <Popover.Root>
      <Popover.Trigger asChild>
        <button
          className="relative p-2 rounded-full hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors focus:outline-none focus:ring-2 focus:ring-blue-500"
          aria-label="Notifications"
        >
          <BellIcon
            className={clsx(
              "w-6 h-6 transition-colors",
              isConnected
                ? "text-gray-700 dark:text-gray-200"
                : "text-gray-400",
            )}
          />

          {/* Status Dot (Online/Offline) */}
          <span
            className={clsx(
              "absolute top-2 right-2 w-2 h-2 rounded-full border border-white dark:border-gray-900",
              isConnected ? "bg-green-500" : "bg-red-500",
            )}
          />

          {/* Badge Counter */}
          {unreadCount > 0 && (
            <span className="absolute -top-1 -right-1 flex h-5 w-5 items-center justify-center rounded-full bg-red-600 text-[10px] font-bold text-white shadow-sm">
              {unreadCount > 9 ? "9+" : unreadCount}
            </span>
          )}
        </button>
      </Popover.Trigger>

      <Popover.Portal>
        <Popover.Content
          className="z-50 w-80 sm:w-96 rounded-lg border border-gray-200 bg-white shadow-xl dark:border-gray-700 dark:bg-gray-900 animate-in fade-in zoom-in-95 duration-200"
          sideOffset={8}
          align="end"
        >
          {/* Header */}
          <div className="flex items-center justify-between border-b px-4 py-3 dark:border-gray-700 bg-gray-50/50 dark:bg-gray-800/50">
            <h4 className="font-semibold text-sm text-gray-900 dark:text-gray-100">
              Notifications
            </h4>
            <span className="text-xs text-gray-500 dark:text-gray-400">
              {unreadCount} unread
            </span>
          </div>

          {/* List Content */}
          <div className="max-h-[400px] overflow-y-auto">
            {notifications.length === 0 ? (
              <div className="flex flex-col items-center justify-center py-10 text-gray-500">
                <BellIcon className="w-10 h-10 mb-3 opacity-20" />
                <p className="text-sm">No new notifications</p>
              </div>
            ) : (
              <ul className="divide-y divide-gray-100 dark:divide-gray-800">
                {notifications.map((notif) => {
                  const title = notif.title || notif.title;
                  const message = notif.message || notif.message;
                  const isRead = notif.is_read;

                  return (
                    <li
                      key={notif.id}
                      onClick={() => {
                        if (!isRead) markAsRead(notif.id);
                      }}
                      className={clsx(
                        "p-4 transition-colors cursor-pointer group relative",
                        isRead
                          ? "bg-white hover:bg-gray-50 dark:bg-gray-900 dark:hover:bg-gray-800"
                          : "bg-blue-50/50 hover:bg-blue-50 dark:bg-blue-900/10 dark:hover:bg-blue-900/20",
                      )}
                    >
                      <div className="flex gap-3 items-start">
                        <div
                          className={clsx(
                            "mt-1.5 h-2 w-2 flex-shrink-0 rounded-full transition-colors",
                            isRead ? "bg-transparent" : "bg-blue-600",
                          )}
                        />

                        <div className="flex-1 space-y-1">
                          <p
                            className={clsx(
                              "text-sm",
                              isRead
                                ? "font-medium text-gray-700 dark:text-gray-300"
                                : "font-semibold text-gray-900 dark:text-white",
                            )}
                          >
                            {title}
                          </p>
                          <p
                            className={clsx(
                              "text-xs leading-relaxed line-clamp-2",
                              isRead
                                ? "text-gray-500 dark:text-gray-500"
                                : "text-gray-600 dark:text-gray-400",
                            )}
                          >
                            {message}
                          </p>

                          {notif.created_at && (
                            <p className="text-[10px] text-gray-400 pt-1">
                              {formatDate(notif.created_at)}
                            </p>
                          )}
                        </div>
                      </div>
                    </li>
                  );
                })}
              </ul>
            )}
          </div>
          <Popover.Arrow className="fill-gray-50/50 dark:fill-gray-800/50 border-gray-200" />
        </Popover.Content>
      </Popover.Portal>
    </Popover.Root>
  );
};
