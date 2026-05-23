import { describe, it, expect, beforeEach } from "vitest";
import { api } from "@/lib/axios";

describe("axios interceptors", () => {
  beforeEach(() => {
    localStorage.clear();
  });

  describe("request interceptor", () => {
    it("should add Authorization header when token exists", async () => {
      localStorage.setItem("token", "test-token-123");

      const config = await api.interceptors.request.handlers[0].fulfilled({
        headers: {} as Record<string, string>,
      });

      expect(config.headers.Authorization).toBe("Bearer test-token-123");
    });

    it("should not add Authorization header when no token", async () => {
      const config = await api.interceptors.request.handlers[0].fulfilled({
        headers: {} as Record<string, string>,
      });

      expect(config.headers.Authorization).toBeUndefined();
    });

    it("should remove Content-Type for FormData", async () => {
      const formData = new FormData();
      const config = await api.interceptors.request.handlers[0].fulfilled({
        headers: { "Content-Type": "application/json" } as Record<string, string>,
        data: formData,
      });

      expect(config.headers["Content-Type"]).toBeUndefined();
    });

    it("should keep Content-Type for non-FormData", async () => {
      const config = await api.interceptors.request.handlers[0].fulfilled({
        headers: { "Content-Type": "application/json" } as Record<string, string>,
        data: { name: "test" },
      });

      expect(config.headers["Content-Type"]).toBe("application/json");
    });
  });

  describe("response interceptor", () => {
    it("should pass through successful responses", async () => {
      const response = { status: 200, data: { message: "ok" } };
      const result = await api.interceptors.response.handlers[0].fulfilled(response);
      expect(result).toEqual(response);
    });

    it("should remove token and redirect on 401", async () => {
      localStorage.setItem("token", "expired-token");
      delete (window as Record<string, unknown>).location;
      window.location = { href: "" } as Location;

      const error = { response: { status: 401 } };

      await expect(
        api.interceptors.response.handlers[0].rejected(error)
      ).rejects.toEqual(error);

      expect(localStorage.getItem("token")).toBeNull();
      expect(window.location.href).toBe("/login");
    });

    it("should not redirect on 401 when no token", async () => {
      delete (window as Record<string, unknown>).location;
      window.location = { href: "" } as Location;

      const error = { response: { status: 401 } };

      await expect(
        api.interceptors.response.handlers[0].rejected(error)
      ).rejects.toEqual(error);

      expect(window.location.href).toBe("");
    });

    it("should pass through non-401 errors", async () => {
      const error = { response: { status: 500 } };

      await expect(
        api.interceptors.response.handlers[0].rejected(error)
      ).rejects.toEqual(error);
    });
  });
});
