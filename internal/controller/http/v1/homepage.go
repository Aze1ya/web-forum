package v1

import (
	"errors"
	"net/http"

	"01.alem.school/git/Taimas/forum/internal/entity"
	"01.alem.school/git/Taimas/forum/pkg/utils"
)

type DataForHomepage struct {
	Username string
	Posts    []entity.Post
}

func (h *Handler) homePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		h.ExecuteErrorPage(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}

	if r.Method != http.MethodGet {
		h.ExecuteErrorPage(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	posts, err := h.useCases.Post.GetAllPostsUseCase()
	if err != nil {
		h.ExecuteErrorPage(w, http.StatusInternalServerError, err.Error())
		return
	}
	data := new(DataForHomepage)
	posts = utils.AddDateFrontToPosts(posts)
	data.Posts = posts
	data.Username, err = h.checkUserAuth(w, r)

	if err != nil && !errors.Is(err, utils.ErrUnauthorized) {
		h.ExecuteErrorPage(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := h.tmpl.ExecuteTemplate(w, "index.html", data); err != nil {
		h.ExecuteErrorPage(w, http.StatusInternalServerError, err.Error())
		return
	}
}
