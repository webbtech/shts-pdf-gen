package pdf

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/webbtech/shts-pdf-gen/config"
	"github.com/webbtech/shts-pdf-gen/model"
	"github.com/webbtech/shts-pdf-gen/mongodb"
)

const (
	estimateNum = 1177 // estimate with 11 different items and several unwanted special characters
	// estimateNum = 1191 // estimate with several special characters and description
)

// to test and preview, do something like: go test -run ^TestPdfSuite$ github.com/webbtech/shts-pdf-gen/pdf && open -a Preview ./tmp/est-1177.pdf

// PdfSuite struct
type PdfSuite struct {
	suite.Suite
	cfg            *config.Config
	estimateRecord *model.Estimate
	requestType    string
	db             model.DbHandler
}

// SetupTest method
func (s *PdfSuite) SetupTest() {
	os.Setenv("Stage", "test")

	s.cfg = &config.Config{IsDefaultsLocal: true}
	s.cfg.DefaultsFilePath = "../config"
	err := s.cfg.Init()
	if err != nil {
		log.Fatalf(err.Error())
	}

	s.db, err = mongodb.NewDb(s.cfg.GetMongoConnectString(), s.cfg.GetDbName())
	if err != nil {
		log.Fatalf(err.Error())
	}

	s.estimateRecord, err = s.db.FetchEstimate(estimateNum)
	s.NoError(err)
}

// TestEstimateToDisk method
func (s *PdfSuite) TestEstimateToDisk() {
	s.requestType = "estimate"
	p, err := New(s.cfg, s.requestType, s.estimateRecord)
	s.NoError(err)
	p.OutputToDisk("../tmp")
}

// TestInvoiceToDisk method
func (s *PdfSuite) TestInvoiceToDisk() {
	s.requestType = "invoice"
	p, err := New(s.cfg, s.requestType, s.estimateRecord)
	s.NoError(err)
	p.OutputToDisk("../tmp")
}

// TestEstimateToS3 method
func (s *PdfSuite) TestEstimateToS3() {
	s.requestType = "estimate"
	p, err := New(s.cfg, s.requestType, s.estimateRecord)
	s.NoError(err)

	err = p.SaveToS3()
	s.NoError(err)
}

// TestInvoiceToS3 method
func (s *PdfSuite) TestInvoiceToS3() {
	s.requestType = "invoice"
	p, err := New(s.cfg, s.requestType, s.estimateRecord)
	s.NoError(err)

	err = p.SaveToS3()
	s.NoError(err)
}

// TestPdfSuite method
func TestPdfSuite(t *testing.T) {
	suite.Run(t, new(PdfSuite))
}
