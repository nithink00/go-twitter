package dto

type (
	LikeResponse struct {
		Message string `json:"message"`
	}

	LikesCountResponse struct {
		Count int `json:"count"`
	}
)
