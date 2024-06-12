package mobile

import "go.mongodb.org/mongo-driver/bson/primitive"

type LocalizationText struct {
	Ar string `json:"ar"`
	En string `json:"en"`
}

type GetMenuGroupItem struct {
	ID       primitive.ObjectID `json:"id"`
	Name     LocalizationText   `json:"name"`
	Image    string             `json:"image"`
	Price    float64            `json:"price"`
	Calories int                `json:"calories"`
	Sort     int                `json:"sort"`
}

type GetMenuGroupItemsResponse struct {
	ID    primitive.ObjectID `json:"id"`
	Name  LocalizationText   `json:"name"`
	Icon  string             `json:"icon"`
	Items []GetMenuGroupItem `json:"items"`
	Sort  int                `json:"sort"`
}
