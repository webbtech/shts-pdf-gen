package pdf

import (
	"fmt"
	"time"

	"github.com/jung-kurt/gofpdf"
	log "github.com/sirupsen/logrus"
	"github.com/webbtech/shts-pdf-gen/model"
)

type invoice struct {
	file   *gofpdf.Fpdf
	record *model.Estimate
	p      *Pdf
}

// Estimate method
func (p *Pdf) Invoice() (file *gofpdf.Fpdf, err error) {

	file = gofpdf.New("P", "mm", "Letter", "")
	p.file = file
	i := &invoice{p: p}

	defFontSize = p.defFontSize
	defLnHt = p.defLnHt
	titleStr := fmt.Sprintf("Quote %d PDF", p.record.Number)

	file.SetDrawColor(200, 200, 200)
	file.SetLineWidth(.35)
	file.SetUnderlineThickness(1)
	file.SetTitle(titleStr, false)
	file.SetAuthor(p.cfg.GetCompanyInfo().Name, false)

	file.SetFooterFunc(func() {
		file.SetTextColor(100, 100, 100)
		file.SetY(-15)
		file.SetFont("Arial", "I", defFontSize)
		file.CellFormat(0, 10, fmt.Sprintf("Page %d of {nb}", file.PageNo()), "", 0, "C", false, 0, "")
	})
	file.AliasNbPages("")

	file.AddPage()
	i.header()
	i.items()
	i.totals()
	i.footer()

	return file, err
}

func (i *invoice) header() {

	file := i.p.file
	companyInfo := i.p.cfg.GetCompanyInfo()
	protocol := "http"

	record := i.p.record
	customer := record.Customer

	customerNm := fmt.Sprintf("%s %s", customer.FirstName, customer.LastName)
	addressLine1 := fmt.Sprintf("%s", customer.Street1)
	addressLine2 := fmt.Sprintf("%s, %s. %s", customer.City, customer.Province, customer.PostalCode)
	invoiceID := fmt.Sprintf("%d", record.Number)
	issueDate := time.Now()

	logoInfo, ok := i.p.getLogo(companyInfo.LogoURI)
	if !ok {
		log.Errorf("Error with GetLogo and url: %s", companyInfo.LogoURI)
		return
	}

	file.ImageOptions(companyInfo.LogoURI, 10, 10, 30, 30, false, logoInfo, 0, fmt.Sprintf("%s://%s", protocol, companyInfo.Domain))

	file.MoveTo(19, 40)
	file.SetTextColor(0, 0, 0)
	file.SetFont("Arial", "", 7)
	file.SetTextColor(100, 100, 100)
	file.CellFormat(0, 3.5, "Est. 1999", "C", 2, "", false, 0, "")

	file.MoveTo(-24, 12)
	file.SetFont("Times", "B", 24)
	file.SetTextColor(60, 98, 61)
	file.CellFormat(0, 10, companyInfo.Name, "", 2, "C", false, 0, "")

	file.SetTextColor(0, 0, 0)
	file.SetFont("Arial", "", defFontSize)
	file.CellFormat(0, defLnHt, "2187989 Ontario Inc.", "", 2, "C", false, 0, "")

	file.MoveTo(150, 16)
	file.SetFont("Arial", "B", 14)
	file.SetTextColor(0, 0, 0)
	file.CellFormat(24, defLnHt, "INVOICE", "", 0, "", false, 0, "")

	file.MoveTo(136, 25)
	file.SetFont("Arial", "", defFontSize)
	file.SetTextColor(100, 100, 100)
	file.CellFormat(24, defLnHt, "From", "", 0, "", false, 0, "")

	file.Line(148, 25, 148, 48)

	file.MoveTo(150, 25)
	file.SetFont("Arial", "B", defFontSize)
	file.SetTextColor(0, 0, 0)
	file.CellFormat(0, defLnHt, companyInfo.Name, "", 2, "", false, 0, "")
	file.SetFont("Arial", "", defFontSize)
	file.CellFormat(0, defLnHt, companyInfo.Address1, "", 2, "", false, 0, "")
	file.CellFormat(0, defLnHt, companyInfo.Address2, "", 2, "", false, 0, "")
	file.SetFont("Arial", "U", defFontSize)
	file.CellFormat(0, defLnHt, fmt.Sprintf("email: %s", companyInfo.Email), "", 2, "", false, 0, fmt.Sprintf("mailto:%s", companyInfo.Email))
	file.SetFont("Arial", "", defFontSize)
	file.CellFormat(0, defLnHt, fmt.Sprintf("phone: %s", companyInfo.Phone), "", 1, "", false, 0, "")

	file.MoveTo(14, 56)
	file.SetTextColor(100, 100, 100)
	file.CellFormat(0, defLnHt, "Invoice For", "", 0, "", false, 0, "")

	file.Line(36, 56, 36, 70)

	file.MoveTo(38, 56)
	file.SetTextColor(0, 0, 0)
	file.SetFont("Arial", "B", defFontSize)
	file.CellFormat(0, defLnHt, customerNm, "", 2, "", false, 0, "")
	file.SetFont("Arial", "", defFontSize)
	file.CellFormat(0, defLnHt, addressLine1, "", 2, "", false, 0, "")
	file.CellFormat(0, defLnHt, addressLine2, "", 2, "", false, 0, "")

	file.MoveTo(128, 56)
	file.SetTextColor(100, 100, 100)
	file.CellFormat(0, 5, "Invoice ID", "", 2, "T", false, 0, "")
	file.CellFormat(0, 7, "Issue Date", "", 2, "", false, 0, "")
	file.CellFormat(0, 7, "Due Date", "", 2, "", false, 0, "")

	file.Line(148, 56, 148, 74)

	file.MoveTo(150, 56)
	file.SetTextColor(0, 0, 0)
	file.SetFont("Arial", "B", defFontSize)
	file.CellFormat(0, 5, invoiceID, "", 2, "T", false, 0, "")
	file.SetFont("Arial", "", defFontSize)
	file.CellFormat(0, 7, issueDate.Format(DATE_FMT_SHRT), "", 2, "", false, 0, "")
	file.CellFormat(0, 7, issueDate.Format(DATE_FMT_SHRT), "", 2, "", false, 0, "")
}

