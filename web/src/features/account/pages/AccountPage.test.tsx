import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { MemoryRouter, Route, Routes } from "react-router-dom";
import { describe, expect, it } from "vitest";
import {
  accountFixture,
  jsonResponse,
  mockFetch,
  postFixture,
  staffFixture,
} from "@/test/fixtures";
import { AccountProvider } from "../session";
import type { Account } from "../types";
import { AccountPage } from "./AccountPage";

function renderAccount() {
  return render(
    <MemoryRouter initialEntries={["/akun"]}>
      <AccountProvider>
        <Routes>
          <Route path="/akun" element={<AccountPage />} />
          <Route path="/" element={<p>beranda</p>} />
        </Routes>
      </AccountProvider>
    </MemoryRouter>,
  );
}

function sessionResponse(account: Account = accountFixture) {
  return jsonResponse({ account, csrf_token: "csrf", expires_at: "2026-11-20T00:00:00Z" });
}

describe("AccountPage", () => {
  it("shows an anonymous call to action", async () => {
    mockFetch(({ url }) =>
      url.includes("/auth/session")
        ? jsonResponse({ error: { code: "unauthenticated", message: "no" } }, 401)
        : jsonResponse({ posts: [] }),
    );
    renderAccount();
    expect(screen.getByText("Memeriksa sesi…")).toBeInTheDocument();
    expect(await screen.findByRole("heading", { name: "Akun komunitas" })).toBeInTheDocument();
    expect(screen.getByRole("link", { name: "Masuk" })).toHaveAttribute("href", "/masuk");
    expect(screen.getByRole("link", { name: "Daftar akun" })).toHaveAttribute("href", "/daftar");
  });

  it("surfaces session errors and retries", async () => {
    let sessions = 0;
    mockFetch(({ url }) => {
      if (url.includes("/auth/session")) {
        sessions += 1;
        if (sessions === 1) {
          return jsonResponse({ error: { code: "internal_error", message: "Server sibuk." } }, 500);
        }
        return sessionResponse();
      }
      return jsonResponse({ posts: [] });
    });
    const user = userEvent.setup();
    renderAccount();
    expect(await screen.findByRole("alert")).toHaveTextContent("Server sibuk.");
    await user.click(screen.getByRole("button", { name: "Coba lagi" }));
    expect(await screen.findByText("Belum divalidasi")).toBeInTheDocument();
  });

  it("shows an unvalidated account without the composer", async () => {
    mockFetch(({ url, method }) => {
      if (url.includes("/auth/session")) return sessionResponse();
      if (url.includes("/api/v1/posts") && method === "GET") {
        return jsonResponse({ posts: [postFixture] });
      }
      return jsonResponse({}, 404);
    });
    renderAccount();
    expect(await screen.findByText("Belum divalidasi")).toBeInTheDocument();
    expect(screen.getByText(/tetap dapat membaca dan menyukai/i)).toBeInTheDocument();
    expect(screen.queryByRole("heading", { name: "Tulis kabar" })).not.toBeInTheDocument();
    expect(
      await screen.findByText("Kajian Maulid di masjid kampung berlangsung selepas magrib."),
    ).toBeInTheDocument();
  });

  it("lets an unvalidated account like a post and reports the result", async () => {
    mockFetch(({ url, method }) => {
      if (url.includes("/auth/session")) return sessionResponse();
      if (url.includes("/like") && method === "PUT") {
        return jsonResponse({ post: { ...postFixture, like_count: 1, liked_by_me: true } });
      }
      if (url.includes("/api/v1/posts")) return jsonResponse({ posts: [postFixture] });
      return jsonResponse({}, 404);
    });
    const user = userEvent.setup();
    renderAccount();
    await screen.findByText("Belum divalidasi");
    await user.click(screen.getByRole("button", { name: "Suka (0)" }));
    expect(await screen.findByText("Kabar disukai.")).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "Batal suka (1)" })).toBeInTheDocument();
  });

  it("reports a failed interaction", async () => {
    mockFetch(({ url, method }) => {
      if (url.includes("/auth/session")) return sessionResponse();
      if (url.includes("/like") && method === "PUT") {
        return jsonResponse({ error: { code: "internal_error", message: "Gagal menyimpan." } }, 500);
      }
      if (url.includes("/api/v1/posts")) return jsonResponse({ posts: [postFixture] });
      return jsonResponse({}, 404);
    });
    const user = userEvent.setup();
    renderAccount();
    await screen.findByText("Belum divalidasi");
    await user.click(screen.getByRole("button", { name: "Suka (0)" }));
    expect(await screen.findByText("Gagal menyimpan.")).toBeInTheDocument();
  });

  it("shows the composer for validated accounts and signs out", async () => {
    mockFetch(({ url }) => {
      if (url.includes("/auth/session")) {
        return sessionResponse({ ...accountFixture, validation_status: "validated" });
      }
      if (url.includes("/auth/logout")) return new Response(null, { status: 204 });
      if (url.includes("/api/v1/posts")) return jsonResponse({ posts: [] });
      return jsonResponse({}, 404);
    });
    const user = userEvent.setup();
    renderAccount();
    expect(await screen.findByText("Sudah divalidasi")).toBeInTheDocument();
    expect(screen.getByRole("heading", { name: "Tulis kabar" })).toBeInTheDocument();
    await user.click(screen.getByRole("button", { name: "Keluar" }));
    expect(await screen.findByText("beranda")).toBeInTheDocument();
  });

  it("links staff to the backoffice", async () => {
    mockFetch(({ url }) =>
      url.includes("/auth/session") ? sessionResponse(staffFixture) : jsonResponse({ posts: [] }),
    );
    renderAccount();
    expect(await screen.findByText("Sudah divalidasi")).toBeInTheDocument();
    expect(screen.getByRole("link", { name: "Buka backoffice staf" })).toHaveAttribute(
      "href",
      "/backoffice",
    );
  });
});
