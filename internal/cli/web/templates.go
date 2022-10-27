package web

import (
	"html/template"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/deeincom/deeincom/pkg/form"
	"github.com/deeincom/deeincom/pkg/models"
	"github.com/deeincom/deeincom/pkg/models/db"
	"github.com/microcosm-cc/bluemonday"
)

var strictPolicy = bluemonday.StrictPolicy()

type templateData struct {
	User       *models.User
	Pagination *db.Pagination

	Localhost  bool
	CurrentURL string

	Flash string

	Form     *form.Form
	Products []*models.Product
	Product  *models.Product
}

var functions = template.FuncMap{
	"__":       translate,
	"upper":    strings.ToUpper,
	"lower":    strings.ToLower,
	"title":    strings.Title,
	"split":    strings.Split,
	"contains": strings.Contains,
	"has_role": hasRole,
}

func hasRole(user *models.User, test string) bool {
	for _, role := range user.Roles {
		if role == test {
			return true
		}
	}
	return false
}

func translate(s string) string {
	return s
}

func use(h http.HandlerFunc, middleware ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for _, m := range middleware {
		h = m(h)
	}
	return h
}

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
