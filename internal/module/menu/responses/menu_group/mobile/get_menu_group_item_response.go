package mobile

import "go.mongodb.org/mongo-driver/bson/primitive"

type MobileGetItemResponse struct {
	ID               primitive.ObjectID      `json:"id"`
	ItemId           primitive.ObjectID      `json:"item_id"`
	Name             LocalizationText        `json:"name"`
	Desc             LocalizationText        `json:"desc"`
	Calories         int                     `json:"calories"`
	Price            float64                 `json:"price"`
	ModifierGroupIds []primitive.ObjectID    `json:"modifier_group_ids"`
	ModifierGroups   []ModifierGroupResponse `json:"modifier_groups"`
	Category         GetItemCategoryResponse `json:"category"`
	Tags             []string                `json:"tags"`
	Image            string                  `json:"image"`
}

type GetItemCategoryResponse struct {
	ID   primitive.ObjectID `json:"id"`
	Name LocalizationText   `json:"name"`
	Icon string             `json:"icon"`
}

type ModifierGroupResponse struct {
	ID         primitive.ObjectID           `json:"id"`
	Name       LocalizationText             `json:"name"`
	Type       string                       `json:"type"`
	Min        int                          `json:"min"`
	Max        int                          `json:"max"`
	ProductIds []primitive.ObjectID         `json:"product_ids"`
	Addons     []MobileGetItemAddonResponse `json:"addons"`
}

type MobileGetItemAddonResponse struct {
	ID       primitive.ObjectID `json:"id"`
	Name     LocalizationText   `json:"name"`
	Type     string             `json:"type"`
	Min      int                `json:"min"`
	Max      int                `json:"max"`
	Calories int                `json:"calories"`
	Price    float64            `json:"price"`
	Image    string             `json:"image"`
}
