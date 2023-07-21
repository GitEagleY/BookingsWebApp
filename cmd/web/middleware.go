package main

import (
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
)

// writes to console every time somoene hits the page
func WriteToConsole(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("someone hited the page")
		next.ServeHTTP(w, r)
	})
}

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
