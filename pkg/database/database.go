package database

import (
	"context"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/postgresql"
	"time"
)

type DBConnector interface {
	Connect() error
	Close() error
}

func NewDB(uri string) DBConnector {
	return &DB{
		uri: uri,
	}
}

type DB struct {
	uri string
	db  *sqlx.DB
	dbx db.Session
}

func (d *DB) Connect() error {
	database, err := sqlx.Open("postgres", d.uri)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if _, err := database.QueryContext(ctx, "SELECT 1"); err != nil {
		return err
	}
	d.db = database
	d.dbx, err = postgresql.New(database.DB)
	if err != nil {
		return err
	}
	return nil
}

func (d *DB) Close() error {
	if err := d.db.Close(); err != nil {
		return err
	}
	return d.dbx.Close()
}
