package dto

type (
	CreatePostRequest struct {
		Title   string `json:"title" validate:"required,min=1,max=100"`
		Content string `json:"content" validate:"required,min=1"`
	}

	CreatePostResponse struct {
		ID int64 `json:"id"`
	}
)

type (
	UpdatePostRequest struct {
		Title   string `json:"title" validate:"required,min=1,max=100"`
		Content string `json:"content" validate:"required,min=1"`
	}

	UpdatePostResponse struct {
		ID int64 `json:"id"`
	}
)

type (
	PostResponse struct {
		ID        int64  `json:"id"`
		UserID    int64  `json:"user_id"`
		Username  string `json:"username"`
		Title     string `json:"title"`
		Content   string `json:"content"`
		LikesCount int   `json:"likes_count"`
		CommentsCount int `json:"comments_count"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}

	PostsResponse struct {
		Posts      []PostResponse `json:"posts"`
		TotalCount int64          `json:"total_count"`
		Page       int            `json:"page"`
		PageSize   int            `json:"page_size"`
		TotalPages int            `json:"total_pages"`
	}
)
