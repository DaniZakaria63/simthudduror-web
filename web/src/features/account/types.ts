export type Role = "member" | "staff";
export type ValidationStatus = "unvalidated" | "validated";

export type Account = {
  id: string;
  email: string;
  name: string;
  role: Role;
  validation_status: ValidationStatus;
  created_at: string;
};

export type SessionPayload = {
  account: Account;
  csrf_token: string;
  expires_at: string;
  token?: string;
};
