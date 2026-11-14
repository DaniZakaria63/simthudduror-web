// Package post owns the public feed, validated publishing and idempotent likes.
package post

import (
	"context"
	"errors"
	"time"
)

// Post is a public community news item. Email and other private account data
// are deliberately absent.
type Post struct {
	ID         string    `json:"id"`
	AuthorID   string    `json:"author_id"`
	AuthorName string    `json:"author_name"`
	Body       string    `json:"body"`
	CreatedAt  time.Time `json:"created_at"`
	LikeCount  int       `json:"like_count"`
	LikedByMe  bool      `json:"liked_by_me"`
}

// Repository is the persistence port owned by the post feature.
type Repository interface {
	Create(ctx context.Context, post Post) error
	List(ctx context.Context, query string, limit int, viewerID string) ([]Post, error)
	FindByID(ctx context.Context, id, viewerID string) (Post, error)
	Like(ctx context.Context, postID, accountID string, at time.Time) error
	Unlike(ctx context.Context, postID, accountID string) error
}

// Authorizer reads the current validation status from trusted persistence.
type Authorizer interface {
	ValidationStatus(ctx context.Context, accountID string) (string, error)
}

var (
	// ErrUnauthenticated is returned when no account is identified.
	ErrUnauthenticated = errors.New("unauthenticated")
	// ErrNotValidated is returned when an authenticated account lacks validation.
	ErrNotValidated = errors.New("account not validated")
	// ErrPostNotFound is returned for unknown post identifiers.
	ErrPostNotFound = errors.New("post not found")
)

// InputError is a safe validation failure for the form layer.
type InputError struct {
	Message string
}

func (e *InputError) Error() string { return e.Message }

// Service implements the feed and publishing rules.
type Service struct {
	repo       Repository
	authorizer Authorizer
	now        func() time.Time
}

// NewService constructs the post service.
func NewService(repo Repository, authorizer Authorizer) *Service {
	return &Service{repo: repo, authorizer: authorizer, now: time.Now}
}

// List returns the public feed, optionally filtered by a visitor's words.
func (s *Service) List(ctx context.Context, query string, limit int, viewerID string) ([]Post, error) {
	query = trimSpace(query)
	if len([]rune(query)) > 100 {
		return nil, &InputError{Message: "Pencarian maksimal 100 karakter."}
	}
	if limit <= 0 {
		limit = 30
	}
	if limit > 100 {
		limit = 100
	}
	return s.repo.List(ctx, query, limit, viewerID)
}

// Create publishes a post only after reloading the author's validation status
// from the database. A stale client claim or revoked validation cannot bypass
// this check.
func (s *Service) Create(ctx context.Context, authorID, body string) (Post, error) {
	if authorID == "" {
		return Post{}, ErrUnauthenticated
	}
	status, err := s.authorizer.ValidationStatus(ctx, authorID)
	if err != nil {
		return Post{}, ErrUnauthenticated
	}
	if status != "validated" {
		return Post{}, ErrNotValidated
	}
	body = trimSpace(body)
	if count := len([]rune(body)); count < 10 || count > 2000 {
		return Post{}, &InputError{Message: "Kabar harus 10 sampai 2.000 karakter."}
	}
	if containsControl(body) {
		return Post{}, &InputError{Message: "Kabar mengandung karakter yang tidak diizinkan."}
	}
	post := Post{
		ID: newID(), AuthorID: authorID, Body: body, CreatedAt: s.now().UTC(),
	}
	if err := s.repo.Create(ctx, post); err != nil {
		return Post{}, err
	}
	return s.repo.FindByID(ctx, post.ID, authorID)
}

// Like adds a like for any authenticated account. Repeating it is harmless.
func (s *Service) Like(ctx context.Context, accountID, postID string) (Post, error) {
	if accountID == "" {
		return Post{}, ErrUnauthenticated
	}
	if _, err := s.repo.FindByID(ctx, postID, accountID); err != nil {
		return Post{}, err
	}
	if err := s.repo.Like(ctx, postID, accountID, s.now().UTC()); err != nil {
		return Post{}, err
	}
	return s.repo.FindByID(ctx, postID, accountID)
}

// Unlike removes a like for any authenticated account, idempotently.
func (s *Service) Unlike(ctx context.Context, accountID, postID string) (Post, error) {
	if accountID == "" {
		return Post{}, ErrUnauthenticated
	}
	if _, err := s.repo.FindByID(ctx, postID, accountID); err != nil {
		return Post{}, err
	}
	if err := s.repo.Unlike(ctx, postID, accountID); err != nil {
		return Post{}, err
	}
	return s.repo.FindByID(ctx, postID, accountID)
}
