package sql

import (
	"context"
	"fmt"

	v1 "github.com/HooYa-Bigdata/rivalservice/genproto/v1"
	"github.com/HooYa-Bigdata/rivalservice/service/sql_store"
	"gorm.io/gorm"
)

type product_details struct {
	db *gorm.DB
}

var _ sql_store.ProductDetailsStore = (*product_details)(nil)

func (pd *product_details) AppendActiveDetail(ctx context.Context, rq *v1.AppendRivalProductActiveDetailRequest) error {
	active_info := sql_store.RivalProdutActiveDetail{
		Asin:            rq.Details.Asin,
		Country:         rq.Details.Country,
		Price:           rq.Details.Price,
		Currency:        rq.Details.Currency,
		Coupon:          rq.Details.Coupon,
		Star:            rq.Details.Star,
		Ratings:         rq.Details.Ratings,
		Image:           rq.Details.Image,
		ParentAsin:      rq.Details.ParentAsin,
		CategoryInfo:    rq.Details.CategoryInfo,
		TopCategoryName: rq.Details.TopCategoryName,
		TopCategoryRank: rq.Details.TopCategoryRank,
		Color:           rq.Details.Color,
		Weight:          rq.Details.Weight,
		WeightUnit:      rq.Details.WeightUnit,
		Dimensions:      rq.Details.Dimensions,
		DimensionsUnit:  rq.Details.DimensionsUnit,
		CreateDate:      rq.Details.CreateDate,
	}

	return pd.db.Create(&active_info).Error
}

func (pd *product_details) DeleteActiveDetail(ctx context.Context, rq *v1.DeleteRivalActiveDetailRequest) error {
	if rq.MinCreateDate != "" {
		return pd.db.Where("create_date < ?", rq.MinCreateDate).Delete(&sql_store.RivalProdutActiveDetail{}).Error
	}
	return fmt.Errorf("min create date参数不可以为空")
}

func (pd *product_details) AppendInactiveDetail(ctx context.Context, rq *v1.AppendRivalProductInactiveDetailRequest) error {
	inactive_info := sql_store.RivalProdutInactiveDetail{
		Asin:         rq.Asin,
		Country:      rq.Country,
		Title:        rq.Title,
		BulletPoints: rq.BulletPoints,
		CreateDate:   rq.CreateDate,
	}

	return pd.db.Create(&inactive_info).Error
}

func (pd *product_details) DeleteInactiveDetail(ctx context.Context, rq *v1.DeleteRivalInactiveDetailRequest) error {
	if rq.MinCreateDate != "" {
		return pd.db.Where("create_date < ?", rq.MinCreateDate).Delete(&sql_store.RivalProdutInactiveDetail{}).Error
	}
	return fmt.Errorf("min create date参数不可以为空")
}
