package handlers

import (
	"fmt"
	"net/http"

	render "github.com/GitEagleY/BookingsWebApp/pkg/Render"
	"github.com/GitEagleY/BookingsWebApp/pkg/config"
	"github.com/GitEagleY/BookingsWebApp/pkg/models"
)

type Repository struct { //defined Repository structure
	App *config.AppConfig //containing data from App
}

var Repo *Repository //making variable

func NewRepo(a *config.AppConfig) *Repository { //strange constructor of Reposytory
	return &Repository{
		App: a,
	}
}
func NewHandlers(r *Repository) { //need for taking here repo to work with
	Repo = r
}
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	TestingStringMap := make(map[string]string) //creating testing map of string
	TestingStringMap["test"] = "hello"          //adding new string to just created map

	remoteIpData := r.RemoteAddr                                                           //taking user ip from session
	m.App.Session.Put(r.Context(), "KEY_TO_TAKE_DATA_FROM_SESSION", remoteIpData)          //putting to session data
	takedRemoteIp := m.App.Session.GetString(r.Context(), "KEY_TO_TAKE_DATA_FROM_SESSION") //taking just putted data from session by key
	TestingStringMap["remote_IP"] = takedRemoteIp                                          //adding just taked remote ip from Session data to TestingStringMap
	//set data to render
	fmt.Println(TestingStringMap) //just make sure everyting ok with map
	render.RenderTemplate(w, "home.page.tmpl", &models.TemplateData{
		StringMap: TestingStringMap, //rendering test sting map
	})

}
