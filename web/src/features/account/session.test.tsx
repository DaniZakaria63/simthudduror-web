import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { useState } from "react";
import { describe, expect, it, vi } from "vitest";
import { getCSRFToken } from "@/shared/api/client";
import {
  accountFixture,
  jsonResponse,
  mockFetch,
} from "@/test/fixtures";
import { AccountProvider, useSession } from "./session";

function Probe() {
  const session = useSession();
  const [error, setError] = useState("");
  return (
    <div>
      <span data-testid="status">{session.status}</span>
      <span data-testid="account">{session.account?.name ?? "-"}</span>
      <span data-testid="error">{session.error ?? "-"}</span>
      <span data-testid="csrf">{session.csrfToken ?? "-"}</span>
      <span data-testid="probe-error">{error}</span>
      <button type="button" onClick={() => void session.reload()}>
        reload
      </button>
      <button
        type="button"
        onClick={() =>
          void session.signIn("rina@example.com", "kata-sandi-aman").catch((caught: Error) =>
            setError(caught.message),
          )
        }
      >
        sign-in
      </button>
      <button
        type="button"
        onClick={() =>
          void session.signUp("rina@example.com", "Rina", "kata-sandi-aman").catch((caught: Error) =>
            setError(caught.message),
          )
        }
      >
        sign-up
      </button>
      <button type="button" onClick={() => void session.signOut()}>
        sign-out
      </button>
    </div>
  );
}

function renderProbe() {
  return render(
    <AccountProvider>
      <Probe />
    </AccountProvider>,
  );
}

describe("AccountProvider", () => {
  it("loads an existing session and exposes the CSRF token", async () => {
    mockFetch(() =>
      jsonResponse({ account: accountFixture, csrf_token: "csrf-session", expires_at: "2026-11-20T00:00:00Z" }),
    );
    renderProbe();
    await waitFor(() => expect(screen.getByTestId("status")).toHaveTextContent("authenticated"));
    expect(screen.getByTestId("account")).toHaveTextContent("Rina");
    expect(screen.getByTestId("csrf")).toHaveTextContent("csrf-session");
    expect(getCSRFToken()).toBe("csrf-session");
  });

  it("treats 401 as anonymous without an error", async () => {
    mockFetch(() => jsonResponse({ error: { code: "unauthenticated", message: "no" } }, 401));
    renderProbe();
    await waitFor(() => expect(screen.getByTestId("status")).toHaveTextContent("anonymous"));
    expect(screen.getByTestId("error")).toHaveTextContent("-");
    expect(getCSRFToken()).toBeNull();
  });

  it("surfaces unexpected session errors and retries", async () => {
    let calls = 0;
    mockFetch(() => {
      calls += 1;
      if (calls === 1) {
        return jsonResponse({ error: { code: "internal_error", message: "Server sibuk." } }, 500);
      }
      return jsonResponse({ account: accountFixture, csrf_token: "csrf-retry", expires_at: "x" });
    });
    renderProbe();
    await waitFor(() => expect(screen.getByTestId("status")).toHaveTextContent("error"));
    expect(screen.getByTestId("error")).toHaveTextContent("Server sibuk.");
    fireEvent.click(screen.getByRole("button", { name: "reload" }));
    await waitFor(() => expect(screen.getByTestId("status")).toHaveTextContent("authenticated"));
    expect(screen.getByTestId("csrf")).toHaveTextContent("csrf-retry");
  });

  it("signs in and signs out with server-managed state", async () => {
    const fetchMock = mockFetch(({ url, method, headers }) => {
      if (url.endsWith("/api/v1/auth/session")) {
        return jsonResponse({ error: { code: "unauthenticated", message: "no" } }, 401);
      }
      if (url.endsWith("/api/v1/auth/login") && method === "POST") {
        return jsonResponse({ account: accountFixture, csrf_token: "csrf-login", expires_at: "x" });
      }
      if (url.endsWith("/api/v1/auth/logout") && method === "POST") {
        expect(headers.get("X-CSRF-Token")).toBe("csrf-login");
        return new Response(null, { status: 204 });
      }
      return jsonResponse({}, 404);
    });
    renderProbe();
    await waitFor(() => expect(screen.getByTestId("status")).toHaveTextContent("anonymous"));

    fireEvent.click(screen.getByRole("button", { name: "sign-in" }));
    await waitFor(() => expect(screen.getByTestId("status")).toHaveTextContent("authenticated"));
    expect(screen.getByTestId("csrf")).toHaveTextContent("csrf-login");

    fireEvent.click(screen.getByRole("button", { name: "sign-out" }));
    await waitFor(() => expect(screen.getByTestId("status")).toHaveTextContent("anonymous"));
    expect(getCSRFToken()).toBeNull();
    expect(fetchMock).toHaveBeenCalled();
  });

  it("registers and reports sign-in failures", async () => {
    mockFetch(({ url }) => {
      if (url.endsWith("/api/v1/auth/session")) {
        return jsonResponse({ error: { code: "unauthenticated", message: "no" } }, 401);
      }
      if (url.endsWith("/api/v1/auth/register")) {
        return jsonResponse({ account: accountFixture, csrf_token: "csrf-register", expires_at: "x" }, 201);
      }
      return jsonResponse({ error: { code: "invalid_credentials", message: "Email atau kata sandi salah." } }, 401);
    });
    renderProbe();
    await waitFor(() => expect(screen.getByTestId("status")).toHaveTextContent("anonymous"));

    fireEvent.click(screen.getByRole("button", { name: "sign-up" }));
    await waitFor(() => expect(screen.getByTestId("status")).toHaveTextContent("authenticated"));
    expect(screen.getByTestId("csrf")).toHaveTextContent("csrf-register");

    fireEvent.click(screen.getByRole("button", { name: "sign-out" }));
    await waitFor(() => expect(screen.getByTestId("status")).toHaveTextContent("anonymous"));
    fireEvent.click(screen.getByRole("button", { name: "sign-in" }));
    await waitFor(() =>
      expect(screen.getByTestId("probe-error")).toHaveTextContent("Email atau kata sandi salah."),
    );
  });

  it("rejects useSession outside the provider", () => {
    const consoleError = vi.spyOn(console, "error").mockImplementation(() => {});
    expect(() => render(<Probe />)).toThrow("useSession harus dipakai di dalam AccountProvider.");
    consoleError.mockRestore();
  });
});
