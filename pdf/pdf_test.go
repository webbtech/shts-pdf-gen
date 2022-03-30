package pdf

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/webbtech/shts-pdf-gen/config"
	"github.com/webbtech/shts-pdf-gen/model"
	"github.com/webbtech/shts-pdf-gen/mongo"
)

const (
	defaultsFp  = "../config/defaults.yml"
	estimateNum = 1011
)

// PdfSuite struct
type PdfSuite struct {
	suite.Suite
	estimateRecord *model.Estimate
	requestType    string
	db             model.DbHandler
}

// SetupTest method
func (s *PdfSuite) SetupTest() {
	cfg := &config.Config{DefaultsFilePath: defaultsFp}
	err := cfg.Init()
	s.NoError(err)

	s.db, err = mongo.NewDb(cfg.GetMongoConnectString(), cfg.GetDbName())
	s.NoError(err)

	s.estimateRecord, err = s.db.FetchEstimate(estimateNum)
	s.NoError(err)
}

// TestInit method
func (s *PdfSuite) TestNew() {
	s.requestType = "estimate"
	p, err := New(s.requestType, s.estimateRecord)
	s.NoError(err)
	// fmt.Printf("p: %+v\n", p)
	p.OutputToDisk("../tmp")
	// suite.IsType(&model.Customer{}, q.)
}

// TestPdfSuite method
func TestPdfSuite(t *testing.T) {
	suite.Run(t, new(PdfSuite))
}
