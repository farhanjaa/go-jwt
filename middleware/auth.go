package middleware

import (
	"context"
	"fmt"
	"go-jwt/helpers"
	"html/template"
	"net/http"
)

type contextKey string

const userInfoKey contextKey = "userinfo"

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accesToken := r.Header.Get("Authorization")

		if accesToken == "" {
			helpers.Response(w, 401, "unauthorized", nil)
			return
		}

		user, err := helpers.ValidateToken(accesToken)
		if err != nil {
			helpers.Response(w, 401, err.Error(), nil)
			return
		}

		ctx := context.WithValue(r.Context(), userInfoKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func SomeProtectedHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(userInfoKey).(helpers.UserInfo)

	if user.Role != "admin" {
		http.Error(w, "Forbidden: Admins only", http.StatusForbidden)
		return
	}

	// Lanjutkan logic khusus admin
	fmt.Fprintf(w, "Welcome Admin %s!", user.Name)
	tmpl := template.Must(template.ParseFiles("views/home/landing.html"))
	tmpl.Execute(w, user)
}
