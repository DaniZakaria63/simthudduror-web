import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { ApiError } from "@/shared/api/client";
import { Button } from "@/shared/components/Button";
import { LoadingState } from "@/shared/components/LoadingState";
import { StatusMessage } from "@/shared/components/StatusMessage";
import { PostComposer, PostFeed, likePost, unlikePost, useFeed } from "@/features/posts";
import type { Post } from "@/features/posts";
import { useSession } from "../session";

export function AccountPage() {
  const { status, account, error: sessionError, reload, signOut } = useSession();
  const navigate = useNavigate();
  const feed = useFeed();
  const [pendingId, setPendingId] = useState<string | null>(null);
  const [feedback, setFeedback] = useState<string | null>(null);
  const [feedbackTone, setFeedbackTone] = useState<"error" | "success">("success");

  async function toggleLike(post: Post) {
    if (!account) {
      navigate("/masuk");
      return;
    }
    setPendingId(post.id);
    setFeedback(null);
    try {
      const updated = post.liked_by_me ? await unlikePost(post.id) : await likePost(post.id);
      feed.updatePost(updated);
      setFeedbackTone("success");
      setFeedback(updated.liked_by_me ? "Kabar disukai." : "Suka pada kabar dibatalkan.");
    } catch (caught) {
      setFeedbackTone("error");
      setFeedback(caught instanceof ApiError ? caught.message : "Interaksi tidak dapat disimpan.");
    } finally {
      setPendingId(null);
    }
  }

  async function handleSignOut() {
    await signOut();
    navigate("/", { replace: true });
  }

  if (status === "loading") {
    return (
      <section className="account-page" aria-labelledby="akun-title">
        <h1 id="akun-title">Akun</h1>
        <LoadingState label="Memeriksa sesi…" />
      </section>
    );
  }

  if (status === "error") {
    return (
      <section className="account-page" aria-labelledby="akun-title">
        <h1 id="akun-title">Akun</h1>
        <StatusMessage tone="error">{sessionError ?? "Sesi tidak dapat diperiksa."}</StatusMessage>
        <Button type="button" onClick={() => void reload()}>
          Coba lagi
        </Button>
      </section>
    );
  }

  if (!account) {
    return (
      <section className="account-page account-page--anonymous" aria-labelledby="akun-title">
        <h1 id="akun-title">Akun komunitas</h1>
        <p>
          Masuk untuk menyukai kabar dan melihat status validasimu. Setiap akun dapat berinteraksi;
          hanya akun yang sudah divalidasi staf yang dapat menerbitkan kabar.
        </p>
        <div className="account-page__actions">
          <Link className="button button--primary" to="/masuk">
            Masuk
          </Link>
          <Link className="button button--quiet" to="/daftar">
            Daftar akun
          </Link>
        </div>
      </section>
    );
  }

  const validated = account.validation_status === "validated";

  return (
    <section className="account-page" aria-labelledby="akun-title">
      <header className="account-page__header">
        <div>
          <h1 id="akun-title">Akun</h1>
          <p className="account-page__identity">
            {account.name} · {account.email}
          </p>
        </div>
        <Button tone="quiet" type="button" onClick={() => void handleSignOut()}>
          Keluar
        </Button>
      </header>

      <div className={`account-page__status ${validated ? "is-validated" : ""}`.trim()}>
        <p className="account-page__status-label">
          {validated ? "Sudah divalidasi" : "Belum divalidasi"}
        </p>
        <p>
          {validated
            ? "Akunmu dapat menerbitkan kabar. Kamu tetap dapat menyukai kabar mana pun."
            : "Kamu tetap dapat membaca dan menyukai kabar. Minta staf komunitas meninjau akunmu untuk membuka penerbitan."}
        </p>
        {account.role === "staff" ? (
          <p>
            <Link to="/backoffice">Buka backoffice staf</Link> untuk meninjau akun lain.
          </p>
        ) : null}
      </div>

      {validated ? <PostComposer onCreated={feed.reload} /> : null}

      <div className="account-page__feed" id="kabar">
        <h2>Kabar komunitas</h2>
        <PostFeed
          feed={feed}
          emptyMessage="Belum ada kabar. Akun yang sudah divalidasi dapat membagikan kabar pertama."
          filteredEmptyMessage="Belum ada kabar yang cocok. Coba kata lain."
          onToggleLike={(post) => void toggleLike(post)}
          likePendingId={pendingId}
          feedback={
            feedback ? <StatusMessage tone={feedbackTone}>{feedback}</StatusMessage> : undefined
          }
        />
      </div>
    </section>
  );
}

import "./AccountPage.css";
