import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { ApiError } from "@/shared/api/client";
import { StatusMessage } from "@/shared/components/StatusMessage";
import { useSession } from "@/features/account";
import { PostFeed, likePost, unlikePost, useFeed } from "@/features/posts";
import type { Post } from "@/features/posts";
import { ContributionSection } from "../components/ContributionSection";
import { Masthead } from "../components/Masthead";

export function LandingPage() {
  const { account } = useSession();
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

  return (
    <>
      <Masthead />
      <section className="feed-section" id="kabar" aria-labelledby="feed-title">
        <div className="feed-section__body">
          <div className="feed-section__rail">
            <h2 className="feed-section__title" id="feed-title" data-reveal>
              Kabar komunitas
            </h2>
            <div className="feed-section__note" data-reveal data-reveal-delay="1">
              <p>
                Semua akun dapat membaca dan menyukai kabar. Penerbitan dibuka setelah staf
                memvalidasi akun.
              </p>
              <p>
                <Link to="/akun">Kelola akunmu</Link> atau <Link to="/daftar">buat akun baru</Link>.
              </p>
            </div>
          </div>
          <div className="feed-section__feed" data-reveal data-reveal-delay="2">
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
        </div>
      </section>
      <ContributionSection />
    </>
  );
}

import "./LandingPage.css";
