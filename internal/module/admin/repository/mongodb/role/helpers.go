package role

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"samm/pkg/database/mongodb"
)

func createIndexes(collection *mongo.Collection) {
	mongodb.CreateIndex(collection, false,
		bson.E{"name", mongodb.IndexType.Text},
	)
}
