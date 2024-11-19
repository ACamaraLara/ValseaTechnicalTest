package main

import (
	"bank-demo-app/internal/bank/memoryBank"
	"bank-demo-app/internal/restServer"
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Inits logger global instance to use it all around the project with same timestamp.
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()

	bankStore := memoryBank.NewBankStore()
	bankRoutes := restServer.InitRestRoutes(bankStore)

	// Create a new HTTP router and set up routes
	router := restServer.NewRouter(bankRoutes)

	server := &http.Server{
		Addr:    ":8080", // Typical test port for REST servers.
		Handler: router,
	}
	go func() {
		log.Info().Msg("Starting server on port 8080...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Server failed")
		}
	}()

	// Correct shutdown handling
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop // Block until a termination signal is received
	log.Info().Msg("Shutting down server...")

	// Maximum of 5 secconds to let the server finish correctly.
	timeout := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("Error during server shutdown")
	} else {
		log.Info().Msg("Server correctly stopped")
	}
}
