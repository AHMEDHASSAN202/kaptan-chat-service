package responses

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CollectionMethod struct {
	Id     primitive.ObjectID `json:"id" bson:"_id"`
	Type   string             `json:"type" bson:"type"`
	Fields map[string]any     `json:"fields" bson:"fields"`
	Values map[string]any     `json:"values" bson:"values"`
}
