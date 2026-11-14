import { useCallback, useDeferredValue, useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { ApiError } from "@/shared/api/client";
import { Button } from "@/shared/components/Button";
import { LoadingState } from "@/shared/components/LoadingState";
import { StatusMessage } from "@/shared/components/StatusMessage";
import { formatDate, formatDateTime } from "@/shared/format";
import { useSession } from "@/features/account";
import { listAccounts, listAudit, setValidation as updateValidation } from "../api";
import type { AuditEntry, StaffAccount, StaffFilter } from "../types";

type PageStatus = "loading" | "ready" | "error";

export function BackofficePage() {
  const { status: sessionStatus, account } = useSession();
  const [accounts, setAccounts] = useState<StaffAccount[]>([]);
  const [audit, setAudit] = useState<AuditEntry[]>([]);
  const [pageStatus, setPageStatus] = useState<PageStatus>("loading");
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);
  const [pendingId, setPendingId] = useState<string | null>(null);
  const [filter, setFilter] = useState<StaffFilter>({ validation_status: "", query: "" });
  const deferredQuery = useDeferredValue(filter.query);
  const [revision, setRevision] = useState(0);

  const load = useCallback(
    async (signal: AbortSignal) => {
      setPageStatus("loading");
      try {
        const [nextAccounts, nextAudit] = await Promise.all([
          listAccounts({ validation_status: filter.validation_status, query: deferredQuery }, signal),
          listAudit(signal),
        ]);
        setAccounts(nextAccounts);
        setAudit(nextAudit);
        setError(null);
        setPageStatus("ready");
      } catch (caught) {
        if (caught instanceof DOMException && caught.name === "AbortError") {
          return;
        }
        setError(caught instanceof ApiError ? caught.message : "Data staf tidak dapat dimuat.");
        setPageStatus("error");
      }
    },
    [deferredQuery, filter.validation_status],
  );

  useEffect(() => {
    if (sessionStatus !== "authenticated" || account?.role !== "staff") {
      return;
    }
    const controller = new AbortController();
    void load(controller.signal);
    return () => controller.abort();
  }, [account?.role, load, revision, sessionStatus]);

  async function changeValidation(target: StaffAccount, nextStatus: "unvalidated" | "validated") {
    setPendingId(target.id);
    setError(null);
    setSuccess(null);
    try {
      const result = await updateValidation(target.id, nextStatus);
      setAccounts((current) =>
        current.map((item) => (item.id === result.account.id ? result.account : item)),
      );
      if (result.audit) {
        setAudit((current) => [result.audit as AuditEntry, ...current]);
      }
      setSuccess(
        `Status ${result.account.name} diperbarui menjadi ${
          nextStatus === "validated" ? "Sudah divalidasi" : "Belum divalidasi"
        }.`,
      );
    } catch (caught) {
      setError(caught instanceof ApiError ? caught.message : "Status tidak dapat diperbarui.");
    } finally {
      setPendingId(null);
    }
  }

  if (sessionStatus === "loading") {
    return (
      <section className="backoffice" aria-labelledby="backoffice-title">
        <h1 id="backoffice-title">Backoffice</h1>
        <LoadingState label="Memeriksa akses staf…" />
      </section>
    );
  }

  if (sessionStatus !== "authenticated" || !account) {
    return (
      <section className="backoffice backoffice--gate" aria-labelledby="backoffice-title">
        <h1 id="backoffice-title">Backoffice</h1>
        <p>Halaman ini memerlukan akun staf yang sudah masuk.</p>
        <p>
          <Link className="button button--primary" to="/masuk">
            Masuk
          </Link>
        </p>
      </section>
    );
  }

  if (account.role !== "staff") {
    return (
      <section className="backoffice backoffice--gate" aria-labelledby="backoffice-title">
        <h1 id="backoffice-title">Backoffice</h1>
        <StatusMessage tone="error">
          Akses staf diperlukan. Akun {account.name} tidak memiliki peran staf.
        </StatusMessage>
        <p>
          <Link to="/akun">Kembali ke akun</Link>
        </p>
      </section>
    );
  }

  return (
    <section className="backoffice" aria-labelledby="backoffice-title">
      <header className="backoffice__header">
        <div>
          <h1 id="backoffice-title">Backoffice</h1>
          <p>
            Tinjau akun dan kelola validasi. Setiap perubahan dicatat dalam audit. Masuk sebagai{" "}
            {account.name}.
          </p>
        </div>
        <Button tone="quiet" type="button" onClick={() => setRevision((value) => value + 1)}>
          Muat ulang
        </Button>
      </header>

      {error ? <StatusMessage tone="error">{error}</StatusMessage> : null}
      {success ? <StatusMessage tone="success">{success}</StatusMessage> : null}

      <div className="backoffice__controls">
        <div className="field">
          <label className="field__label" htmlFor="staff-filter">
            Status validasi
          </label>
          <select
            id="staff-filter"
            value={filter.validation_status}
            onChange={(event) =>
              setFilter((current) => ({
                ...current,
                validation_status: event.target.value as StaffFilter["validation_status"],
              }))
            }
          >
            <option value="">Semua status</option>
            <option value="unvalidated">Belum divalidasi</option>
            <option value="validated">Sudah divalidasi</option>
          </select>
        </div>
        <div className="field">
          <label className="field__label" htmlFor="staff-search">
            Cari akun
          </label>
          <input
            id="staff-search"
            type="search"
            value={filter.query}
            onChange={(event) =>
              setFilter((current) => ({ ...current, query: event.target.value }))
            }
            placeholder="Nama atau email"
            maxLength={100}
            autoComplete="off"
          />
        </div>
      </div>

      {pageStatus === "loading" && accounts.length === 0 ? (
        <LoadingState label="Memuat akun…" />
      ) : null}
      {pageStatus === "error" && accounts.length === 0 ? (
        <Button tone="quiet" type="button" onClick={() => setRevision((value) => value + 1)}>
          Coba lagi
        </Button>
      ) : null}

      <p className="backoffice__count" aria-live="polite">
        {accounts.length} akun pada tampilan ini.
      </p>

      {accounts.length === 0 && pageStatus === "ready" ? (
        <p className="backoffice__empty">Tidak ada akun yang cocok dengan filter ini.</p>
      ) : (
        <div className="backoffice__table-wrap">
          <table className="backoffice__table">
            <caption className="visually-hidden">Daftar akun dan status validasinya</caption>
            <thead>
              <tr>
                <th scope="col">Nama</th>
                <th scope="col">Email</th>
                <th scope="col">Peran</th>
                <th scope="col">Status validasi</th>
                <th scope="col">Bergabung</th>
                <th scope="col">Tindakan</th>
              </tr>
            </thead>
            <tbody>
              {accounts.map((item) => {
                const targetStatus = item.validation_status === "validated" ? "unvalidated" : "validated";
                return (
                  <tr key={item.id}>
                    <td data-label="Nama">{item.name}</td>
                    <td data-label="Email">{item.email}</td>
                    <td data-label="Peran">{item.role === "staff" ? "Staf" : "Anggota"}</td>
                    <td data-label="Status validasi">
                      {item.validation_status === "validated"
                        ? "Sudah divalidasi"
                        : "Belum divalidasi"}
                    </td>
                    <td data-label="Bergabung">{formatDate(item.created_at)}</td>
                    <td data-label="Tindakan">
                      <Button
                        tone={targetStatus === "validated" ? "primary" : "danger"}
                        type="button"
                        disabled={pendingId === item.id}
                        onClick={() => void changeValidation(item, targetStatus)}
                      >
                        {pendingId === item.id
                          ? "Menyimpan…"
                          : targetStatus === "validated"
                            ? "Validasi"
                            : "Cabut validasi"}
                      </Button>
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
      )}

      <section className="backoffice__audit" aria-labelledby="audit-title">
        <h2 id="audit-title">Audit perubahan</h2>
        {audit.length === 0 ? (
          <p className="backoffice__empty">Belum ada perubahan status yang tercatat.</p>
        ) : (
          <ol className="backoffice__audit-list">
            {audit.map((entry) => (
              <li key={entry.id}>
                <p>
                  <strong>{entry.actor_name}</strong> mengubah {entry.target_email}:{" "}
                  {entry.from_status === "validated" ? "Sudah divalidasi" : "Belum divalidasi"} →{" "}
                  {entry.to_status === "validated" ? "Sudah divalidasi" : "Belum divalidasi"}.
                </p>
                <time dateTime={entry.created_at}>{formatDateTime(entry.created_at)}</time>
              </li>
            ))}
          </ol>
        )}
      </section>
    </section>
  );
}

import "./BackofficePage.css";
