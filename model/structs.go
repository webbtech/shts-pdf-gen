package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PdfRequest struct {
	EstimateNumber *int    `json:"number"`
	FileType       *string `json:"type"`
}

type Estimate struct {
	ID primitive.ObjectID `bson:"_id" json:"id"`
}
