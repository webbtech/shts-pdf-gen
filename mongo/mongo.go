package mongo

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/webbtech/shts-pdf-gen/model"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Mdb struct
type Mdb struct {
	client *mongo.Client
	// db     *mongo.Database
}

// NewDb sets up a new Mdb struct
func NewDb(connection string) (model.DbHandler, error) {

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

	log.Println("Connected to MongoDB!")

	return &Mdb{
		client: client,
		// db: client.Database
	}, nil
}

// FetchEstimate method
func (db *Mdb) FetchEstimate(estimateId string) (*model.Estimate, error) {

	// Initialize
	q := &model.Estimate{}

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
