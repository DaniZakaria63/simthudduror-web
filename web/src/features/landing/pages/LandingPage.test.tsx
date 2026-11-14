import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { MemoryRouter, Route, Routes } from "react-router-dom";
import { describe, expect, it } from "vitest";
import {
  accountFixture,
  jsonResponse,
  mockFetch,
  postFixture,
} from "@/test/fixtures";
import { AccountProvider } from "@/features/account";
import { LandingPage } from "./LandingPage";

function renderLanding() {
  return render(
    <MemoryRouter initialEntries={["/"]}>
      <AccountProvider>
        <Routes>
          <Route path="/" element={<LandingPage />} />
          <Route path="/masuk" element={<p>halaman masuk</p>} />
        </Routes>
      </AccountProvider>
    </MemoryRouter>,
  );
}

describe("LandingPage", () => {
  it("renders the editorial masthead and contribution argument", async () => {
    mockFetch(({ url }) =>
      url.includes("/auth/session")
        ? jsonResponse({ error: { code: "unauthenticated", message: "no" } }, 401)
        : jsonResponse({ posts: [] }),
    );
    renderLanding();
    expect(
      screen.getByRole("heading", { name: "Kabar Maulid kita" }),
    ).toBeInTheDocument();
    expect(screen.getByRole("link", { name: "Baca kabar" })).toHaveAttribute("href", "/#kabar");
    expect(screen.getByRole("heading", { name: "Pertama, buat akun." })).toBeInTheDocument();
    expect(screen.getByRole("link", { name: "Buat akun" })).toHaveAttribute("href", "/daftar");
    expect(await screen.findByText(/Belum ada kabar\. Akun yang sudah divalidasi/)).toBeInTheDocument();
  });

  it("filters the real feed by the visitor's words", async () => {
    const fetchMock = mockFetch(({ url }) => {
      if (url.includes("/auth/session")) {
        return jsonResponse({ error: { code: "unauthenticated", message: "no" } }, 401);
      }
      if (url.includes("q=bandung")) {
        return jsonResponse({ posts: [{ ...postFixture, id: "post-2", body: "Kabar dari Bandung." }] });
      }
      return jsonResponse({ posts: [postFixture] });
    });
    const user = userEvent.setup();
    renderLanding();
    expect(await screen.findByText(postFixture.body)).toBeInTheDocument();
    await user.type(screen.getByLabelText("Cari kabar"), "bandung");
    expect(await screen.findByText("Kabar dari Bandung.")).toBeInTheDocument();
    expect(fetchMock.mock.calls.some(([url]) => String(url).includes("q=bandung"))).toBe(true);
  });

  it("sends anonymous readers to sign in before liking", async () => {
    mockFetch(({ url }) => {
      if (url.includes("/auth/session")) {
        return jsonResponse({ error: { code: "unauthenticated", message: "no" } }, 401);
      }
      return jsonResponse({ posts: [postFixture] });
    });
    const user = userEvent.setup();
    renderLanding();
    await user.click(await screen.findByRole("button", { name: "Suka (0)" }));
    expect(await screen.findByText("halaman masuk")).toBeInTheDocument();
  });

  it("lets an authenticated reader like a post", async () => {
    mockFetch(({ url, method }) => {
      if (url.includes("/auth/session")) {
        return jsonResponse({ account: accountFixture, csrf_token: "csrf", expires_at: "x" });
      }
      if (url.includes("/like") && method === "PUT") {
        return jsonResponse({ post: { ...postFixture, like_count: 1, liked_by_me: true } });
      }
      return jsonResponse({ posts: [postFixture] });
    });
    const user = userEvent.setup();
    renderLanding();
    await user.click(await screen.findByRole("button", { name: "Suka (0)" }));
    expect(await screen.findByText("Kabar disukai.")).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "Batal suka (1)" })).toBeInTheDocument();
  });

  it("reports an interaction failure", async () => {
    mockFetch(({ url, method }) => {
      if (url.includes("/auth/session")) {
        return jsonResponse({ account: accountFixture, csrf_token: "csrf", expires_at: "x" });
      }
      if (url.includes("/like") && method === "PUT") {
        return jsonResponse({ error: { code: "internal_error", message: "Gagal menyimpan." } }, 500);
      }
      return jsonResponse({ posts: [postFixture] });
    });
    const user = userEvent.setup();
    renderLanding();
    await user.click(await screen.findByRole("button", { name: "Suka (0)" }));
    expect(await screen.findByText("Gagal menyimpan.")).toBeInTheDocument();
  });

  it("offers a retry when the feed fails", async () => {
    let calls = 0;
    mockFetch(({ url }) => {
      if (url.includes("/auth/session")) {
        return jsonResponse({ error: { code: "unauthenticated", message: "no" } }, 401);
      }
      calls += 1;
      if (calls === 1) {
        return jsonResponse({ error: { code: "internal_error", message: "Kabar gagal dimuat." } }, 500);
      }
      return jsonResponse({ posts: [postFixture] });
    });
    const user = userEvent.setup();
    renderLanding();
    expect(await screen.findByRole("alert")).toHaveTextContent("Kabar gagal dimuat.");
    await user.click(screen.getByRole("button", { name: "Coba lagi" }));
    expect(await screen.findByText(postFixture.body)).toBeInTheDocument();
  });

  it("shows no matching result for an empty filter", async () => {
    mockFetch(({ url }) => {
      if (url.includes("/auth/session")) {
        return jsonResponse({ error: { code: "unauthenticated", message: "no" } }, 401);
      }
      if (url.includes("q=")) return jsonResponse({ posts: [] });
      return jsonResponse({ posts: [postFixture] });
    });
    const user = userEvent.setup();
    renderLanding();
    await screen.findByText(postFixture.body);
    await user.type(screen.getByLabelText("Cari kabar"), "tidakada");
    await waitFor(() =>
      expect(screen.getByText("Belum ada kabar yang cocok. Coba kata lain.")).toBeInTheDocument(),
    );
  });
});
