package main

import (
	"net/http"

	handlers "github.com/GitEagleY/BookingsWebApp/pkg/Handlers"
	"github.com/GitEagleY/BookingsWebApp/pkg/config"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func routes(app *config.AppConfig) http.Handler {
	mx := chi.NewRouter()
	mx.Use(middleware.Recoverer) //need so app doesnt fall because of some random error
	mx.Use(NoSurf)               //CSRF attacs protection (just why not)
	mx.Use(SessionLoad)          //save and load user session data

	mx.Get("/", handlers.Repo.Home) //make site & middleware work
	mx.Get("/about", handlers.Repo.About)
	mx.Get("/book_now", handlers.Repo.Book_now)
	mx.Get("/contact", handlers.Repo.Contact)
	mx.Get("/login", handlers.Repo.Login)
	mx.Get("/make_reservation", handlers.Repo.Make_reservation)
	mx.Get("/quarters", handlers.Repo.Quarters)
	mx.Get("/register", handlers.Repo.Register)
	mx.Get("/suite", handlers.Repo.Suite)

	fileserver := http.FileServer(http.Dir("./static/"))
	mx.Handle("/static/*", http.StripPrefix("/static", fileserver))
	return mx //return setted up mixer to main
}
