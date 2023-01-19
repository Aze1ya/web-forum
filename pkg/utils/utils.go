package utils

import (
	"errors"
	"regexp"
	"strings"
	"unicode"

	"01.alem.school/git/Taimas/forum/internal/entity"
)

var (
	ErrUsernameNotUnique      = errors.New("username is not unique")
	ErrEmailNotUnique         = errors.New("email is not unique")
	ErrInvalidUsernameFormat  = errors.New("invalid username")
	ErrInvalidEmailFormat     = errors.New("invalid email")
	ErrInvalidPasswordFormat  = errors.New("invalid password")
	ErrUserNotFound           = errors.New("user does not exist ")
	ErrIncorrectPassword      = errors.New("password is incorrect")
	ErrEmailFieldNotFilled    = errors.New("email field not filled")
	ErrUsernameFieldNotFilled = errors.New("username field not filled")
	ErrPasswordFieldNotFilled = errors.New("password field not filled")
	ErrPasswordFieldNotSame   = errors.New("passwords are not same")
	ErrInvalidPostTitle       = errors.New("post title not filled")
	ErrInvalidPostBody        = errors.New("post body not filled")
	ErrInvalidPostCategories  = errors.New("post categories not filled")
	ErrInvalidPostId          = errors.New("invalid post id")
	ErrInvalidCommentId       = errors.New("invalid comment/post id")
	ErrInvalidCreator         = errors.New("invalid creator")
	ErrInvalidCommentBody     = errors.New("invalid comment body")
	ErrSqlNotFound            = errors.New("row not found in DB")
	ErrUnauthorized           = errors.New("unauthorized")
)

func ValidateSignUpParams(username, email, password, passwordRepeat string) error {
	if password != passwordRepeat {
		return ErrPasswordFieldNotSame
	}
	if len(username) < 4 || len(username) >= 12 {
		return ErrInvalidUsernameFormat
	}
	for _, char := range username {
		if char <= 32 || char >= 127 {
			return ErrInvalidUsernameFormat
		}
	}

	validEmail, err := regexp.MatchString(`[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`, email)
	if err != nil {
		return err
	}
	if !validEmail {
		return ErrInvalidEmailFormat
	}

	if !validatePassword(password) {
		return ErrInvalidPasswordFormat
	}

	return nil
}

func validatePassword(s string) bool {
	var (
		hasMinLen  = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)
	if len(s) >= 7 && len(s) <=20 {
		hasMinLen = true
	}
	for _, char := range s {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}

func ValidateCreateCommentParams(postId int, creator, body string) error {
	if postId == 0 {
		return ErrInvalidPostId
	}
	if creator == "" {
		return ErrInvalidCreator
	}

	if body == "" || body == " " {
		return ErrInvalidCommentBody
	}
	return nil
}

func ValidateCreatePostParams(title, body string, categories []string) error {
	if title == "" || title == " " {
		return ErrInvalidPostTitle
	}
	body = strings.Join(strings.Fields(body), " ")
	if body == "" || body == " " {
		return ErrInvalidPostBody
	}
	if len(categories) == 0 {
		return ErrInvalidPostCategories
	}
	for i:=0; i<len(categories); i++{
		if len(categories[i]) > 30 {
			return ErrInvalidPostCategories
		}
	}

	return nil
}

func AddDateFrontToPosts(posts []entity.Post) []entity.Post {
	for i := 0; i < len(posts); i++ {
		posts[i].CreationDateFront = posts[i].CreationDate.Format("02.01.2006 15:04")
	}
	return posts
}

func AddDateFrontToComments(comments []entity.Comment) []entity.Comment {
	for i := 0; i < len(comments); i++ {
		comments[i].CreationDateFront = comments[i].CreationDate.Format("02.01.2006 15:04")
	}
	return comments
}
