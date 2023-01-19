package usecase

import (
	"errors"
	"fmt"
	"strings"

	"01.alem.school/git/Taimas/forum/internal/entity"
	"01.alem.school/git/Taimas/forum/pkg/utils"
)

type Comment interface {
	CreateCommentUseCase(*entity.Comment) error
	// GetCommentByIdUseCase(id int) (*entity.Comment, error)
}

type CommentRepository interface {
	CreateComment(comment *entity.Comment) error
	GetPostById(id int) (*entity.Post, error)
	GetCommentById(id int) (*entity.Comment, error)
}

type CommentUseCase struct {
	repo CommentRepository
}

// Constructor
func NewCommentUseCase(r CommentRepository) *CommentUseCase {
	return &CommentUseCase{
		repo: r,
	}
}

func (uc *CommentUseCase) CreateCommentUseCase(comment *entity.Comment) error {
	comment.Body = strings.Join(strings.Fields(comment.Body), " ")
	if err := utils.ValidateCreateCommentParams(comment.PostID, comment.Creator, comment.Body); err != nil {
		return fmt.Errorf("use case: create comment - %w", err)
	}

	if _, err := uc.repo.GetPostById(comment.PostID); err != nil {
		if !errors.Is(err, utils.ErrSqlNotFound) {
			return fmt.Errorf("use case: create comment - %w", err)
		} else {
			return utils.ErrInvalidPostId
		}
	}

	err := uc.repo.CreateComment(comment)
	if err != nil {
		return err
	}
	return nil
}

func (uc *CommentUseCase) GetCommentByIdUseCase(id int) (*entity.Comment, error) {
	if id == 0 {
		return nil, utils.ErrInvalidPostId
	}
	comment, err := uc.repo.GetCommentById(id)
	if err != nil {
		if !errors.Is(err, utils.ErrSqlNotFound) {
			return nil, err
		} else {
			return nil, utils.ErrInvalidPostId
		}
	}
	return comment, nil
}
