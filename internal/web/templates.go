package web

import (
	"html/template"
	"path/filepath"

	"github.com/deeincom/deeincom/pkg/models"
	"github.com/deeincom/deeincom/pkg/models/db"
	"github.com/microcosm-cc/bluemonday"
)

var strictPolicy = bluemonday.StrictPolicy()

type templateData struct {
	User       *models.User
	Pagination *db.Pagination
	HomePage   bool
	Localhost  bool
	CurrentURL string
}

var functions = template.FuncMap{}

func parseHTML(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	pages, err := filepath.Glob(filepath.Join(dir, "pages", "*.page.html"))
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		name := filepath.Base(page)

		ts := template.New(name).Funcs(functions)

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.html"))
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "partials", "*.partial.html"))
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}
	return cache, nil
}
