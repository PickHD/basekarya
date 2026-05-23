import { describe, it, expect } from "vitest";
import { renderHook, act } from "@testing-library/react";
import { useDebounce } from "@/hooks/useDebounce";

describe("useDebounce", () => {
  it("should return initial value immediately", () => {
    const { result } = renderHook(() => useDebounce("hello", 500));
    expect(result.current).toBe("hello");
  });

  it("should return debounced value after delay", () => {
    vi.useFakeTimers();
    const { result, rerender } = renderHook(
      ({ value, delay }) => useDebounce(value, delay),
      { initialProps: { value: "hello", delay: 500 } }
    );

    rerender({ value: "world", delay: 500 });
    expect(result.current).toBe("hello");

    act(() => {
      vi.advanceTimersByTime(500);
    });
    expect(result.current).toBe("world");

    vi.useRealTimers();
  });

  it("should use default delay of 500ms", () => {
    vi.useFakeTimers();
    const { result, rerender } = renderHook(
      ({ value }) => useDebounce(value),
      { initialProps: { value: "initial" } }
    );

    rerender({ value: "updated" });
    expect(result.current).toBe("initial");

    act(() => {
      vi.advanceTimersByTime(500);
    });
    expect(result.current).toBe("updated");

    vi.useRealTimers();
  });

  it("should debounce rapid changes and only use last value", () => {
    vi.useFakeTimers();
    const { result, rerender } = renderHook(
      ({ value, delay }) => useDebounce(value, delay),
      { initialProps: { value: "a", delay: 300 } }
    );

    rerender({ value: "b", delay: 300 });
    act(() => {
      vi.advanceTimersByTime(150);
    });

    rerender({ value: "c", delay: 300 });
    expect(result.current).toBe("a");

    act(() => {
      vi.advanceTimersByTime(300);
    });
    expect(result.current).toBe("c");

    vi.useRealTimers();
  });
});
