package pdf

import (
	"fmt"

	"github.com/jung-kurt/gofpdf"
	log "github.com/sirupsen/logrus"
	"github.com/webbtech/shts-pdf-gen/model"
)

type estimate struct {
	file   *gofpdf.Fpdf
	record *model.Estimate
	p      *Pdf
}

var (
	defFontSize float64
	defLnHt     float64
)

const (
	PAYMENT_DUE_NOTE = "Payment is due upon completion unless otherwise established."
	ESTIMATE_VALID   = "Estimate is valid until %s"
	DATE_FMT_SHRT    = "Jan 2, 2006"
)

// Estimate method
func (p *Pdf) Estimate() (file *gofpdf.Fpdf, err error) {

	file = gofpdf.New("P", "mm", "Letter", "")
	p.file = file
	e := &estimate{p: p}
	defFontSize = p.defFontSize
	defLnHt = p.defLnHt

	titleStr := fmt.Sprintf("Quote %d PDF", p.record.Number)

	file.SetTitle(titleStr, false)
	file.SetAuthor("Shorthills Tree Service", false)

	file.SetFooterFunc(func() {
		file.SetY(-15)
		file.SetFont("Arial", "I", defFontSize)
		file.CellFormat(0, 10, fmt.Sprintf("Page %d of {nb}", file.PageNo()), "", 0, "C", false, 0, "")
	})
	file.AliasNbPages("")

	file.AddPage()
	e.header()
	e.items()
	e.footer()

	return file, err
}

func (e *estimate) header() {

	record := e.p.record
	file := e.p.file
	companyInfo := e.p.cfg.GetCompanyInfo()
	protocol := "http"
	customer := record.Customer

	customerNm := fmt.Sprintf("%s %s", customer.FirstName, customer.LastName)
	addressLine1 := fmt.Sprintf("%s", customer.Street1)
	addressLine2 := fmt.Sprintf("%s, %s. %s", customer.City, customer.Province, customer.PostalCode)
	estimateNo := fmt.Sprintf("%d", record.Number)

	logoInfo, ok := e.p.GetLogo(companyInfo.LogoURI)
	if !ok {
		log.Errorf("Error with GetLogo and url: %s", companyInfo.LogoURI)
		return
	}

	file.SetUnderlineThickness(1)

	file.ImageOptions(companyInfo.LogoURI, 9, 11, 40, 40, false, logoInfo, 0, fmt.Sprintf("%s://%s", protocol, companyInfo.Domain))
	file.MoveTo(22, 52)
	file.SetTextColor(0, 0, 0)
	file.SetFont("Arial", "", 8)
	file.CellFormat(0, 3.5, "Est. 1999", "C", 2, "", false, 0, "")

	file.MoveTo(6, 12)
	file.SetFont("Times", "B", 26)
	file.SetTextColor(60, 98, 61)
	file.CellFormat(0, 10, "Shorthills Tree Service", "", 2, "CT", false, 0, "")

	file.SetFont("Arial", "B", defFontSize)
	file.SetTextColor(0, 0, 0)
	file.CellFormat(0, defLnHt, companyInfo.Phone, "", 2, "C", false, 0, "")

	file.SetFont("Arial", "U", defFontSize)
	file.CellFormat(0, defLnHt, companyInfo.Domain, "", 2, "C", false, 0, fmt.Sprintf("%s://%s", protocol, companyInfo.Domain))
	file.CellFormat(0, defLnHt, companyInfo.Email, "", 2, "C", false, 0, fmt.Sprintf("mailto:%s", companyInfo.Email))

	// credentials
	file.MoveTo(84, 38)
	file.SetTextColor(0, 0, 0)
	file.SetFont("Arial", "", 7.5)
	file.CellFormat(0, 3.5, "Certified ISA Arborist ON1111-A", "", 2, "", false, 0, "")
	file.CellFormat(0, 3.5, "Certified Tree Risk Assessor #1859", "", 2, "", false, 0, "")
	file.CellFormat(0, 3.5, "Certified Utility Arborist #400145204", "", 2, "", false, 0, "")
	file.CellFormat(0, 3.5, "Certified Ontario Arborist #400143375", "", 2, "", false, 0, "")

	file.MoveTo(162, 14)
	file.SetFont("Arial", "B", 12)
	file.SetTextColor(0, 0, 0)
	file.CellFormat(24, defLnHt, "ESTIMATE", "", 0, "T", false, 0, "")
	file.SetFont("Arial", "", 12)
	file.SetTextColor(200, 0, 0)
	file.CellFormat(0, defLnHt, fmt.Sprintf("%s", estimateNo), "", 2, "T", false, 0, "")

	file.MoveTo(162, 19)
	file.SetFont("Arial", "", defFontSize)
	file.SetTextColor(0, 0, 0)
	file.CellFormat(0, 5, "Mar 15th, 2022", "", 2, "", false, 0, "")

	file.MoveTo(162, 28)

	file.SetFont("Arial", "", defFontSize)
	file.CellFormat(0, defLnHt, customerNm, "", 2, "", false, 0, "")
	file.CellFormat(0, defLnHt, addressLine1, "", 2, "", false, 0, "")
	file.CellFormat(0, defLnHt, addressLine2, "", 2, "", false, 0, "")
	file.CellFormat(0, defLnHt, customer.Phone, "", 2, "", false, 0, "")
	if customer.Email != "" {
		file.SetFont("Arial", "U", defFontSize)
		file.CellFormat(0, defLnHt, customer.Email, "", 2, "", false, 0, fmt.Sprintf("mailto:%s", customer.Email))
	}

	file.MoveTo(152, 28)
	file.SetFont("Arial", "B", defFontSize)
	file.CellFormat(0, defLnHt, "For:", "", 2, "M", false, 0, "")

	// Items heading
	file.SetFillColor(60, 98, 61)
	file.Rect(10, 58, 196, 1, "F")
	file.SetFont("Arial", "", 7.5)
	file.SetFillColor(200, 200, 200)
	file.SetTextColor(0, 0, 0)

	file.RoundedRect(10, 64, 196, 6, 0.75, "1234", "F")
	file.MoveTo(11, 65.25)
	file.CellFormat(69, 4, "ITEM", "", 0, "", false, 0, "")
	file.CellFormat(86, 4, "DESCRIPTION", "", 0, "", false, 0, "")
	file.CellFormat(9, 4, "QTY", "", 0, "", false, 0, "")
	file.CellFormat(13, 4, "UNIT", "", 0, "", false, 0, "")
	file.CellFormat(20, 4, "EXTENDED", "", 0, "", false, 0, "")

	file.Ln(7)
}

