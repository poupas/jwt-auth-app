package middleware_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
	"ubik.org/jwt-auth-app/middleware"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
)

func TestJWTAuthMiddleware(t *testing.T) {
	// Load the secret key
	secretKeyPath := "../secret.key"
	secretKey, err := middleware.LoadSecretKey(secretKeyPath)
	if err != nil {
		t.Fatalf("Failed to load secret key: %v", err)
	}

	// Create a test router with the middleware
	r := mux.NewRouter()
	r.Use(func(next http.Handler) http.Handler {
		return middleware.JWTAuthMiddleware(secretKey, next)
	})
	r.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Success"))
	}).Methods("GET")

	// Helper function to generate JWT tokens
	generateToken := func(iat int64) (string, error) {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"iat": iat,
		})
		return token.SignedString(secretKey)
	}

	// Define test cases
	tests := []struct {
		name           string
		tokenGenerator func() (string, error)
		expectedStatus int
	}{
		{
			name: "Valid Token",
			tokenGenerator: func() (string, error) {
				return generateToken(time.Now().Unix())
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Missing Authorization Header",
			tokenGenerator: func() (string, error) {
				return "", nil
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Invalid Token Format",
			tokenGenerator: func() (string, error) {
				return "InvalidToken", nil
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Expired Token",
			tokenGenerator: func() (string, error) {
				return generateToken(time.Now().Unix() - 2*middleware.TokenWindowSeconds)
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Future Token",
			tokenGenerator: func() (string, error) {
				return generateToken(time.Now().Unix() + 2*middleware.TokenWindowSeconds)
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Generate the token
			token, err := tt.tokenGenerator()
			if err != nil && tt.expectedStatus == http.StatusOK {
				t.Fatalf("Failed to generate token: %v", err)
			}

			// Create a new HTTP request
			req, err := http.NewRequest("GET", "/test", nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			// Add the Authorization header if token is provided
			if token != "" {
				req.Header.Add("Authorization", middleware.BearerPrefix+token)
			}

			// Record the response
			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			// Check the status code
			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}

// Helper function to create a temporary secret key for testing
func TestMain(m *testing.M) {
	// Create a temporary secret key file
	secretKey := []byte("testsecretkey")
	err := os.WriteFile("../secret.key", secretKey, 0600)
	if err != nil {
		log.Fatalf("Failed to create secret key file: %v", err)
	}

	// Run tests
	code := m.Run()

	// Clean up
	os.Remove("../secret.key")

	os.Exit(code)
}
