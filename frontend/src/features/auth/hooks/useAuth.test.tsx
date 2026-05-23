import { describe, it, expect, beforeEach, vi } from "vitest";
import { renderHook, waitFor } from "@testing-library/react";
import { createQueryWrapper } from "@/test/utils";
import { useLogin, useLogout, useForgotPassword, useRegister } from "@/features/auth/hooks/useAuth";

vi.mock("@/lib/axios", () => ({
  api: {
    post: vi.fn(),
    get: vi.fn(),
  },
}));

vi.mock("sonner", () => ({
  toast: { success: vi.fn(), error: vi.fn() },
}));

vi.mock("react-router-dom", () => ({
  useNavigate: () => vi.fn(),
}));

import { api } from "@/lib/axios";
import { toast } from "sonner";

describe("useLogin", () => {
  beforeEach(() => {
    vi.clearAllMocks();
    localStorage.clear();
  });

  it("should save token to localStorage on success", async () => {
    const mockToken = "jwt-token-123";
    (api.post as ReturnType<typeof vi.fn>).mockResolvedValue({
      data: { message: "ok", data: { token: mockToken, must_change_password: false } },
    });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useLogin(), { wrapper });

    result.current.mutate({ username: "admin", password: "password123" });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(localStorage.getItem("token")).toBe(mockToken);
  });

  it("should show toast on error", async () => {
    (api.post as ReturnType<typeof vi.fn>).mockRejectedValue({
      response: { data: { message: "Invalid credentials" } },
    });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useLogin(), { wrapper });

    result.current.mutate({ username: "admin", password: "wrong" });

    await waitFor(() => expect(result.current.isError).toBe(true));
    expect(toast.error).toHaveBeenCalled();
  });

  it("should handle validation errors", async () => {
    (api.post as ReturnType<typeof vi.fn>).mockRejectedValue({
      response: {
        data: {
          message: "Validation failed",
          error: {
            errors: [{ message: "Username is required" }, { message: "Password is required" }],
          },
        },
      },
    });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useLogin(), { wrapper });

    result.current.mutate({ username: "", password: "" });

    await waitFor(() => expect(result.current.isError).toBe(true));
    expect(toast.error).toHaveBeenCalledWith("Validation Failed", {
      description: "Username is required, Password is required",
    });
  });
});

describe("useLogout", () => {
  beforeEach(() => {
    vi.clearAllMocks();
    localStorage.clear();
  });

  it("should remove token and show toast", async () => {
    localStorage.setItem("token", "some-token");

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useLogout(), { wrapper });

    result.current.logout();

    expect(localStorage.getItem("token")).toBeNull();
    expect(toast.success).toHaveBeenCalledWith("Logout successful", {
      description: "You have been logged out successfully",
    });
  });
});

describe("useForgotPassword", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("should call forgot password API", async () => {
    (api.post as ReturnType<typeof vi.fn>).mockResolvedValue({
      data: { message: "OTP sent" },
    });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useForgotPassword(), { wrapper });

    result.current.mutate({ email: "test@example.com" });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(api.post).toHaveBeenCalledWith("/auth/forgot-password", {
      email: "test@example.com",
    });
  });
});

describe("useRegister", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("should call register API and return data", async () => {
    (api.post as ReturnType<typeof vi.fn>).mockResolvedValue({
      data: { message: "ok", data: { username: "newuser" } },
    });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useRegister(), { wrapper });

    const payload = {
      company_name: "Test Co",
      admin_name: "Admin",
      admin_email: "admin@test.com",
      password: "password123",
      phone_number: "081234567890",
      plan_slug: "starter",
    };

    result.current.mutate(payload);

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(api.post).toHaveBeenCalledWith("/auth/register", payload);
  });
});
