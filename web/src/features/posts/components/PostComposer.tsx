import { useState, type FormEvent } from "react";
import { ApiError } from "@/shared/api/client";
import { Button } from "@/shared/components/Button";
import { Field, describedBy } from "@/shared/components/Field";
import { StatusMessage } from "@/shared/components/StatusMessage";
import { createPost } from "../api";

type PostComposerProps = {
  onCreated: () => void;
};

export function PostComposer({ onCreated }: PostComposerProps) {
  const [body, setBody] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);
  const [pending, setPending] = useState(false);
  const remaining = 2000 - body.length;

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setError(null);
    setSuccess(null);
    if (body.trim().length < 10) {
      setError("Kabar harus minimal 10 karakter.");
      return;
    }
    setPending(true);
    try {
      await createPost(body.trim());
      setBody("");
      setSuccess("Kabar diterbitkan.");
      onCreated();
    } catch (caught) {
      setError(caught instanceof ApiError ? caught.message : "Kabar tidak dapat diterbitkan.");
    } finally {
      setPending(false);
    }
  }

  return (
    <section className="post-composer" aria-labelledby="composer-title">
      <h2 id="composer-title">Tulis kabar</h2>
      <form onSubmit={handleSubmit} noValidate>
        {success ? <StatusMessage tone="success">{success}</StatusMessage> : null}
        <Field
          id="post-body"
          label="Isi kabar"
          hint="Tulis 10–2.000 karakter. Sebutkan tempat dan waktu bila relevan."
          error={error ?? undefined}
        >
          <textarea
            id="post-body"
            name="body"
            value={body}
            onChange={(event) => {
              setBody(event.target.value);
              setError(null);
            }}
            maxLength={2000}
            minLength={10}
            aria-describedby={describedBy(
              "post-body",
              "Tulis 10–2.000 karakter. Sebutkan tempat dan waktu bila relevan.",
              error ?? undefined,
            )}
            aria-invalid={error ? true : undefined}
            required
          />
        </Field>
        <div className="post-composer__footer">
          <span className="post-composer__count" aria-live="polite">
            {remaining} karakter tersisa
          </span>
          <Button type="submit" disabled={pending}>
            {pending ? "Menerbitkan…" : "Terbitkan kabar"}
          </Button>
        </div>
      </form>
    </section>
  );
}

import "./PostComposer.css";
