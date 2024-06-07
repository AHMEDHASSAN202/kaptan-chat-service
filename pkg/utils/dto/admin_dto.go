package dto

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AdminDetails struct {
	Id        primitive.ObjectID `bson:"id"`
	Name      string             `bson:"name"`
	Operation string             `bson:"operation"`
	UpdatedAt time.Time          `bson:"updatedAt"`
}
