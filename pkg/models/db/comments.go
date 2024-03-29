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

type CommentModel struct {
	DB         *sql.DB
	Pagination *Pagination
}

var commentUserColumes = []string{
	"comments.id",
	"comments.parent_id",
	"comments.slug",
	"comments.message",
	"comments.is_censorship",
	"comments.created_at",
	"comments.updated_at",
	"users.id AS user_id",
	"users.first_name",
	"users.last_name",
}

func scanComment(r scanner, o *models.Comment) error {
	o.Poster = models.User{}
	if err := r.Scan(
		&o.ID,
		&o.ParrentId,
		&o.Slug,
		&o.Message,
		&o.IsCensorship,
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
	return fmt.Sprintf(`SELECT %s FROM comments LEFT JOIN users ON comments.user_id = users.id %s`, strings.Join(commentUserColumes, ","), s)
}

func (m *CommentModel) count(s string) string {
	return fmt.Sprintf(`SELECT count(*) FROM comments %s`, s)
}

func (m *CommentModel) Find() ([]*models.Comment, error) {
	q := m.query("WHERE comments.is_censorship = False ORDER BY comments.id desc")
	count := m.count("WHERE comments.is_censorship = False")

	if err := m.Pagination.Count(count); err != nil {
		return nil, err
	}

	rows, err := m.DB.Query(m.Pagination.Generate(q))
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

func (m *CommentModel) Slug(slug string) ([]*models.Comment, error) {
	q := m.query(fmt.Sprintf("WHERE comments.slug = '%s' AND comments.is_censorship = true", slug))

	rows, err := m.DB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := []*models.Comment{}
	commentMap := make(map[int]*models.Comment)
	for rows.Next() {
		o := &models.Comment{}
		if err := scanComment(rows, o); err != nil {
			log.Println(err)
		}

		if o.ParrentId == nil {
			v, has := commentMap[o.ID]
			if has {
				o.ChildComments = v.ChildComments
			}
			commentMap[o.ID] = o

			list = append(list, o)
		} else {
			v, has := commentMap[*o.ParrentId]
			if has {
				v.ChildComments = append(v.ChildComments, o)
			} else {
				parentComment := &models.Comment{}
				parentComment.ChildComments = []*models.Comment{o}
			}
		}
	}

	return list, nil
}

func (m *CommentModel) Create(f *form.Form) (*models.Comment, error) {
	q := `
		INSERT INTO comments (user_id, parent_id, slug, message, is_censorship)
		VALUES ($1, $2, $3, $4, false)
		RETURNING id`
	parrentId := f.GetInt("ParentId")

	var row *sql.Row
	if parrentId > 0 {
		row = m.DB.QueryRow(q,
			f.GetInt("UserId"),
			parrentId,
			f.Get("Slug"),
			f.Get("Message"),
		)
	} else {
		row = m.DB.QueryRow(q,
			f.GetInt("UserId"),
			nil,
			f.Get("Slug"),
			f.Get("Message"),
		)
	}

	o := new(models.Comment)
	if err := row.Scan(&o.ID); err != nil {
		return nil, errors.Wrap(err, "Comment.Create")
	}

	return o, nil
}

func (m *CommentModel) Remove(id string) error {
	q := "DELETE FROM comments WHERE comments.id = $1;"
	_, err := m.DB.Exec(q,
		id,
	)

	return err
}

func (m *CommentModel) ChangeCensorship(id string) error {
	q := `
		UPDATE comments SET
			updated_at = now(),
			is_censorship = not(is_censorship)
		WHERE id = $1`

	_, err := m.DB.Exec(q,
		id,
	)

	return err
}
