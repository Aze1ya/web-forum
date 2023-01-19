package v1

import (
	"errors"
	"net/http"
	"strconv"

	"01.alem.school/git/Taimas/forum/pkg/utils"
)

func (h *Handler) createLike(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.ExecuteErrorPage(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	creator := r.Context().Value(ctxKeyUser).(string)
	postId, err := strconv.Atoi(r.URL.Query().Get("post_id"))
	if err != nil {
		h.ExecuteErrorPage(w, http.StatusBadRequest, err.Error())
		return
	}
	comment_id := r.URL.Query().Get("comment_id")
	commentId := 0
	if comment_id != "" {
		commentId, err = strconv.Atoi(comment_id)
		if err != nil {
			h.ExecuteErrorPage(w, http.StatusBadRequest, err.Error())
			return
		}
	}
	err = h.useCases.LikeDisLike.CreateLikeUseCase(creator, postId, commentId)
	if err != nil {
		if errors.Is(err, utils.ErrInvalidPostId) || errors.Is(err, utils.ErrInvalidCommentId) {
			h.ExecuteErrorPage(w, http.StatusBadRequest, err.Error())
			return
		}
		h.ExecuteErrorPage(w, http.StatusInternalServerError, err.Error())
		return
	}
	http.Redirect(w, r, "/posts?post_id="+strconv.Itoa(postId), http.StatusSeeOther)
}

func (h *Handler) createDisLike(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.ExecuteErrorPage(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	creator := r.Context().Value(ctxKeyUser).(string)

	postId, err := strconv.Atoi(r.URL.Query().Get("post_id"))
	if err != nil {
		h.ExecuteErrorPage(w, http.StatusBadRequest, err.Error())
		return
	}

	comment_id := r.URL.Query().Get("comment_id")
	commentId := 0
	if comment_id != "" {
		commentId, err = strconv.Atoi(comment_id)
		if err != nil {
			h.ExecuteErrorPage(w, http.StatusBadRequest, err.Error())
			return
		}
	}
	err = h.useCases.LikeDisLike.CreateDisLikeUseCase(creator, postId, commentId)
	if err != nil {
		if errors.Is(err, utils.ErrInvalidPostId) || errors.Is(err, utils.ErrInvalidCommentId) {
			h.ExecuteErrorPage(w, http.StatusBadRequest, err.Error())
			return
		}
		h.ExecuteErrorPage(w, http.StatusInternalServerError, err.Error())
		return
	}

	http.Redirect(w, r, "/posts?post_id="+strconv.Itoa(postId), http.StatusSeeOther)
}
