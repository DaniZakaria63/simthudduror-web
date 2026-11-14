import { render, screen, waitFor, within } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { MemoryRouter, Route, Routes } from "react-router-dom";
import { describe, expect, it } from "vitest";
import {
  accountFixture,
  jsonResponse,
  mockFetch,
  staffFixture,
} from "@/test/fixtures";
import { AccountProvider } from "@/features/account";
import type { Account } from "@/features/account";
import { BackofficePage } from "./BackofficePage";

const memberRow = {
  id: accountFixture.id,
  email: accountFixture.email,
  name: accountFixture.name,
  role: "member" as const,
  validation_status: "unvalidated" as const,
  created_at: accountFixture.created_at,
};

const auditRow = {
  id: "audit-1",
  actor_id: staffFixture.id,
  actor_name: "Staf",
  target_account_id: memberRow.id,
  target_email: memberRow.email,
  action: "validation_status_changed",
  from_status: "unvalidated",
  to_status: "validated",
  created_at: "2026-11-14T09:00:00.000000000Z",
};

function renderBackoffice() {
  return render(
    <MemoryRouter initialEntries={["/backoffice"]}>
      <AccountProvider>
        <Routes>
          <Route path="/backoffice" element={<BackofficePage />} />
        </Routes>
      </AccountProvider>
    </MemoryRouter>,
  );
}

function session(account: Account = staffFixture) {
  return jsonResponse({ account, csrf_token: "csrf", expires_at: "2026-11-20T00:00:00Z" });
}

