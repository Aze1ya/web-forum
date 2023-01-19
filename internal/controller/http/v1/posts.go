package v1

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"01.alem.school/git/Taimas/forum/internal/entity"
	"01.alem.school/git/Taimas/forum/pkg/utils"
)

type DataForPostPage struct {
	Username string
	Post     entity.Post
}

type DataForCategoryPostsPage struct {
	Username string
	Posts    []entity.Post
}

func (h *Handler) postCR(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// get post creation form
		if r.URL.RawQuery == "" {
			if r.URL.Path != "/posts/" {
				h.ExecuteErrorPage(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
				return
			}
			_, err := h.checkUserAuth(w, r)
			if err != nil {
				if errors.Is(err, utils.ErrUnauthorized) {
					h.ExecuteErrorPage(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
					return
				} else {
					h.ExecuteErrorPage(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
					return
				}
			}
			if err := h.tmpl.ExecuteTemplate(w, "createPost.html", nil); err != nil {
				h.ExecuteErrorPage(w, http.StatusInternalServerError, err.Error())
				return
			}
			return
		}
		// get posts by category
		category := r.URL.Query().Get("category")
		if category != "" {
			posts, err := h.useCases.Post.GetPostsByCategoryNameUseCase(category)
			if err != nil {
				h.ExecuteErrorPage(w, http.StatusInternalServerError, err.Error())
				return
			}
			data := new(DataForCategoryPostsPage)
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
			return
		}
		// get post by id
		postId, err := strconv.Atoi(r.URL.Query().Get("post_id"))
		if err != nil || postId == 0 {
			h.ExecuteErrorPage(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
			return
		}

		post, err := h.useCases.Post.GetPostByIdUseCase(postId)
		if err != nil {
			if errors.Is(err, utils.ErrSqlNotFound) {
				h.ExecuteErrorPage(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
				return
			} else {
				h.ExecuteErrorPage(w, http.StatusInternalServerError, err.Error())
				return
			}
		}
		data := new(DataForPostPage)
		post.CreationDateFront = post.CreationDate.Format("02.01.2006 15:04")
		post.Comments = utils.AddDateFrontToComments(post.Comments)
		data.Post = *post
		data.Username, err = h.checkUserAuth(w, r)
		if err != nil && !errors.Is(err, utils.ErrUnauthorized) {
			h.ExecuteErrorPage(w, http.StatusInternalServerError, err.Error())
			return
		}
		if err := h.tmpl.ExecuteTemplate(w, "post.html", data); err != nil {
			h.ExecuteErrorPage(w, http.StatusInternalServerError, err.Error())
			return
		}
		return
	}
	// create post
	if r.Method != http.MethodPost {
		h.ExecuteErrorPage(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	creator, err := h.checkUserAuth(w, r)
	if err != nil {
		if errors.Is(err, utils.ErrUnauthorized) {
			h.ExecuteErrorPage(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			return
		} else {
			h.ExecuteErrorPage(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
	}
	if r.URL.RawQuery != "" {
		h.ExecuteErrorPage(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	post := &entity.Post{
		Creator:      creator,
		Title:        r.FormValue("title"),
		Body:         r.FormValue("body"),
		CreationDate: time.Now(),
		Categories:   r.Form["categories"],
	}

	postId, err := h.useCases.Post.CreatePostUseCase(post)
	if err != nil {
		if errors.Is(err, utils.ErrInvalidPostTitle) || errors.Is(err, utils.ErrInvalidPostBody) ||
			errors.Is(err, utils.ErrInvalidPostCategories) {
			h.ExecuteErrorPage(w, http.StatusBadRequest, err.Error())
			return
		}
		h.ExecuteErrorPage(w, http.StatusInternalServerError, err.Error())
		return
	}
	http.Redirect(w, r, "/posts?post_id="+strconv.Itoa(postId), http.StatusSeeOther)
}
