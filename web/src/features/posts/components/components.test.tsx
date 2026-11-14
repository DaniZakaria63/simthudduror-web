import { fireEvent, render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, expect, it, vi } from "vitest";
import { jsonResponse, mockFetch, postFixture } from "@/test/fixtures";
import type { FeedState } from "../types";
import { PostCard } from "./PostCard";
import { PostComposer } from "./PostComposer";
import { PostFeed } from "./PostFeed";

function feedState(overrides: Partial<FeedState> = {}): FeedState {
  return {
    posts: [],
    status: "ready",
    error: null,
    query: "",
    setQuery: vi.fn(),
    reload: vi.fn(),
    updatePost: vi.fn(),
    ...overrides,
  };
}

describe("PostCard", () => {
  it("renders author, date, body paragraphs and like state", () => {
    const onToggleLike = vi.fn();
    const post = { ...postFixture, body: "Baris pertama.\nBaris kedua.", like_count: 2, liked_by_me: true };
    render(<PostCard post={post} onToggleLike={onToggleLike} />);
    expect(screen.getByText("Bima")).toBeInTheDocument();
    expect(screen.getByText("Baris pertama.")).toBeInTheDocument();
    expect(screen.getByText("Baris kedua.")).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "Batal suka (2)" })).toHaveAttribute(
      "aria-pressed",
      "true",
    );
    fireEvent.click(screen.getByRole("button", { name: "Batal suka (2)" }));
    expect(onToggleLike).toHaveBeenCalledWith(post);
  });

  it("disables the like button while saving", () => {
    render(<PostCard post={postFixture} onToggleLike={vi.fn()} likePending />);
    expect(screen.getByRole("button", { name: "Menyimpan…" })).toBeDisabled();
  });
});

describe("PostFeed", () => {
  it("shows loading, error and retry states", () => {
    const reload = vi.fn();
    const errorFeed = feedState({ status: "error", error: "Server sibuk.", reload });
    const { rerender } = render(
      <PostFeed
        feed={feedState({ status: "loading" })}
        emptyMessage="kosong"
        filteredEmptyMessage="tidak cocok"
        onToggleLike={vi.fn()}
        likePendingId={null}
      />,
    );
    expect(screen.getByText("Memuat kabar…")).toBeInTheDocument();

    rerender(
      <PostFeed
        feed={errorFeed}
        emptyMessage="kosong"
        filteredEmptyMessage="tidak cocok"
        onToggleLike={vi.fn()}
        likePendingId={null}
      />,
    );
    expect(screen.getByRole("alert")).toHaveTextContent("Server sibuk.");
    fireEvent.click(screen.getByRole("button", { name: "Coba lagi" }));
    expect(reload).toHaveBeenCalled();
  });

  it("shows the right empty message and clears the filter", () => {
    const setQuery = vi.fn();
    const { rerender } = render(
      <PostFeed
        feed={feedState()}
        emptyMessage="Belum ada kabar."
        filteredEmptyMessage="Tidak cocok."
        onToggleLike={vi.fn()}
        likePendingId={null}
      />,
    );
    expect(screen.getByText("Belum ada kabar.")).toBeInTheDocument();
    expect(screen.queryByRole("button", { name: "Hapus" })).not.toBeInTheDocument();

    rerender(
      <PostFeed
        feed={feedState({ query: "jakarta", setQuery })}
        emptyMessage="Belum ada kabar."
        filteredEmptyMessage="Tidak cocok."
        onToggleLike={vi.fn()}
        likePendingId={null}
      />,
    );
    expect(screen.getByText("Tidak cocok.")).toBeInTheDocument();
    fireEvent.click(screen.getByRole("button", { name: "Hapus" }));
    expect(setQuery).toHaveBeenCalledWith("");
  });

  it("renders posts and a loading-more status during refetch", () => {
    const { rerender } = render(
      <PostFeed
        feed={feedState({ posts: [postFixture] })}
        emptyMessage="kosong"
        filteredEmptyMessage="tidak cocok"
        onToggleLike={vi.fn()}
        likePendingId={null}
        feedback={<p>Kabar disukai.</p>}
      />,
    );
    expect(screen.getByText(postFixture.body)).toBeInTheDocument();
    expect(screen.getByText("Kabar disukai.")).toBeInTheDocument();

    rerender(
      <PostFeed
        feed={feedState({ posts: [postFixture], status: "loading" })}
        emptyMessage="kosong"
        filteredEmptyMessage="tidak cocok"
        onToggleLike={vi.fn()}
        likePendingId={postFixture.id}
      />,
    );
    expect(screen.getByText("Memuat hasil terbaru…")).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "Menyimpan…" })).toBeDisabled();
  });
});

describe("PostComposer", () => {
  it("requires a minimum length", async () => {
    const user = userEvent.setup();
    const onCreated = vi.fn();
    render(<PostComposer onCreated={onCreated} />);
    await user.type(screen.getByLabelText("Isi kabar"), "pendek");
    await user.click(screen.getByRole("button", { name: "Terbitkan kabar" }));
    expect(screen.getByRole("alert")).toHaveTextContent("Kabar harus minimal 10 karakter.");
    expect(onCreated).not.toHaveBeenCalled();
  });

  it("publishes a post and clears the form", async () => {
    mockFetch(() =>
      jsonResponse(
        {
          post: {
            ...postFixture,
            body: "Kabar baru yang cukup panjang.",
          },
        },
        201,
      ),
    );
    const user = userEvent.setup();
    const onCreated = vi.fn();
    render(<PostComposer onCreated={onCreated} />);
    await user.type(screen.getByLabelText("Isi kabar"), "Kabar baru yang cukup panjang.");
    await user.click(screen.getByRole("button", { name: "Terbitkan kabar" }));
    expect(await screen.findByText("Kabar diterbitkan.")).toBeInTheDocument();
    expect(onCreated).toHaveBeenCalledOnce();
    expect(screen.getByLabelText("Isi kabar")).toHaveValue("");
  });

  it("reports a server failure", async () => {
    mockFetch(() =>
      jsonResponse(
        { error: { code: "validation_required", message: "Akun perlu divalidasi." } },
        403,
      ),
    );
    const user = userEvent.setup();
    render(<PostComposer onCreated={vi.fn()} />);
    await user.type(screen.getByLabelText("Isi kabar"), "Kabar baru yang cukup panjang.");
    await user.click(screen.getByRole("button", { name: "Terbitkan kabar" }));
    expect(await screen.findByRole("alert")).toHaveTextContent("Akun perlu divalidasi.");
  });
});
