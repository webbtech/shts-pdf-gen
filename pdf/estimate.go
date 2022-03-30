package pdf

import (
	"fmt"

	"github.com/jung-kurt/gofpdf"
	"github.com/webbtech/shts-pdf-gen/model"
)

type estimate struct {
	file   *gofpdf.Fpdf
	record *model.Estimate
}

// Estimate method
func (p *Pdf) Estimate() (file *gofpdf.Fpdf, err error) {

	file = gofpdf.New("P", "mm", "Letter", "")
	e := &estimate{file: file, record: p.record}

	titleStr := fmt.Sprintf("Quote %d PDF", p.record.Number)

	file.SetTitle(titleStr, false)
	file.SetAuthor("Shorthills Tree Service", false)

	file.SetFooterFunc(func() {
		file.SetY(-15)
		file.SetFont("Arial", "I", 9)
		file.CellFormat(0, 10, fmt.Sprintf("Page %d of {nb}", file.PageNo()), "", 0, "C", false, 0, "")
	})
	file.AliasNbPages("")

	file.AddPage()

	e.estimateTitle()

	return file, err
}

func (e *estimate) estimateTitle() {

	customerNm := fmt.Sprintf("%s %s", e.record.Customer.FirstName, e.record.Customer.LastName)
	e.file.Ln(1)
	e.file.SetFont("Arial", "", 12)
	e.file.CellFormat(0, 5.5, customerNm, "", 2, "", false, 0, "")
	e.file.Ln(1)
	// pdf.CellFormat(0, 5.5, addressLine1, "", 2, "", false, 0, "")
	// pdf.CellFormat(0, 4, addressLine2, "", 2, "", false, 0, "")
}
