import { describe, it, expect, beforeEach, vi } from "vitest";
import { renderHook, waitFor } from "@testing-library/react";
import { createQueryWrapper } from "@/test/utils";
import { useNotifications, useMarkAsRead } from "@/features/notification/hooks/useNotification";

vi.mock("@/lib/axios", () => ({
  api: { get: vi.fn(), put: vi.fn() },
}));

import { api } from "@/lib/axios";

describe("useNotifications", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should fetch notifications", async () => {
    const mockNotifs = {
      data: { data: [{ id: 1, title: "Test", is_read: false }] },
    };
    (api.get as ReturnType<typeof vi.fn>).mockResolvedValue(mockNotifs);

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useNotifications(), { wrapper });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data).toEqual(mockNotifs.data.data);
  });
});

describe("useMarkAsRead", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should mark notification as read and update cache", async () => {
    (api.put as ReturnType<typeof vi.fn>).mockResolvedValue({});

    const { wrapper, queryClient } = createQueryWrapper();
    const initialData = [
      { id: 1, title: "Test", is_read: false },
      { id: 2, title: "Other", is_read: false },
    ];

    const { result: notifResult } = renderHook(() => useNotifications(), { wrapper });

    await waitFor(() => expect(notifResult.current.isSuccess).toBe(true));

    queryClient.setQueryData(["notifications"], initialData);

    const { result } = renderHook(() => useMarkAsRead(), { wrapper });
    result.current.mutate(1);

    await waitFor(() => expect(result.current.isSuccess).toBe(true));

    const updated = queryClient.getQueryData(["notifications"]) as Array<{
      id: number;
      is_read: boolean;
    }>;
    expect(updated).toBeDefined();
    expect(updated[0].is_read).toBe(true);
    expect(updated[1].is_read).toBe(false);
  });

  it("should handle empty cache gracefully", async () => {
    (api.put as ReturnType<typeof vi.fn>).mockResolvedValue({});

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useMarkAsRead(), { wrapper });

    result.current.mutate(1);

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
  });
});
