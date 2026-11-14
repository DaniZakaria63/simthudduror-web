import { useCallback, useDeferredValue, useEffect, useState } from "react";
import { ApiError } from "@/shared/api/client";
import { listPosts } from "./api";
import type { FeedState, Post } from "./types";

export function useFeed(): FeedState {
  const [query, setQuery] = useState("");
  const deferredQuery = useDeferredValue(query);
  const [posts, setPosts] = useState<Post[]>([]);
  const [status, setStatus] = useState<FeedState["status"]>("loading");
  const [error, setError] = useState<string | null>(null);
  const [revision, setRevision] = useState(0);

  useEffect(() => {
    const controller = new AbortController();
    setStatus("loading");
    listPosts(deferredQuery, controller.signal)
      .then((next) => {
        setPosts(next);
        setError(null);
        setStatus("ready");
      })
      .catch((caught: unknown) => {
        if (caught instanceof DOMException && caught.name === "AbortError") {
          return;
        }
        setError(caught instanceof ApiError ? caught.message : "Kabar tidak dapat dimuat.");
        setStatus("error");
      });
    return () => controller.abort();
  }, [deferredQuery, revision]);

  const reload = useCallback(() => setRevision((current) => current + 1), []);
  const updatePost = useCallback((updated: Post) => {
    setPosts((current) => current.map((post) => (post.id === updated.id ? updated : post)));
  }, []);

  return { posts, status, error, query, setQuery, reload, updatePost };
}
