package cmd

import (
	"errors"
	"html/template"
	"io"
	"log"
	"net/http"
	"path/filepath"

	"github.com/deeincom/deeincom/app/repositories"
	configuration "github.com/deeincom/deeincom/config"
	"github.com/deeincom/deeincom/database"
	"github.com/deeincom/deeincom/routes"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/cobra"
)

type App struct {
	*echo.Echo
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "serve http server",
	Long:  "serve http server",
	Run:   serve,
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

type TemplateRenderer struct {
	templates map[string]*template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	tmpl, ok := t.templates[name]
	if !ok {
		err := errors.New("Template not found -> " + name)
		return err
	}
	return tmpl.ExecuteTemplate(w, name, data)
}

func serve(cmd *cobra.Command, args []string) {
	config := configuration.New()
	db := database.NewDB(config.GetString("DB_URI"))
	if err := db.Connect(); err != nil {
		log.Fatal("unable connect to db", err)
	}
	app := App{
		Echo: echo.New(),
	}
	app.Use(middleware.Logger())
	app.Use(middleware.Recover())
	app.Static("/", "./public")
	renderer := &TemplateRenderer{
		templates: getRenderer(),
	}

	app.Renderer = renderer
	routes.RegisterWeb(app.Echo, repositories.New(db.GetSession()))
	routes.RegisterAdmin(app.Echo)
	app.Echo.GET("/test", func(c echo.Context) error {
		return c.JSON(http.StatusOK, app.Routes())
	})

	app.Logger.Fatal(app.Start(config.GetString("APP_ADDR")))
}

func getRenderer() map[string]*template.Template {
	t, err := parseHTML(filepath.Join("resources", "views"))
	if err != nil {
		panic(err)
	}
	ats, err := parseHTML(filepath.Join("resources", "views", "admin"))
	if err != nil {
		panic(err)
	}
	for n, at := range ats {
		t[n] = at
	}
	ats, err = parseHTML(filepath.Join("resources", "views", "web"))
	if err != nil {
		panic(err)
	}
	for n, at := range ats {
		t[n] = at
	}
	return t
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
