package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ubik.org/jwt-auth-app/middleware"

	"github.com/gorilla/mux"
)

func main() {
	// Command-line flags
	secretKeyPath := flag.String("secret", "secret.key", "Path to the JWT secret key file")
	addr := flag.String("addr", ":8080", "HTTP network address")
	flag.Parse()

	secretKey, err := middleware.LoadSecretKey(*secretKeyPath)
	if err != nil {
		log.Fatalf("Failed to load secret key: %v", err)
	}

	r := mux.NewRouter()

	// Apply the JWT middleware to all routes.
	r.Use(func(next http.Handler) http.Handler {
		return middleware.JWTAuthMiddleware(secretKey, next)
	})

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Access granted")
	}).Methods("GET")

	// Additional routes can be defined here and will be protected by JWT middleware.

	srv := &http.Server{
		Addr:    *addr,
		Handler: r,
	}

	serverErrors := make(chan error, 1)

	// Start the server.
	go func() {
		log.Printf("Server starting on %s", *addr)
		serverErrors <- srv.ListenAndServe()
	}()

	// Trap signals for graceful shutdown.
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt, syscall.SIGTERM)

	// Blocking select to wait for either an error or a shutdown signal.
	select {
	case err := <-serverErrors:
		log.Fatalf("Could not start server: %v", err)

	case sig := <-sigchan:
		log.Printf("Received signal %v. Initiating graceful shutdown...", sig)

		// Create a context with timeout for the shutdown process.
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Attempt graceful shutdown.
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatalf("Graceful shutdown failed: %v", err)
		}

		log.Println("Server gracefully stopped.")
	}
}
