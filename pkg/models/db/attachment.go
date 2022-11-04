package db

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/deeincom/deeincom/pkg/form"
	"github.com/deeincom/deeincom/pkg/models"
	"github.com/pkg/errors"
)

type AttachmentModel struct {
	DB         *sql.DB
	Pagination *Pagination
}

var attachmentColumes = []string{
	"attachments.id",
	"attachments.title",
	"attachments.content_type",
	"attachments.mime_type",
	"attachments.link",
	"attachments.width",
	"attachments.height",
	"attachments.size",
	"attachments.video_thumbnail",
	"attachments.video_length",
	"attachments.product_id",
}

func (m *AttachmentModel) query(s string) string {
	return fmt.Sprintf(`SELECT %s FROM attachments %s`, strings.Join(attachmentColumes, ","), s)
}

func (m *AttachmentModel) count(s string) string {
	return fmt.Sprintf(`SELECT count(*) FROM attachments %s`, s)
}

func scanAttachment(r scanner, o *models.Attachment) error {
	if err := r.Scan(
		&o.ID,
		&o.Title,
		&o.ContentType,
		&o.MineType,
		&o.Link,
		&o.Width,
		&o.Height,
		&o.Size,
		&o.VideoThumbnail,
		&o.VideoLength,
		&o.Product.ID,
	); err != nil {
		return errors.Wrap(err, "scanAttachment")
	}

	return nil
}

func (m *AttachmentModel) ID(id string) (*models.Attachment, error) {
	q := m.query(`where attachments.id = $1`)
	row := m.DB.QueryRow(q, id)
	o := new(models.Attachment)
	if err := scanAttachment(row, o); err != nil {
		return nil, errors.Wrap(err, "AttachmentModel.ID")
	}
	return o, nil
}

func (m *AttachmentModel) Find(product *models.Product) ([]*models.Attachment, error) {
	q := m.query("where product_id = $1 and updated_at > '0001-01-01 00:00:00+00'::date order by id desc")
	count := m.count("where product_id = $1 and updated_at > '0001-01-01 00:00:00+00'::date")

	if err := m.Pagination.Count(count, product.ID); err != nil {
		return nil, err
	}
	rows, err := m.DB.Query(m.Pagination.Generate(q), product.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := []*models.Attachment{}
	for rows.Next() {
		o := &models.Attachment{}
		if err := scanAttachment(rows, o); err != nil {
			log.Println(err)
		}
		list = append(list, o)
	}
	return list, nil
}

func (m *AttachmentModel) Create(product *models.Product, t string) (*models.Attachment, error) {
	q := `insert into attachments (product_id, content_type) values ($1, $2) returning id`
	row := m.DB.QueryRow(q, product.ID, t)
	o := new(models.Attachment)
	if err := row.Scan(&o.ID); err != nil {
		return nil, errors.Wrap(err, "AttachmentModel.Create")
	}
	return o, nil
}

func (m *AttachmentModel) Update(o *models.Attachment, f *form.Form) error {
	q := `
		update
			attachments
		set
			updated_at = now(),
			title = $2,
			link = $3,
			size = $4
		where
			id = $1
	`
	_, err := m.DB.Exec(q,
		o.ID,
		f.Get("Title"),
		f.Get("Link"),
		f.Get("Size"),
	)

	return err
}
