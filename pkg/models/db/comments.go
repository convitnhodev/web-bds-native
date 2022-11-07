package db

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/pkg/errors"
	"github.com/deeincom/deeincom/pkg/models"
)

type CommentModel struct {
	DB         *sql.DB
	Pagination *Pagination
}

var commentUserColumes = []string{
	"comments.id",
	"comments.parent_id",
	"comments.slug",
	"comments.message",
	"comments.created_at",
	"comments.updated_at",
	"users.id AS user_id",
	"users.first_name",
	"users.last_name",
}

func scanComment(r scanner, o *models.Comment) error {
	o.Poster = models.UserInfo{}
	if err := r.Scan(
		&o.ID,
		&o.Slug,
		&o.Message,
		&o.CreatedAt,
		&o.UpdatedAt,
		&o.Poster.ID,
		&o.Poster.FirstName,
		&o.Poster.LastName,
	); err != nil {
		return errors.Wrap(err, "scanComment")
	}

	return nil
}

func (m *CommentModel) query(s string) string {
	return fmt.Sprintf(`SELECT %s FROM comments %s LEFT JOIN users ON comments.user_id = users.id`, strings.Join(commentUserColumes, ","), s)
}

func (m *CommentModel) count(s string) string {
	return fmt.Sprintf(`SELECT count(*) FROM comments %s`, s)
}

func (m *CommentModel) Find(slug string) ([]*models.Comment, error) {
	q := m.query(fmt.Sprintf("where comments.slug = '%s'", slug))

	rows, err := m.DB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := []*models.Comment{}
	for rows.Next() {
		o := &models.Comment{}
		if err := scanComment(rows, o); err != nil {
			log.Println(err)
		}
		list = append(list, o)
	}
	return list, nil
}
