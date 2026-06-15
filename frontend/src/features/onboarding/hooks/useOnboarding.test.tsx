import { describe, it, expect, beforeEach, vi } from "vitest";
import { renderHook, waitFor } from "@testing-library/react";
import { createQueryWrapper } from "@/test/utils";
import {
  useOnboardingWorkflows,
  useOnboardingWorkflowDetail,
  useCreateWorkflow,
  useCompleteTask,
} from "@/features/onboarding/hooks/useOnboarding";

vi.mock("@/lib/axios", () => ({
  api: { get: vi.fn(), post: vi.fn(), put: vi.fn(), delete: vi.fn() },
}));

vi.mock("sonner", () => ({
  toast: { success: vi.fn(), error: vi.fn() },
}));

import { api } from "@/lib/axios";
import { toast } from "sonner";

describe("useOnboardingWorkflows", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should fetch workflows", async () => {
    const mockData = { data: [{ id: 1, status: "in_progress" }] };
    (api.get as ReturnType<typeof vi.fn>).mockResolvedValue({ data: mockData });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useOnboardingWorkflows(), { wrapper });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data).toEqual(mockData);
  });
});

describe("useOnboardingWorkflowDetail", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should fetch workflow detail", async () => {
    const mockDetail = { id: 1, status: "in_progress", tasks: [] };
    (api.get as ReturnType<typeof vi.fn>).mockResolvedValue({ data: { data: mockDetail } });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useOnboardingWorkflowDetail(1), { wrapper });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data).toEqual(mockDetail);
  });

  it("should not fetch when id is null", () => {
    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useOnboardingWorkflowDetail(null), { wrapper });
    expect(result.current.fetchStatus).toBe("idle");
  });
});

describe("useCreateWorkflow", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should create workflow", async () => {
    (api.post as ReturnType<typeof vi.fn>).mockResolvedValue({ data: { message: "ok" } });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useCreateWorkflow(), { wrapper });

    result.current.mutate({ employee_id: 1, template_id: 1 });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(toast.success).toHaveBeenCalledWith("Onboarding workflow created");
  });
});

describe("useCompleteTask", () => {
  beforeEach(() => vi.clearAllMocks());

  it("should complete task", async () => {
    (api.put as ReturnType<typeof vi.fn>).mockResolvedValue({ data: { message: "ok" } });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useCompleteTask(), { wrapper });

    result.current.mutate({ id: 1, notes: "Done" });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(toast.success).toHaveBeenCalledWith("Task completed!");
    expect(api.put).toHaveBeenCalledWith("/onboarding/tasks/1/complete", { notes: "Done" });
  });

  it("should complete task without notes", async () => {
    (api.put as ReturnType<typeof vi.fn>).mockResolvedValue({ data: { message: "ok" } });

    const { wrapper } = createQueryWrapper();
    const { result } = renderHook(() => useCompleteTask(), { wrapper });

    result.current.mutate({ id: 2 });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(api.put).toHaveBeenCalledWith("/onboarding/tasks/2/complete", { notes: "" });
  });
});
