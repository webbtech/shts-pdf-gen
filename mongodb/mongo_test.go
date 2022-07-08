package mongodb

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/webbtech/shts-pdf-gen/config"
	"github.com/webbtech/shts-pdf-gen/model"
)

const (
	estimateNum = 1191
)

// MongoSuite struct
type MongoSuite struct {
	suite.Suite
	cfg *config.Config
}

func (suite *MongoSuite) SetupTest() {
	os.Setenv("Stage", "test")
	suite.cfg = &config.Config{IsDefaultsLocal: true}
	suite.cfg.DefaultsFilePath = "../config"
	err := suite.cfg.Init()
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func (suite *MongoSuite) TestNewDb() {
	_, err := NewDb(suite.cfg.GetMongoConnectString(), suite.cfg.GetDbName())
	suite.NoError(err)
}

// not exactly a comprehensive test suite... but it covers the bases
func (suite *MongoSuite) TestFetchEstimate() {
	db, _ := NewDb(suite.cfg.GetMongoConnectString(), suite.cfg.GetDbName())
	q, err := db.FetchEstimate(estimateNum)
	suite.NoError(err)
	suite.IsType(&model.Customer{}, q.Customer)
	suite.NotEmpty(q.Customer.FirstName)

	suite.Equal(len(q.ItemIds), len(q.Items))
}

func TestMongoSuite(t *testing.T) {
	suite.Run(t, new(MongoSuite))
}
