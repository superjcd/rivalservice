package sql_store

import (
	v1 "github.com/superjcd/rivalservice/genproto/v1"
	"gorm.io/gorm"
)

// UserRivals
type Rival struct {
	gorm.Model
	User      string `json:"user" gorm:"column:user;index:idx_user_rivals,priority:1"`
	Asin      string `json:"asin" gorm:"column:asin;index:idx_user_rivals,priority:2"`
	Country   string `json:"country" gorm:"column:country;index:idx_user_rival,priority:3"`
	RivalAsin string `json:"rival_asin" gorm:"column:rival_asin;index:idx_user_rival, priority:4"`
}

type RivalList struct {
	TotalCount int64    `json:"totalCount"`
	Items      []*Rival `json:"items"`
}

func (rl *RivalList) ConvertToListRivalResponse(msg string, status v1.Status) v1.ListRivalResponse {
	rivals := make([]*v1.Rival, 0, 8)

	for _, rival := range rl.Items {
		rivals = append(rivals, &v1.Rival{
			User:      rival.User,
			Country:   rival.Country,
			Asin:      rival.Asin,
			RivalAsin: rival.RivalAsin,
		})
	}
	return v1.ListRivalResponse{
		Msg:    msg,
		Status: status,
		Rivals: rivals,
	}
}

// + my asin      join rival_asin = asin and country =country
type RivalChange struct {
	gorm.Model
	Asin       string `json:"asin" gorm:"column:asin;index:idx_rival_changes,priority:1"`
	Country    string `json:"country" gorm:"column:country;index:idx_rival_changes,priority:2"`
	Field      string `json:"field" gorm:"column:field"`
	OldValue   string `json:"old_value" gorm:"column:old_value"`
	NewValue   string `json:"new_value" gorm:"column:new_value"`
	CreateDate string `json:"create_date" gorm:"column:create_date;index: idx_rival_changes,priority:3"`
}

type UserRivalChange struct {
	gorm.Model
	MyAsin     string
	RivalAsin  string
	Country    string
	Field      string
	OldValue   string
	NewValue   string
	CreateDate string
}

type UserRivalChangeList struct {
	TotalCount int               `json:"totalCount"`
	Items      []UserRivalChange `json:"items"`
}

func (rcl *UserRivalChangeList) ConvertToListRivalChangeResponse(msg string, status v1.Status) v1.ListRivalChangesResponse {
	rcs := make([]*v1.RivalChange, 0, 8)

	for _, rc := range rcl.Items {
		rcs = append(rcs, &v1.RivalChange{
			Country:   rc.Country,
			MyAsin:    rc.MyAsin,
			RivalAsin: rc.RivalAsin,
			Field:     rc.Field,
			OldValue:  rc.OldValue,
			NewValue:  rc.NewValue,
		})
	}
	return v1.ListRivalChangesResponse{
		Msg:          msg,
		Status:       status,
		RivalChanges: rcs,
	}
}

type RivalProdutInactiveDetail struct {
	gorm.Model
	Asin         string `json:"asin" gorm:"column:asin;index:idx_inactive_product_details,priority:1"`
	Country      string `json:"country" gorm:"column:country;index:idx_inactive_product_details,priority:2"`
	Title        string `json:"title" gorm:"column:title"`
	BulletPoints string `json:"bullet_points" gorm:"column:bullet_points"`
	CreateDate   string `json:"create_date" gorm:"column:create_date;index:idx_inactive_product_details,priority:3"`
}

type RivalProdutActiveDetail struct {
	gorm.Model
	Asin            string `json:"asin" gorm:"column:asin;index:idx_product_details,priority:1"`
	Country         string `json:"country" gorm:"column:country;index:idx_product_details,priority:2"`
	Price           string `json:"price" gorm:"column:price"`
	Currency        string `json:"currency" gorm:"column:currency"`
	Coupon          string `json:"coupon" gorm:"column:coupon"`
	Star            string `json:"star" gorm:"column:star"`
	Ratings         uint32 `json:"ratings" gorm:"column:ratings"`
	Image           string `json:"image" gorm:"column:image"`
	ParentAsin      string `json:"parent_asin" gorm:"column:parent_asin"`
	CategoryInfo    string `json:"category_info" gorm:"column:category_info"`
	TopCategoryName string `json:"top_category_name" gorm:"column:top_category_name"`
	TopCategoryRank uint32 `json:"top_category_rank" gorm:"column:top_category_rank"`
	Color           string `json:"color" gorm:"column:color"`
	Weight          string `json:"weight" gorm:"column:weight"`
	WeightUnit      string `json:"weight_unit" gorm:"column:weight_unit"`
	Dimensions      string `json:"dimensions" gorm:"column:dimensions"`
	DimensionsUnit  string `json:"dimensions_unit" gorm:"column:dimensions_unit"`
	CreateDate      string `json:"create_date" gorm:"column:create_date;index:idx_product_details,priority:3"`
}

func MigrateDatabase(db *gorm.DB) error {
	if err := db.AutoMigrate(Rival{}, RivalChange{}, RivalProdutActiveDetail{}, RivalProdutInactiveDetail{}); err != nil {
		return err
	}
	return nil
}
