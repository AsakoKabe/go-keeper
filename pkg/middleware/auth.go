package middleware

import (
	"net/http"

	"go-keeper/internal/context"
	"go-keeper/internal/utils/jwt"
)

// CookieName Ключ под которым хранится кука
const CookieName = "jwt"

// Auth Middleware для аутентификации пользователя и получения userID
func Auth(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(CookieName)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			userID, err := jwt.GetUserID(cookie.Value)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			ctx := context.SetUserID(r.Context(), userID)

			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}(next)
}
