package deein

import (
	"database/sql"
	"log"

	"github.com/deeincom/deeincom/config"
	"github.com/deeincom/deeincom/pkg/files"
	"github.com/deeincom/deeincom/pkg/models/db"
	"github.com/deeincom/deeincom/pkg/telegram"
	"github.com/pkg/errors"

	// register pq with database/sql
	_ "github.com/lib/pq"
)

// App represent an app
type App struct {
	Config           *config.Config
	Telegram         *telegram.Telegram
	Migration        *db.MigrationModel
	Users            *db.UserModel
	AdminUsers       *db.UserModel
	Products         *db.ProductModel
	AdminProducts    *db.ProductModel
	Attachments      *db.AttachmentModel
	AdminAttachments *db.AttachmentModel
	Posts            *db.PostModel
	Comments         *db.CommentModel
	Files            *db.FileModel
	KYC              *db.KYCModel
	Partner          *db.PartnerModel
	Log              *db.LogModel
	Invoice          *db.InvoiceModel
	InvoiceItem      *db.InvoiceItemModel
	Payment          *db.PaymentModel
	LocalFile        *files.LocalFile
	B2Scheduler      *files.LocalToB2
}

// New return an app instance
func New(c *config.Config) (*App, error) {
	conn, err := sql.Open("postgres", c.DB)
	if err != nil {
		return nil, errors.Wrap(err, "postgres")
	} else {
		log.Println("DATABASE: OK")
	}

	Files := db.FileModel{
		DB: conn,
		Pagination: &db.Pagination{
			DB:      conn,
			Min:     25,
			Max:     100,
			Default: 25,
			Data:    &db.PaginationData{Limit: 25},
		},
	}

	lf := files.LocalFile{
		Root:                   c.UploadingRoot,
		MappingUploadLocalLink: c.MappingUploadLocalLink,
		Files:                  &Files,
	}

	B2Scheduler, err := files.NewB2Scheduler(
		c.B2AccountId,
		c.B2AccountKey,
		c.B2BucketName,
		c.UploadToB2At,
		c.B2Prefix,
		c.UploadingRoot,
		c.MappingUploadLocalLink,
		&Files,
	)

	if err != nil {
		return nil, errors.Wrap(err, "Backblaze scheduler")
	}

	B2Scheduler.StartScheduler()

	app := &App{
		B2Scheduler: B2Scheduler,
		LocalFile:   &lf,
		Config:      c,
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
		Products: &db.ProductModel{
			DB: conn,
			Pagination: &db.Pagination{
				DB:      conn,
				Min:     25,
				Max:     100,
				Default: 25,
				Data:    &db.PaginationData{Limit: 25},
			},
		},
		Posts: &db.PostModel{
			DB: conn,
			Pagination: &db.Pagination{
				DB:      conn,
				Min:     25,
				Max:     100,
				Default: 25,
				Data:    &db.PaginationData{Limit: 25},
			},
		},
		Comments: &db.CommentModel{
			DB: conn,
			Pagination: &db.Pagination{
				DB:      conn,
				Min:     25,
				Max:     100,
				Default: 25,
				Data:    &db.PaginationData{Limit: 25},
			},
		},
		AdminProducts: &db.ProductModel{
			DB: conn,
			Pagination: &db.Pagination{
				DB:      conn,
				Min:     25,
				Max:     100,
				Default: 25,
				Data:    &db.PaginationData{Limit: 25},
			},
		},
		AdminUsers: &db.UserModel{
			DB: conn,
			Pagination: &db.Pagination{
				DB:      conn,
				Min:     25,
				Max:     100,
				Default: 25,
				Data:    &db.PaginationData{Limit: 25},
			},
		},
		Attachments: &db.AttachmentModel{
			DB: conn,
			Pagination: &db.Pagination{
				DB:      conn,
				Min:     25,
				Max:     100,
				Default: 25,
				Data:    &db.PaginationData{Limit: 25},
			},
		},
		AdminAttachments: &db.AttachmentModel{
			DB: conn,
			Pagination: &db.Pagination{
				DB:      conn,
				Min:     25,
				Max:     100,
				Default: 25,
				Data:    &db.PaginationData{Limit: 25},
			},
		},
		Files: &Files,
		KYC: &db.KYCModel{
			DB: conn,
			Pagination: &db.Pagination{
				DB:      conn,
				Min:     25,
				Max:     100,
				Default: 25,
				Data:    &db.PaginationData{Limit: 25},
			},
		},
		Partner: &db.PartnerModel{
			DB: conn,
			Pagination: &db.Pagination{
				DB:      conn,
				Min:     25,
				Max:     100,
				Default: 25,
				Data:    &db.PaginationData{Limit: 25},
			},
		},
		Log: &db.LogModel{
			DB: conn,
			Pagination: &db.Pagination{
				DB:      conn,
				Min:     25,
				Max:     100,
				Default: 25,
				Data:    &db.PaginationData{Limit: 25},
			},
		},
		Invoice: &db.InvoiceModel{
			DB: conn,
			Pagination: &db.Pagination{
				DB:      conn,
				Min:     25,
				Max:     100,
				Default: 25,
				Data:    &db.PaginationData{Limit: 25},
			},
		},
		InvoiceItem: &db.InvoiceItemModel{
			DB: conn,
			Pagination: &db.Pagination{
				DB:      conn,
				Min:     25,
				Max:     100,
				Default: 25,
				Data:    &db.PaginationData{Limit: 25},
			},
		},
		Payment: &db.PaymentModel{
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
