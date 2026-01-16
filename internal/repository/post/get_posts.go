package post

import (
	"context"
	"go-twitter/internal/model"
)

func (r *postRepository) GetPosts(ctx context.Context, limit, offset int) ([]*model.PostModel, error) {
	query := `SELECT id, user_id, title, content, deleted_at, created_at, updated_at FROM posts WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT ? OFFSET ?`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*model.PostModel
	for rows.Next() {
		var post model.PostModel
		err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.DeletedAt, &post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}
	return posts, nil
}

func (r *postRepository) GetPostsWithUserInfo(ctx context.Context, limit, offset int) ([]*model.PostModel, []string, error) {
	query := `
		SELECT p.id, p.user_id, p.title, p.content, p.deleted_at, p.created_at, p.updated_at, u.username
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE p.deleted_at IS NULL
		ORDER BY p.created_at DESC
		LIMIT ? OFFSET ?
	`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var posts []*model.PostModel
	var usernames []string
	for rows.Next() {
		var post model.PostModel
		var username string
		err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.DeletedAt, &post.CreatedAt, &post.UpdatedAt, &username)
		if err != nil {
			return nil, nil, err
		}
		posts = append(posts, &post)
		usernames = append(usernames, username)
	}
	return posts, usernames, nil
}

func (r *postRepository) GetPostsByUserID(ctx context.Context, userID int64, limit, offset int) ([]*model.PostModel, error) {
	query := `SELECT id, user_id, title, content, deleted_at, created_at, updated_at FROM posts WHERE user_id = ? AND deleted_at IS NULL ORDER BY created_at DESC LIMIT ? OFFSET ?`
	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*model.PostModel
	for rows.Next() {
		var post model.PostModel
		err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.DeletedAt, &post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}
	return posts, nil
}

func (r *postRepository) GetPostsCount(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM posts WHERE deleted_at IS NULL`
	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
