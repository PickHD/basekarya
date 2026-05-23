import { describe, it, expect, beforeEach, vi } from "vitest";
import { renderHook, waitFor } from "@testing-library/react";
import { createQueryWrapper } from "@/test/utils";
import { usePublishAnnouncement } from "@/features/announcement/hooks/useAnnouncement";

vi.mock("@/lib/axios", () => ({
  api: { post: vi.fn() },
}));

vi.mock("sonner", () => ({
  toast: { success: vi.fn(), error: vi.fn() },
}));

import { api } from "@/lib/axios";
import { toast } from "sonner";

describe("usePublishAnnouncement", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should publish announcement", async () => {
    (api.post as ReturnType<typeof vi.fn>).mockResolvedValue({
      data: { message: "ok", data: { id: 1 } },
    });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => usePublishAnnouncement(), { wrapper });

    result.current.mutate({ title: "Test", content: "Hello" });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(api.post).toHaveBeenCalledWith("/announcements/publish", {
      title: "Test",
      content: "Hello",
    });
  });

  it("should show error toast on failure", async () => {
    (api.post as ReturnType<typeof vi.fn>).mockRejectedValue({
      response: { data: { message: "Unauthorized" } },
    });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => usePublishAnnouncement(), { wrapper });

    result.current.mutate({ title: "Test", content: "Hello" });

    await waitFor(() => expect(result.current.isError).toBe(true));
    expect(toast.error).toHaveBeenCalledWith("Gagal", {
      description: "Unauthorized",
    });
  });

  it("should handle string error", async () => {
    (api.post as ReturnType<typeof vi.fn>).mockRejectedValue({
      response: { data: { error: "Some error string" } },
    });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => usePublishAnnouncement(), { wrapper });

    result.current.mutate({ title: "Test", content: "Hello" });

    await waitFor(() => expect(result.current.isError).toBe(true));
    expect(toast.error).toHaveBeenCalled();
  });
});
