package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	config "github.com/GitEagleY/BookingsWebApp/internal/config"
	"github.com/GitEagleY/BookingsWebApp/internal/models"
	"github.com/justinas/nosurf"
)

var functions = template.FuncMap{} // Custom template functions.

var app *config.AppConfig // Holds the application configuration.

var pathToTemplates = "./templates" // Path to the template files.

// NewRenderer sets the configuration for the template package.
func NewRenderer(a *config.AppConfig) {
	app = a // Assign the provided AppConfig to the app variable for access in template functions.
}

// AddDefaultData adds common data to the template data.
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	// Retrieve flash messages from the session and add them to the template data.
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.Error = app.Session.PopString(r.Context(), "error")

	// Add CSRF token to the template data.

	if app.Session.Exists(r.Context(), "user_id") {
		td.IsAuthenticated = 1
	}
	td.CSRFToken = nosurf.Token(r)
	return td
}

// Template renders a template
func Template(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {
	var tc map[string]*template.Template

	if app.UseCache {
		// get the template cache from the app config
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("Could not get template from template cache")
		return errors.New("cant get error from cache")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td, r)

	_ = t.Execute(buf, td)

	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("error writing template to browser", err)
		return err
	}
	return nil

}

// CreateTemplateCache creates a template cache as a map
func CreateTemplateCache() (map[string]*template.Template, error) {

	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}

	return myCache, nil
}
