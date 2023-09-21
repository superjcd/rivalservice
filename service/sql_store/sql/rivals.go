package sql

import (
	"context"

	v1 "github.com/superjcd/rivalservice/genproto/v1"
	"github.com/superjcd/rivalservice/service/sql_store"
	"gorm.io/gorm"
)

type rivals struct {
	db *gorm.DB
}

var _ sql_store.RivalStore = (*rivals)(nil)

func (ur *rivals) Create(ctx context.Context, rq *v1.CreateRivalRequest) error {
	var err error
	for _, rivalAsin := range rq.Rivals {
		rival := sql_store.Rival{
			User:      rq.User,
			Asin:      rq.Asin,
			Country:   rq.Country,
			RivalAsin: rivalAsin,
		}
		err = ur.db.Create(&rival).Error

		if err != nil {
			return err
		}
	}

	return nil
}

func (ur *rivals) List(ctx context.Context, rq *v1.ListRivalRequest) (*sql_store.RivalList, error) {
	result := &sql_store.RivalList{}

	tx := ur.db

	if rq.User != "" {
		tx = tx.Where("user = ?", rq.User)
	}

	if rq.Country != "" {
		tx = tx.Where("country = ?", rq.Country)
	}

	if rq.Asin != "" {
		tx = tx.Where("asin = ?", rq.Asin)
	}

	d := tx.
		Offset(int(rq.Offset)).
		Limit(int(rq.Limit)).
		Find(&result.Items).
		Offset(-1).
		Limit(-1).
		Count(&result.TotalCount)

	return result, d.Error
}

func (ur *rivals) Delete(ctx context.Context, rq *v1.DeleteRivalRequest) error {
	tx := ur.db

	if rq.User != "" {
		tx = tx.Where("user = ?", rq.User)
	}

	if rq.Country != "" {
		tx = tx.Where("country = ?", rq.Country)
	}

	if rq.Asin != "" {
		tx = tx.Where("asin = ?", rq.Asin)
	}

	return tx.Delete(&sql_store.Rival{}).Error

}
