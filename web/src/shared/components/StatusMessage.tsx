import type { ReactNode } from "react";

type StatusMessageProps = {
  tone?: "info" | "success" | "error";
  children: ReactNode;
};

export function StatusMessage({ tone = "info", children }: StatusMessageProps) {
  return (
    <p
      className={`status-message status-message--${tone}`}
      role={tone === "error" ? "alert" : "status"}
    >
      {children}
    </p>
  );
}

import "./StatusMessage.css";
