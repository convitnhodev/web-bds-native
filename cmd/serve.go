package cmd

import (
	configuration "github.com/deeincom/deeincom/config"
	"github.com/deeincom/deeincom/database"
	"github.com/deeincom/deeincom/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/spf13/cobra"
	"log"
)

type App struct {
	*fiber.App
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

func serve(cmd *cobra.Command, args []string) {
	config := configuration.New()
	database.Connect()
	app := App{
		App: fiber.New(*config.GetFiberConfig()),
	}
	app.Use(recover.New())
	app.Use(logger.New())

	v1 := app.Group("/api/v1")

	v1.Get("/users", handlers.UserList)
	v1.Post("/users", handlers.UserCreate)

	app.Static("/", "./static/public")

	app.Use(handlers.NotFound)
	log.Fatal(app.Listen(config.GetString("APP_ADDR")))
}
