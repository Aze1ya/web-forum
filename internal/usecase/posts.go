package usecase

import (
	"fmt"
	"strings"

	"01.alem.school/git/Taimas/forum/internal/entity"
	"01.alem.school/git/Taimas/forum/pkg/utils"
)

type Post interface {
	CreatePostUseCase(post *entity.Post) (int, error)
	GetPostByIdUseCase(id int) (*entity.Post, error)
	GetAllPostsUseCase() ([]entity.Post, error)
	GetPostsByCategoryNameUseCase(category string) ([]entity.Post, error)
}

type PostRepository interface {
	CreatePost(post *entity.Post) (int64, error)
	GetPostById(id int) (*entity.Post, error)
	GetAllPosts() ([]entity.Post, error)
	GetPostsByCategoryName(category string) ([]entity.Post, error)
}

type PostUseCase struct {
	repo PostRepository
}

// Constructor
func NewPostUseCase(r PostRepository) *PostUseCase {
	return &PostUseCase{
		repo: r,
	}
}

// use cases

func (uc *PostUseCase) CreatePostUseCase(post *entity.Post) (int, error) {
	post.Title = strings.Join(strings.Fields(post.Title), " ")
	for i := 0; i < len(post.Categories); i++ {
		post.Categories[i] = strings.Join(strings.Fields(post.Categories[i]), " ")
	}
	err := utils.ValidateCreatePostParams(post.Title, post.Body, post.Categories)
	if err != nil {
		return 0, fmt.Errorf("use case: create post - %w", err)
	}

	postId, err := uc.repo.CreatePost(post)
	if err != nil {
		return 0, fmt.Errorf("use case: create post - %w", err)
	}
	return int(postId), nil
}

func (uc *PostUseCase) GetPostByIdUseCase(id int) (*entity.Post, error) {
	post, err := uc.repo.GetPostById(id)
	if err != nil {
		return nil, fmt.Errorf("use case: get post by id - %w", err)
	}

	return post, nil
}

func (uc *PostUseCase) GetAllPostsUseCase() ([]entity.Post, error) {
	posts, err := uc.repo.GetAllPosts()
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (uc *PostUseCase) GetPostsByCategoryNameUseCase(category string) ([]entity.Post, error) {
	posts, err := uc.repo.GetPostsByCategoryName(category)
	if err != nil {
		return nil, err
	}
	return posts, nil
}
