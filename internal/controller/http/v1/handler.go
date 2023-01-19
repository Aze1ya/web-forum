package v1

import (
	"context"
	"html/template"
	"net/http"

	"01.alem.school/git/Taimas/forum/internal/usecase"
)

type Handler struct {
	useCases *usecase.UseCases
	ctx      context.Context
	tmpl     *template.Template
	Mux      *http.ServeMux
}

func NewHandler(useCases *usecase.UseCases) *Handler {
	return &Handler{
		Mux:      http.NewServeMux(),
		useCases: useCases,
		ctx:      context.Background(),
		tmpl:     template.Must(template.ParseGlob("./ui/html/*.html")),
	}
}

func (h *Handler) RegisterHTTPEndpoints() {
	h.Mux.HandleFunc("/", h.homePage)

	h.Mux.HandleFunc("/auth/signin", h.signIn)
	h.Mux.HandleFunc("/auth/signup", h.signUp)
	h.Mux.Handle("/auth/logout", h.authenticateUser(http.HandlerFunc(h.logOut)))

	h.Mux.Handle("/posts/", http.HandlerFunc(h.postCR))
	h.Mux.Handle("/comments/", h.authenticateUser(http.HandlerFunc(h.commentCreate)))

	h.Mux.Handle("/users/", h.authenticateUser(http.HandlerFunc(h.userProfile)))

	h.Mux.Handle("/likes/", h.authenticateUser(http.HandlerFunc(h.createLike)))
	h.Mux.Handle("/dislikes/", h.authenticateUser(http.HandlerFunc(h.createDisLike)))

	fs := http.FileServer(http.Dir("./ui/static"))
	h.Mux.Handle("/static/", http.StripPrefix("/static", fs))
}
