package repositories

import (
	"github.com/deeincom/deeincom/app/models"
	"time"
)

type SMSRepository repository

func (r *SMSRepository) Create(d *models.SMS) error {
	d.CreatedAt = time.Now()
	d.UpdatedAt = time.Now()
	if err := r.r.db.Collection(models.SMSTable).InsertReturning(d); err != nil {
		return err
	}
	return nil
}
