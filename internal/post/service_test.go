package post

import (
	"context"
	"errors"
	"strings"
	"testing"
)

func TestCreateRequiresCurrentValidation(t *testing.T) {
	service, _, db := newTestService(t)
	seedAccount(t, db, "validated-id", "Valid")
	seedAccount(t, db, "unvalidated-id", "Belum")
	seedAccount(t, db, "revoked-id", "Dicabut")
	authorizer := fakeAuthorizer{statuses: map[string]string{
		"validated-id":   "validated",
		"unvalidated-id": "unvalidated",
		"revoked-id":     "validated",
	}}
	service.authorizer = authorizer
	ctx := context.Background()

	if _, err := service.Create(ctx, "", "Kabar komunitas yang layak dibaca."); !errors.Is(err, ErrUnauthenticated) {
		t.Fatalf("anonymous create error = %v", err)
	}
	if _, err := service.Create(ctx, "unvalidated-id", "Kabar komunitas yang layak dibaca."); !errors.Is(err, ErrNotValidated) {
		t.Fatalf("unvalidated create error = %v", err)
	}

	created, err := service.Create(ctx, "validated-id", "  Kajian Maulid di masjid kampung berlangsung santai.  ")
	if err != nil {
		t.Fatalf("validated create: %v", err)
	}
	if created.ID == "" || created.AuthorName != "Valid" || strings.HasPrefix(created.Body, " ") {
		t.Fatalf("unexpected post: %+v", created)
	}

	// Revocation takes effect on the next call, with no cached claim.
	authorizer.statuses["revoked-id"] = "unvalidated"
	service.authorizer = authorizer
	if _, err := service.Create(ctx, "revoked-id", "Kabar lain yang seharusnya ditolak."); !errors.Is(err, ErrNotValidated) {
		t.Fatalf("revoked create error = %v", err)
	}
}

func TestCreateRejectsInvalidBodies(t *testing.T) {
	service, _, db := newTestService(t)
	seedAccount(t, db, "valid", "Valid")
	service.authorizer = fakeAuthorizer{statuses: map[string]string{"valid": "validated"}}
	ctx := context.Background()
	for _, body := range []string{"pendek", strings.Repeat("a", 2001), "baris\x00terlarang dalam kabar"} {
		var input *InputError
		if _, err := service.Create(ctx, "valid", body); !errors.As(err, &input) {
			t.Fatalf("body %q error = %v", body[:min(10, len(body))], err)
		}
	}
}

func TestCreateMapsAuthorizerFailureToUnauthenticated(t *testing.T) {
	service, _, _ := newTestService(t)
	service.authorizer = fakeAuthorizer{err: errors.New("database down")}
	if _, err := service.Create(context.Background(), "someone", "Kabar yang cukup panjang."); !errors.Is(err, ErrUnauthenticated) {
		t.Fatalf("authorizer failure error = %v", err)
	}
}

func TestLikeAndUnlikeAreIdempotentForAnyAccount(t *testing.T) {
	service, _, db := newTestService(t)
	seedAccount(t, db, "author", "Penulis")
	seedAccount(t, db, "reader", "Pembaca")
	service.authorizer = fakeAuthorizer{statuses: map[string]string{"author": "validated"}}
	ctx := context.Background()

	created, err := service.Create(ctx, "author", "Kabar Maulid yang bisa disukai semua akun.")
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	// The reader is unvalidated yet may interact.
	liked, err := service.Like(ctx, "reader", created.ID)
	if err != nil {
		t.Fatalf("like: %v", err)
	}
	if liked.LikeCount != 1 || !liked.LikedByMe {
		t.Fatalf("unexpected like result: %+v", liked)
	}
	liked, err = service.Like(ctx, "reader", created.ID)
	if err != nil || liked.LikeCount != 1 {
		t.Fatalf("repeated like should be idempotent: %+v %v", liked, err)
	}
	unliked, err := service.Unlike(ctx, "reader", created.ID)
	if err != nil || unliked.LikeCount != 0 || unliked.LikedByMe {
		t.Fatalf("unlike failed: %+v %v", unliked, err)
	}
	unliked, err = service.Unlike(ctx, "reader", created.ID)
	if err != nil || unliked.LikeCount != 0 {
		t.Fatalf("repeated unlike should be idempotent: %+v %v", unliked, err)
	}
}

func TestInteractionRejectsAnonymousAndMissingPost(t *testing.T) {
	service, _, _ := newTestService(t)
	ctx := context.Background()
	if _, err := service.Like(ctx, "", "post"); !errors.Is(err, ErrUnauthenticated) {
		t.Fatalf("anonymous like error = %v", err)
	}
	if _, err := service.Unlike(ctx, "", "post"); !errors.Is(err, ErrUnauthenticated) {
		t.Fatalf("anonymous unlike error = %v", err)
	}
	if _, err := service.Like(ctx, "reader", "missing"); !errors.Is(err, ErrPostNotFound) {
		t.Fatalf("missing post like error = %v", err)
	}
	if _, err := service.Unlike(ctx, "reader", "missing"); !errors.Is(err, ErrPostNotFound) {
		t.Fatalf("missing post unlike error = %v", err)
	}
}

func TestListFilterAndViewerState(t *testing.T) {
	service, _, db := newTestService(t)
	seedAccount(t, db, "author", "Penulis")
	seedAccount(t, db, "viewer", "Pengunjung")
	service.authorizer = fakeAuthorizer{statuses: map[string]string{"author": "validated"}}
	ctx := context.Background()
	first, _ := service.Create(ctx, "author", "Kabar pertama tentang Maulid di Surabaya.")
	if _, err := service.Create(ctx, "author", "Kabar kedua tentang kajian di Bandung."); err != nil {
		t.Fatalf("create second: %v", err)
	}
	if _, err := service.Like(ctx, "viewer", first.ID); err != nil {
		t.Fatalf("like: %v", err)
	}

	all, err := service.List(ctx, "", 30, "viewer")
	if err != nil || len(all) != 2 {
		t.Fatalf("list all = %d, %v", len(all), err)
	}
	found := false
	for _, item := range all {
		if item.ID == first.ID && item.LikedByMe && item.LikeCount == 1 {
			found = true
		}
	}
	if !found {
		t.Fatal("viewer like state missing from feed")
	}

	filtered, err := service.List(ctx, "bandung", 10, "")
	if err != nil || len(filtered) != 1 || !strings.Contains(strings.ToLower(filtered[0].Body), "bandung") {
		t.Fatalf("filter failed: %+v %v", filtered, err)
	}
	anonymous, err := service.List(ctx, "", 30, "")
	if err != nil || anonymous[0].LikedByMe {
		t.Fatalf("anonymous feed leaked like state: %+v %v", anonymous, err)
	}
	if _, err := service.List(ctx, strings.Repeat("a", 101), 10, ""); err == nil {
		t.Fatal("oversized search accepted")
	}
}
