package v1

import (
	"context"
	"errors"
	"net/http"
	"time"

	"01.alem.school/git/Taimas/forum/pkg/utils"
)

const (
	ctxKeyUser = "session_username"
)

func (h *Handler) authenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, err := h.checkUserAuth(w, r)
		if err != nil {
			if errors.Is(err, utils.ErrUnauthorized) {
				h.ExecuteErrorPage(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
				return
			} else {
				h.ExecuteErrorPage(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
				return
			}
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, username)))
	})
}

func (h *Handler) checkUserAuth(w http.ResponseWriter, r *http.Request) (string, error) {
	c, err := r.Cookie("session_token")
	if err != nil {
		return "", utils.ErrUnauthorized
	}

	session, err := h.useCases.Authorization.ParseToken(c.Value)
	if err != nil {
		if errors.Is(err, utils.ErrSqlNotFound) {
			return "", utils.ErrUnauthorized
		} else {
			return "", err
		}
	}
	if session.TokenExpDate.Before(time.Now()) {
		if err := h.useCases.Authorization.LogOut(session.Token); err != nil {
			return "", err
		}
		return "", utils.ErrUnauthorized
	}

	return session.Username, nil
}
