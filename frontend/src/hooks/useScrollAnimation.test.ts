import { describe, it, expect } from "vitest";
import { renderHook } from "@testing-library/react";
import { useScrollAnimation } from "@/hooks/useScrollAnimation";

describe("useScrollAnimation", () => {
  it("should return a ref and isVisible as false initially", () => {
    const { result } = renderHook(() => useScrollAnimation());
    expect(result.current.ref.current).toBeNull();
    expect(result.current.isVisible).toBe(false);
  });

  it("should set isVisible to true when element intersects", () => {
    const { result, rerender } = renderHook(() => useScrollAnimation());

    const mockElement = document.createElement("div");
    result.current.ref.current = mockElement;

    rerender();

    const observerInstances = (globalThis.IntersectionObserver as unknown as typeof IntersectionObserver & { __instances: InstanceType<typeof IntersectionObserver>[] }).__instances;
  });

  it("should accept custom threshold", () => {
    const { result } = renderHook(() => useScrollAnimation(0.5));
    expect(result.current.isVisible).toBe(false);
  });
});
