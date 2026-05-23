import { describe, it, expect, beforeEach, vi } from "vitest";
import { renderHook, waitFor } from "@testing-library/react";
import { createQueryWrapper } from "@/test/utils";
import { useCompanyProfile, useUpdateCompanyProfile } from "@/features/company/hooks/useCompany";

vi.mock("@/lib/axios", () => ({
  api: {
    get: vi.fn(),
    put: vi.fn(),
  },
}));

vi.mock("sonner", () => ({
  toast: { success: vi.fn(), error: vi.fn() },
}));

import { api } from "@/lib/axios";
import { toast } from "sonner";

const mockProfile = {
  id: 1,
  name: "Test Company",
  email: "test@company.com",
  phone: "081234567890",
  address: "Jakarta",
  logo_url: null,
  industry: null,
  description: null,
  employee_count: 10,
  subscription: null,
};

describe("useCompanyProfile", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("should fetch company profile", async () => {
    (api.get as ReturnType<typeof vi.fn>).mockResolvedValue({
      data: { data: mockProfile },
    });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useCompanyProfile(), { wrapper });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data).toEqual(mockProfile);
    expect(api.get).toHaveBeenCalledWith("/companies/profile");
  });

  it("should handle fetch error", async () => {
    (api.get as ReturnType<typeof vi.fn>).mockRejectedValue(new Error("Network error"));

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useCompanyProfile(), { wrapper });

    await waitFor(() => expect(result.current.isError).toBe(true));
  });
});

describe("useUpdateCompanyProfile", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("should update company profile and show success toast", async () => {
    (api.put as ReturnType<typeof vi.fn>).mockResolvedValue({ data: { message: "ok" } });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useUpdateCompanyProfile(), { wrapper });

    const formData = new FormData();
    formData.append("name", "Updated Company");
    result.current.mutate(formData);

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(toast.success).toHaveBeenCalledWith("Company profile updated successfully");
  });

  it("should show error toast on failure", async () => {
    (api.put as ReturnType<typeof vi.fn>).mockRejectedValue({
      response: { data: { message: "Update failed" } },
    });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useUpdateCompanyProfile(), { wrapper });

    result.current.mutate(new FormData());

    await waitFor(() => expect(result.current.isError).toBe(true));
    expect(toast.error).toHaveBeenCalled();
  });
});
