import { useState, type FormEvent } from "react";
import { Link, useNavigate } from "react-router-dom";
import { ApiError } from "@/shared/api/client";
import { Button } from "@/shared/components/Button";
import { Field, describedBy } from "@/shared/components/Field";
import { LoadingState } from "@/shared/components/LoadingState";
import { StatusMessage } from "@/shared/components/StatusMessage";
import { useSession } from "../session";

type FormErrors = {
  name?: string;
  email?: string;
  password?: string;
  confirm?: string;
};

export function RegisterPage() {
  const { status, account, signUp } = useSession();
  const navigate = useNavigate();
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirm, setConfirm] = useState("");
  const [errors, setErrors] = useState<FormErrors>({});
  const [requestError, setRequestError] = useState<string | null>(null);
  const [pending, setPending] = useState(false);

  function validate(): boolean {
    const next: FormErrors = {};
    if (name.trim().length < 2) {
      next.name = "Nama harus 2 sampai 80 karakter.";
    }
    if (!email.includes("@")) {
      next.email = "Masukkan alamat email yang valid.";
    }
    if (password.length < 10) {
      next.password = "Kata sandi harus minimal 10 karakter.";
    }
    if (confirm !== password) {
      next.confirm = "Ulangi kata sandi dengan tepat.";
    }
    setErrors(next);
    return Object.keys(next).length === 0;
  }

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setRequestError(null);
    if (!validate()) {
      return;
    }
    setPending(true);
    try {
      await signUp(email.trim(), name.trim(), password);
      navigate("/akun", { replace: true });
    } catch (caught) {
      setRequestError(
        caught instanceof ApiError ? caught.message : "Pendaftaran tidak dapat diproses. Coba lagi.",
      );
    } finally {
      setPending(false);
    }
  }

  if (status === "loading") {
    return (
      <section className="auth-page" aria-labelledby="daftar-title">
        <h1 id="daftar-title">Daftar</h1>
        <LoadingState label="Memeriksa sesi…" />
      </section>
    );
  }

  if (account) {
    return (
      <section className="auth-page" aria-labelledby="daftar-title">
        <h1 id="daftar-title">Sudah punya sesi</h1>
        <StatusMessage tone="success">
          Kamu masuk sebagai {account.name}. Satu perangkat memakai satu sesi akun pada satu waktu.
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
    <section className="auth-page" aria-labelledby="daftar-title">
      <div className="auth-page__intro">
        <h1 id="daftar-title">Daftar</h1>
        <p>
          Buat akun untuk membaca, menyukai, dan kelak menerbitkan kabar. Akun baru berstatus belum
          divalidasi sampai staf meninjaunya.
        </p>
      </div>
      <form className="auth-page__form" onSubmit={handleSubmit} noValidate>
        {requestError ? <StatusMessage tone="error">{requestError}</StatusMessage> : null}
        <Field
          id="register-name"
          label="Nama"
          hint="Nama yang tampil di samping kabar yang kamu terbitkan."
          error={errors.name}
        >
          <input
            id="register-name"
            name="name"
            type="text"
            autoComplete="name"
            value={name}
            onChange={(event) => setName(event.target.value)}
            aria-describedby={describedBy(
              "register-name",
              "Nama yang tampil di samping kabar yang kamu terbitkan.",
              errors.name,
            )}
            aria-invalid={errors.name ? true : undefined}
            maxLength={80}
            required
          />
        </Field>
        <Field id="register-email" label="Email" error={errors.email}>
          <input
            id="register-email"
            name="email"
            type="email"
            autoComplete="email"
            value={email}
            onChange={(event) => setEmail(event.target.value)}
            aria-describedby={describedBy("register-email", undefined, errors.email)}
            aria-invalid={errors.email ? true : undefined}
            maxLength={254}
            required
          />
        </Field>
        <Field
          id="register-password"
          label="Kata sandi"
          hint="Minimal 10 karakter. Gunakan frasa yang mudah kamu ingat dan sulit ditebak."
          error={errors.password}
        >
          <input
            id="register-password"
            name="password"
            type="password"
            autoComplete="new-password"
            value={password}
            onChange={(event) => setPassword(event.target.value)}
            aria-describedby={describedBy(
              "register-password",
              "Minimal 10 karakter. Gunakan frasa yang mudah kamu ingat dan sulit ditebak.",
              errors.password,
            )}
            aria-invalid={errors.password ? true : undefined}
            minLength={10}
            maxLength={128}
            required
          />
        </Field>
        <Field id="register-confirm" label="Ulangi kata sandi" error={errors.confirm}>
          <input
            id="register-confirm"
            name="confirm"
            type="password"
            autoComplete="new-password"
            value={confirm}
            onChange={(event) => setConfirm(event.target.value)}
            aria-describedby={describedBy("register-confirm", undefined, errors.confirm)}
            aria-invalid={errors.confirm ? true : undefined}
            maxLength={128}
            required
          />
        </Field>
        <Button type="submit" disabled={pending}>
          {pending ? "Mendaftar…" : "Buat akun"}
        </Button>
        <p className="auth-page__switch">
          Sudah punya akun? <Link to="/masuk">Masuk</Link>
        </p>
      </form>
    </section>
  );
}

import "./AuthPage.css";
