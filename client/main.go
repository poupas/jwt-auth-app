package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const (
	BearerPrefix = "Bearer "
)

func loadSecretKey(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func main() {
	// Command-line flags.
	secretKeyPath := flag.String("secret", "secret.key", "Path to the JWT secret key file")
	serverURL := flag.String("url", "http://localhost:8080/", "URL of the server")
	flag.Parse()

	// Load the shared secret key from a file.
	secretKey, err := loadSecretKey(*secretKeyPath)
	if err != nil {
		log.Fatalf("Failed to load secret key: %v", err)
	}

	// Create a new JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iat": time.Now().Unix(),
		// Additional claims can be added here.
	})

	// Authentication the token with the secret key.
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		log.Fatalf("Failed to sign token: %v", err)
	}

	// Create a new HTTP request.
	req, err := http.NewRequest("GET", *serverURL, nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	// Add the Authorization header.
	req.Header.Add("Authorization", BearerPrefix+tokenString)

	// Show the JWT.
	fmt.Printf("Authorization header: %s\n", req.Header["Authorization"][0])

	// Send the request.
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	// Read and print the response.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response: %v", err)
	}

	fmt.Printf("Response Status: %s\n", resp.Status)
	fmt.Printf("Response Body: %s\n", string(body))
}
