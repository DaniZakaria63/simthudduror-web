import { describe, expect, it } from "vitest";
import { ApiError, apiRequest, getCSRFToken, setCSRFToken } from "./client";
import { jsonResponse } from "@/test/fixtures";

describe("apiRequest", () => {
  it("returns parsed JSON for successful GET", async () => {
    const captured: { headers: Headers | null } = { headers: null };
    globalThis.fetch = async (_input, init) => {
      captured.headers = new Headers(init?.headers);
      return jsonResponse({ ok: true });
    };
    const result = await apiRequest<{ ok: boolean }>("/api/v1/posts");
    expect(result.ok).toBe(true);
    expect(captured.headers?.get("Accept")).toBe("application/json");
    expect(captured.headers?.get("X-CSRF-Token")).toBeNull();
  });

  it("adds the in-memory CSRF token to unsafe requests", async () => {
    setCSRFToken("csrf-123");
    expect(getCSRFToken()).toBe("csrf-123");
    let requestInit: RequestInit | undefined;
    globalThis.fetch = async (_input, init) => {
      requestInit = init;
      return jsonResponse({ created: true });
    };
    await apiRequest("/api/v1/posts", { method: "POST", body: { body: "kabar" } });
    expect(new Headers(requestInit?.headers).get("X-CSRF-Token")).toBe("csrf-123");
    expect(new Headers(requestInit?.headers).get("Content-Type")).toBe("application/json");
    expect(requestInit?.body).toBe('{"body":"kabar"}');
    expect(requestInit?.credentials).toBe("same-origin");
  });

  it("returns undefined for an empty 204 response", async () => {
    globalThis.fetch = async () => new Response(null, { status: 204 });
    await expect(apiRequest("/api/v1/auth/logout", { method: "POST" })).resolves.toBeUndefined();
  });

  it("maps an API error envelope to ApiError", async () => {
    globalThis.fetch = async () =>
      jsonResponse({ error: { code: "validation_required", message: "Perlu validasi." } }, 403);
    await expect(apiRequest("/api/v1/posts", { method: "POST" })).rejects.toMatchObject({
      name: "ApiError",
      status: 403,
      code: "validation_required",
      message: "Perlu validasi.",
    });
  });

  it("uses a safe fallback for non-JSON failures", async () => {
    globalThis.fetch = async () => new Response("bad gateway", { status: 502 });
    await expect(apiRequest("/api/v1/posts")).rejects.toMatchObject({
      status: 502,
      code: "request_failed",
      message: "Permintaan tidak dapat diproses. Coba lagi.",
    });
  });

  it("reports network failures without leaking the cause", async () => {
    globalThis.fetch = async () => {
      throw new TypeError("socket hang up with secret detail");
    };
    const failure = apiRequest("/api/v1/posts");
    await expect(failure).rejects.toBeInstanceOf(ApiError);
    await expect(failure).rejects.not.toThrow(/secret detail/);
  });

  it("rethrows abort errors unchanged", async () => {
    const abort = new DOMException("aborted", "AbortError");
    globalThis.fetch = async () => {
      throw abort;
    };
    await expect(apiRequest("/api/v1/posts")).rejects.toBe(abort);
  });
});
