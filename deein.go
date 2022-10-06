package deein

import (
	"database/sql"
	"log"

	"github.com/deeincom/deeincom/config"
	"github.com/deeincom/deeincom/pkg/models/db"
	"github.com/deeincom/deeincom/pkg/telegram"
	"github.com/pkg/errors"

	// register pq with database/sql
	_ "github.com/lib/pq"
)

// App represent an app
type App struct {
	Config    *config.Config
	Telegram  *telegram.Telegram
	Migration *db.MigrationModel
	Users     *db.UserModel
}

// New return an app instance
func New(c *config.Config) (*App, error) {
	conn, err := sql.Open("postgres", c.DB)
	if err != nil {
		return nil, errors.Wrap(err, "postgres")
	} else {
		log.Println("DATABASE: OK")
	}

	app := &App{
		Config: c,
		Migration: &db.MigrationModel{
			DB: conn,
		},
		Users: &db.UserModel{
			DB: conn,
			Pagination: &db.Pagination{
				DB:      conn,
				Min:     25,
				Max:     100,
				Default: 25,
				Data:    &db.PaginationData{Limit: 25},
			},
		},
		Telegram: &telegram.Telegram{
			Token:  c.TelegramToken,
			ChatID: c.TelegramChatID,
		},
	}

	return app, nil
}
