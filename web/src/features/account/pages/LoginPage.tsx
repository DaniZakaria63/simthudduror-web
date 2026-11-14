import { useState, type FormEvent } from "react";
import { Link, useNavigate } from "react-router-dom";
import { ApiError } from "@/shared/api/client";
import { Button } from "@/shared/components/Button";
import { Field, describedBy } from "@/shared/components/Field";
import { LoadingState } from "@/shared/components/LoadingState";
import { StatusMessage } from "@/shared/components/StatusMessage";
import { useSession } from "../session";

export function LoginPage() {
  const { status, account, signIn } = useSession();
  const navigate = useNavigate();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [pending, setPending] = useState(false);

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setError(null);
    if (!email.trim() || !password) {
      setError("Isi email dan kata sandi terlebih dahulu.");
      return;
    }
    setPending(true);
    try {
      await signIn(email.trim(), password);
      navigate("/akun", { replace: true });
    } catch (caught) {
      setError(caught instanceof ApiError ? caught.message : "Tidak dapat masuk. Coba lagi.");
    } finally {
      setPending(false);
    }
  }

  if (status === "loading") {
    return (
      <section className="auth-page" aria-labelledby="masuk-title">
        <h1 id="masuk-title">Masuk</h1>
        <LoadingState label="Memeriksa sesi…" />
      </section>
    );
  }

  if (account) {
    return (
      <section className="auth-page" aria-labelledby="masuk-title">
        <h1 id="masuk-title">Sudah masuk</h1>
        <StatusMessage tone="success">
          Kamu masuk sebagai {account.name}. Lanjutkan ke halaman akun untuk membaca dan
          berinteraksi.
        </StatusMessage>
        <p>
          <Link className="auth-page__link" to="/akun">
            Buka halaman akun
          </Link>
        </p>
      </section>
    );
  }

  return (
    <section className="auth-page" aria-labelledby="masuk-title">
      <div className="auth-page__intro">
        <h1 id="masuk-title">Masuk</h1>
        <p>
          Gunakan akun yang sudah terdaftar. Semua akun dapat menyukai kabar; penerbitan memerlukan
          validasi staf.
        </p>
      </div>
      <form className="auth-page__form" onSubmit={handleSubmit} noValidate>
        {error ? <StatusMessage tone="error">{error}</StatusMessage> : null}
        <Field id="login-email" label="Email" hint="Alamat email yang dipakai saat mendaftar.">
          <input
            id="login-email"
            name="email"
            type="email"
            autoComplete="email"
            value={email}
            onChange={(event) => setEmail(event.target.value)}
            aria-describedby={describedBy("login-email", "Alamat email yang dipakai saat mendaftar.")}
            required
          />
        </Field>
        <Field id="login-password" label="Kata sandi">
          <input
            id="login-password"
            name="password"
            type="password"
            autoComplete="current-password"
            value={password}
            onChange={(event) => setPassword(event.target.value)}
            required
          />
        </Field>
        <Button type="submit" disabled={pending}>
          {pending ? "Sedang masuk…" : "Masuk"}
        </Button>
        <p className="auth-page__switch">
          Belum punya akun? <Link to="/daftar">Daftar</Link>
        </p>
      </form>
    </section>
  );
}

import "./AuthPage.css";
