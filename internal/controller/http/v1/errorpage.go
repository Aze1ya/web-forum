package v1

import (
	"net/http"
)

type errorData struct {
	Status int
	Text   string
}

func (h *Handler) ExecuteErrorPage(w http.ResponseWriter, status int, text string) {
	w.WriteHeader(status)
	data := errorData{
		Status: status,
		Text:   text,
	}

	if err := h.tmpl.ExecuteTemplate(w, "error.html", data); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
