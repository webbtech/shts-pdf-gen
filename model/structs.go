package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PdfRequest struct {
	EstimateNumber *int    `bson:"number" json:"number" validate:"required"`
	RequestType    *string `bson:"requestType" json:"requestType" validate:"required"`
}

type Customer struct {
	ID           primitive.ObjectID `bson:"_id" json:"id"`
	Active       bool               `bson:"active" json:"active"`
	City         string             `bson:"addrCity" json:"city"`
	Location     []float64          `bson:"addrLocation" json:"location"`
	PostalCode   string             `bson:"addrPostalCode" json:"postalCode"`
	Province     string             `bson:"addrProvinceCode" json:"province"`
	Street1      string             `bson:"addrStreet1" json:"street1"`
	Email        string             `bson:"email" json:"email"`
	FirstName    string             `bson:"nameFirst" json:"firstName"`
	LastName     string             `bson:"nameLast" json:"lastName"`
	Notes        string             `bson:"notes" json:"notes"`
	Phone        string             `bson:"phone" json:"phone"`
	PhoneMsgOK   bool               `bson:"phoneMsgOK" json:"phoneMsgOK"`
	PreferredCom string             `bson:"preferredCommunication" json:"preferredCom"`
}

type Estimate struct {
	ID            primitive.ObjectID   `bson:"_id" json:"id"`
	CreatedAt     time.Time            `bson:"createdAt" json:"createdAt"`
	Customer      *Customer            `bson:"customerRecord" json:"customerRecord"`
	CustomerId    primitive.ObjectID   `bson:"customer" json:"customerId"`
	CustomerNotes string               `bson:"customerNotes" json:"customerNotes"`
	Discount      float64              `bson:"discount" json:"discount"`
	HST           int                  `bson:"HST" json:"HST"`
	ItemIds       []primitive.ObjectID `bson:"estimateItems" json:"itemIds"`
	Items         []EstimateItem       `json:"items"`
	ItemsCost     float64              `bson:"itemsCost" json:"itemsCost"`
	ItemsCostNet  float64              `bson:"itemsCostNet" json:"itemsCostNet"`
	JobNotes      string               `bson:"jobNotes" json:"jobNotes"`
	Number        int64                `bson:"number" json:"number"`
	PaidStatus    string               `bson:"paidStatus" json:"paidStatus"`
	Tax           float64              `bson:"tax" json:"tax"`
	TotalCost     float64              `bson:"totalCost" json:"totalCost"`
	UpdatedAt     time.Time            `bson:"updatedAt" json:"updatedAt"`
}

type EstimateItem struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	Cost        float64            `bson:"cost" json:"cost"`
	Description string             `bson:"descrip" json:"description"`
	Notes       string             `bson:"notes" json:"notes"`
	Qty         int                `bson:"qty" json:"qty"`
	Total       float64            `bson:"total" json:"total"`
}
