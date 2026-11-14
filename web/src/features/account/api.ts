import { apiRequest } from "@/shared/api/client";
import type { SessionPayload } from "./types";

export function fetchSession(signal?: AbortSignal): Promise<SessionPayload> {
  return apiRequest<SessionPayload>("/api/v1/auth/session", { signal });
}

export function login(email: string, password: string): Promise<SessionPayload> {
  return apiRequest<SessionPayload>("/api/v1/auth/login", {
    method: "POST",
    body: { email, password },
  });
}

export function register(email: string, name: string, password: string): Promise<SessionPayload> {
  return apiRequest<SessionPayload>("/api/v1/auth/register", {
    method: "POST",
    body: { email, name, password },
  });
}

export function logout(): Promise<void> {
  return apiRequest<void>("/api/v1/auth/logout", { method: "POST" });
}