func (e *estimate) items() {

	var xPos, yPos float64
	record := e.p.record
	file := e.p.file
	descripW := 70.0
	notesW := 85.0

	// this should in practice never happen, but still need to protect against
	if len(record.Items) <= 0 {
		return
	}

	// Test page break for multi-page estimate
	/* record.Items = append(record.Items, record.Items[0])
	record.Items = append(record.Items, record.Items[1])
	record.Items = append(record.Items, record.Items[0])
	record.Items = append(record.Items, record.Items[1])
	record.Items = append(record.Items, record.Items[0])
	record.Items = append(record.Items, record.Items[1])
	record.Items = append(record.Items, record.Items[0])
	record.Items = append(record.Items, record.Items[1])
	record.Items = append(record.Items, record.Items[0])
	record.Items = append(record.Items, record.Items[1]) */

	file.SetFont("Arial", "", defFontSize)
	file.SetTextColor(0, 0, 0)
	file.SetFillColor(100, 100, 100)

	for idx, i := range record.Items {
		file.Ln(3)
		xPos = file.GetX()
		yPos = file.GetY()
		file.MultiCell(descripW, defLnHt, i.Description, "", "T", false)
		file.MoveTo(xPos+descripW, yPos)

		xPos = file.GetX()
		yPos = file.GetY()
		file.MultiCell(notesW, defLnHt, i.Notes, "", "T", false)
		file.MoveTo(xPos+notesW, yPos)
		file.CellFormat(8, defLnHt, fmt.Sprintf("%d", i.Qty), "", 0, "TC", false, 0, "")
		file.CellFormat(15, defLnHt, fmt.Sprintf("%.2f", i.Cost), "", 0, "TR", false, 0, "")
		file.SetFont("Arial", "B", defFontSize)
		file.CellFormat(15, defLnHt, fmt.Sprintf("%.2f", i.Total), "", 0, "TR", false, 0, "")

		file.SetDrawColor(200, 200, 200)
		file.SetFont("Arial", "", defFontSize)
		file.Ln(14)
		file.CellFormat(0, 6, "", "B", 1, "", false, 0, "")

		if (idx+1)%8 == 0 {
			file.AddPage()
			file.CellFormat(0, 6, "", "B", 1, "", false, 0, "")
		}
	}
}

