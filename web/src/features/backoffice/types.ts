export type StaffAccount = {
  id: string;
  email: string;
  name: string;
  role: "member" | "staff";
  validation_status: "unvalidated" | "validated";
  created_at: string;
};

export type AuditEntry = {
  id: string;
  actor_id: string;
  actor_name: string;
  target_account_id: string;
  target_email: string;
  action: string;
  from_status: string;
  to_status: string;
  created_at: string;
};

export type StaffFilter = {
  validation_status: "" | "unvalidated" | "validated";
  query: string;
};
