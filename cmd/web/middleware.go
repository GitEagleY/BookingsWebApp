package main

import (
	"net/http"

	"github.com/justinas/nosurf"
)

// CSRF attacs protection
func NoSurf(next http.Handler) http.Handler {

	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.Production,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}

// loads and saves session every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}
