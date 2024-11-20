package main

import (
	"bank-demo-app/internal/bank/dbBank"
	"bank-demo-app/internal/bank/memoryBank"
	"bank-demo-app/internal/inputParams"
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

const (
	serverPort      = ":8080"
	shutdownTimeout = 5 * time.Second
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // This will allow us to cancel the context on app

	// Inits logger global instance to use it all around the project with same timestamp.
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()

	config, err := inputParams.ParseInputParams()
	if err != nil {
		log.Error().Err(err).Msg("Finishing application")
		return
	}

	// Initialize BankStore type.
	bankStore := initBankStore(ctx, config)

	if err := startServer(ctx, bankStore); err != nil {
		log.Fatal().Err(err).Msg("Application terminated with error")
	}
}

/// BootStrap helper functions would be moved to different package, but not needed for this technical test.

// initBankStore initializes the appropriate BankStore based on configuration.
func initBankStore(ctx context.Context, config *inputParams.AppConfig) restServer.BankStore {
	if config.InMemory {
		return memoryBank.NewBankStore()
	}
	return dbBank.NewBankStore(ctx, &config.MongoConf)
}

// startServer starts the HTTP server and handles graceful shutdown.
func startServer(ctx context.Context, bankStore restServer.BankStore) error {
	router := restServer.NewRouter(restServer.InitRestRoutes(bankStore))
	server := &http.Server{
		Addr:    serverPort,
		Handler: router,
	}

	// Run server in a goroutine
	go func() {
		log.Info().Msgf("Server listening on port %s...", serverPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Server failed")
		}
	}()

	// Wait for termination signal
	signalCtx := handleSignals()
	<-signalCtx.Done()

	// Shutdown server gracefully
	log.Info().Msg("Shutting down server...")
	ctxWithTimeout, cancel := context.WithTimeout(ctx, shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctxWithTimeout); err != nil {
		log.Error().Err(err).Msg("Error during server shutdown")
		return err
	}

	log.Info().Msg("Server stopped gracefully")
	return nil
}

// handleSignals listens for termination signals and returns a cancellable context.
func handleSignals() context.Context {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-stop
		cancel()
	}()
	return ctx
}
