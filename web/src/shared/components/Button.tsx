import type { ButtonHTMLAttributes, ReactNode } from "react";

type ButtonProps = ButtonHTMLAttributes<HTMLButtonElement> & {
  tone?: "primary" | "quiet" | "danger";
  children: ReactNode;
};

export function Button({ tone = "primary", className = "", children, ...rest }: ButtonProps) {
  return (
    <button className={`button button--${tone} ${className}`.trim()} {...rest}>
      {children}
    </button>
  );
}

import "./Button.css";
