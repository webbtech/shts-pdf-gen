package pdf

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/webbtech/shts-pdf-gen/config"
	"github.com/webbtech/shts-pdf-gen/model"
	"github.com/webbtech/shts-pdf-gen/mongodb"
)

const (
	defaultsFp  = "../config/defaults.yml"
	estimateNum = 1011
)

// to test and preview, do something like: go test -run ^TestPdfSuite$ github.com/webbtech/shts-pdf-gen/pdf && open -a Preview ./tmp/estimate-1011.pdf

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

	s.cfg = &config.Config{DefaultsFilePath: defaultsFp}
	err := s.cfg.Init()
	s.NoError(err)

	s.db, err = mongodb.NewDb(s.cfg.GetMongoConnectString(), s.cfg.GetDbName())
	s.NoError(err)

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
	// suite.IsType(&model.Customer{}, q.)
}

// TestEstimateToS3 method
func (s *PdfSuite) TestEstimateToS3() {
	s.requestType = "estimate"
	p, err := New(s.cfg, s.requestType, s.estimateRecord)
	s.NoError(err)

	l, err := p.SaveToS3()
	s.NoError(err)

	expectLocation := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s/%s", s.cfg.S3Bucket, s.cfg.AwsRegion, s.requestType, p.outputName)
	s.Equal(expectLocation, l)
}

// TestInvoiceToS3 method
func (s *PdfSuite) TestInvoiceToS3() {
	s.requestType = "invoice"
	p, err := New(s.cfg, s.requestType, s.estimateRecord)
	s.NoError(err)

	l, err := p.SaveToS3()
	s.NoError(err)

	expectLocation := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s/%s", s.cfg.S3Bucket, s.cfg.AwsRegion, s.requestType, p.outputName)
	s.Equal(expectLocation, l)
}

// TestPdfSuite method
func TestPdfSuite(t *testing.T) {
	suite.Run(t, new(PdfSuite))
}
