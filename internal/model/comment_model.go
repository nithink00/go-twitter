package model

import (
	"database/sql"
	"time"
)

type CommentModel struct {
	ID        int64        `db:"id"`
	PostID    int64        `db:"post_id"`
	UserID    int64        `db:"user_id"`
	Content   string       `db:"content"`
	DeletedAt sql.NullTime `db:"deleted_at"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt time.Time    `db:"updated_at"`
}
