package pdf

import (
	"fmt"
	"path"

	"github.com/jung-kurt/gofpdf"
	"github.com/webbtech/shts-pdf-gen/model"
)

// Pdf struct
type Pdf struct {
	outputName  string
	file        *gofpdf.Fpdf
	record      *model.Estimate
	requestType string
}

// New function
func New(requestType string, record *model.Estimate) (p *Pdf, err error) {

	p = &Pdf{
		record:      record,
		requestType: requestType,
	}

	switch p.requestType {
	case "estimate":
		p.file, err = p.Estimate()
	case "invoice":
		p.file, err = p.Invoice()
	}

	p.outputName = fmt.Sprintf("%s-%d.pdf", p.requestType, p.record.Number)

	return p, err
}

// ========================== Public Methods =============================== //

// OutputToDisk method
func (p *Pdf) OutputToDisk(dir string) (err error) {
	outputPath := path.Join(dir, p.outputName)
	err = p.file.OutputFileAndClose(outputPath)
	return err
}

// SaveToS3 method

// ========================== Private Methods =============================== //