func (i *invoice) items() {

	record := i.p.record
	file := i.p.file
	descripW := 148.0

	// Green separator
	file.SetFillColor(60, 98, 61)
	file.Rect(10, 78, 196, 1, "F")
	file.SetFont("Arial", "", 7.5)
	file.SetFillColor(200, 200, 200)
	file.SetTextColor(0, 0, 0)

	// Items heading
	file.SetFont("Arial", "B", 7.0)
	file.SetFillColor(220, 220, 220)
	file.SetTextColor(0, 0, 0)

	file.RoundedRect(10, 82, 196, 6, 0.75, "1234", "F")
	file.MoveTo(11, 83.25)
	file.CellFormat(146, 4, "ITEM", "", 0, "", false, 0, "")
	file.CellFormat(6, 4, "QTY", "", 0, "R", false, 0, "")
	file.CellFormat(19, 4, "UNIT", "", 0, "R", false, 0, "")
	file.CellFormat(22, 4, "EXTENDED", "", 0, "R", false, 0, "")

	file.Ln(6)
	file.SetFont("Arial", "", defFontSize)

	for idx, i := range record.Items {
		file.Ln(2.5)
		file.CellFormat(descripW, defLnHt, i.Description, "", 0, "", false, 0, "")
		file.CellFormat(6, defLnHt, fmt.Sprintf("%d", i.Qty), "", 0, "C", false, 0, "")
		file.CellFormat(18, defLnHt, fmt.Sprintf("%.2f", i.Cost), "", 0, "R", false, 0, "")
		file.SetFont("Arial", "B", defFontSize)
		file.CellFormat(22, defLnHt, fmt.Sprintf("%.2f", i.Total), "", 1, "R", false, 0, "")
		file.SetFont("Arial", "", defFontSize)
		file.CellFormat(0, 2.5, "", "B", 1, "", false, 0, "")

		if (idx+1)%16 == 0 { // TODO: test this
			file.AddPage()
			file.CellFormat(0, 6, "", "B", 1, "", false, 0, "")
		}
	}
}

