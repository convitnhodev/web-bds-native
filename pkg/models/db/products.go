package db

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/deeincom/deeincom/pkg/models"
	"github.com/pkg/errors"
)

type ProductModel struct {
	DB         *sql.DB
	Pagination *Pagination
}

var productColumes = []string{
	"products.product_id",
}

func (m *ProductModel) query(s string) string {
	return fmt.Sprintf(`SELECT %s FROM products %s`, strings.Join(productColumes, ","), s)
}

func (m *ProductModel) count(s string) string {
	return fmt.Sprintf(`SELECT count(*) FROM products %s`, s)
}

func scanProduct(r scanner, o *models.Product) error {
	if err := r.Scan(
		&o.ID,
	); err != nil {
		return errors.Wrap(err, "scanProduct")
	}

	return nil
}

func (m *ProductModel) Find() ([]*models.Product, error) {
	q := m.query("order by product_id desc")
	count := m.count("")

	if err := m.Pagination.Count(count); err != nil {
		return nil, err
	}
	rows, err := m.DB.Query(m.Pagination.Generate(q))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := []*models.Product{}
	for rows.Next() {
		o := &models.Product{}
		if err := scanProduct(rows, o); err != nil {
			log.Println(err)
		}
		list = append(list, o)
	}
	return list, nil
}

func (m *ProductModel) GetBySlug(slug string) (*models.Product, error) {
	q := m.query(`where products.product_slug = $1`)
	row := m.DB.QueryRow(q, slug)
	o := new(models.Product)
	if err := scanProduct(row, o); err != nil {
		return nil, errors.Wrap(err, "Products.GetBySlug")
	}
	return o, nil
}
