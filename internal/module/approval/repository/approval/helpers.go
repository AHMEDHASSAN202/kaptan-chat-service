package approval

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"samm/pkg/database/mongodb"
)

func createApprovalIndexes(collection *mongo.Collection) {
	mongodb.CreateIndex(collection, false,
		bson.E{"country_id", mongodb.IndexType.Asc},
		bson.E{"status", mongodb.IndexType.Asc},
		bson.E{"type", mongodb.IndexType.Asc},
	)

	mongodb.CreateIndex(collection, false,
		bson.E{"entity_id", mongodb.IndexType.Asc},
		bson.E{"entity_type", mongodb.IndexType.Asc},
		bson.E{"status", mongodb.IndexType.Asc},
	)

	mongodb.CreateIndex(collection, false,
		bson.E{"created_at", mongodb.IndexType.Desc},
	)

	oneMonth := int32(30 * 24 * 60 * 60)
	mongodb.CreateIndexWithTTL(collection, false, bson.E{"dates.approved_at", mongodb.IndexType.Desc}, oneMonth)
	mongodb.CreateIndexWithTTL(collection, false, bson.E{"dates.rejected_at", mongodb.IndexType.Desc}, oneMonth)
}