func (i *invoice) totals() {

	var yPos float64

	record := i.p.record
	file := i.p.file
	totalsCellH := 5.0
	totalsCellW := 22.0
	totalsCellX := 160.0

	file.Ln(4)
	yPos = file.GetY()
	file.MoveTo(totalsCellX, yPos)

	file.SetFont("Arial", "", defFontSize)
	file.SetTextColor(100, 100, 100)
	file.CellFormat(totalsCellW, totalsCellH, "Item Cost", "", 0, "MR", false, 0, "")
	file.SetFont("Arial", "B", defFontSize)
	file.SetTextColor(0, 0, 0)
	file.CellFormat(totalsCellW, totalsCellH, fmt.Sprintf("%.2f", record.ItemsCost), "", 0, "R", false, 0, "")

	file.MoveTo(totalsCellX, yPos+totalsCellH)

	file.SetFont("Arial", "", defFontSize)
	file.SetTextColor(100, 100, 100)
	file.CellFormat(totalsCellW, totalsCellH, "Discount", "", 0, "MR", false, 0, "")
	file.SetFont("Arial", "B", defFontSize)
	file.SetTextColor(0, 0, 0)
	file.CellFormat(totalsCellW, totalsCellH, fmt.Sprintf("-%.2f", record.Discount), "", 0, "MR", false, 0, "")

	file.MoveTo(totalsCellX, yPos+(2*totalsCellH))

	file.SetFont("Arial", "", defFontSize)
	file.SetTextColor(100, 100, 100)
	file.CellFormat(totalsCellW, totalsCellH, "Subtotal", "", 0, "MR", false, 0, "")
	file.SetFont("Arial", "B", defFontSize)
	file.SetTextColor(0, 0, 0)
	file.CellFormat(totalsCellW, totalsCellH, fmt.Sprintf("%.2f", record.ItemsCostNet), "", 0, "MR", false, 0, "")

	file.MoveTo(totalsCellX, yPos+(3*totalsCellH))

	file.SetFont("Arial", "", defFontSize)
	file.SetTextColor(100, 100, 100)
	file.CellFormat(totalsCellW, totalsCellH, fmt.Sprintf("HST (%d%%)", record.HST), "", 0, "MR", false, 0, "")
	file.SetFont("Arial", "B", defFontSize)
	file.SetTextColor(0, 0, 0)
	file.CellFormat(totalsCellW, totalsCellH, fmt.Sprintf("%.2f", record.Tax), "", 0, "MR", false, 0, "")

	file.MoveTo(totalsCellX, yPos+2+(4*totalsCellH))

	file.SetFont("Arial", "B", 10)
	file.CellFormat(totalsCellW, totalsCellH+1, "Amount Due", "", 0, "MR", false, 0, "")
	file.CellFormat(totalsCellW, totalsCellH+1, fmt.Sprintf("$%.2f", record.TotalCost), "", 0, "MR", false, 0, "")
}

func (i *invoice) footer() {

	companyInfo := i.p.cfg.GetCompanyInfo()
	file := i.p.file

	file.Ln(16)
	file.CellFormat(0, 2, "", "T", 1, "", false, 0, "")

	file.SetFont("Arial", "B", defFontSize)
	file.CellFormat(0, defLnHt, "NOTES", "", 1, "", false, 0, "")
	file.Ln(2)

	file.SetFont("Arial", "", defFontSize)
	file.CellFormat(0, defLnHt, fmt.Sprintf("Payable to: %s", companyInfo.Name), "", 1, "", false, 0, "")
	file.CellFormat(0, defLnHt, fmt.Sprintf("HST #: %s", companyInfo.HST), "", 1, "", false, 0, "")
	file.CellFormat(0, defLnHt, fmt.Sprintf("e-Transfer to: %s (auto deposit)", companyInfo.Email), "", 1, "", false, 0, "")
}
