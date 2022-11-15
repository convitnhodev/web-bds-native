package db

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/deeincom/deeincom/pkg/models"
	"github.com/pkg/errors"
)

type FileModel struct {
	DB         *sql.DB
	Pagination *Pagination
}

var fileUserColumes = []string{
	"files.id",
	"files.local_path",
	"files.cloud_link",
	"files.status",
	"files.created_at",
	"files.updated_at",
}

func scanFile(r scanner, o *models.File) error {
	if err := r.Scan(
		&o.ID,
		&o.LocalPath,
		&o.CloudLink,
		&o.Status,
		&o.CreatedAt,
		&o.UpdatedAt,
	); err != nil {
		return errors.Wrap(err, "scanFile")
	}

	return nil
}

func (m *FileModel) query(s string) string {
	return fmt.Sprintf(`SELECT %s FROM files %s`, strings.Join(fileUserColumes, ","), s)
}

func (m *FileModel) count(s string) string {
	return fmt.Sprintf(`SELECT count(*) FROM files %s`, s)
}

func (m *FileModel) Local(s string) (*models.File, error) {
	q := m.query(`WHERE files.local_path = $1`)
	row := m.DB.QueryRow(q, s)
	o := new(models.File)
	if err := scanFile(row, o); err != nil {
		return nil, errors.Wrap(err, "Files.LocalPath")
	}
	return o, nil
}

func (m *FileModel) Create(localPath string) (*models.File, error) {
	q := `INSERT INTO files (local_path) VALUES ($1) returning id`
	row := m.DB.QueryRow(q, localPath)
	o := new(models.File)
	if err := row.Scan(&o.ID); err != nil {
		return nil, errors.Wrap(err, "FileModel.Create")
	}
	return o, nil
}

func (m *FileModel) UploadCloudLink(localPath string, cloudLink string) error {
	q := "UPDATE files SET cloud_link = $2, status = 'sync' WHERE files.local_path = $1"
	_, err := m.DB.Exec(q, localPath, cloudLink)
	return err
}

func (m *FileModel) NotUpload() ([]*models.File, error) {
	q := m.query(`WHERE files.cloud_link = ''`)

	rows, err := m.DB.Query(q)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := []*models.File{}
	for rows.Next() {
		o := &models.File{}
		if err := scanFile(rows, o); err != nil {
			log.Println(err)
		}
		list = append(list, o)
	}

	return list, nil
}
