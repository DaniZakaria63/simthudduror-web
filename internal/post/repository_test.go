package post

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestRepositoryOrderFilterAndLikes(t *testing.T) {
	db := newTestDB(t)
	repo := NewRepository(db)
	seedAccount(t, db, "a1", "Anisa")
	seedAccount(t, db, "a2", "Bima")
	ctx := context.Background()
	base := time.Date(2026, 11, 14, 8, 0, 0, 0, time.UTC)
	for i, body := range []string{"Kabar lama 100% jujur", "Kabar baru dengan garis_bawah"} {
		post := Post{ID: "p" + string(rune('1'+i)), AuthorID: "a1", Body: body, CreatedAt: base.Add(time.Duration(i) * time.Hour)}
		if err := repo.Create(ctx, post); err != nil {
			t.Fatalf("create: %v", err)
		}
	}
	if err := repo.Create(ctx, Post{ID: "p3", AuthorID: "a2", Body: "Kabar ketiga", CreatedAt: base.Add(2 * time.Hour)}); err != nil {
		t.Fatalf("create: %v", err)
	}
	if err := repo.Create(ctx, Post{ID: "bad", AuthorID: "missing", Body: "x", CreatedAt: base}); err == nil {
		t.Fatal("foreign-key violation accepted")
	}

	posts, err := repo.List(ctx, "", 2, "")
	if err != nil || len(posts) != 2 || posts[0].ID != "p3" {
		t.Fatalf("list = %+v, %v", posts, err)
	}

	literalPercent, err := repo.List(ctx, "100%", 10, "")
	if err != nil || len(literalPercent) != 1 || literalPercent[0].ID != "p1" {
		t.Fatalf("literal percent search = %+v, %v", literalPercent, err)
	}
	literalUnderscore, err := repo.List(ctx, "garis_", 10, "")
	if err != nil || len(literalUnderscore) != 1 || literalUnderscore[0].ID != "p2" {
		t.Fatalf("literal underscore search = %+v, %v", literalUnderscore, err)
	}

	if err := repo.Like(ctx, "p3", "a2", base); err != nil {
		t.Fatalf("like: %v", err)
	}
	if err := repo.Like(ctx, "p3", "a2", base); err != nil {
		t.Fatalf("repeated like: %v", err)
	}
	found, err := repo.FindByID(ctx, "p3", "a2")
	if err != nil || found.LikeCount != 1 || !found.LikedByMe || found.AuthorName != "Bima" {
		t.Fatalf("find liked = %+v, %v", found, err)
	}
	if err := repo.Unlike(ctx, "p3", "a2"); err != nil {
		t.Fatalf("unlike: %v", err)
	}
	found, err = repo.FindByID(ctx, "p3", "a2")
	if err != nil || found.LikeCount != 0 || found.LikedByMe {
		t.Fatalf("find unliked = %+v, %v", found, err)
	}
	if err := repo.Like(ctx, "missing", "a2", base); !errors.Is(err, ErrPostNotFound) {
		t.Fatalf("like missing error = %v", err)
	}
	if _, err := repo.FindByID(ctx, "missing", ""); !errors.Is(err, ErrPostNotFound) {
		t.Fatalf("find missing error = %v", err)
	}
}

func TestRepositoryRejectsMalformedPostTimestamp(t *testing.T) {
	db := newTestDB(t)
	repo := NewRepository(db)
	seedAccount(t, db, "a1", "Anisa")
	if _, err := db.Exec(`INSERT INTO posts (id, author_id, body, created_at) VALUES ('bad', 'a1', 'body', 'nope')`); err != nil {
		t.Fatalf("insert fixture: %v", err)
	}
	if _, err := repo.FindByID(context.Background(), "bad", ""); err == nil {
		t.Fatal("malformed timestamp accepted")
	}
}
