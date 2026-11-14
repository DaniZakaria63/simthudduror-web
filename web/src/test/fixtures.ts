import { vi } from "vitest";

export type MockRequest = {
  url: string;
  method: string;
  headers: Headers;
  body: unknown;
};

type Handler = (request: MockRequest) => Response | Promise<Response>;

export function jsonResponse(data: unknown, status = 200): Response {
  return new Response(JSON.stringify(data), {
    status,
    headers: { "Content-Type": "application/json" },
  });
}

export function emptyResponse(status = 204): Response {
  return new Response(null, { status });
}

export function mockFetch(handler: Handler) {
  const mock = vi.fn(async (input: RequestInfo | URL, init?: RequestInit) => {
    const url = typeof input === "string" ? input : input.toString();
    const method = init?.method ?? "GET";
    const headers = new Headers(init?.headers);
    const body = typeof init?.body === "string" ? JSON.parse(init.body) : undefined;
    return handler({ url, method, headers, body });
  });
  vi.stubGlobal("fetch", mock);
  return mock;
}

export const accountFixture = {
  id: "acc-1",
  email: "rina@example.com",
  name: "Rina",
  role: "member" as const,
  validation_status: "unvalidated" as const,
  created_at: "2026-11-14T08:00:00.000000000Z",
};

export const staffFixture = {
  ...accountFixture,
  id: "staff-1",
  email: "staf@example.com",
  name: "Staf",
  role: "staff" as const,
  validation_status: "validated" as const,
};

export const postFixture = {
  id: "post-1",
  author_id: "acc-2",
  author_name: "Bima",
  body: "Kajian Maulid di masjid kampung berlangsung selepas magrib.",
  created_at: "2026-11-13T12:00:00.000000000Z",
  like_count: 0,
  liked_by_me: false,
};
