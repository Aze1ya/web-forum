package usecase

import (
	"fmt"

	"01.alem.school/git/Taimas/forum/internal/entity"
	"01.alem.school/git/Taimas/forum/pkg/utils"
)

type LikeDisLikeRepository interface {
	CreateLikePost(like *entity.Like) error
	CreateDisLikePost(dislike *entity.DisLike) error
	CreateLikeComment(like *entity.Like) error
	CreateDisLikeComment(dislike *entity.DisLike) error
}

type LikeDisLike interface {
	CreateLikeUseCase(creator string, postID int, commentID int) error
	CreateDisLikeUseCase(creator string, postID int, commentID int) error
}

type LikeDisLikeUseCase struct {
	repo LikeDisLikeRepository
}

// Constructor
func NewLikeDisLikeUseCase(r LikeDisLikeRepository) *LikeDisLikeUseCase {
	return &LikeDisLikeUseCase{
		repo: r,
	}
}

func (uc *LikeDisLikeUseCase) CreateLikeUseCase(creator string, postID int, commentID int) error {
	if postID == 0 {
		return utils.ErrInvalidPostId
	}
	like := entity.Like{
		Creator:   creator,
		PostID:    postID,
		CommentID: commentID,
	}

	if commentID == 0 {
		err := uc.repo.CreateLikePost(&like)
		if err != nil {
			return fmt.Errorf("use case: create likepost - %w", err)
		}
	} else {
		err := uc.repo.CreateLikeComment(&like)
		if err != nil {
			return fmt.Errorf("use case: create likecomment - %w", err)
		}
	}

	return nil
}

func (uc *LikeDisLikeUseCase) CreateDisLikeUseCase(creator string, postID int, commentID int) error {
	if postID == 0 {
		return utils.ErrInvalidPostId
	}

	dislike := entity.DisLike{
		Creator:   creator,
		PostID:    postID,
		CommentID: commentID,
	}

	if commentID == 0 {
		err := uc.repo.CreateDisLikePost(&dislike)
		if err != nil {
			return fmt.Errorf("use case: create dislikepost - %w", err)
		}
	} else {
		err := uc.repo.CreateDisLikeComment(&dislike)
		if err != nil {
			return fmt.Errorf("use case: create dislikecomment - %w", err)
		}
	}
	return nil
}
