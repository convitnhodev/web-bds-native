package cmd

import (
	"errors"
	"github.com/deeincom/deeincom/app/repositories"
	configuration "github.com/deeincom/deeincom/config"
	"github.com/deeincom/deeincom/database"
	"github.com/deeincom/deeincom/routes"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/cobra"
	"html/template"
	"io"
	"log"
	"net/http"
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
	renderer, err := config.GetEmbedRender()
	if err != nil {
		log.Fatal("unable to get template", err)
	}
	e := echo.New()
	e.Renderer = renderer
	app := App{
		Echo: e,
	}
	app.Use(middleware.Logger())
	app.Use(middleware.Recover())
	app.Static("/", "./public")
	routes.RegisterWeb(app.Echo, repositories.New(db.GetSession()))
	routes.RegisterAdmin(app.Echo)
	app.Echo.GET("/test", func(c echo.Context) error {
		return c.JSON(http.StatusOK, app.Routes())
	})

	app.Logger.Fatal(app.Start(config.GetString("APP_ADDR")))
}
