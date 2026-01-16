package post

import (
	"context"
	"go-twitter/internal/dto"
	"math"
	"net/http"
)

func (s *postService) GetPosts(ctx context.Context, page, pageSize int) (*dto.PostsResponse, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	posts, usernames, err := s.postRepo.GetPostsWithUserInfo(ctx, pageSize, offset)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	totalCount, err := s.postRepo.GetPostsCount(ctx)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	var postResponses []dto.PostResponse
	for i, post := range posts {
		likesCount, err := s.getPostLikesCount(ctx, post.ID)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}

		commentsCount, err := s.getPostCommentsCount(ctx, post.ID)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}

		postResponses = append(postResponses, dto.PostResponse{
			ID:            post.ID,
			UserID:        post.UserID,
			Username:      usernames[i],
			Title:         post.Title,
			Content:       post.Content,
			LikesCount:    likesCount,
			CommentsCount: commentsCount,
			CreatedAt:     post.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:     post.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))

	response := &dto.PostsResponse{
		Posts:      postResponses,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}

	return response, http.StatusOK, nil
}

func (s *postService) GetPostsByUserID(ctx context.Context, userID int64, page, pageSize int) (*dto.PostsResponse, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	posts, err := s.postRepo.GetPostsByUserID(ctx, userID, pageSize, offset)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	totalCount, err := s.postRepo.GetPostsCount(ctx)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	var postResponses []dto.PostResponse
	for _, post := range posts {
		likesCount, err := s.getPostLikesCount(ctx, post.ID)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}

		commentsCount, err := s.getPostCommentsCount(ctx, post.ID)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}

		postResponses = append(postResponses, dto.PostResponse{
			ID:            post.ID,
			UserID:        post.UserID,
			Title:         post.Title,
			Content:       post.Content,
			LikesCount:    likesCount,
			CommentsCount: commentsCount,
			CreatedAt:     post.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:     post.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))

	response := &dto.PostsResponse{
		Posts:      postResponses,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}

	return response, http.StatusOK, nil
}
