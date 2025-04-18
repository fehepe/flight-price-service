package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/fehepe/flight-price-service/internal/handlers"
	"github.com/fehepe/flight-price-service/internal/middleware"
	"github.com/gorilla/mux"
)

// NewRouter wires up routes and middleware.
func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.Use(middleware.Logging)
	r.StrictSlash(true)

	r.HandleFunc("/health", handlers.HealthCheck).Methods(http.MethodGet)
	r.HandleFunc("/flights/search", handlers.GetFlights).Methods(http.MethodGet)

	return r
}

// Run starts the HTTP server and handles graceful shutdown.
func Run(addr string) error {
	srv := &http.Server{
		Addr:           addr,
		Handler:        NewRouter(),
		ReadTimeout:    time.Duration(getEnvInt("READ_TIMEOUT", 5)) * time.Second,
		WriteTimeout:   time.Duration(getEnvInt("WRITE_TIMEOUT", 10)) * time.Second,
		IdleTimeout:    time.Duration(getEnvInt("IDLE_TIMEOUT", 120)) * time.Second,
		MaxHeaderBytes: 1 << 20,
		ErrorLog:       log.New(os.Stdout, "server: ", log.LstdFlags),
	}

	go func() {
		log.Printf("Server listening on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return srv.Shutdown(ctx)
}

// getEnvInt reads an env var into int or returns fallback.
func getEnvInt(key string, fallback int) int {
	if v, ok := os.LookupEnv(key); ok {
		if iv, err := strconv.Atoi(v); err == nil {
			return iv
		}
	}
	return fallback
}
