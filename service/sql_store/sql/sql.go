package sql

import (
	"fmt"
	"sync"

	"github.com/HooYa-Bigdata/rivalservice/service/sql_store"
	"gorm.io/gorm"
)

type datastore struct {
	db *gorm.DB
}

func (ds *datastore) Rivals() sql_store.RivalStore {
	return &rivals{db: ds.db}
}

func (ds *datastore) RivalChanges() sql_store.RivalChangeStore {
	return &rival_changes{db: ds.db}
}

func (ds *datastore) ProductDetails() sql_store.ProductDetailsStore {
	return &product_details{db: ds.db}
}

var (
	DB         *gorm.DB
	sqlFactory sql_store.SqlFactory
	once       sync.Once
)

func NewSqlStoreFactory(db *gorm.DB) (sql_store.SqlFactory, error) {
	if db == nil && sqlFactory == nil {
		return nil, fmt.Errorf("failed to get sql store fatory")
	}
	once.Do(func() {
		DB = db
		sql_store.MigrateDatabase(db)
		sqlFactory = &datastore{db: db}
	})

	return sqlFactory, nil
}

func (ds *datastore) Close() error {
	db, _ := ds.db.DB()

	if err := db.Close(); err != nil {
		return err
	}
	return nil
}
