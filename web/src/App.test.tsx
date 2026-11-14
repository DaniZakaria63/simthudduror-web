import { render, screen } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import { describe, expect, it } from "vitest";
import App from "./App";
import { jsonResponse, mockFetch } from "@/test/fixtures";

function renderApp(path = "/") {
  return render(
    <MemoryRouter initialEntries={[path]}>
      <App />
    </MemoryRouter>,
  );
}

describe("App", () => {
  it("renders the public page with chrome", async () => {
    mockFetch(({ url }) =>
      url.includes("/auth/session")
        ? jsonResponse({ error: { code: "unauthenticated", message: "no" } }, 401)
        : jsonResponse({ posts: [] }),
    );
    renderApp("/");
    expect(screen.getByRole("link", { name: "Lewati ke konten" })).toHaveAttribute("href", "#main");
    expect(screen.getByRole("banner")).toBeInTheDocument();
    expect(screen.getByRole("contentinfo")).toBeInTheDocument();
    expect(await screen.findByRole("heading", { name: "Kabar Maulid kita" })).toBeInTheDocument();
  });

  it("renders the not-found route", async () => {
    mockFetch(({ url }) =>
      url.includes("/auth/session")
        ? jsonResponse({ error: { code: "unauthenticated", message: "no" } }, 401)
        : jsonResponse({ posts: [] }),
    );
    renderApp("/tidak-ada");
    expect(await screen.findByRole("heading", { name: "Halaman tidak ditemukan" })).toBeInTheDocument();
  });
});
