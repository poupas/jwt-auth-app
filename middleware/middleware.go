package middleware

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const (
	TokenWindowSeconds = 60
	BearerPrefix       = "Bearer "
)

func LoadSecretKey(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// JWTAuthMiddleware validates JWT tokens and protects routes.
func JWTAuthMiddleware(secretKey []byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the Authorization header.
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		// Check if the Authorization header has the Bearer prefix.
		if !strings.HasPrefix(authHeader, BearerPrefix) {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		// Extract the token string.
		tokenString := strings.TrimPrefix(authHeader, BearerPrefix)

		// Parse the token.
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Ensure the signing method is HMAC.
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return secretKey, nil
		})

		if err != nil {
			log.Printf("Token parsing error: %v", err)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Validate the token and extract claims.
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		// Extract 'iat' claim.
		iatFloat, ok := claims["iat"].(float64)
		if !ok {
			http.Error(w, "Invalid 'iat' claim", http.StatusUnauthorized)
			return
		}
		iat := int64(iatFloat)
		currentTime := time.Now().Unix()

		if currentTime < iat-TokenWindowSeconds || currentTime > iat+TokenWindowSeconds {
			http.Error(w, "Token expired", http.StatusUnauthorized)
			return
		}

		// Token is valid; proceed to the next handler.
		next.ServeHTTP(w, r)
	})
}
