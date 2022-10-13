package pdf

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"path"
	"regexp"
	"strings"

	"github.com/jung-kurt/gofpdf"
	"github.com/webbtech/shts-pdf-gen/config"
	"github.com/webbtech/shts-pdf-gen/model"
	"github.com/webbtech/shts-pdf-gen/services"
)

// Pdf struct
type Pdf struct {
	cfg         *config.Config
	outputName  string
	file        *gofpdf.Fpdf
	record      *model.Estimate
	requestType string
	defFontSize float64
	defLnHt     float64
}

const (
	DEFAULT_FONT_SIZE   = 8.5
	DEFAULT_LINE_HEIGHT = 4.5
)

var (
	defFontSize float64
	defLnHt     float64
)

// New function
func New(cfg *config.Config, requestType string, record *model.Estimate) (p *Pdf, err error) {

	p = &Pdf{
		cfg:         cfg,
		defFontSize: DEFAULT_FONT_SIZE,
		defLnHt:     DEFAULT_LINE_HEIGHT,
		record:      record,
		requestType: requestType,
	}

	switch p.requestType {
	case "estimate":
		p.outputName = fmt.Sprintf("%s-%d.pdf", "est", p.record.Number)
		p.file, err = p.Estimate()
	case "invoice":
		p.outputName = fmt.Sprintf("%s-%d.pdf", "inv", p.record.Number)
		p.file, err = p.Invoice()
	}

	return p, err
}

// ========================== Public Methods =============================== //

// OutputToDisk method
func (p *Pdf) OutputToDisk(dir string) (err error) {
	outputPath := path.Join(dir, p.outputName)
	return p.file.OutputFileAndClose(outputPath)
}

// SaveToS3 method
func (p *Pdf) SaveToS3() (err error) {

	fileObject := path.Join(p.requestType, p.outputName)
	var buf bytes.Buffer
	if err = p.file.Output(&buf); err != nil {
		return err
	}

	return services.UploadS3Object(&buf, fileObject, p.cfg.AwsRegion, p.cfg.S3Bucket)
}

// ========================== Private Methods =============================== //

// getLogo method
func (p *Pdf) getLogo(url string) (gofpdf.ImageOptions, bool) {
	var (
		rsp     *http.Response
		tp      string
		imgInfo gofpdf.ImageOptions
	)
	rsp, err := http.Get(url)

	if err == nil {
		tp = p.file.ImageTypeFromMime(rsp.Header["Content-Type"][0])
		if p.file.Err() { // tp produced error because of invalid image so we need to clear and create something that makes a little more sense
			p.file.ClearError()
			p.file.SetError(fmt.Errorf("Invalid or missing filepath: %s", url))
			return imgInfo, false
		}
		imgInfo = gofpdf.ImageOptions{ImageType: tp}
		p.file.RegisterImageReader(url, tp, rsp.Body)
	} else {
		p.file.SetError(err)
		return imgInfo, false
	}
	return imgInfo, true
}

// ========================== Helper Functions ============================= //

func cleanStr(text string) string {

	text = strings.ReplaceAll(text, "’", "'")
	text = strings.ReplaceAll(text, "“", "\"")
	text = strings.ReplaceAll(text, "”", "\"")
	text = strings.ReplaceAll(text, "…", "...")

	// the above should handle most characters,
	// but various unknown characters, possibly cntrl,
	// command, or similar are finding their way in
	// and so we remove them here
	reg, err := regexp.Compile("[^0-9\\w\\s'\"!.,]+")
	if err != nil {
		log.Fatal(err)
	}

	return reg.ReplaceAllString(text, "")
}

func setItemRowHeight(item model.EstimateItem, cellWidth int, compare bool) float64 {

	var numLines = 0.0
	descripLen := len(item.Description)
	longestStr := descripLen

	if compare {
		notesLen := len(item.Notes)
		if notesLen > descripLen {
			longestStr = notesLen
		}
	}
	if longestStr > cellWidth+3 { // we don't need a line break if the string is only marginally longer
		numLines = float64(longestStr / cellWidth)
	}

	return (numLines * defLnHt) + 2
}
