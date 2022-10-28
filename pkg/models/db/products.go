package db

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/Machiel/slugify"
	"github.com/deeincom/deeincom/pkg/form"
	"github.com/deeincom/deeincom/pkg/models"
	"github.com/pkg/errors"
)

type ProductModel struct {
	DB         *sql.DB
	Pagination *Pagination
}

var productColumes = []string{
	"products.id",
	"products.title",
	"products.short_description",
	"products.full_description",
	"products.city",
	"products.district",
	"products.ward",
	"products.address_number",
	"products.street",
	"products.house_direction",
	"products.balcony_direction",
	"products.financial_plan",
	"products.furniture",
	"products.bedroom",
	"products.toilet",
	"products.floor",
	"products.front_width",
	"products.street_width",
	"products.pavement_width",
	"products.created_at",
	"products.updated_at",
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
		&o.Title,
		&o.Short,
		&o.Full,
		&o.City,
		&o.District,
		&o.Ward,
		&o.AddressNumber,
		&o.Street,
		&o.HouseDirection,
		&o.BalconyDirection,
		&o.FinancialPlan,
		&o.Furniture,
		&o.Bedroom,
		&o.Toilet,
		&o.Floor,
		&o.FrontWidth,
		&o.StreetWidth,
		&o.PavementWidth,
		&o.CreatedAt,
		&o.UpdatedAt,
	); err != nil {
		return errors.Wrap(err, "scanProduct")
	}

	return nil
}

func (m *ProductModel) Find() ([]*models.Product, error) {
	q := m.query("order by id desc")
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
	q := m.query(`where products.slug = $1`)
	row := m.DB.QueryRow(q, slug)
	o := new(models.Product)
	if err := scanProduct(row, o); err != nil {
		return nil, errors.Wrap(err, "Products.GetBySlug")
	}
	return o, nil
}

func (m *ProductModel) Create(f *form.Form) (*models.Product, error) {
	q := `insert into products (title) values ('') returning id`
	row := m.DB.QueryRow(q)
	o := new(models.Product)
	if err := row.Scan(&o.ID); err != nil {
		return nil, errors.Wrap(err, "ProductModel.Create")
	}
	return o, nil
}

func (m *ProductModel) Update(o *models.Product, f *form.Form) error {
	q := `
		update
			products
		set
			title = $2,
			short_description = $3,
			full_description = $4,
			city = $5,
			district = $6,
			ward = $7,
			address_number = $8,
			street = $9,
			house_direction = $10,
			balcony_direction = $11,
			financial_plan = $12,
			furniture = $13,
			bedroom = $14,
			toilet = $15,
			floor = $16,
			front_width = $17,
			street_width = $18,
			pavement_width = $19,
			slug = $20
		where
			id = $1
	`
	_, err := m.DB.Exec(q,
		o.ID,
		f.Get("Title"),
		f.Get("Short"),
		f.Get("Full"),
		f.Get("City"),
		f.Get("District"),
		f.Get("Ward"),
		f.Get("AddressNumber"),
		f.Get("Street"),
		f.Get("HouseDirection"),
		f.Get("BalconyDirection"),
		f.Get("FinancialPlan"),
		f.Get("Furniture"),
		f.GetInt("Bedroom"),
		f.GetInt("Toilet"),
		f.GetInt("Floor"),
		f.GetInt("FrontWidth"),
		f.GetInt("StreetWidth"),
		f.GetInt("PavementWidth"),
		slugify.Slugify(f.Get("Title")),
	)

	return err
}
