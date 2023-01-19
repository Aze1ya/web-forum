package v1

import (
	"errors"
	"net/http"

	"01.alem.school/git/Taimas/forum/internal/entity"
	"01.alem.school/git/Taimas/forum/pkg/utils"
)

type DataForUserpage struct {
	User  entity.User
	Posts []entity.Post
}

func (h *Handler) userProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.ExecuteErrorPage(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	username := ""
	if r.URL.RawQuery == "" {
		username = r.Context().Value(ctxKeyUser).(string)
	} else {
		username = r.URL.Query().Get("username")
	}
	data := new(DataForUserpage)
	var err error

	switch r.URL.Path {
	case "/users/profile":
		user, err := h.useCases.User.GetUserProfileUseCase(username)
		if err != nil {
			if errors.Is(err, utils.ErrUserNotFound) {
				h.ExecuteErrorPage(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
				return
			} else {
				h.ExecuteErrorPage(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
				return
			}
		}
		user.CreationDateFront = user.CreationDate.Format("02.01.2006 15:04")
		data.User = *user
		if err := h.tmpl.ExecuteTemplate(w, "userProfile.html", data); err != nil {
			h.ExecuteErrorPage(w, http.StatusInternalServerError, err.Error())
			return
		}
	case "/users/posts":
		data.Posts, err = h.useCases.User.GetUserPostsUseCase(username)
		if err != nil {
			if errors.Is(err, utils.ErrUserNotFound) {
				h.ExecuteErrorPage(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
				return
			} else {
				h.ExecuteErrorPage(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
				return
			}
		}
		data.Posts = utils.AddDateFrontToPosts(data.Posts)
		data.User.Username = username
		if err := h.tmpl.ExecuteTemplate(w, "userPosts.html", data); err != nil {
			h.ExecuteErrorPage(w, http.StatusInternalServerError, err.Error())
			return
		}
	case "/users/comments":
		data.Posts, err = h.useCases.User.GetUsersCommentedPostsUseCase(username)
		if err != nil {
			if errors.Is(err, utils.ErrUserNotFound) {
				h.ExecuteErrorPage(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
				return
			} else {
				h.ExecuteErrorPage(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
				return
			}
		}
		data.Posts = utils.AddDateFrontToPosts(data.Posts)
		data.User.Username = username
		if err := h.tmpl.ExecuteTemplate(w, "userCommentedPosts.html", data); err != nil {
			h.ExecuteErrorPage(w, http.StatusInternalServerError, err.Error())
			return
		}
	case "/users/likes":
		data.Posts, err = h.useCases.User.GetUsersLikedPostsUseCase(username)
		if err != nil {
			if errors.Is(err, utils.ErrUserNotFound) {
				h.ExecuteErrorPage(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
				return
			} else {
				h.ExecuteErrorPage(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
				return
			}
		}
		data.Posts = utils.AddDateFrontToPosts(data.Posts)
		data.User.Username = username
		if err := h.tmpl.ExecuteTemplate(w, "userLikedPosts.html", data); err != nil {
			h.ExecuteErrorPage(w, http.StatusInternalServerError, err.Error())
			return
		}
	case "/users/dislikes":
		data.Posts, err = h.useCases.User.GetUsersDislikedPostsUseCase(username)
		if err != nil {
			if errors.Is(err, utils.ErrUserNotFound) {
				h.ExecuteErrorPage(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
				return
			} else {
				h.ExecuteErrorPage(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
				return
			}
		}
		data.Posts = utils.AddDateFrontToPosts(data.Posts)
		data.User.Username = username
		if err := h.tmpl.ExecuteTemplate(w, "userDislikedPosts.html", data); err != nil {
			h.ExecuteErrorPage(w, http.StatusInternalServerError, err.Error())
			return
		}
	default:
		h.ExecuteErrorPage(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
}
