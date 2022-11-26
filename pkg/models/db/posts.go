package db

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/Machiel/slugify"
	"github.com/deeincom/deeincom/pkg/form"
	"github.com/deeincom/deeincom/pkg/models"
)

type PostModel struct {
	DB         *sql.DB
	Pagination *Pagination
}

var postUserColumes = []string{
	"posts.id",
	"posts.title",
	"posts.thumbnail",
	"posts.post_type",
	"posts.slug",
	"posts.tags",
	"posts.short",
	"posts.content",
	"posts.published_at",
	"posts.created_at",
	"posts.updated_at",
	"users.id AS user_id",
	"users.first_name",
	"users.last_name",
}

func scanPost(r scanner, o *models.Post, includePoster bool) error {
	if includePoster {
		o.Poster = models.User{}
		if err := r.Scan(
			&o.ID,
			&o.Title,
			&o.Thumbnail,
			&o.PostType,
			&o.Slug,
			pq.Array(&o.Tags),
			&o.Short,
			&o.Content,
			&o.PublishedAt,
			&o.CreatedAt,
			&o.UpdatedAt,
			&o.Poster.ID,
			&o.Poster.FirstName,
			&o.Poster.LastName,
		); err != nil {
			return errors.Wrap(err, "scanPost")
		}
	} else {
		if err := r.Scan(
			&o.ID,
			&o.Title,
			&o.Thumbnail,
			&o.PostType,
			&o.Slug,
			pq.Array(&o.Tags),
			&o.Short,
			&o.Content,
			&o.PublishedAt,
			&o.CreatedAt,
			&o.UpdatedAt,
		); err != nil {
			return errors.Wrap(err, "scanPost")
		}
	}

	return nil
}

func (m *PostModel) query(s string, includePoster bool) string {
	posterSql := ""
	columes := postUserColumes[0:11]
	if includePoster {
		columes = postUserColumes
		posterSql = "LEFT JOIN users ON posts.poster_id = users.id"
	}

	return fmt.Sprintf(`SELECT %s FROM posts %s  %s`, strings.Join(columes, ","), posterSql, s)
}

func (m *PostModel) count(s string) string {
	return fmt.Sprintf(`SELECT count(*) FROM posts %s`, s)
}

func (m *PostModel) Find() ([]*models.Post, error) {
	q := m.query("order by posts.id desc", true)
	count := m.count("")

	if err := m.Pagination.Count(count); err != nil {
		return nil, err
	}

	rows, err := m.DB.Query(m.Pagination.Generate(q))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := []*models.Post{}
	for rows.Next() {
		o := &models.Post{}
		if err := scanPost(rows, o, true); err != nil {
			log.Println(err)
		}
		list = append(list, o)
	}
	return list, nil
}

func (m *PostModel) Published() ([]*models.Post, error) {
	q := m.query("where posts.published_at <= now() order by posts.id desc", true)
	count := m.count("")

	if err := m.Pagination.Count(count); err != nil {
		return nil, err
	}

	rows, err := m.DB.Query(m.Pagination.Generate(q))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := []*models.Post{}
	for rows.Next() {
		o := &models.Post{}
		if err := scanPost(rows, o, true); err != nil {
			log.Println(err)
		}
		list = append(list, o)
	}
	return list, nil
}

func (m *PostModel) Tags(tags []string) []*models.Post {
	q := m.query(
		fmt.Sprintf(
			"where posts.tags @> '{%s}' order by posts.id desc",
			strings.Join(tags, ","),
		),
		false,
	)

	rows, err := m.DB.Query(q)
	if err != nil {
		return nil
	}

	defer rows.Close()

	list := []*models.Post{}
	for rows.Next() {
		o := &models.Post{}
		if err := scanPost(rows, o, false); err != nil {
			log.Println(err)
		}
		list = append(list, o)
	}
	return list
}

func (m *PostModel) Create(userId int, postType string) (*models.Post, error) {
	if userId == 0 {
		return nil, errors.New("err_user_id_empty")
	}

	if postType == "" {
		return nil, errors.New("err_post_type_empty")
	}

	q := fmt.Sprintf(`INSERT INTO posts (poster_id, post_type) VALUES (%d, '%s') RETURNING id`, userId, postType)
	row := m.DB.QueryRow(q)

	o := new(models.Post)
	if err := row.Scan(&o.ID); err != nil {
		return nil, errors.Wrap(err, "Post.Create")
	}

	return o, nil
}

func (m *PostModel) Update(o *models.Post, f *form.Form) error {
	q := `
		UPDATE
			posts
		SET
			updated_at = now(),
			title = $2,
			short = $3,
			slug = $4,
			"content" = $5,
			published_at = $6,
			tags = $7,
			thumbnail = $8,
			post_type = $9
		WHERE
			id = $1
	`

	tags := []string{}
	for _, t := range strings.Split(f.Get("Tags"), ",") {
		t = strings.TrimSpace(t)
		if t != "" {
			tags = append(tags, t)
		}
	}

	var publishedAt *string = nil
	if len(f.Get("PublishedAt")) > 0 {
		time := fmt.Sprintf("{%s} +0700", f.Get("PublishedAt"))
		publishedAt = &time
	}

	_, err := m.DB.Exec(q,
		o.ID,
		f.Get("Title"),
		f.Get("Short"),
		slugify.Slugify(f.Get("Title")),
		f.Get("Content"),
		&publishedAt,
		fmt.Sprintf("{%s}", strings.Join(tags, ",")),
		f.Get("Thumbnail"),
		f.Get("PostType"),
	)

	return err
}

func (m *PostModel) ID(id string) (*models.Post, error) {
	if id == "" {
		return nil, errors.New("err_id_empty")
	}

	q := m.query(`WHERE posts.id = $1`, false)
	row := m.DB.QueryRow(q, id)
	o := new(models.Post)

	if err := scanPost(row, o, false); err != nil {
		return nil, errors.Wrap(err, "user.ID")
	}

	return o, nil
}

func (m *PostModel) Remove(id string) error {
	o, err := m.ID(id)

	if err != nil {
		return err
	}

	q := "DELETE FROM posts WHERE posts.id = $1;"
	_, err = m.DB.Exec(q,
		id,
	)

	if err != nil {
		return err
	}

	q = "DELETE FROM comments WHERE comments.slug = $1;"
	_, err = m.DB.Exec(q,
		o.Slug,
	)

	return err
}

func (m *PostModel) GetBySlug(slug string) (*models.Post, error) {
	q := m.query(`where posts.slug = $1`, true)
	row := m.DB.QueryRow(q, slug)
	o := new(models.Post)

	if err := scanPost(row, o, true); err != nil {
		return nil, errors.Wrap(err, "Post.GetBySlug")
	}

	return o, nil
}
