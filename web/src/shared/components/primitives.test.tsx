import { fireEvent, render, screen } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import { describe, expect, it, vi } from "vitest";
import { Button } from "./Button";
import { Field, describedBy } from "./Field";
import { LoadingState } from "./LoadingState";
import { NotFoundPage } from "./NotFoundPage";
import { SiteFooter } from "./SiteFooter";
import { SiteHeader } from "./SiteHeader";
import { StatusMessage } from "./StatusMessage";

describe("shared primitives", () => {
  it("renders and activates buttons", () => {
    const onClick = vi.fn();
    render(<Button onClick={onClick}>Simpan</Button>);
    fireEvent.click(screen.getByRole("button", { name: "Simpan" }));
    expect(onClick).toHaveBeenCalledOnce();
  });

  it("links field hints and errors through describedBy", () => {
    expect(describedBy("email")).toBeUndefined();
    expect(describedBy("email", "hint")).toBe("email-hint");
    expect(describedBy("email", "hint", "error")).toBe("email-hint email-error");
    render(
      <Field id="email" label="Email" hint="Pakai email aktif." error="Email salah.">
        <input id="email" />
      </Field>,
    );
    expect(screen.getByLabelText("Email")).toBeInTheDocument();
    expect(screen.getByText("Pakai email aktif.")).toHaveAttribute("id", "email-hint");
    expect(screen.getByRole("alert")).toHaveTextContent("Email salah.");
  });

  it("announces status and loading states", () => {
    render(
      <>
        <StatusMessage tone="success">Berhasil.</StatusMessage>
        <StatusMessage tone="error">Gagal.</StatusMessage>
        <LoadingState label="Memuat…" />
      </>,
    );
    expect(screen.getAllByRole("status").length).toBeGreaterThanOrEqual(2);
    expect(screen.getByText("Berhasil.")).toHaveAttribute("role", "status");
    expect(screen.getByText("Gagal.")).toHaveAttribute("role", "alert");
    expect(screen.getByText("Memuat…")).toBeInTheDocument();
  });

  it("shows the correct header action for anonymous, member and staff", () => {
    const { rerender } = render(
      <MemoryRouter>
        <SiteHeader account={null} />
      </MemoryRouter>,
    );
    expect(screen.getByRole("link", { name: "Masuk" })).toHaveAttribute("href", "/masuk");

    rerender(
      <MemoryRouter>
        <SiteHeader account={{ name: "Rina", role: "member" }} />
      </MemoryRouter>,
    );
    expect(screen.getByRole("link", { name: "Rina" })).toHaveAttribute("href", "/akun");

    rerender(
      <MemoryRouter>
        <SiteHeader account={{ name: "Staf", role: "staff" }} overlay />
      </MemoryRouter>,
    );
    expect(screen.getByRole("link", { name: "Backoffice" })).toHaveAttribute("href", "/backoffice");
    expect(document.querySelector(".site-header--overlay")).not.toBeNull();

    rerender(
      <MemoryRouter>
        <SiteHeader account={null} loading />
      </MemoryRouter>,
    );
    expect(screen.queryByRole("link", { name: "Masuk" })).not.toBeInTheDocument();
  });

  it("renders footer navigation, purpose and honest credits", () => {
    render(
      <MemoryRouter>
        <SiteFooter />
      </MemoryRouter>,
    );
    expect(screen.getByRole("link", { name: "Backoffice staf" })).toHaveAttribute(
      "href",
      "/backoffice",
    );
    expect(screen.getByText(/bukan dokumentasi kegiatan maulid\.id/i)).toBeInTheDocument();
    expect(screen.getByRole("link", { name: "Masjid Istiqlal Jakarta" })).toBeInTheDocument();
    expect(screen.getAllByRole("link", { name: "CC BY-SA 4.0" })).toHaveLength(2);
  });

  it("renders a useful not-found page", () => {
    render(
      <MemoryRouter>
        <NotFoundPage />
      </MemoryRouter>,
    );
    expect(screen.getByRole("heading", { name: "Halaman tidak ditemukan" })).toBeInTheDocument();
    expect(screen.getByRole("link", { name: "Buka beranda" })).toHaveAttribute("href", "/");
  });
});
