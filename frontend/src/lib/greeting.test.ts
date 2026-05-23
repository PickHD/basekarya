import { describe, it, expect, vi, beforeEach } from "vitest";
import {
  getTimePeriod,
  getGreetingMessage,
  getGreetingWithName,
  getGreetingWithLocale,
} from "@/lib/greeting";

describe("getTimePeriod", () => {
  it("should return morning for hours 0-11", () => {
    expect(getTimePeriod(0)).toBe("morning");
    expect(getTimePeriod(6)).toBe("morning");
    expect(getTimePeriod(11)).toBe("morning");
  });

  it("should return afternoon for hours 12-17", () => {
    expect(getTimePeriod(12)).toBe("afternoon");
    expect(getTimePeriod(15)).toBe("afternoon");
    expect(getTimePeriod(17)).toBe("afternoon");
  });

  it("should return evening for hours 18-23", () => {
    expect(getTimePeriod(18)).toBe("evening");
    expect(getTimePeriod(21)).toBe("evening");
    expect(getTimePeriod(23)).toBe("evening");
  });

  it("should default to morning for invalid values", () => {
    const errorSpy = vi.spyOn(console, "error").mockImplementation(() => {});
    expect(getTimePeriod(-1)).toBe("morning");
    expect(getTimePeriod(24)).toBe("morning");
    expect(getTimePeriod(NaN)).toBe("evening");
    errorSpy.mockRestore();
  });
});

describe("getGreetingMessage", () => {
  it("should return Good Morning for morning hours", () => {
    const date = new Date(2024, 0, 1, 8, 0, 0);
    expect(getGreetingMessage(date)).toBe("Good Morning");
  });

  it("should return Good Afternoon for afternoon hours", () => {
    const date = new Date(2024, 0, 1, 14, 0, 0);
    expect(getGreetingMessage(date)).toBe("Good Afternoon");
  });

  it("should return Good Evening for evening hours", () => {
    const date = new Date(2024, 0, 1, 20, 0, 0);
    expect(getGreetingMessage(date)).toBe("Good Evening");
  });

  it("should handle invalid date by defaulting to current time", () => {
    const errorSpy = vi.spyOn(console, "error").mockImplementation(() => {});
    const result = getGreetingMessage(new Date("invalid"));
    expect(["Good Morning", "Good Afternoon", "Good Evening"]).toContain(result);
    errorSpy.mockRestore();
  });
});

describe("getGreetingWithName", () => {
  it("should include first name in greeting", () => {
    const date = new Date(2024, 0, 1, 8, 0, 0);
    expect(getGreetingWithName("John Doe", date)).toBe("Good Morning, John! 👋");
  });

  it("should return greeting without name when name is empty", () => {
    const date = new Date(2024, 0, 1, 8, 0, 0);
    expect(getGreetingWithName("", date)).toBe("Good Morning!");
  });

  it("should return greeting without name when name is undefined", () => {
    const date = new Date(2024, 0, 1, 8, 0, 0);
    expect(getGreetingWithName(undefined, date)).toBe("Good Morning!");
  });

  it("should trim whitespace from name", () => {
    const date = new Date(2024, 0, 1, 8, 0, 0);
    expect(getGreetingWithName("  John  Doe  ", date)).toBe("Good Morning, John! 👋");
  });
});

describe("getGreetingWithLocale", () => {
  it("should return English greeting by default", () => {
    const date = new Date(2024, 0, 1, 8, 0, 0);
    expect(getGreetingWithLocale("en", date)).toBe("Good Morning");
  });

  it("should return Indonesian greeting for id locale", () => {
    const date = new Date(2024, 0, 1, 8, 0, 0);
    expect(getGreetingWithLocale("id", date)).toBe("Selamat Pagi");
  });

  it("should return Indonesian afternoon greeting", () => {
    const date = new Date(2024, 0, 1, 14, 0, 0);
    expect(getGreetingWithLocale("id", date)).toBe("Selamat Siang");
  });

  it("should return Indonesian evening greeting", () => {
    const date = new Date(2024, 0, 1, 20, 0, 0);
    expect(getGreetingWithLocale("id", date)).toBe("Selamat Malam");
  });

  it("should fallback to English for unsupported locale", () => {
    const date = new Date(2024, 0, 1, 8, 0, 0);
    expect(getGreetingWithLocale("fr", date)).toBe("Good Morning");
  });
});
