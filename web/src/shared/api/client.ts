let csrfToken: string | null = null;

/** CSRF token lives in memory only, never localStorage. */
export function setCSRFToken(token: string | null): void {
  csrfToken = token;
}

export function getCSRFToken(): string | null {
  return csrfToken;
}

export class ApiError extends Error {
  readonly status: number;
  readonly code: string;

  constructor(message: string, status = 0, code = "request_failed") {
    super(message);
    this.name = "ApiError";
    this.status = status;
    this.code = code;
  }
}

type RequestOptions = {
  method?: "GET" | "POST" | "PUT" | "PATCH" | "DELETE";
  body?: unknown;
  signal?: AbortSignal;
};

type ErrorEnvelope = {
  error?: { code?: string; message?: string };
};

/** Typed fetch wrapper for the maulid.id JSON API. */
export async function apiRequest<T>(path: string, options: RequestOptions = {}): Promise<T> {
  const method = options.method ?? "GET";
  const headers: Record<string, string> = { Accept: "application/json" };
  if (options.body !== undefined) {
    headers["Content-Type"] = "application/json";
  }
  if (csrfToken && method !== "GET") {
    headers["X-CSRF-Token"] = csrfToken;
  }

  let response: Response;
  try {
    response = await fetch(path, {
      method,
      headers,
      credentials: "same-origin",
      body: options.body === undefined ? undefined : JSON.stringify(options.body),
      signal: options.signal,
    });
  } catch (error) {
    if (error instanceof DOMException && error.name === "AbortError") {
      throw error;
    }
    throw new ApiError("Tidak dapat menghubungi server. Periksa koneksi lalu coba lagi.");
  }

  const contentType = response.headers.get("Content-Type") ?? "";
  const payload = contentType.includes("application/json") ? await response.json() : null;

  if (!response.ok) {
    const envelope = payload as ErrorEnvelope | null;
    throw new ApiError(
      envelope?.error?.message ?? "Permintaan tidak dapat diproses. Coba lagi.",
      response.status,
      envelope?.error?.code ?? "request_failed",
    );
  }
  if (response.status === 204 || payload === null) {
    return undefined as T;
  }
  return payload as T;
}
