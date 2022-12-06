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
	"products.slug",
	"products.full_content",
	"products.product_type",
	"products.business_advantage",
	"products.legal",
	"products.area",
	"products.poster_link",
	"products.house_certificate_link",
	"products.finance_plan_link",
	"products.num_of_slot",
	"products.cost_per_slot",
	"products.deposit_percent",
	"products.created_by",
	"products.is_censorship",
	"products.censored_at",
	"products.remain_of_slot",
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
		&o.Slug,
		&o.FullContent,
		&o.Type,
		&o.BusinessAdvantage,
		&o.Legal,
		&o.Area,
		&o.PosterLink,
		&o.HouseCertificateLink,
		&o.FinancePlanLink,
		&o.NumOfSlot,
		&o.CostPerSlot,
		&o.DepositPercent,
		&o.CreatedBy,
		&o.IsCensorship,
		&o.CensoredAt,
		&o.RemainOfSlot,
	); err != nil {
		return errors.Wrap(err, "scanProduct")
	}

	return nil
}

func (m *ProductModel) Find() ([]*models.Product, error) {
	q := m.query("where is_deleted = false and updated_at > '0001-01-01 00:00:00+00'::date order by id desc")
	count := m.count("where is_deleted = false and updated_at > '0001-01-01 00:00:00+00'::date")

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

func (m *ProductModel) Published() ([]*models.Product, error) {
	q := m.query("where is_censorship != false and is_deleted = false and updated_at > '0001-01-01 00:00:00+00'::date order by id desc")
	count := m.count("where is_censorship != false and is_deleted = false and updated_at > '0001-01-01 00:00:00+00'::date")

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

func (m *ProductModel) CreatedBy(userId int) ([]*models.Product, error) {
	q := m.query("where created_by = $1 and is_deleted = false and updated_at > '0001-01-01 00:00:00+00'::date order by id desc")
	count := m.count("where created_by = $1 and is_deleted = false and updated_at > '0001-01-01 00:00:00+00'::date")

	if err := m.Pagination.Count(count, userId); err != nil {
		return nil, err
	}
	rows, err := m.DB.Query(m.Pagination.Generate(q), userId)
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

func (m *ProductModel) ID(id string) (*models.Product, error) {
	q := m.query(`where products.id = $1`)
	row := m.DB.QueryRow(q, id)
	o := new(models.Product)
	if err := scanProduct(row, o); err != nil {
		return nil, errors.Wrap(err, "Products.ID")
	}
	return o, nil
}

func (m *ProductModel) Create(userId int, isCensorship bool) (*models.Product, error) {
	q := `insert into products (title, created_by, is_censorship) values ('', $1, $2) returning id`
	row := m.DB.QueryRow(q, userId, isCensorship)
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
			updated_at = now(),
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
			slug = $20,
			full_content = $21,
			business_advantage = $22,
			product_type = $23,
			legal = $24,
			area = $25,
			num_of_slot = $26,
			cost_per_slot = $27,
			deposit_percent = $28,
			remain_of_slot = (CASE WHEN remain_of_slot = num_of_slot THEN $29 ELSE remain_of_slot END)
		where
			id = $1
	`

	// Xử lí slug
	slug := slugify.Slugify(f.Get("Title"))
	if slug != o.Slug {
		productCount := 0
		count := m.count("WHERE id != $1 And slug like $2")

		row := m.DB.QueryRow(count, o.ID, fmt.Sprintf("%s%s", slug, "%"))
		if err := row.Scan(&productCount); err != nil {
			return err
		}

		if productCount > 0 {
			slug = fmt.Sprintf("%s-%d", slug, productCount)
		}
	}

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
		slug,
		f.Get("FullContent"),
		f.Get("BusinessAdvantage"),
		f.Get("Type"),
		f.Get("Legal"),
		f.Get("Area"),
		f.GetInt("NumOfSlot"),
		f.GetInt("CostPerSlot"),
		f.GetFloat("DepositPercent"),
		f.GetInt("NumOfSlot"),
	)

	if err != nil {
		return err
	}

	if slug != o.Slug {
		prefix := "/real-estate"
		_, err = m.DB.Exec(
			"UPDATE comments SET slug = $2 WHERE slug = $1;",
			fmt.Sprintf("%s/%s", prefix, o.Slug),
			fmt.Sprintf("%s/%s", prefix, slug),
		)
	}

	return err
}

func (m *ProductModel) Set(id string, key string, value string) error {
	q := fmt.Sprintf(`
		UPDATE
			products
		SET
			%s = $2
		WHERE id = $1
	`, key)

	_, err := m.DB.Exec(q, id, value)

	return err
}

func (m *ProductModel) Remove(id string) error {
	q := `update products set is_deleted = true where id = $1`
	_, err := m.DB.Exec(q, id)
	return err
}

func (m *ProductModel) Approve(id string) error {
	q := `update products set is_censorship = true, censored_at = now() where id = $1`
	_, err := m.DB.Exec(q, id)
	return err
}
