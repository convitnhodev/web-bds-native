package repositories

import (
	"github.com/deeincom/deeincom/app/models"
	"github.com/upper/db/v4"
)

type ProductRepository repository

func (r *ProductRepository) ListProducts() ([]models.Product, error) {
	ps := make([]models.Product, 0)
	var p models.Product
	res := r.r.db.Collection(models.ProductTable).Find().Limit(10)
	defer res.Close()

	for res.Next(&p) {
		ps = append(ps, p)
	}
	if err := res.Err(); err != nil {
		return nil, err
	}
	return ps, nil
}

func (r *ProductRepository) FindByID(id string) (*models.Product, error) {
	var p models.Product
	if err := r.Find(db.Cond{"id": id}).One(&p); err != nil {
		return nil, err
	}
	return &p, nil
}
