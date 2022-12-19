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

type Location struct {
	WardID       string `db:"ward_id"`
	WardName     string `db:"ward_name"`
	DistrictID   string `db:"district_id"`
	DistrictName string `db:"district_name"`
	ProvinceID   string `db:"province_id"`
	ProvinceName string `db:"province_name"`
}

func (r *LocationRepository) ListLocationByWardIDs(ids []string) (map[string]Location, error) {
	rms := make([]Location, 0)

	rmm := make(map[string]Location, 0)

	is := make([]interface{}, len(ids))
	for _, id := range ids {
		is = append(is, id)
	}

	err := r.r.db.SQL().Select(
		"w.code as ward_id",
		"w.full_name as ward_name",
		"d.code as district_id",
		"d.full_name as district_name",
		"p.code as province_id",
		"p.full_name as province_name",
	).
		From("wards AS w").
		Join("districts as d").On("w.district_code = d.code").
		Join("provinces as p").On("d.province_code = p.code").
		Where(db.Cond{"w.code": db.In(is...)}).
		All(&rms)
	if err != nil {
		return nil, err
	}
	for _, rm := range rms {
		rmm[rm.WardID] = rm
	}
	return rmm, nil
}
