package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	log "github.com/sirupsen/logrus"
	"github.com/webbtech/shts-pdf-gen/config"
	"github.com/webbtech/shts-pdf-gen/handlers"
	"github.com/webbtech/shts-pdf-gen/model"
	"github.com/webbtech/shts-pdf-gen/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
)

// see link below for best practices in connection pooling for lambda
// https://www.mongodb.com/docs/atlas/best-practices-connecting-from-aws-lambda/
// although the example above shows example with node.js, I believe the principal is the same when using the init function below

var (
	cfg    *config.Config
	db     model.DbHandler
	client *mongo.Client
)

// init isn't called for each invocation, so we take advantage and only setup cfg and db for (I'm assuming) cold starts
func init() {

	log.Info("calling config.Config.init in main")
	cfg = &config.Config{}
	err := cfg.Init()
	if err != nil {
		log.Fatal(err)
		return
	}

	db, err = mongodb.NewDb(cfg.GetMongoConnectString(), cfg.GetDbName())
	if err != nil {
		log.Fatal(err)
	}
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var h handlers.Handler

	switch request.HTTPMethod {
	case "DELETE":
		h = &handlers.Delete{Cfg: cfg}
	case "POST":
		h = &handlers.Post{Cfg: cfg, Db: db}
	case "GET":
		h = &handlers.Get{}
	default:
		h = &handlers.Any{}
	}

	return h.Response(request)
}

func main() {
	lambda.Start(handler)
}
