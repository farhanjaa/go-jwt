package middleware

import (
	"go-jwt/helpers"
	"net/http"
)

// AdminOnly membatasi akses hanya untuk user dengan role "admin"
func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(userInfoKey).(helpers.UserInfo)
		if !ok || user.Role != "admin" {
			http.Error(w, "Forbidden - Admins only", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
