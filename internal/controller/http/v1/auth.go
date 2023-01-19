package v1

import (
	"errors"
	"net/http"

	"01.alem.school/git/Taimas/forum/pkg/utils"
)

func (h *Handler) signUp(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/auth/signup" {
		h.ExecuteErrorPage(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}

	if r.Method == http.MethodGet {
		if err := h.tmpl.ExecuteTemplate(w, "signup.html", nil); err != nil {
			h.ExecuteErrorPage(w, http.StatusInternalServerError, err.Error())
			return
		}
		return
	}
	if r.Method != http.MethodPost {
		h.ExecuteErrorPage(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	if err := r.ParseForm(); err != nil {
		h.ExecuteErrorPage(w, http.StatusInternalServerError, err.Error())
		return
	}
	email, ok := r.Form["email"]
	if !ok {
		h.ExecuteErrorPage(w, http.StatusBadRequest, utils.ErrEmailFieldNotFilled.Error())
		return
	}
	username, ok := r.Form["username"]
	if !ok {
		h.ExecuteErrorPage(w, http.StatusBadRequest, utils.ErrUsernameFieldNotFilled.Error())
		return
	}
	password, ok := r.Form["password"]
	if !ok {
		h.ExecuteErrorPage(w, http.StatusBadRequest, utils.ErrPasswordFieldNotFilled.Error())
		return
	}

	passwordRepeat, ok := r.Form["passwordRepeat"]
	if !ok {
		h.ExecuteErrorPage(w, http.StatusBadRequest, utils.ErrPasswordFieldNotFilled.Error())
		return
	}

	err := h.useCases.Authorization.SignUp(email[0], username[0], password[0], passwordRepeat[0])
	if err != nil {
		if errors.Is(err, utils.ErrUsernameNotUnique) || errors.Is(err, utils.ErrEmailNotUnique) || errors.Is(err, utils.ErrPasswordFieldNotSame) || errors.Is(err, utils.ErrInvalidUsernameFormat) || errors.Is(err, utils.ErrInvalidEmailFormat) || errors.Is(err, utils.ErrInvalidPasswordFormat) {
			h.ExecuteErrorPage(w, http.StatusBadRequest, err.Error())
		} else {
			h.ExecuteErrorPage(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	http.Redirect(w, r, "/auth/signin", http.StatusSeeOther)
}

func (h *Handler) signIn(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/auth/signin" {
		h.ExecuteErrorPage(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	if r.Method == http.MethodGet {
		if err := h.tmpl.ExecuteTemplate(w, "signin.html", nil); err != nil {
			h.ExecuteErrorPage(w, http.StatusInternalServerError, err.Error())
			return
		}
		return
	}
	if r.Method != http.MethodPost {
		h.ExecuteErrorPage(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	if err := r.ParseForm(); err != nil {
		h.ExecuteErrorPage(w, http.StatusInternalServerError, err.Error())
		return
	}
	username, ok := r.Form["username"]
	if !ok {
		h.ExecuteErrorPage(w, http.StatusBadRequest, utils.ErrUsernameFieldNotFilled.Error())
		return
	}
	password, ok := r.Form["password"]
	if !ok {
		h.ExecuteErrorPage(w, http.StatusBadRequest, utils.ErrPasswordFieldNotFilled.Error())
		return
	}
	token, expirationDate, err := h.useCases.Authorization.SignIn(username[0], password[0])
	if err != nil {
		if errors.Is(err, utils.ErrUserNotFound) || errors.Is(err, utils.ErrIncorrectPassword) {
			h.ExecuteErrorPage(w, http.StatusBadRequest, err.Error())
			return
		}
		h.ExecuteErrorPage(w, http.StatusInternalServerError, err.Error())
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   token,
		Path:    "/",
		Expires: expirationDate,
	})
	// fmt.Println(r.Method)
	// r.Method = http.MethodGet
	// fmt.Println(r.Method)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) logOut(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/auth/logout" {
		h.ExecuteErrorPage(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	if r.Method != http.MethodPost {
		h.ExecuteErrorPage(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			h.ExecuteErrorPage(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			return
		}
		h.ExecuteErrorPage(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := h.useCases.Authorization.LogOut(c.Value); err != nil {
		h.ExecuteErrorPage(w, http.StatusInternalServerError, err.Error())
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:  "session_token",
		Value: "",
		Path:  "/",
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
