package menu_group

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"samm/pkg/database/mongodb"
)

func createIndexes(collection *mongo.Collection) {
	mongodb.CreateIndex(collection, false,
		bson.E{"account_id", mongodb.IndexType.Asc},
		bson.E{"branch_ids", mongodb.IndexType.Asc},
		bson.E{"status", mongodb.IndexType.Asc},
	)

	mongodb.CreateIndex(collection, false,
		bson.E{"name.ar", mongodb.IndexType.Text},
		bson.E{"name.en", mongodb.IndexType.Text},
	)

	mongodb.CreateIndex(collection, false,
		bson.E{"categories._id", mongodb.IndexType.Asc},
		bson.E{"categories.sort", mongodb.IndexType.Asc},
		bson.E{"categories.status", mongodb.IndexType.Asc},
	)

	mongodb.CreateIndex(collection, false,
		bson.E{"branch_ids", mongodb.IndexType.Asc},
		bson.E{"status", mongodb.IndexType.Asc},
	)
}
