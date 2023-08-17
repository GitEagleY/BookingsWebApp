package main

import (
	"net/http"

	"github.com/GitEagleY/BookingsWebApp/internal/helpers"
	"github.com/justinas/nosurf"
)

// NoSurf is the CSRF protection middleware.
func NoSurf(next http.Handler) http.Handler {
	// Create a new CSRF handler using the nosurf package.
	csrfHandler := nosurf.New(next)

	// Configure the base CSRF cookie settings.
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,                 // Ensure the cookie is accessible only via HTTP (not JavaScript).
		Path:     "/",                  // Set the cookie path to the root.
		Secure:   app.InProduction,     // Enable "Secure" flag if in production mode.
		SameSite: http.SameSiteLaxMode, // Set the SameSite attribute to Lax mode for better cross-site request protection.
	})

	return csrfHandler
}

// SessionLoad loads and saves session data for the current request.
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

func Auth(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if !helpers.IsAuthenticated(r) {
			session.Put(r.Context(), "error", "Log in First")
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}
