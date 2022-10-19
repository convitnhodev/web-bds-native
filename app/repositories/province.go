package repositories

import (
	"github.com/deeincom/deeincom/app/models"
	"github.com/upper/db/v4"
)

type LocationRepository repository

func (r *LocationRepository) ListProvinces() ([]models.Province, error) {
	provinces := make([]models.Province, 0)
	var province models.Province
	res := r.r.db.Collection(models.ProvinceTable).Find()
	defer res.Close()

	for res.Next(&province) {
		provinces = append(provinces, province)
	}

	if err := res.Err(); err != nil {
		return nil, err
	}
	return provinces, nil
}

func (r *LocationRepository) ListDistrictsByProvinceID(id string) ([]models.District, error) {
	ds := make([]models.District, 0)
	var d models.District
	res := r.r.db.Collection(models.DistrictTable).Find(db.Cond{"province_code": id})
	defer res.Close()

	for res.Next(&d) {
		ds = append(ds, d)
	}

	if err := res.Err(); err != nil {
		return nil, err
	}
	return ds, nil
}

func (r *LocationRepository) ListWardsByDistrictID(id string) ([]models.Ward, error) {
	ds := make([]models.Ward, 0)
	var d models.Ward
	res := r.r.db.Collection(models.WardTable).Find(db.Cond{"district_code": id})
	defer res.Close()

	for res.Next(&d) {
		ds = append(ds, d)
	}

	if err := res.Err(); err != nil {
		return nil, err
	}
	return ds, nil
}
