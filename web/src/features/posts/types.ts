export type Post = {
  id: string;
  author_id: string;
  author_name: string;
  body: string;
  created_at: string;
  like_count: number;
  liked_by_me: boolean;
};

export type FeedStatus = "loading" | "ready" | "error";

export type FeedState = {
  posts: Post[];
  status: FeedStatus;
  error: string | null;
  query: string;
  setQuery: (value: string) => void;
  reload: () => void;
  updatePost: (post: Post) => void;
};
