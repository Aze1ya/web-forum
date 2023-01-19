package usecase

import (
	"fmt"

	"01.alem.school/git/Taimas/forum/internal/entity"
)

type User interface {
	GetUserProfileUseCase(username string) (*entity.User, error)
	GetUserPostsUseCase(username string) ([]entity.Post, error)
	GetUsersCommentedPostsUseCase(username string) ([]entity.Post, error)
	GetUsersLikedPostsUseCase(username string) ([]entity.Post, error)
	GetUsersDislikedPostsUseCase(username string) ([]entity.Post, error)
}

type UserRepository interface {
	GetUserProfile(username string) (*entity.User, error)
	GetUserPosts(username string) ([]entity.Post, error)
	GetUsersCommentedPosts(username string) ([]entity.Post, error)
	GetUsersLikedPosts(username string) ([]entity.Post, error)
	GetUsersDislikedPosts(username string) ([]entity.Post, error)
}

type UserUseCase struct {
	repo UserRepository
}

func NewUserUseCase(r UserRepository) *UserUseCase {
	return &UserUseCase{
		repo: r,
	}
}

func (uc *UserUseCase) GetUserProfileUseCase(username string) (*entity.User, error) {
	user, err := uc.repo.GetUserProfile(username)
	if err != nil {
		return nil, fmt.Errorf("use case: get user profile - %w", err)
	}
	return user, nil
}

func (uc *UserUseCase) GetUserPostsUseCase(username string) ([]entity.Post, error) {
	_, err := uc.repo.GetUserProfile(username)
	if err != nil {
		return nil, fmt.Errorf("use case: get user profile - %w", err)
	}
	posts, err := uc.repo.GetUserPosts(username)
	if err != nil {
		return nil, fmt.Errorf("use case: get user posts - %w", err)
	}
	return posts, nil
}

func (uc *UserUseCase) GetUsersCommentedPostsUseCase(username string) ([]entity.Post, error) {
	_, err := uc.repo.GetUserProfile(username)
	if err != nil {
		return nil, fmt.Errorf("use case: get user profile - %w", err)
	}
	posts, err := uc.repo.GetUsersCommentedPosts(username)
	if err != nil {
		return nil, fmt.Errorf("use case: get user commented posts - %w", err)
	}
	return posts, nil
}

func (uc *UserUseCase) GetUsersLikedPostsUseCase(username string) ([]entity.Post, error) {
	_, err := uc.repo.GetUserProfile(username)
	if err != nil {
		return nil, fmt.Errorf("use case: get user profile - %w", err)
	}
	posts, err := uc.repo.GetUsersLikedPosts(username)
	if err != nil {
		return nil, fmt.Errorf("use case: get user liked posts - %w", err)
	}
	return posts, nil
}

func (uc *UserUseCase) GetUsersDislikedPostsUseCase(username string) ([]entity.Post, error) {
	_, err := uc.repo.GetUserProfile(username)
	if err != nil {
		return nil, fmt.Errorf("use case: get user profile - %w", err)
	}
	posts, err := uc.repo.GetUsersDislikedPosts(username)
	if err != nil {
		return nil, fmt.Errorf("use case: get user disliked posts - %w", err)
	}
	return posts, nil
}
