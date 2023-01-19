package usecase

import "01.alem.school/git/Taimas/forum/internal/repo"

type UseCases struct {
	Authorization Authorization
	Post          Post
	LikeDisLike   LikeDisLike
	Comment       Comment
	User          User
}

func NewUseCases(repo *repo.Repository) *UseCases {
	return &UseCases{
		Authorization: NewAuthUseCase(repo),
		Post:          NewPostUseCase(repo),
		LikeDisLike:   NewLikeDisLikeUseCase(repo),
		Comment:       NewCommentUseCase(repo),
		User:          NewUserUseCase(repo),
	}
}
