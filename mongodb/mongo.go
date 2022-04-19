package mongodb

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/webbtech/shts-pdf-gen/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DB and Collection constants
const (
	colCustomer      = "Customer"
	colEstimate      = "Estimate"
	colEstimateItems = "EstimateItem"
)

// Mdb struct
type Mdb struct {
	client *mongo.Client
	dbName string
	db     *mongo.Database
}

// ========================== Public Methods =============================== //

// NewDb sets up a new Mdb struct
func NewDb(connection string, dbNm string) (model.DbHandler, error) {

	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().ApplyURI(connection).SetServerAPIOptions(serverAPIOptions)
	if err := clientOptions.Validate(); err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	stage := os.Getenv("Stage")
	if stage != "test" {
		log.Println("Connected to MongoDB!")
	}

	return &Mdb{
		client: client,
		dbName: dbNm,
		db:     client.Database(dbNm),
	}, nil
}

// FetchEstimate method
func (db *Mdb) FetchEstimate(estimateNum int) (*model.Estimate, error) {

	// Initialize
	q := &model.Estimate{}

	// Fetch estimate
	if err := db.getEstimate(q, estimateNum); err != nil {
		return q, err
	}

	// Fetch customer
	if err := db.getCustomer(q); err != nil {
		return q, err
	}

	// Fetch Items
	if err := db.getItems(q); err != nil {
		return q, err
	}

	return q, nil
}

// Close method
func (db *Mdb) Close() {
	err := db.client.Disconnect(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")
}

// ========================== Private Methods =============================== //

func (db *Mdb) getEstimate(e *model.Estimate, estimateNum int) (err error) {

	if estimateNum <= 0 {
		return errors.New("missing estimateNum string")
	}

	col := db.db.Collection(colEstimate)
	filter := bson.D{primitive.E{Key: "number", Value: estimateNum}}

	if err = col.FindOne(context.Background(), filter).Decode(&e); err != nil {
		return err
	}

	return err
}

func (db *Mdb) getCustomer(e *model.Estimate) (err error) {

	if e.CustomerId.IsZero() {
		return errors.New("invalid customer id")
	}

	col := db.db.Collection(colCustomer)
	filter := bson.D{primitive.E{Key: "_id", Value: e.CustomerId}}

	err = col.FindOne(context.Background(), filter).Decode(&e.Customer)
	if err != nil {
		return err
	}

	return err
}

func (db *Mdb) getItems(e *model.Estimate) (err error) {

	if len(e.ItemIds) <= 0 {
		return nil
	}

	findOptions := options.Find()

	col := db.db.Collection(colEstimateItems)
	filter := bson.D{primitive.E{Key: "_id", Value: bson.D{primitive.E{Key: "$in", Value: e.ItemIds}}}}
	cur, err := col.Find(context.Background(), filter, findOptions)
	defer cur.Close(context.Background())
	if err != nil {
		return err
	}

	// var tmpItems []model.EstimateItem
	if err = cur.All(context.TODO(), &e.Items); err != nil {
		return err
	}

	// fmt.Printf("tmpItem: %+v\n", tmpItems)

	return err
}
