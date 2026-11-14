import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { MemoryRouter, Route, Routes } from "react-router-dom";
import { describe, expect, it } from "vitest";
import { accountFixture, jsonResponse, mockFetch } from "@/test/fixtures";
import { AccountProvider } from "../session";
import { LoginPage } from "./LoginPage";

function renderLogin() {
  return render(
    <MemoryRouter initialEntries={["/masuk"]}>
      <AccountProvider>
        <Routes>
          <Route path="/masuk" element={<LoginPage />} />
          <Route path="/akun" element={<p>halaman akun</p>} />
        </Routes>
      </AccountProvider>
    </MemoryRouter>,
  );
}

describe("LoginPage", () => {
  it("validates required fields before calling the API", async () => {
    const fetchMock = mockFetch(() => jsonResponse({}, 401));
    renderLogin();
    await screen.findByRole("heading", { name: "Masuk" });
    await waitFor(() => expect(screen.queryByText("Memeriksa sesi…")).not.toBeInTheDocument());
    fireEvent.click(screen.getByRole("button", { name: "Masuk" }));
    expect(screen.getByRole("alert")).toHaveTextContent("Isi email dan kata sandi terlebih dahulu.");
    expect(fetchMock.mock.calls.filter(([url]) => String(url).includes("/login"))).toHaveLength(0);
  });

  it("shows the server error for invalid credentials", async () => {
    mockFetch(({ url }) => {
      if (url.endsWith("/auth/session")) {
        return jsonResponse({ error: { code: "unauthenticated", message: "no" } }, 401);
      }
      return jsonResponse({ error: { code: "invalid_credentials", message: "Email atau kata sandi salah." } }, 401);
    });
    const user = userEvent.setup();
    renderLogin();
    await screen.findByLabelText("Email");
    await user.type(screen.getByLabelText("Email"), "rina@example.com");
    await user.type(screen.getByLabelText("Kata sandi"), "kata-sandi-salah");
    await user.click(screen.getByRole("button", { name: "Masuk" }));
    expect(await screen.findByRole("alert")).toHaveTextContent("Email atau kata sandi salah.");
  });

  it("signs in and navigates to the account page", async () => {
    mockFetch(({ url }) => {
      if (url.endsWith("/auth/session")) {
        return jsonResponse({ error: { code: "unauthenticated", message: "no" } }, 401);
      }
      return jsonResponse({ account: accountFixture, csrf_token: "csrf", expires_at: "x" });
    });
    const user = userEvent.setup();
    renderLogin();
    await screen.findByLabelText("Email");
    await user.type(screen.getByLabelText("Email"), "rina@example.com");
    await user.type(screen.getByLabelText("Kata sandi"), "kata-sandi-aman");
    await user.click(screen.getByRole("button", { name: "Masuk" }));
    expect(await screen.findByText("halaman akun")).toBeInTheDocument();
  });

  it("shows the current session instead of the form", async () => {
    mockFetch(() => jsonResponse({ account: accountFixture, csrf_token: "csrf", expires_at: "x" }));
    renderLogin();
    expect(await screen.findByRole("heading", { name: "Sudah masuk" })).toBeInTheDocument();
    expect(screen.getByRole("link", { name: "Buka halaman akun" })).toHaveAttribute("href", "/akun");
  });
});
