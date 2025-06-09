package helpers

import (
	"errors"
	"go-jwt/models"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var mySigningKey = []byte("mysecretkey")

type MyCustomClaims struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"` // Tambahkan role di sini
	jwt.RegisteredClaims
}

type UserInfo struct {
	ID    int
	Name  string
	Email string
	Role  string
}

func CreateToken(user *models.User) (string, error) {
	// Ambil TTL dari ENV
	ttlStr := os.Getenv("TOKEN_TTL")
	ttl, err := strconv.Atoi(ttlStr)
	if err != nil {
		ttl = 1800 // default 30 menit jika gagal ambil env
	}

	claims := jwt.MapClaims{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
		"role":  user.Role,
		"exp":   time.Now().Add(time.Duration(ttl) * time.Second).Unix(),
		"iat":   time.Now().Unix(),
		"nbf":   time.Now().Unix(),
	}

	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ValidateToken(tokenStr string) (UserInfo, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		return UserInfo{}, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*MyCustomClaims)
	if !ok {
		return UserInfo{}, errors.New("invalid claims")
	}

	return UserInfo{
		ID:    claims.ID,
		Name:  claims.Name,
		Email: claims.Email,
		Role:  claims.Role,
	}, nil
}
