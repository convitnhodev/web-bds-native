package db

import (
	"database/sql"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"

	"github.com/deeincom/deeincom/pkg/helper"
	"github.com/lib/pq"
)

type MigrationModel struct {
	DB *sql.DB
}

func (m *MigrationModel) Find() ([]string, error) {
	var result []string
	q := `select migration_version from migrations order by migration_id desc;`
	rows, err := m.DB.Query(q)
	if err != nil {
		return result, err
	}
	defer rows.Close()
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return result, err
		}
		result = append(result, version)
	}
	return result, nil
}

func (m *MigrationModel) Create(version string) error {
	q := `insert into migrations (migration_version) values($1);`
	_, err := m.DB.Exec(q, version)
	return err
}

func (m *MigrationModel) Migrate() error {
	completed, err := m.Find()
	// only allow undefined_table, all other err will result in panic
	if err != nil {
		if err, ok := err.(*pq.Error); !ok {
			return err
		} else if err.Code.Name() != "undefined_table" {
			return err
		}
	}

	files, err := filepath.Glob("./sql/*.sql")
	if err != nil {
		return err
	}
	sort.Strings(files)
	for _, file := range files {
		filename := filepath.Base(file)

		// skip completed
		if helper.Contains(filename, completed) {
			continue
		}
		// exec pending sql file into database
		if err := m.Run(file); err != nil {
			return err
		}
		// create a new rows in migrations table to record the process
		if err := m.Create(filename); err != nil {
			return err
		}
	}
	return nil
}

// Run execute content of file into database
func (m *MigrationModel) Run(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	_, err = m.DB.Exec(string(b))
	return err
}
