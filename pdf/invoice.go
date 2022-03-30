package pdf

import "github.com/jung-kurt/gofpdf"

// Estimate method
func (p *Pdf) Invoice() (file *gofpdf.Fpdf, err error) {

	file = gofpdf.New("P", "mm", "Letter", "")

	return file, err
}
