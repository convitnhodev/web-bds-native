package database

import (
	"context"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/postgresql"
	"time"
)

type DBConnector interface {
	Connect() error
	Close() error
	Migrate() error
	GetSession() db.Session
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
	db.LC().SetLevel(db.LogLevelTrace)

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

func (d *DB) Migrate() error {
	driver, err := postgres.WithInstance(d.db.DB, &postgres.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://database/migrations",
		"postgres", driver)
	if err != nil {
		return err
	}
	return m.Up()
}

func (d *DB) GetSession() db.Session {
	return d.dbx
}
