package cmd

import (
	configuration "github.com/deeincom/deeincom/config"
	"github.com/deeincom/deeincom/pkg/database"
	"github.com/deeincom/deeincom/routes"
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
	db := database.NewDB(config.GetString("DB_URI"))
	if err := db.Connect(); err != nil {
		log.Fatal("unable connect to db", err)
	}
	app := App{
		App: fiber.New(*config.GetFiberConfig()),
	}
	app.Use(recover.New())
	app.Use(logger.New())
	routes.RegisterWeb(app)
	log.Fatal(app.Listen(config.GetString("APP_ADDR")))
}
