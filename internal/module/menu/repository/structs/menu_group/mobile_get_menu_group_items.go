package menu_group

import "go.mongodb.org/mongo-driver/bson/primitive"

type MobileGetMenuGroupItem struct {
	ID       primitive.ObjectID `bson:"_id" json:"id"`
	Name     LocalizationText   `bson:"name" json:"name"`
	Image    string             `bson:"image" json:"image"`
	Price    float64            `bson:"price" json:"price"`
	Calories int                `bson:"calories" json:"calories"`
	Sort     int                `bson:"sort" json:"sort"`
}

type MobileGetMenuGroupItems struct {
	ID    primitive.ObjectID       `bson:"_id" json:"id"`
	Name  LocalizationText         `bson:"name" json:"name"`
	Icon  string                   `bson:"icon" json:"icon"`
	Items []MobileGetMenuGroupItem `bson:"items" json:"items"`
	Sort  int                      `bson:"sort" json:"sort"`
}
