import { apiRequest } from "@/shared/api/client";
import type { AuditEntry, StaffAccount, StaffFilter } from "./types";

type AccountsResponse = { accounts: StaffAccount[] };
type AuditResponse = { audit: AuditEntry[] };
type UpdateResponse = { account: StaffAccount; audit: AuditEntry | null };

export async function listAccounts(
  filter: StaffFilter,
  signal?: AbortSignal,
): Promise<StaffAccount[]> {
  const params = new URLSearchParams();
  if (filter.validation_status) {
    params.set("validation_status", filter.validation_status);
  }
  if (filter.query.trim()) {
    params.set("q", filter.query.trim());
  }
  const suffix = params.toString() ? `?${params.toString()}` : "";
  const payload = await apiRequest<AccountsResponse>(`/api/v1/staff/accounts${suffix}`, {
    signal,
  });
  return payload.accounts;
}

export async function setValidation(
  accountID: string,
  status: "unvalidated" | "validated",
): Promise<UpdateResponse> {
  return apiRequest<UpdateResponse>(
    `/api/v1/staff/accounts/${encodeURIComponent(accountID)}/validation`,
    { method: "PATCH", body: { validation_status: status } },
  );
}

export async function listAudit(signal?: AbortSignal): Promise<AuditEntry[]> {
  const payload = await apiRequest<AuditResponse>("/api/v1/staff/audit", { signal });
  return payload.audit;
}
