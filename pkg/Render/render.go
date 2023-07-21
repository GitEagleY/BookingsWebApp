package render

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	config "github.com/GitEagleY/BookingsWebApp/pkg/config"
	"github.com/GitEagleY/BookingsWebApp/pkg/models"
)

//var functions = template.FuncMap{}

var app *config.AppConfig
var templateCache map[string]*template.Template

func AddDefaultData(td *models.TemplateData) *models.TemplateData { //default data avalivable to every template
	return td
}
func NewTemplates(a *config.AppConfig) { //template constructor
	app = a
}
func RenderTemplate(w http.ResponseWriter, templateName string, td *models.TemplateData) error {
	if app.Production { //if not in developing mode use templates each time from cache
		templateCache = app.TemplateCache
	} else { //or if in development mode use from disk
		templateCache, _ = CacheTemplate()
	}

	templateToRender, ok := templateCache[templateName] //make template to render
	if !ok {
		log.Fatal(ok)
	}
	td = AddDefaultData(td) //adding deafult data availivable to every template

	templateToRender.Execute(w, td) //render template
	return nil
}
func CacheTemplate() (map[string]*template.Template, error) {
	Cache := map[string]*template.Template{}

	// get all of the files named *.page.tmpl from ./templates
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return Cache, err
	}

	// range through all files ending with *.page.tmpl
	for _, page := range pages {
		name := filepath.Base(page)                            //take page file name
		parsdTmplt, err := template.New(name).ParseFiles(page) //allocating template and create template by parsing page file
		if err != nil {
			return Cache, err
		}
		// get all of the files named *.layout.tmpl from ./templates
		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return Cache, err
		}

		if len(matches) > 0 {
			parsdTmplt, err = parsdTmplt.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return Cache, err
			}
		}

		Cache[name] = parsdTmplt
	}

	return Cache, nil
}
