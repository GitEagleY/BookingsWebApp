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
	mx.Use(middleware.Recoverer)    //need so app doesnt fall because of some random error
	mx.Use(WriteToConsole)          //just to check how many times user hitted page
	mx.Use(NoSurf)                  //CSRF attacs protection (just why not)
	mx.Use(SessionLoad)             //save and load user session data
	mx.Get("/", handlers.Repo.Home) //make site & middleware work

	return mx //return setted up mixer to main
}
