package usecase

import (
	"errors"
	"fmt"
	"time"

	"01.alem.school/git/Taimas/forum/internal/entity"
	"01.alem.school/git/Taimas/forum/pkg/utils"
	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Authorization interface {
	SignUp(email, username, password, passwordRepeat string) error
	SignIn(username, password string) (string, time.Time, error)
	ParseToken(token string) (*entity.Session, error)
	LogOut(token string) error
}

type AuthRepository interface {
	CreateUser(user *entity.User) error
	GetUserByUsername(username string) (*entity.User, error)
	GetSessionByToken(token string) (*entity.Session, error)
	CreateSession(session *entity.Session) error
	DeleteSession(token string) error
}

type AuthUseCase struct {
	repo AuthRepository
}

// Constructor
func NewAuthUseCase(r AuthRepository) *AuthUseCase {
	return &AuthUseCase{
		repo: r,
	}
}

// SignIn - sign in to forum
func (uc *AuthUseCase) SignUp(email, username, password, passwordRepeat string) error {
	// validate sign up params
	if err := utils.ValidateSignUpParams(username, email, password, passwordRepeat); err != nil {
		return fmt.Errorf("use case: sign up - %w", err)
	}

	user := &entity.User{
		Email:        email,
		Username:     username,
		CreationDate: time.Now(),
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("use case: sign up - %w", err)
	}
	user.Password = string(hash)
	err = uc.repo.CreateUser(user)
	if err != nil {
		return fmt.Errorf("use case: sign up - %w", err)
	}
	return nil
}

// SignUp - sign up to forum
func (uc *AuthUseCase) SignIn(username, password string) (string, time.Time, error) {
	user, err := uc.repo.GetUserByUsername(username)
	if err != nil {
		if !errors.Is(err, utils.ErrSqlNotFound) {
			return "", time.Now(), fmt.Errorf("use case: sign in - %w", err)
		} else {
			return "", time.Now(), fmt.Errorf("use case: sign in - %w", utils.ErrUserNotFound)
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", time.Now(), fmt.Errorf("use case: sign in - %w", utils.ErrIncorrectPassword)
	}
	uuid1, err := uuid.NewV4()
	if err != nil {
		return "", time.Now(), fmt.Errorf("use case: sign in - %w", err)
	}

	session := &entity.Session{
		Username:     user.Username,
		Token:        uuid1.String(),
		TokenExpDate: time.Now().Add(time.Hour * 24),
	}
	if err := uc.repo.CreateSession(session); err != nil {
		return "", time.Now(), fmt.Errorf("use case: sing in - %w", err)
	}

	return session.Token, session.TokenExpDate, nil
}

func (uc *AuthUseCase) LogOut(token string) error {
	err := uc.repo.DeleteSession(token)
	if err != nil {
		return fmt.Errorf("use case: log out - %w", err)
	}
	return nil
}

func (uc *AuthUseCase) ParseToken(token string) (*entity.Session, error) {
	session, err := uc.repo.GetSessionByToken(token)
	if err != nil {
		if errors.Is(err, utils.ErrSqlNotFound) {
			return nil, fmt.Errorf("use case: parse token - %w", utils.ErrSqlNotFound)
		} else {
			return nil, fmt.Errorf("use case: parse token - %w", err)
		}
	}

	return session, nil
}
