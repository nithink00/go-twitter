package dto

type (
	CreateCommentRequest struct {
		Content string `json:"content" validate:"required,min=1"`
	}

	CreateCommentResponse struct {
		ID int64 `json:"id"`
	}
)

type (
	UpdateCommentRequest struct {
		Content string `json:"content" validate:"required,min=1"`
	}

	UpdateCommentResponse struct {
		ID int64 `json:"id"`
	}
)

type (
	CommentResponse struct {
		ID         int64  `json:"id"`
		PostID     int64  `json:"post_id"`
		UserID     int64  `json:"user_id"`
		Username   string `json:"username"`
		Content    string `json:"content"`
		LikesCount int    `json:"likes_count"`
		CreatedAt  string `json:"created_at"`
		UpdatedAt  string `json:"updated_at"`
	}

	CommentsResponse struct {
		Comments   []CommentResponse `json:"comments"`
		TotalCount int64             `json:"total_count"`
		Page       int               `json:"page"`
		PageSize   int               `json:"page_size"`
		TotalPages int               `json:"total_pages"`
	}
)
