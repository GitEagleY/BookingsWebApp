package helpers

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/GitEagleY/BookingsWebApp/internal/config"
)

var app *config.AppConfig

// NewHelpers sets up the app config for helpers.
func NewHelpers(a *config.AppConfig) {
	app = a // Assign the provided AppConfig to the app variable for access in helper functions.
}

// ClientError logs a client error and responds with the specified status.
func ClientError(w http.ResponseWriter, status int) {
	// Log the client error with the provided status.
	app.InfoLog.Println("Client error with status of:", status)
	// Respond with the appropriate status and message.
	http.Error(w, http.StatusText(status), status)
}

// ServerError logs a server error, including a stack trace, and responds with a 500 status.
func ServerError(w http.ResponseWriter, err error) {
	// Generate a stack trace and log the server error.
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Println(trace)
	// Respond with a 500 status and generic server error message.
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
