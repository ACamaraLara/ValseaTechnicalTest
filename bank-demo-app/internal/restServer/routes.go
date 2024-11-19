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

func NewRouter() *gin.Engine {
	// Avoid GIN verbose messages.
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard
	// Create the muxer of our Rest server that will
	// route each request and response to correspondent
	// declared route.
	router := gin.Default()

	serverRoutes := initRestRoutes()

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

// This is the default status handler that will be used to check if the REST server is up.
func statusHandler(c *gin.Context) {
	log.Info().Msg("Called GET status method.")
	if c.Request.Method != http.MethodGet {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Expected GET method!",
		})
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
}

func initRestRoutes() Routes {
	serverRoutes := Routes{
		// Route to get if server is up.
		{
			Method:  http.MethodGet,
			Pattern: "/status",
			Handler: statusHandler,
		},
		// Create a new bank account with an initial balance.
		{
			Method:  http.MethodPost,
			Pattern: "/accounts",
			Handler: func(c *gin.Context) {
			},
		},
		// Retrieve a list of all bank accounts.
		{
			Method:  http.MethodGet,
			Pattern: "/accounts",
			Handler: func(c *gin.Context) {
			},
		},
		// Retrieve details of a specific account by ID.
		{
			Method:  http.MethodGet,
			Pattern: "/accounts/:id",
			Handler: func(c *gin.Context) {
			},
		},
		// Create a deposit or withdrawal transaction for a specific account.
		{
			Method:  http.MethodPost,
			Pattern: "/accounts/:id/transactions",
			Handler: func(c *gin.Context) {
			},
		},
		// Retrieve all transactions associated with a specific account.
		{
			Method:  http.MethodGet,
			Pattern: "/accounts/:id/transactions",
			Handler: func(c *gin.Context) {
			},
		},
		// Transfer funds from one account to another.
		{
			Method:  http.MethodPost,
			Pattern: "/transfer",
			Handler: func(c *gin.Context) {
			},
		},
	}
	return serverRoutes
}
