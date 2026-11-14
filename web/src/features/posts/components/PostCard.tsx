import { authorInitial, formatDate } from "@/shared/format";
import type { Post } from "../types";

type PostCardProps = {
  post: Post;
  onToggleLike: (post: Post) => void;
  likePending?: boolean;
};

export function PostCard({ post, onToggleLike, likePending = false }: PostCardProps) {
  const paragraphs = post.body
    .split(/\n+/)
    .map((paragraph) => paragraph.trim())
    .filter(Boolean);
  const likeLabel = post.liked_by_me ? "Batal suka" : "Suka";

  return (
    <article className="post-card">
      <header className="post-card__meta">
        <span className="post-card__avatar" aria-hidden="true">
          {authorInitial(post.author_name)}
        </span>
        <span className="post-card__author">{post.author_name}</span>
        <time className="post-card__date" dateTime={post.created_at}>
          {formatDate(post.created_at)}
        </time>
      </header>
      <div className="post-card__body">
        {paragraphs.map((paragraph, index) => (
          <p key={`${post.id}-${index}`}>{paragraph}</p>
        ))}
      </div>
      <div className="post-card__actions">
        <button
          type="button"
          className="post-card__like"
          aria-pressed={post.liked_by_me}
          onClick={() => onToggleLike(post)}
          disabled={likePending}
        >
          {likePending ? "Menyimpan…" : `${likeLabel} (${post.like_count})`}
        </button>
      </div>
    </article>
  );
}

import "./PostCard.css";
