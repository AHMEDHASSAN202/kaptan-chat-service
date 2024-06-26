package admin

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"samm/pkg/database/mongodb"
)

func createIndexes(collection *mongo.Collection) {
	mongodb.CreateIndex(collection, false,
		bson.E{"deleted_at", mongodb.IndexType.Asc},
		bson.E{"country_ids", mongodb.IndexType.Asc},
		bson.E{"status", mongodb.IndexType.Asc},
		bson.E{"type", mongodb.IndexType.Asc},
		bson.E{"role", mongodb.IndexType.Asc},
	)

	mongodb.CreateIndex(collection, false,
		bson.E{"deleted_at", mongodb.IndexType.Asc},
		bson.E{"name", mongodb.IndexType.Text},
		bson.E{"email", mongodb.IndexType.Text},
	)

	mongodb.CreateIndex(collection, true,
		bson.E{"deleted_at", mongodb.IndexType.Asc},
		bson.E{"email", mongodb.IndexType.Asc},
	)

	mongodb.CreateIndex(collection, false,
		bson.E{"deleted_at", mongodb.IndexType.Asc},
		bson.E{"status", mongodb.IndexType.Asc},
	)
}
