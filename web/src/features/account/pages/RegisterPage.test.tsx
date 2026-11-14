import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { MemoryRouter, Route, Routes } from "react-router-dom";
import { describe, expect, it } from "vitest";
import { accountFixture, jsonResponse, mockFetch } from "@/test/fixtures";
import { AccountProvider } from "../session";
import { RegisterPage } from "./RegisterPage";

function renderRegister() {
  return render(
    <MemoryRouter initialEntries={["/daftar"]}>
      <AccountProvider>
        <Routes>
          <Route path="/daftar" element={<RegisterPage />} />
          <Route path="/akun" element={<p>halaman akun</p>} />
        </Routes>
      </AccountProvider>
    </MemoryRouter>,
  );
}

describe("RegisterPage", () => {
  it("validates names, emails, passwords and confirmation", async () => {
    mockFetch(() => jsonResponse({ error: { code: "unauthenticated", message: "no" } }, 401));
    const user = userEvent.setup();
    renderRegister();
    expect(screen.getByText("Memeriksa sesi…")).toBeInTheDocument();
    await screen.findByLabelText("Nama");
    await user.type(screen.getByLabelText("Nama"), "R");
    await user.type(screen.getByLabelText("Email"), "bukan-email");
    await user.type(screen.getByLabelText("Kata sandi"), "pendek");
    await user.type(screen.getByLabelText("Ulangi kata sandi"), "beda");
    await user.click(screen.getByRole("button", { name: "Buat akun" }));
    expect(screen.getByText("Nama harus 2 sampai 80 karakter.")).toBeInTheDocument();
    expect(screen.getByText("Masukkan alamat email yang valid.")).toBeInTheDocument();
    expect(screen.getByText("Kata sandi harus minimal 10 karakter.")).toBeInTheDocument();
    expect(screen.getByText("Ulangi kata sandi dengan tepat.")).toBeInTheDocument();
  });

  it("registers and navigates to the account page", async () => {
    mockFetch(({ url }) => {
      if (url.endsWith("/auth/session")) {
        return jsonResponse({ error: { code: "unauthenticated", message: "no" } }, 401);
      }
      return jsonResponse({ account: accountFixture, csrf_token: "csrf", expires_at: "x" }, 201);
    });
    const user = userEvent.setup();
    renderRegister();
    await screen.findByLabelText("Nama");
    await user.type(screen.getByLabelText("Nama"), "Rina");
    await user.type(screen.getByLabelText("Email"), "rina@example.com");
    await user.type(screen.getByLabelText("Kata sandi"), "kata-sandi-aman");
    await user.type(screen.getByLabelText("Ulangi kata sandi"), "kata-sandi-aman");
    await user.click(screen.getByRole("button", { name: "Buat akun" }));
    expect(await screen.findByText("halaman akun")).toBeInTheDocument();
  });

  it("shows a duplicate-email or server failure", async () => {
    mockFetch(({ url }) => {
      if (url.endsWith("/auth/session")) {
        return jsonResponse({ error: { code: "unauthenticated", message: "no" } }, 401);
      }
      return jsonResponse({ error: { code: "email_unavailable", message: "Email tidak dapat digunakan." } }, 409);
    });
    const user = userEvent.setup();
    renderRegister();
    await screen.findByLabelText("Nama");
    await user.type(screen.getByLabelText("Nama"), "Rina");
    await user.type(screen.getByLabelText("Email"), "rina@example.com");
    await user.type(screen.getByLabelText("Kata sandi"), "kata-sandi-aman");
    await user.type(screen.getByLabelText("Ulangi kata sandi"), "kata-sandi-aman");
    await user.click(screen.getByRole("button", { name: "Buat akun" }));
    expect(await screen.findByRole("alert")).toHaveTextContent("Email tidak dapat digunakan.");
  });

  it("shows the existing session instead of the form", async () => {
    mockFetch(() => jsonResponse({ account: accountFixture, csrf_token: "csrf", expires_at: "x" }));
    renderRegister();
    expect(await screen.findByRole("heading", { name: "Sudah punya sesi" })).toBeInTheDocument();
  });
});
