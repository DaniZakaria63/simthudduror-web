import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
  type ReactNode,
} from "react";
import { ApiError, setCSRFToken } from "@/shared/api/client";
import * as accountApi from "./api";
import type { Account } from "./types";

type SessionStatus = "loading" | "anonymous" | "authenticated" | "error";

type SessionValue = {
  status: SessionStatus;
  account: Account | null;
  csrfToken: string | null;
  error: string | null;
  reload: () => Promise<void>;
  signIn: (email: string, password: string) => Promise<Account>;
  signUp: (email: string, name: string, password: string) => Promise<Account>;
  signOut: () => Promise<void>;
};

const SessionContext = createContext<SessionValue | null>(null);

export function AccountProvider({ children }: { children: ReactNode }) {
  const [status, setStatus] = useState<SessionStatus>("loading");
  const [account, setAccount] = useState<Account | null>(null);
  const [csrfToken, setCsrfTokenState] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

  const applySession = useCallback((nextAccount: Account, nextCSRF: string) => {
    setCSRFToken(nextCSRF);
    setCsrfTokenState(nextCSRF);
    setAccount(nextAccount);
    setError(null);
    setStatus("authenticated");
  }, []);

  const clearSession = useCallback(() => {
    setCSRFToken(null);
    setCsrfTokenState(null);
    setAccount(null);
    setError(null);
    setStatus("anonymous");
  }, []);

  const reload = useCallback(async () => {
    setStatus("loading");
    try {
      const payload = await accountApi.fetchSession();
      applySession(payload.account, payload.csrf_token);
    } catch (caught) {
      if (caught instanceof ApiError && caught.status === 401) {
        clearSession();
        return;
      }
      setError(caught instanceof Error ? caught.message : "Sesi tidak dapat diperiksa.");
      setStatus("error");
    }
  }, [applySession, clearSession]);

  useEffect(() => {
    const controller = new AbortController();
    setStatus("loading");
    accountApi
      .fetchSession(controller.signal)
      .then((payload) => {
        applySession(payload.account, payload.csrf_token);
      })
      .catch((caught: unknown) => {
        if (caught instanceof DOMException && caught.name === "AbortError") {
          return;
        }
        if (caught instanceof ApiError && caught.status === 401) {
          clearSession();
          return;
        }
        setError(caught instanceof Error ? caught.message : "Sesi tidak dapat diperiksa.");
        setStatus("error");
      });
    return () => controller.abort();
  }, [applySession, clearSession]);

  const signIn = useCallback(
    async (email: string, password: string) => {
      const payload = await accountApi.login(email, password);
      applySession(payload.account, payload.csrf_token);
      return payload.account;
    },
    [applySession],
  );

  const signUp = useCallback(
    async (email: string, name: string, password: string) => {
      const payload = await accountApi.register(email, name, password);
      applySession(payload.account, payload.csrf_token);
      return payload.account;
    },
    [applySession],
  );

  const signOut = useCallback(async () => {
    try {
      await accountApi.logout();
    } catch {
      // Local state is cleared even when the server is unreachable; the
      // credential will expire server-side regardless.
    } finally {
      clearSession();
    }
  }, [clearSession]);

  const value = useMemo(
    () => ({ status, account, csrfToken, error, reload, signIn, signUp, signOut }),
    [status, account, csrfToken, error, reload, signIn, signUp, signOut],
  );

  return <SessionContext.Provider value={value}>{children}</SessionContext.Provider>;
}

export function useSession(): SessionValue {
  const value = useContext(SessionContext);
  if (!value) {
    throw new Error("useSession harus dipakai di dalam AccountProvider.");
  }
  return value;
}
