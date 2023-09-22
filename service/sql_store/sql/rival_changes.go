package sql

import (
	"context"
	"fmt"

	v1 "github.com/superjcd/rivalservice/genproto/v1"
	"github.com/superjcd/rivalservice/service/sql_store"
	"gorm.io/gorm"
)

type rival_changes struct {
	db *gorm.DB
}

var _ sql_store.RivalChangeStore = (*rival_changes)(nil)

func (rc *rival_changes) Append(ctx context.Context, rq *v1.AppendRivalChangesRequest) error {
	rivalChanges := make([]sql_store.RivalChange, 0, 512)
	sql := fmt.Sprintf(`
	        SELECT
				t1.country country,
				t1.asin asin,
				'%s' as field,
				t1.price as old_value,
				t2.price as new_value,
				'%s' as create_date
			FROM
			(
			Select
				country,
				asin,
				%s
			FROM rival_produt_active_details
			WHERE create_date = '%s'
				and %s != ''
			)t1 LEFT JOIN
			(
				Select
				country,
				asin,
				%s
			FROM rival_produt_active_details
			WHERE create_date = '%s'
				and %s != ''
			)t2
			on t1.asin = t2.asin
			and t1.country = t2.country
			where t1.price != t2.price
		`, rq.Field, rq.NewDate, rq.Field, rq.OldDate, rq.Field, rq.Field, rq.NewDate, rq.Field)

	d := rc.db.Raw(sql).Scan(&rivalChanges)
	if d.Error != nil {
		return d.Error
	}
	if len(rivalChanges) > 0 {
		return rc.db.Create(&rivalChanges).Error
	}

	return nil

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
			AND field = '%s'
		) t2 ON t1.rival_asin = t2.asin 
		  AND t1.country = t2.country	
		
	`
	sql = fmt.Sprintf(sql, rq.User, rq.Country, rq.Country, rq.CreateDate, rq.Field)

	rc.db.Raw(sql).Scan(&userChanges)

	return &sql_store.UserRivalChangeList{
		TotalCount: len(userChanges),
		Items:      userChanges,
	}, nil

}

func (rc *rival_changes) Delete(ctx context.Context, rq *v1.DeleteRivalChangesRequest) error {
	if rq.MinCreateDate != "" {
		return rc.db.Unscoped().Where("create_date < ?", rq.MinCreateDate).Delete(&sql_store.RivalChange{}).Error
	}

	return fmt.Errorf("min_create_date不可谓空字段")
}
