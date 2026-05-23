import { describe, it, expect, beforeEach, vi } from "vitest";
import { renderHook, waitFor } from "@testing-library/react";
import { createQueryWrapper } from "@/test/utils";
import { useProfile, useUpdateProfile, useChangePassword } from "@/features/user/hooks/useProfile";

vi.mock("@/lib/axios", () => ({
  api: { get: vi.fn(), put: vi.fn() },
}));

vi.mock("sonner", () => ({
  toast: { success: vi.fn(), error: vi.fn() },
}));

import { api } from "@/lib/axios";
import { toast } from "sonner";

const mockProfile = {
  id: 1,
  full_name: "John Doe",
  email: "john@test.com",
  phone: "081234567890",
  avatar_url: null,
  role: "admin",
  department: "Engineering",
};

describe("useProfile", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should fetch user profile", async () => {
    (api.get as ReturnType<typeof vi.fn>).mockResolvedValue({ data: { data: mockProfile } });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useProfile(), { wrapper });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data).toEqual(mockProfile);
    expect(api.get).toHaveBeenCalledWith("/users/me");
  });
});

describe("useUpdateProfile", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should update profile and show toast", async () => {
    (api.put as ReturnType<typeof vi.fn>).mockResolvedValue({ data: { message: "ok" } });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useUpdateProfile(), { wrapper });

    const formData = new FormData();
    formData.append("full_name", "Jane Doe");
    result.current.mutate(formData);

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(toast.success).toHaveBeenCalledWith("Profile updated successfully");
  });

  it("should show error on failure", async () => {
    (api.put as ReturnType<typeof vi.fn>).mockRejectedValue({
      response: { data: { message: "Invalid email" } },
    });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useUpdateProfile(), { wrapper });

    result.current.mutate(new FormData());

    await waitFor(() => expect(result.current.isError).toBe(true));
    expect(toast.error).toHaveBeenCalledWith("Update Profile Failed", {
      description: "Invalid email",
    });
  });

  it("should handle validation errors", async () => {
    (api.put as ReturnType<typeof vi.fn>).mockRejectedValue({
      response: {
        data: {
          error: {
            errors: [{ message: "Email is required" }, { message: "Name is required" }],
          },
        },
      },
    });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useUpdateProfile(), { wrapper });

    result.current.mutate(new FormData());

    await waitFor(() => expect(result.current.isError).toBe(true));
    expect(toast.error).toHaveBeenCalledWith("Validation Failed", {
      description: "Email is required, Name is required",
    });
  });
});

describe("useChangePassword", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should change password and show toast", async () => {
    (api.put as ReturnType<typeof vi.fn>).mockResolvedValue({ data: { message: "ok" } });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useChangePassword(), { wrapper });

    result.current.mutate({ old_password: "old123", new_password: "new456" });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(toast.success).toHaveBeenCalledWith("Password changed successfully");
    expect(api.put).toHaveBeenCalledWith("/users/change-password", {
      old_password: "old123",
      new_password: "new456",
    });
  });

  it("should show error on wrong old password", async () => {
    (api.put as ReturnType<typeof vi.fn>).mockRejectedValue({
      response: { data: { message: "Wrong password" } },
    });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useChangePassword(), { wrapper });

    result.current.mutate({ old_password: "wrong", new_password: "new456" });

    await waitFor(() => expect(result.current.isError).toBe(true));
    expect(toast.error).toHaveBeenCalledWith("Change Password Failed", {
      description: "Wrong password",
    });
  });
});
