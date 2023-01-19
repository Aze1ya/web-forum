package v1

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"01.alem.school/git/Taimas/forum/internal/entity"
	"01.alem.school/git/Taimas/forum/pkg/utils"
)

func (h *Handler) commentCreate(w http.ResponseWriter, r *http.Request) {
	// create comment for post
	if r.Method != http.MethodPost {
		h.ExecuteErrorPage(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	if r.URL.RawQuery == "" {
		h.ExecuteErrorPage(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	postId, err := strconv.Atoi(r.URL.Query().Get("post_id"))
	if err != nil || postId == 0 {
		h.ExecuteErrorPage(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	if err := r.ParseForm(); err != nil {
		h.ExecuteErrorPage(w, http.StatusInternalServerError, err.Error())
		return
	}

	body, ok := r.Form["body"]
	if !ok {
		h.ExecuteErrorPage(w, http.StatusBadRequest, utils.ErrInvalidCommentBody.Error())
		return
	}

	comment := &entity.Comment{
		PostID:       postId,
		Creator:      r.Context().Value(ctxKeyUser).(string),
		Body:         body[0],
		CreationDate: time.Now(),
	}

	err = h.useCases.Comment.CreateCommentUseCase(comment)
	if err != nil {
		if errors.Is(err, utils.ErrInvalidPostId) || errors.Is(err, utils.ErrInvalidCreator) ||
			errors.Is(err, utils.ErrInvalidCommentBody) {
			h.ExecuteErrorPage(w, http.StatusBadRequest, err.Error())
			return
		}
		h.ExecuteErrorPage(w, http.StatusInternalServerError, err.Error())
		return
	}

	http.Redirect(w, r, "/posts?post_id="+strconv.Itoa(postId), http.StatusSeeOther)
}
