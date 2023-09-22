package sql

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	v1 "github.com/superjcd/rivalservice/genproto/v1"
	"github.com/superjcd/rivalservice/service/sql_store"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var dbFile = "fake.db"

type FakeStoreTestSuite struct {
	suite.Suite
	Dbfile      string
	FakeFactory sql_store.SqlFactory
}

func (suite *FakeStoreTestSuite) SetupSuite() {
	file, err := os.Create(dbFile)
	assert.Nil(suite.T(), err)
	defer file.Close()

	suite.Dbfile = dbFile

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:        time.Second, // Slow SQL threshold
			LogLevel:             logger.Warn, // Log level
			ParameterizedQueries: true,        // Don't include params in the SQL log
			Colorful:             true,        // Disable color
		},
	)

	db, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{Logger: newLogger})
	assert.Nil(suite.T(), err)

	factory, err := NewSqlStoreFactory(db)
	assert.Nil(suite.T(), err)
	suite.FakeFactory = factory
}

func (suite *FakeStoreTestSuite) TearDownSuite() {
	var err error
	err = suite.FakeFactory.Close()
	assert.Nil(suite.T(), err)
	err = os.Remove(dbFile)
	assert.Nil(suite.T(), err)
}

func (suite *FakeStoreTestSuite) TestCreateRival() {
	rq := &v1.CreateRivalRequest{
		User:    "Jack",
		Asin:    "B001",
		Country: "US",
		Rivals:  []string{"B1001", "B1002"},
	}

	err := suite.FakeFactory.Rivals().Create(context.Background(), rq)
	assert.Nil(suite.T(), err)
}

func (suite *FakeStoreTestSuite) TestListRival() {
	rq := &v1.ListRivalRequest{
		User:    "Jack",
		Country: "US",
		Asin:    "B001",
		Offset:  0,
		Limit:   10,
	}
	rivals, err := suite.FakeFactory.Rivals().List(context.Background(), rq)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 2, int(rivals.TotalCount))
	assert.Equal(suite.T(), 2, len(rivals.Items))
}

func (suite *FakeStoreTestSuite) TestAppendDetail() {
	// 导入active
	rq := &v1.AppendRivalProductActiveDetailRequest{
		Details: []*v1.AmzProductActiveDetail{
			{
				Asin:       "B1001",
				Country:    "US",
				Price:      "100",
				CreateDate: "2022-01-01",
			},
		},
	}

	err := suite.FakeFactory.ProductDetails().AppendActiveDetail(context.Background(), rq)
	assert.Nil(suite.T(), err)

	rq2 := &v1.AppendRivalProductInactiveDetailRequest{
		Details: []*v1.AmzProductInactivateDetail{
			{
				Asin:         "B1001",
				Country:      "US",
				Title:        "Iphone 15",
				BulletPoints: "1 good 2 cheap",
				CreateDate:   "2022-01-01",
			},
		},
	}

	err2 := suite.FakeFactory.ProductDetails().AppendInactiveDetail(context.Background(), rq2)
	assert.Nil(suite.T(), err2)

}

func (suite *FakeStoreTestSuite) TestAppendRivalChanges() {
	rq := &v1.AppendRivalChangesRequest{
		OldDate: "2022-01-01",
		NewDate: "2022-01-02",
		Field:   "price",
	}

	err := suite.FakeFactory.RivalChanges().Append(context.Background(), rq)
	assert.Nil(suite.T(), err)
}

func (suite *FakeStoreTestSuite) TestListUserRivalChanges() {
	rq := &v1.ListRivalChangesRequest{
		User:       "Jack",
		Country:    "US",
		CreateDate: "2022-01-02",
		Field:      "price",
		Offset:     0,
		Limit:      10,
	}
	result, err := suite.FakeFactory.RivalChanges().List(context.Background(), rq)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 2, len(result.Items))

}

func (suite *FakeStoreTestSuite) TestXDeleteRivals() {
	rq := &v1.DeleteRivalRequest{
		User:    "Jack",
		Country: "US",
	}

	err := suite.FakeFactory.Rivals().Delete(context.Background(), rq)
	assert.Nil(suite.T(), err)

	rq2 := &v1.ListRivalRequest{
		User:    "Jack",
		Country: "US",
		Asin:    "B001",
		Offset:  0,
		Limit:   10,
	}
	rivals, err := suite.FakeFactory.Rivals().List(context.Background(), rq2)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 0, int(rivals.TotalCount))
}

func (suite *FakeStoreTestSuite) TestYDeleteRivalChanges() {
	rq := &v1.DeleteRivalChangesRequest{
		MinCreateDate: "2023-10-01",
	}

	err := suite.FakeFactory.RivalChanges().Delete(context.Background(), rq)
	assert.Nil(suite.T(), err)

}

func (suite *FakeStoreTestSuite) TestZDeleteActiveProductDetail() {
	rq := &v1.DeleteRivalActiveDetailRequest{
		MinCreateDate: "2023-10-01",
	}

	err := suite.FakeFactory.ProductDetails().DeleteActiveDetail(context.Background(), rq)
	assert.Nil(suite.T(), err)
}

func (suite *FakeStoreTestSuite) TestZDeleteInactiveProductDetail() {
	rq := &v1.DeleteRivalInactiveDetailRequest{
		MinCreateDate: "2023-10-01",
	}

	err := suite.FakeFactory.ProductDetails().DeleteInactiveDetail(context.Background(), rq)
	assert.Nil(suite.T(), err)
}

func TestFakeStoreSuite(t *testing.T) {
	suite.Run(t, new(FakeStoreTestSuite))
}
