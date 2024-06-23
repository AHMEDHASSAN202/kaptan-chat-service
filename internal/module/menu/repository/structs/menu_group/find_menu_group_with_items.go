package menu_group

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type LocalizationText struct {
	Ar string `json:"ar" bson:"ar"`
	En string `json:"en" bson:"en"`
}

type Availability struct {
	Day  string `json:"day" bson:"day"`
	From string `json:"from" bson:"from"`
	To   string `json:"to" bson:"to"`
}

type MenuItem struct {
	Id       string           `json:"id" bson:"_id"`
	ItemId   string           `json:"item_id" bson:"item_id"`
	Sort     int              `json:"sort" bson:"sort"`
	Name     LocalizationText `json:"name" bson:"name"`
	Calories int              `json:"calories" bson:"calories"`
	Price    float64          `json:"price" bson:"price"`
	Image    string           `json:"image" bson:"image"`
	Status   string           `json:"status" bson:"status"`
}

type MenuCategory struct {
	ID     primitive.ObjectID `json:"id" bson:"_id"`
	Name   LocalizationText   `json:"name" bson:"name"`
	Icon   string             `json:"icon" bson:"icon"`
	Sort   int                `json:"sort" bson:"sort"`
	Status string             `json:"status" bson:"status"`
	Items  []MenuItem         `json:"menu_items" bson:"menu_items"`
}

type ItemMenuGroup struct {
	ID             primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
	BranchIds      []primitive.ObjectID `json:"branch_ids" bson:"branch_ids"`
	Availabilities []Availability       `json:"availabilities" bson:"availabilities"`
	Status         string               `json:"status" bson:"status"`
}

type Branch struct {
	ID   primitive.ObjectID `json:"id"`
	Name LocalizationText   `json:"name"`
}

type FindMenuGroupWithItems struct {
	ID             primitive.ObjectID   `json:"id" bson:"_id"`
	AccountId      primitive.ObjectID   `json:"account_id" bson:"account_id"`
	Name           LocalizationText     `json:"name" bson:"name"`
	BranchIds      []primitive.ObjectID `json:"branch_ids" bson:"branch_ids"`
	Branches       []Branch             `json:"branches" bson:"branches"`
	Categories     []MenuCategory       `json:"categories" bson:"categories"`
	Availabilities []Availability       `json:"availabilities" bson:"availabilities"`
	Status         string               `json:"status" bson:"status"`
	CreatedAt      time.Time            `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time            `json:"updated_at" bson:"updated_at"`
}
