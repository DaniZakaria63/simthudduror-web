# Component shape — annotated patterns

## A presentational component (`PostCard`)

```tsx
import { authorInitial, formatDate } from "@/shared/format";
import type { Post } from "../types";

type PostCardProps = {                 // 1. imports
  post: Post;
  onToggleLike: (post: Post) => void;
  likePending?: boolean;
};

export function PostCard({ post, onToggleLike, likePending = false }: PostCardProps) {
  const paragraphs = post.body.split(/\n+/).map((p) => p.trim()).filter(Boolean);
  const likeLabel = post.liked_by_me ? "Batal suka" : "Suka";   // 2. variables

  return (                                                      // 3. UI
    <article className="post-card">
      <header className="post-card__meta">…</header>
      <div className="post-card__body">
        {paragraphs.map((paragraph, index) => (
          <p key={`${post.id}-${index}`}>{paragraph}</p>
        ))}
      </div>
      <button
        type="button"
        className="post-card__like"
        aria-pressed={post.liked_by_me}
        onClick={() => onToggleLike(post)}
        disabled={likePending}
      >
        {likePending ? "Menyimpan…" : `${likeLabel} (${post.like_count})`}
      </button>
    </article>
  );
}

import "./PostCard.css";                                        // 4. styles
```

Why it is shaped this way:

- the component renders text, never HTML, so post bodies cannot inject markup;
- state (pending) is owned by the page and passed down; the card stays reusable;
- the accessible name of the button changes with state (`Suka`/`Batal suka`) and
  `aria-pressed` exposes the toggle;
- the stylesheet import is last and is the only one.

## A feature hook (`useFeed`)

```ts
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
        if (caught instanceof DOMException && caught.name === "AbortError") return;
        setError(caught instanceof ApiError ? caught.message : "Kabar tidak dapat dimuat.");
        setStatus("error");
      });
    return () => controller.abort();
  }, [deferredQuery, revision]);

  const reload = useCallback(() => setRevision((value) => value + 1), []);
  const updatePost = useCallback((updated: Post) => {
    setPosts((current) => current.map((post) => (post.id === updated.id ? updated : post)));
  }, []);

  return { posts, status, error, query, setQuery, reload, updatePost };
}
```

Points to copy:

- the hook owns fetching, cancellation, status and error text;
- `useDeferredValue` keeps typing responsive without a second store;
- `updatePost` applies the server's authoritative post after a like instead of
  guessing the new count;
- components consume the returned state; they do not import `api.ts` themselves
  unless they perform a write (composer) and even then the page passes the handler.

## Feature entry point

```ts
export { PostFeed } from "./components/PostFeed";
export { PostComposer } from "./components/PostComposer";
export { useFeed } from "./useFeed";
export { createPost, likePost, unlikePost, listPosts } from "./api";
export type { FeedState, FeedStatus, Post } from "./types";
```

If another feature needs something not exported here, add it deliberately to
`index.ts`; do not reach into `components/` from outside.

## Anti-patterns the lint and review reject

- `import { PostCard } from "@/features/posts/components/PostCard"` from another
  feature — use the index.
- `fetch("/api/v1/posts")` inside a component — put it in `api.ts` and call it from a
  hook/page.
- `dangerouslySetInnerHTML` for post bodies — render text; React escapes it.
- `localStorage.setItem("csrf", …)` — the token stays in the in-memory client.
- Styles above imports, multiple stylesheet imports, or `style={{ color: "#fff" }}` —
  use the trailing CSS import and tokens.
- A “global” `PostCard` in `shared/` — it is domain-specific; shared is for Button,
  Field, LoadingState, StatusMessage, header/footer.
