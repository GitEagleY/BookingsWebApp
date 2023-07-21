package main

import (
	"fmt"
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
	//session.Cookie.Name="lol"
	//session.Cookie.Path="/"
	fmt.Println(session.Cookie) //just to make sure everything ok with cookies
	app.Session = session       //giving session to app data

	///////////////////////////CACHING////////////////////////////
	tc, err := render.CacheTemplate() //creating template cache
	if err != nil {
		log.Fatal("cannot create template cache")
	}
	app.TemplateCache = tc //giving just created cache to app config

	//////////////////RENDER//////////////////////////////////////
	render.NewTemplates(&app) //giving app config to renders

	//////////////////WORK WITH HANDLERS AND REPO/////////////////
	repo := handlers.NewRepo(&app) //making new repo
	handlers.NewHandlers(repo)     //giving new just created repo to Handlers

	//http.HandleFunc("/", handlers.Repo.Home)
	//////////////////////////SERVER////////////////////////////////////
	srv := &http.Server{ //server
		Addr:    portnum,
		Handler: routes(&app),
	}
	err = srv.ListenAndServe() //check if there no errros while running the server
	if err != nil {
		log.Fatal(err)
	}
	//http.ListenAndServe(portnum, nil) //running server
}
