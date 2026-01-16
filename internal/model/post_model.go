package model

import (
	"database/sql"
	"time"
)

type PostModel struct {
	ID        int64
	UserID    int64
	Title     string
	Content   string
	DeletedAt sql.NullTime
	CreatedAt time.Time
	UpdatedAt time.Time
}