func (e *estimate) footer() {

	var yPos float64
	record := e.p.record
	file := e.p.file
	totalsCellH := defLnHt
	totalsCellW := 22.0
	totalsCellX := 160.0

	expireDate := record.Date.AddDate(0, 0, 90)

	file.Ln(6)
	yPos = file.GetY() // Save this for totals section

	file.SetFont("Arial", "B", defFontSize)
	file.CellFormat(0, 6, "JOB NOTES", "", 1, "T", false, 0, "")
	file.SetFont("Arial", "", defFontSize)
	file.MultiCell(120, defLnHt, record.JobNotes, "", "T", false)

	file.CellFormat(0, 6, "", "B", 1, "", false, 0, "")
	file.Ln(4)
	file.SetFont("Arial", "B", defFontSize)
	file.CellFormat(0, 6, "FINANCE NOTES", "", 1, "T", false, 0, "")

	file.SetFont("Arial", "", defFontSize)
	file.CellFormat(0, defLnHt, "Incorporation: 2187989 Ontario Inc.", "", 1, "", false, 0, "")
	file.CellFormat(0, defLnHt, fmt.Sprintf("HST #: %s", e.p.cfg.GetCompanyInfo().HST), "", 1, "", false, 0, "")
	file.CellFormat(0, defLnHt, PAYMENT_DUE_NOTE, "", 1, "", false, 0, "")
	file.CellFormat(0, defLnHt, fmt.Sprintf(ESTIMATE_VALID, expireDate.Format(DATE_FMT_SHRT)), "", 1, "", false, 0, "")

	// totals section
	file.MoveTo(totalsCellX, yPos)

	file.SetFont("Arial", "", defFontSize)
	file.CellFormat(totalsCellW, totalsCellH, "Item Cost", "", 0, "MR", false, 0, "")
	file.SetFont("Arial", "B", defFontSize)
	file.CellFormat(totalsCellW, totalsCellH, fmt.Sprintf("$%.2f", record.ItemsCost), "", 0, "MR", false, 0, "")

	file.MoveTo(totalsCellX, yPos+totalsCellH)

	file.SetFont("Arial", "", defFontSize)
	file.CellFormat(totalsCellW, totalsCellH, "Discount", "", 0, "MR", false, 0, "")
	file.SetFont("Arial", "B", defFontSize)
	file.CellFormat(totalsCellW, totalsCellH, fmt.Sprintf("-$%.2f", record.Discount), "", 0, "MR", false, 0, "")

	file.MoveTo(totalsCellX, yPos+(2*totalsCellH))

	file.SetFont("Arial", "", defFontSize)
	file.CellFormat(totalsCellW, totalsCellH, "Subtotal", "", 0, "MR", false, 0, "")
	file.SetFont("Arial", "B", defFontSize)
	file.CellFormat(totalsCellW, totalsCellH, fmt.Sprintf("$%.2f", record.ItemsCostNet), "", 0, "MR", false, 0, "")

	file.MoveTo(totalsCellX, yPos+(3*totalsCellH))

	file.SetFont("Arial", "", defFontSize)
	file.CellFormat(totalsCellW, totalsCellH, fmt.Sprintf("HST (%d)%%", record.HST), "", 0, "MR", false, 0, "")
	file.SetFont("Arial", "B", defFontSize)
	file.CellFormat(totalsCellW, totalsCellH, fmt.Sprintf("$%.2f", record.Tax), "", 0, "MR", false, 0, "")

	file.MoveTo(totalsCellX, yPos+2+(4*totalsCellH))

	file.SetFont("Arial", "B", 10)
	file.CellFormat(totalsCellW, totalsCellH+1, "Total", "", 0, "MR", false, 0, "")
	file.CellFormat(totalsCellW, totalsCellH+1, fmt.Sprintf("$%.2f", record.TotalCost), "", 0, "MR", false, 0, "")

	file.MoveTo(totalsCellX, yPos+3+(5*totalsCellH))

}
