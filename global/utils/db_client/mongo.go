package db_client

import (
	pkgMongo "tool/pkg/mongo"

	"go.mongodb.org/mongo-driver/mongo"
)

func MongoLocal() *mongo.Database {

	return pkgMongo.NewClient("Local")
}
