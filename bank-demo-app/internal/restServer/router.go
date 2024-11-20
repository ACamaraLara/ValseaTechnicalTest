package restServer

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// Struct that stores information of a single
// route for current service.
type Route struct {
	Method  string
	Pattern string
	Handler gin.HandlerFunc
}

// Vector to store declared routes.
type Routes []Route

func NewRouter(serverRoutes Routes) *gin.Engine {
	// Avoid GIN verbose messages.
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard
	// Create the muxer of our Rest server that will
	// route each request and response to correspondent
	// declared route.
	router := gin.Default()

	for _, route := range serverRoutes {
		addRoute(router, route)
	}

	// Return muxer with all its added routes.
	return router
}

func addRoute(router *gin.Engine, route Route) {
	switch route.Method {
	case http.MethodGet:
		router.GET(route.Pattern, route.Handler)
	case http.MethodPost:
		router.POST(route.Pattern, route.Handler)
	default:
		log.Warn().Msg("Invalid HTTP method specified: " + route.Method)
	}
}
