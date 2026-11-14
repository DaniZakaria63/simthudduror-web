package post

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

// SQLRepository is the SQLite adapter for posts and likes.
type SQLRepository struct {
	db *sql.DB
}

// NewRepository returns the post repository.
func NewRepository(db *sql.DB) *SQLRepository {
	return &SQLRepository{db: db}
}

const selectPost = `
	SELECT p.id, p.author_id, a.name, p.body, p.created_at,
	       (SELECT COUNT(*) FROM likes l WHERE l.post_id = p.id) AS like_count,
	       EXISTS (SELECT 1 FROM likes l WHERE l.post_id = p.id AND l.account_id = ?) AS liked_by_me
	FROM posts p
	JOIN accounts a ON a.id = p.author_id`

func (r *SQLRepository) Create(ctx context.Context, post Post) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO posts (id, author_id, body, created_at) VALUES (?, ?, ?, ?)`,
		post.ID, post.AuthorID, post.Body, formatTime(post.CreatedAt))
	if err != nil {
		return fmt.Errorf("insert post: %w", err)
	}
	return nil
}

func (r *SQLRepository) List(ctx context.Context, query string, limit int, viewerID string) ([]Post, error) {
	pattern := "%" + escapeLike(query) + "%"
	rows, err := r.db.QueryContext(ctx, selectPost+`
		WHERE (? = '' OR p.body LIKE ? ESCAPE '\')
		ORDER BY p.created_at DESC, p.id DESC
		LIMIT ?`, viewerID, query, pattern, limit)
	if err != nil {
		return nil, fmt.Errorf("list posts: %w", err)
	}
	defer rows.Close()

	posts := make([]Post, 0)
	for rows.Next() {
		post, err := scanPost(rows)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate posts: %w", err)
	}
	return posts, nil
}

func (r *SQLRepository) FindByID(ctx context.Context, id, viewerID string) (Post, error) {
	row := r.db.QueryRowContext(ctx, selectPost+` WHERE p.id = ?`, viewerID, id)
	post, err := scanPost(row)
	if errors.Is(err, sql.ErrNoRows) {
		return Post{}, ErrPostNotFound
	}
	return post, err
}

func (r *SQLRepository) Like(ctx context.Context, postID, accountID string, at time.Time) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT OR IGNORE INTO likes (post_id, account_id, created_at) VALUES (?, ?, ?)`,
		postID, accountID, formatTime(at))
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "foreign key constraint") {
			return ErrPostNotFound
		}
		return fmt.Errorf("insert like: %w", err)
	}
	return nil
}

func (r *SQLRepository) Unlike(ctx context.Context, postID, accountID string) error {
	if _, err := r.db.ExecContext(ctx, `DELETE FROM likes WHERE post_id = ? AND account_id = ?`, postID, accountID); err != nil {
		return fmt.Errorf("delete like: %w", err)
	}
	return nil
}

type scanner interface {
	Scan(dest ...any) error
}

func scanPost(row scanner) (Post, error) {
	var post Post
	var createdAt string
	var liked int
	if err := row.Scan(&post.ID, &post.AuthorID, &post.AuthorName, &post.Body, &createdAt, &post.LikeCount, &liked); err != nil {
		return Post{}, err
	}
	parsed, err := time.Parse(persistedTimeLayout, createdAt)
	if err != nil {
		return Post{}, fmt.Errorf("parse post timestamp: %w", err)
	}
	post.CreatedAt = parsed
	post.LikedByMe = liked != 0
	return post, nil
}

const persistedTimeLayout = "2006-01-02T15:04:05.000000000Z"

func formatTime(t time.Time) string { return t.UTC().Format(persistedTimeLayout) }

func escapeLike(value string) string {
	replacer := strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)
	return replacer.Replace(value)
}
