import { apiRequest } from "@/shared/api/client";
import type { Post } from "./types";

type PostsResponse = { posts: Post[] };
type PostResponse = { post: Post };

export async function listPosts(query: string, signal?: AbortSignal): Promise<Post[]> {
  const params = new URLSearchParams();
  if (query.trim()) {
    params.set("q", query.trim());
  }
  const suffix = params.toString() ? `?${params.toString()}` : "";
  const payload = await apiRequest<PostsResponse>(`/api/v1/posts${suffix}`, { signal });
  return payload.posts;
}

export async function createPost(body: string): Promise<Post> {
  const payload = await apiRequest<PostResponse>("/api/v1/posts", {
    method: "POST",
    body: { body },
  });
  return payload.post;
}

export async function likePost(id: string): Promise<Post> {
  const payload = await apiRequest<PostResponse>(`/api/v1/posts/${encodeURIComponent(id)}/like`, {
    method: "PUT",
  });
  return payload.post;
}

export async function unlikePost(id: string): Promise<Post> {
  const payload = await apiRequest<PostResponse>(`/api/v1/posts/${encodeURIComponent(id)}/like`, {
    method: "DELETE",
  });
  return payload.post;
}
