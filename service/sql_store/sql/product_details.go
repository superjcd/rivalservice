package sql

import (
	"context"
	"fmt"

	v1 "github.com/superjcd/rivalservice/genproto/v1"
	"github.com/superjcd/rivalservice/service/sql_store"
	"gorm.io/gorm"
)

type product_details struct {
	db *gorm.DB
}

var _ sql_store.ProductDetailsStore = (*product_details)(nil)

func (pd *product_details) AppendActiveDetail(ctx context.Context, rq *v1.AppendRivalProductActiveDetailRequest) error {
	activeInfos := make([]sql_store.RivalProdutActiveDetail, 0, 16)

	for _, d := range rq.Details {
		activeInfo := sql_store.RivalProdutActiveDetail{
			Asin:            d.Asin,
			Country:         d.Country,
			Price:           d.Price,
			Currency:        d.Currency,
			Coupon:          d.Coupon,
			Star:            d.Star,
			Ratings:         d.Ratings,
			Image:           d.Image,
			ParentAsin:      d.ParentAsin,
			CategoryInfo:    d.CategoryInfo,
			TopCategoryName: d.TopCategoryName,
			TopCategoryRank: d.TopCategoryRank,
			Color:           d.Color,
			Weight:          d.Weight,
			WeightUnit:      d.WeightUnit,
			Dimensions:      d.Dimensions,
			DimensionsUnit:  d.DimensionsUnit,
			CreateDate:      d.CreateDate,
		}
		activeInfos = append(activeInfos, activeInfo)
	}

	return pd.db.Create(&activeInfos).Error
}

func (pd *product_details) DeleteActiveDetail(ctx context.Context, rq *v1.DeleteRivalActiveDetailRequest) error {
	if rq.MinCreateDate != "" {
		return pd.db.Unscoped().Where("create_date < ?", rq.MinCreateDate).Delete(&sql_store.RivalProdutActiveDetail{}).Error
	}
	return fmt.Errorf("min create date参数不可以为空")
}

func (pd *product_details) AppendInactiveDetail(ctx context.Context, rq *v1.AppendRivalProductInactiveDetailRequest) error {
	inactiveInfos := make([]sql_store.RivalProdutInactiveDetail, 0, 16)

	for _, d := range rq.Details {
		inactiveInfo := sql_store.RivalProdutInactiveDetail{
			Asin:         d.Asin,
			Country:      d.Country,
			Title:        d.Title,
			BulletPoints: d.BulletPoints,
			CreateDate:   d.CreateDate,
		}

		inactiveInfos = append(inactiveInfos, inactiveInfo)
	}

	return pd.db.Create(&inactiveInfos).Error
}

func (pd *product_details) DeleteInactiveDetail(ctx context.Context, rq *v1.DeleteRivalInactiveDetailRequest) error {
	if rq.MinCreateDate != "" {
		return pd.db.Unscoped().Where("create_date < ?", rq.MinCreateDate).Delete(&sql_store.RivalProdutInactiveDetail{}).Error
	}
	return fmt.Errorf("min create date参数不可以为空")
}
