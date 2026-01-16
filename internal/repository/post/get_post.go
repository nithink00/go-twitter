package post

import (
	"context"
	"database/sql"
	"go-twitter/internal/model"
)

func (r *postRepository) GetPostByID(ctx context.Context, id int64) (*model.PostModel, error) {
	query := `SELECT id, user_id, title, content, deleted_at, created_at, updated_at FROM posts WHERE id = ? AND deleted_at IS NULL`
	row := r.db.QueryRowContext(ctx, query, id)

	var post model.PostModel
	err := row.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.DeletedAt, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &post, nil
}

func (r *postRepository) GetPostWithUserInfo(ctx context.Context, id int64) (*model.PostModel, string, error) {
	query := `
		SELECT p.id, p.user_id, p.title, p.content, p.deleted_at, p.created_at, p.updated_at, u.username
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE p.id = ? AND p.deleted_at IS NULL
	`
	row := r.db.QueryRowContext(ctx, query, id)

	var post model.PostModel
	var username string
	err := row.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.DeletedAt, &post.CreatedAt, &post.UpdatedAt, &username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, "", nil
		}
		return nil, "", err
	}
	return &post, username, nil
}
