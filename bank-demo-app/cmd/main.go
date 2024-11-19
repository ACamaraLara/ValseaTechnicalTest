package main

import (
	"bank-demo-app/internal/restServer"

	"net/http"

	"github.com/rs/zerolog/log"
)

func main() {
	router := restServer.NewRouter()
	// Starts listening for HTTP requests.
	log.Fatal().Msg(http.ListenAndServe(":8080", router).Error())
	log.Info().Msg("Exiting...")
}
