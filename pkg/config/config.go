package config

import (
	"html/template"

	"github.com/alexedwards/scs/v2"
)

// holds app config
type AppConfig struct {
	TemplateCache map[string]*template.Template

	Production bool //if yes ur in producton mode
	Session    *scs.SessionManager
}
