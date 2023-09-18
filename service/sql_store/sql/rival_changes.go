package sql

import (
	"context"
	"fmt"

	v1 "github.com/HooYa-Bigdata/rivalservice/genproto/v1"
	"github.com/HooYa-Bigdata/rivalservice/service/sql_store"
	"gorm.io/gorm"
)

type rival_changes struct {
	db *gorm.DB
}

var _ sql_store.RivalChangeStore = (*rival_changes)(nil)

func (rc *rival_changes) Append(ctx context.Context, rq *v1.AppendRivalChangesRequest) error {
	rivalChanges := make([]sql_store.RivalChange, 0, 512)

	rc.db.Raw(`SELECT
				t1.country country,
				t1.asin asin,
				'price' as field,
				t1.price as old_value,
				t2.price as new_value,
				? as create_date
			FROM
			(
			Select
				country,
				asin,
				price
			FROM rival_produt_active_details
			WHERE create_date = ?
				and price != ''
			)t1 LEFT JOIN
			(
				Select
				country,
				asin,
				price
			FROM rival_produt_active_details
			WHERE create_date = ?
				and price != ''
			)t2
			on t1.asin = t2.asin
			and t1.country = t2.country
			where t1.price != t2.price`, rq.NewDate, rq.OldDate, rq.NewDate).Scan(&rivalChanges)

	return rc.db.Create(rivalChanges).Error

}

func (rc *rival_changes) List(ctx context.Context, rq *v1.ListRivalChangesRequest) (*sql_store.UserRivalChangeList, error) {
	userChanges := make([]sql_store.UserRivalChange, 0, 32)
	sql := `
	   SELECT 
	     t1.asin as my_asin,
		 t1.country as country,
		 t2.asin as rival_asin,
		 t2.field field,
		 t2.old_value,
		 t2.new_value,
		 t2.create_date
		FROM
		(
		SELECT 
		  asin,
		  rival_asin,
		  country
		FROM rivals 
		WHERE user = '%s'
		  AND country = '%s'
		) t1 
		LEFT JOIN (
		  SELECT 
			asin,
			country,
			field,
			old_value, 
			new_value,
			create_date
		  FROM rival_changes
		  WHERE  country = '%s'		
			AND create_date = '%s'
		) t2 ON t1.rival_asin = t2.asin 
		  AND t1.country = t2.country	
	`
	sql = fmt.Sprintf(sql, rq.User, rq.Country, rq.Country, rq.CreateDate)

	rc.db.Raw(sql).Scan(&userChanges)

	return &sql_store.UserRivalChangeList{
		TotalCount: len(userChanges),
		Items:      userChanges,
	}, nil

}

func (rc *rival_changes) Delete(ctx context.Context, rq *v1.DeleteRivalChangesRequest) error {
	if rq.MinCreateDate != "" {
		return rc.db.Where("create_date < ?", rq.MinCreateDate).Delete(&sql_store.RivalChange{}).Error
	}

	return fmt.Errorf("min_create_date不可谓空字段")
}
