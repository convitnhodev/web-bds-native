package models

import "github.com/upper/db/v4"

var ProvinceTable = "provinces"

type Province struct {
	Code       string `json:"code" db:"code"`
	Name       string `json:"name" db:"name"`
	NameEn     string `json:"name_en" db:"name_en"`
	FullName   string `json:"full_name" db:"full_name"`
	FullNameEn string `json:"full_name_en" db:"full_name_en"`
	CodeName   string `json:"code_name" db:"code_name"`
}

func (p *Province) Store(sess db.Session) db.Store {
	return sess.Collection(ProvinceTable)
}

var _ = db.Record(&Province{})
