package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type AdminDetails struct {
	Id        primitive.ObjectID `bson:"id"`
	Name      string             `bson:"name"`
	UpdatedAt string             `bson:"updatedAt"`
}
