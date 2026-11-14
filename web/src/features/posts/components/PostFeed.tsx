import type { ReactNode } from "react";
import { Button } from "@/shared/components/Button";
import { LoadingState } from "@/shared/components/LoadingState";
import { StatusMessage } from "@/shared/components/StatusMessage";
import { PostCard } from "./PostCard";
import type { FeedState, Post } from "../types";

type PostFeedProps = {
  feed: FeedState;
  emptyMessage: string;
  filteredEmptyMessage: string;
  onToggleLike: (post: Post) => void;
  likePendingId: string | null;
  feedback?: ReactNode;
};

export function PostFeed({
  feed,
  emptyMessage,
  filteredEmptyMessage,
  onToggleLike,
  likePendingId,
  feedback,
}: PostFeedProps) {
  const searching = feed.query.trim().length > 0;

  return (
    <div className="post-feed">
      <div className="post-feed__filter">
        <label htmlFor="feed-search">Cari kabar</label>
        <div className="post-feed__filter-row">
          <input
            id="feed-search"
            name="q"
            type="search"
            value={feed.query}
            onChange={(event) => feed.setQuery(event.target.value)}
            placeholder="Kata kunci, nama masjid, atau kota"
            autoComplete="off"
            maxLength={100}
          />
          {searching ? (
            <button type="button" className="post-feed__clear" onClick={() => feed.setQuery("")}>
              Hapus
            </button>
          ) : null}
        </div>
      </div>

      {feedback}

      {feed.status === "loading" && feed.posts.length === 0 ? (
        <LoadingState label="Memuat kabar…" />
      ) : null}

      {feed.status === "error" ? (
        <div className="post-feed__error">
          <StatusMessage tone="error">{feed.error ?? "Kabar tidak dapat dimuat."}</StatusMessage>
          <Button tone="quiet" type="button" onClick={feed.reload}>
            Coba lagi
          </Button>
        </div>
      ) : null}

      {feed.status !== "error" && feed.posts.length === 0 && feed.status === "ready" ? (
        <p className="post-feed__empty">{searching ? filteredEmptyMessage : emptyMessage}</p>
      ) : null}

      {feed.status === "loading" && feed.posts.length > 0 ? (
        <p className="post-feed__loading-more" role="status">
          Memuat hasil terbaru…
        </p>
      ) : null}

      <ul className="post-feed__list">
        {feed.posts.map((post) => (
          <li key={post.id}>
            <PostCard
              post={post}
              onToggleLike={onToggleLike}
              likePending={likePendingId === post.id}
            />
          </li>
        ))}
      </ul>
    </div>
  );
}

import "./PostFeed.css";
