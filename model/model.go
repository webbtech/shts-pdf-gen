package model

type PdfRequest struct {
	EstimateNumber *int    `json:"number"`
	FileType       *string `json:"type"`
}