describe("BackofficePage", () => {
  it("shows a loading state while the session is checked", () => {
    mockFetch(() => new Promise<Response>(() => {}));
    renderBackoffice();
    expect(screen.getByText("Memeriksa akses staf…")).toBeInTheDocument();
  });

  it("gates anonymous visitors", async () => {
    mockFetch(() => jsonResponse({ error: { code: "unauthenticated", message: "no" } }, 401));
    renderBackoffice();
    expect(
      await screen.findByText("Halaman ini memerlukan akun staf yang sudah masuk."),
    ).toBeInTheDocument();
    expect(screen.getByRole("link", { name: "Masuk" })).toHaveAttribute("href", "/masuk");
  });

  it("gates non-staff accounts", async () => {
    mockFetch(() => session(accountFixture));
    renderBackoffice();
    expect(await screen.findByRole("alert")).toHaveTextContent("Akses staf diperlukan");
  });

  it("lists accounts and audit entries for staff", async () => {
    const fetchMock = mockFetch(({ url }) => {
      if (url.includes("/auth/session")) return session();
      if (url.includes("/staff/accounts")) return jsonResponse({ accounts: [memberRow] });
      if (url.includes("/staff/audit")) return jsonResponse({ audit: [auditRow] });
      return jsonResponse({}, 404);
    });
    renderBackoffice();
    expect(await screen.findByText("rina@example.com")).toBeInTheDocument();
    expect(screen.getByRole("heading", { name: "Backoffice" })).toBeInTheDocument();
    const table = screen.getByRole("table");
    expect(within(table).getByText("Belum divalidasi")).toBeInTheDocument();
    expect(screen.getByText(/mengubah rina@example\.com/)).toBeInTheDocument();
    expect(screen.getByText("1 akun pada tampilan ini.")).toBeInTheDocument();
    expect(fetchMock).toHaveBeenCalled();
  });

  it("validates an account and appends the audit entry", async () => {
    mockFetch(({ url, method, body }) => {
      if (url.includes("/auth/session")) return session();
      if (url.includes("/staff/audit")) return jsonResponse({ audit: [] });
      if (url.includes("/staff/accounts/") && method === "PATCH") {
        expect(body).toEqual({ validation_status: "validated" });
        return jsonResponse({
          account: { ...memberRow, validation_status: "validated" },
          audit: auditRow,
        });
      }
      if (url.includes("/staff/accounts")) return jsonResponse({ accounts: [memberRow] });
      return jsonResponse({}, 404);
    });
    const user = userEvent.setup();
    renderBackoffice();
    await user.click(await screen.findByRole("button", { name: "Validasi" }));
    expect(await screen.findByText("Status Rina diperbarui menjadi Sudah divalidasi.")).toBeInTheDocument();
    const table = screen.getByRole("table");
    expect(within(table).getByText("Sudah divalidasi")).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "Cabut validasi" })).toBeInTheDocument();
    expect(screen.getByText(/mengubah rina@example\.com/)).toBeInTheDocument();
  });

  it("reports a failed validation update", async () => {
    mockFetch(({ url, method }) => {
      if (url.includes("/auth/session")) return session();
      if (url.includes("/staff/audit")) return jsonResponse({ audit: [] });
      if (url.includes("/staff/accounts/") && method === "PATCH") {
        return jsonResponse({ error: { code: "internal_error", message: "Tidak dapat menyimpan." } }, 500);
      }
      if (url.includes("/staff/accounts")) return jsonResponse({ accounts: [memberRow] });
      return jsonResponse({}, 404);
    });
    const user = userEvent.setup();
    renderBackoffice();
    await user.click(await screen.findByRole("button", { name: "Validasi" }));
    expect(await screen.findByRole("alert")).toHaveTextContent("Tidak dapat menyimpan.");
  });

  it("filters by status and search query", async () => {
    const fetchMock = mockFetch(({ url }) => {
      if (url.includes("/auth/session")) return session();
      if (url.includes("/staff/audit")) return jsonResponse({ audit: [] });
      if (url.includes("validation_status=validated")) return jsonResponse({ accounts: [] });
      if (url.includes("q=rina")) return jsonResponse({ accounts: [memberRow] });
      if (url.includes("/staff/accounts")) return jsonResponse({ accounts: [memberRow] });
      return jsonResponse({}, 404);
    });
    const user = userEvent.setup();
    renderBackoffice();
    await screen.findByText("rina@example.com");
    await user.selectOptions(screen.getByLabelText("Status validasi"), "validated");
    await waitFor(() =>
      expect(fetchMock.mock.calls.some(([url]) => String(url).includes("validation_status=validated"))).toBe(
        true,
      ),
    );
    expect(await screen.findByText("Tidak ada akun yang cocok dengan filter ini.")).toBeInTheDocument();
    await user.selectOptions(screen.getByLabelText("Status validasi"), "");
    await screen.findByText("rina@example.com");
    await user.type(screen.getByLabelText("Cari akun"), "rina");
    await waitFor(() =>
      expect(fetchMock.mock.calls.some(([url]) => String(url).includes("q=rina"))).toBe(true),
    );
  });

  it("shows an empty audit trail and reloads", async () => {
    mockFetch(({ url }) => {
      if (url.includes("/auth/session")) return session();
      if (url.includes("/staff/accounts")) return jsonResponse({ accounts: [memberRow] });
      if (url.includes("/staff/audit")) return jsonResponse({ audit: [] });
      return jsonResponse({}, 404);
    });
    const user = userEvent.setup();
    renderBackoffice();
    expect(await screen.findByText("Belum ada perubahan status yang tercatat.")).toBeInTheDocument();
    await user.click(screen.getByRole("button", { name: "Muat ulang" }));
    expect(await screen.findByText("rina@example.com")).toBeInTheDocument();
  });

  it("shows a retry action when staff data fails", async () => {
    mockFetch(({ url }) => {
      if (url.includes("/auth/session")) return session();
      return jsonResponse({ error: { code: "internal_error", message: "Data gagal." } }, 500);
    });
    renderBackoffice();
    expect(await screen.findByRole("alert")).toHaveTextContent("Data gagal.");
    expect(screen.getByRole("button", { name: "Coba lagi" })).toBeInTheDocument();
  });
});
