package config

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"html/template"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type TemplateRenderer struct {
	templates map[string]*template.Template
}

var functions = template.FuncMap{
	"currentPathMatch": func(c echo.Context, p string) bool {
		return c.Request().URL.String() == p
	},
	"currentPathPrefixMatch": func(c echo.Context, p string) bool {
		return strings.HasPrefix(c.Request().URL.String(), p)
	},
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	tmpl, ok := t.templates[name]
	if !ok {
		return fmt.Errorf("template: %s not found", name)
	}

	return tmpl.ExecuteTemplate(w, name, data)
}

func (c *Config) GetEmbedRender() (*TemplateRenderer, error) {
	cache := map[string]*template.Template{}
	viewPaths := []string{
		filepath.Join("resources", "views", "web"),
		filepath.Join("resources", "views", "admin"),
	}
	for _, path := range viewPaths {
		templates, err := c.getEmbedTemplate(path)
		if err != nil {
			return nil, err
		}
		for n, t := range templates {
			cache[n] = t
		}
	}
	return &TemplateRenderer{templates: cache}, nil
}

func (c *Config) getEmbedTemplate(path string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	fd := os.DirFS(path)
	pages, _ := fs.Glob(fd, filepath.Join("pages", "*.page.html"))
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFS(fd, "pages/*.page.html", "*.layout.html", "partials/*.partial.html")
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}
	return cache, nil
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

	for s, t := range cache {
		fmt.Println("name:", s, "template:", t.Name())
	}

	return cache, nil
}
