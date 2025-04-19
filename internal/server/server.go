package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fehepe/flight-price-service/internal/config"
	"github.com/fehepe/flight-price-service/internal/handlers"
	"github.com/fehepe/flight-price-service/internal/middleware"
	"github.com/fehepe/flight-price-service/internal/providers"
	"github.com/gorilla/mux"
)

// NewRouter sets up routes, applying logging globally and auth on protected endpoints.
func NewRouter(providerList []providers.Provider) *mux.Router {

	r := mux.NewRouter()
	r.Use(middleware.Logging)
	r.StrictSlash(true)

	r.HandleFunc("/health", handlers.HealthCheck).Methods(http.MethodGet)
	r.HandleFunc("/auth/token", handlers.GenerateToken).Methods(http.MethodPost)

	flights := r.PathPrefix("/flights").Subrouter()
	flights.Use(middleware.Auth)
	fh := handlers.NewFlightHandler(providerList)
	flights.HandleFunc("/search", fh.GetFlights).Methods(http.MethodGet)

	return r
}

func Run(addr string) error {
	return RunWithProvider(addr, MustLoadProviders())
}

func RunWithProvider(addr string, provider providers.Provider) error {
	srv := &http.Server{
		Addr:           addr,
		Handler:        NewRouter([]providers.Provider{provider}),
		ReadTimeout:    time.Duration(config.GetEnvInt("READ_TIMEOUT", 5)) * time.Second,
		WriteTimeout:   time.Duration(config.GetEnvInt("WRITE_TIMEOUT", 10)) * time.Second,
		IdleTimeout:    time.Duration(config.GetEnvInt("IDLE_TIMEOUT", 120)) * time.Second,
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
