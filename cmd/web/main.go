package main

import (
	"log"
	"net/http"
	"time"

	render "github.com/GitEagleY/BookingsWebApp/pkg/Render"
	"github.com/GitEagleY/BookingsWebApp/pkg/config"
	"github.com/alexedwards/scs/v2"

	handlers "github.com/GitEagleY/BookingsWebApp/pkg/Handlers"
)

const portnum = ":8080"

var app config.AppConfig
var session *scs.SessionManager

func main() {

	app.Production = false
	// /////////////////////SESSION////////////////////
	session = scs.New()                            //declaring new session
	session.Lifetime = 2 * time.Hour               //liftemi of a session
	session.Cookie.Persist = false                 //will cookie exist after user leves page
	session.Cookie.SameSite = http.SameSiteLaxMode //idk some of the site mods
	session.Cookie.Secure = app.Production         //if in prodution than use https
	app.Session = session                          //giving session to app data

	tc, err := render.CacheTemplate() //creating template cache
	if err != nil {
		log.Fatal("cannot create template cache")
	}
	app.TemplateCache = tc //giving just created cache to app config

	render.NewTemplates(&app) //giving app config to renders

	repo := handlers.NewRepo(&app) //making new repo
	handlers.NewHandlers(repo)     //giving new just created repo to Handlers

	srv := &http.Server{ //server
		Addr:    portnum,
		Handler: routes(&app),
	}
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
