import { act, renderHook, waitFor } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import { jsonResponse, mockFetch, postFixture } from "@/test/fixtures";
import { useFeed } from "./useFeed";

describe("useFeed", () => {
  it("loads the public feed and filters with the query", async () => {
    const fetchMock = mockFetch(({ url }) => {
      if (url.toLowerCase().includes("q=bandung")) {
        return jsonResponse({ posts: [{ ...postFixture, id: "post-2", body: "Kabar dari Bandung." }] });
      }
      return jsonResponse({ posts: [postFixture] });
    });
    const { result } = renderHook(() => useFeed());
    await waitFor(() => expect(result.current.status).toBe("ready"));
    expect(result.current.posts).toHaveLength(1);

    act(() => result.current.setQuery("Bandung"));
    await waitFor(() => expect(result.current.posts[0]?.id).toBe("post-2"));
    expect(fetchMock.mock.calls.some(([url]) => String(url).includes("q=Bandung"))).toBe(true);
  });

  it("updates one post in place", async () => {
    mockFetch(() => jsonResponse({ posts: [postFixture] }));
    const { result } = renderHook(() => useFeed());
    await waitFor(() => expect(result.current.status).toBe("ready"));
    act(() =>
      result.current.updatePost({ ...postFixture, like_count: 3, liked_by_me: true }),
    );
    expect(result.current.posts[0]?.like_count).toBe(3);
    expect(result.current.posts[0]?.liked_by_me).toBe(true);
  });

  it("reports errors and retries", async () => {
    let attempts = 0;
    mockFetch(() => {
      attempts += 1;
      if (attempts === 1) {
        return jsonResponse({ error: { code: "internal_error", message: "Server sibuk." } }, 500);
      }
      return jsonResponse({ posts: [postFixture] });
    });
    const { result } = renderHook(() => useFeed());
    await waitFor(() => expect(result.current.status).toBe("error"));
    expect(result.current.error).toBe("Server sibuk.");
    act(() => result.current.reload());
    await waitFor(() => expect(result.current.status).toBe("ready"));
    expect(result.current.posts).toHaveLength(1);
  });
});
