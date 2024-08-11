package kitchen

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type LocalizationText struct {
	Ar string `json:"ar" bson:"ar"`
	En string `json:"en" bson:"en"`
}

type MetaData struct {
	HasMissingItems  bool                 `json:"has_missing_items" bson:"has_missing_items"`
	TargetKitchenIds []primitive.ObjectID `bson:"target_kitchen_ids" json:"target_kitchen_ids"`
}
type KitchenListOrdersResponse struct {
	ID               primitive.ObjectID `json:"id" bson:"_id"`
	SerialNum        string             `json:"serial_num" bson:"serial_num"`
	User             User               `json:"user" bson:"user"`
	Items            []Item             `json:"items" bson:"items"`
	Location         Location           `json:"location" bson:"location"`
	PreparationTime  int                `json:"preparation_time" bson:"preparation_time"`
	PriceSummary     OrderPriceSummary  `json:"price_summary" bson:"price_summary"`
	Status           string             `json:"status" bson:"status"`
	IsFavourite      bool               `json:"is_favourite" bson:"is_favourite"`
	AcceptedAt       *time.Time         `json:"accepted_at" bson:"accepted_at"`
	PaidAt           *time.Time         `json:"paid_at" bson:"paid_at"`
	ArrivedAt        *time.Time         `json:"arrived_at" bson:"arrived_at"`
	PickedUpAt       *time.Time         `json:"pickedup_at" bson:"pickedup_at"`
	ReadyForPickUpAt *time.Time         `json:"ready_for_pickup_at" bson:"ready_for_pickup_at"`
	CancelledAt      *time.Time         `json:"cancelled_at" bson:"cancelled_at"`
	RejectedAt       *time.Time         `json:"rejected_at" bson:"rejected_at"`
	NoShowAt         *time.Time         `json:"no_show_at" bson:"no_show_at"`
	Cancelled        *Rejected          `json:"cancelled,omitempty" bson:"cancelled,omitempty"`
	Rejected         *Rejected          `json:"rejected,omitempty" bson:"rejected,omitempty"`
	Notes            string             `json:"notes" bson:"notes"`
	MetaData         MetaData           `json:"meta_data" bson:"meta_data"`
}
