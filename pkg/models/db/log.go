package db

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/deeincom/deeincom/pkg/models"
	"github.com/pkg/errors"
)

type LogModel struct {
	DB         *sql.DB
	Pagination *Pagination
}

var logColumes = []string{
	"logs.id",
	"logs.user_id",
	"logs.content",
	"logs.created_at",
	"users.id",
	"users.first_name",
	"users.last_name",
}

func scanLog(r scanner, o *models.Log) error {
	o.Actor = models.User{}
	if err := r.Scan(
		&o.ID,
		&o.UserId,
		&o.Content,
		&o.CreatedAt,
		&o.Actor.ID,
		&o.Actor.FirstName,
		&o.Actor.LastName,
	); err != nil {
		return errors.Wrap(err, "scanFile")
	}

	return nil
}

func (m *LogModel) query(s string) string {
	return fmt.Sprintf(`SELECT %s FROM logs JOIN users ON logs.user_id = users.id  %s`, strings.Join(logColumes, ","), s)
}

func (m *LogModel) count(s string) string {
	return fmt.Sprintf(`SELECT count(*) FROM logs %s`, s)
}

func (m *LogModel) Add(userId string, content string) error {
	q := `INSERT INTO logs (user_id, content) VALUES ($1, $2)`

	_, err := m.DB.Exec(q,
		userId,
		content,
	)

	return err
}

func (m *LogModel) Find(userId string, date string) ([]*models.Log, error) {
	wheres := []string{}
	values := []string{}

	if userId != "" {
		wheres = append(wheres, "logs.user_id = $1")
		values = append(values, userId)
	}

	if date != "" {
		loc := len(wheres) + 1
		wheres = append(wheres, fmt.Sprintf("logs.created_at = $%d", loc))
		values = append(values, date)
	}
	orderStm := "ORDER BY logs.id DESC"
	queryWhere := ""
	if len(wheres) >= 1 {
		queryWhere = fmt.Sprintf("WHERE %s", strings.Join(wheres, " AND "))
	}

	q := m.query(fmt.Sprintf("%s %s", queryWhere, orderStm))
	count := m.count(queryWhere)

	inputs := make([]interface{}, len(values))

	for i, v := range values {
		inputs[i] = v
	}

	if err := m.Pagination.Count(count, inputs...); err != nil {
		return nil, err
	}

	rows, err := m.DB.Query(m.Pagination.Generate(q), inputs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := []*models.Log{}
	for rows.Next() {
		o := &models.Log{}
		if err := scanLog(rows, o); err != nil {
			log.Println(err)
		}
		list = append(list, o)
	}
	return list, nil
}
