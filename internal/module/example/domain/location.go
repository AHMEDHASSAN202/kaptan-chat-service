package domain

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
	"time"
)

type Location struct {
	// DefaultModel adds _id, created_at and updated_at fields to the Model.
	mgm.DefaultModel  `bson:",inline"`
	ID                primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Details           Items              `json:"details" bson:"details"`
	Status            string             `json:"status" bson:"status"`
	ChannelLocationID string             `json:"channelLocationId" bson:"channel_location_id"`
	ChannelLinkID     string             `json:"channelLinkId" bson:"channel_link_id"`
	LocationID        string             `json:"locationId" bson:"location_id"`
	ChannelLinkName   string             `json:"channelLinkName" bson:"channel_link_name"`
	NgoBranchRefId    string             `json:"ngo_branch_ref_id" bson:"ngo_branch_ref_id"`
}
type Dummy struct {
	AnonymizeCustomer          bool        `json:"anonymizeCustomer" bson:"anonymize_customer"`
	ValidateProducts           bool        `json:"validateProducts" bson:"validate_products"`
	SendDiscount               bool        `json:"sendDiscount" bson:"send_discount"`
	SendDeliveryFee            bool        `json:"sendDeliveryFee" bson:"send_delivery_fee"`
	SendServiceCharge          bool        `json:"sendServiceCharge" bson:"send_service_charge"`
	SendDeliveryFeeCondition   int         `json:"sendDeliveryFeeCondition" bson:"send_delivery_fee_condition"`
	PosOrdersAreAlwaysPaid     bool        `json:"posOrdersAreAlwaysPaid" bson:"pos_orders_are_always_paid"`
	SortModifiers              int         `json:"sortModifiers" bson:"sort_modifiers"`
	BufferOrders               interface{} `json:"bufferOrders" bson:"buffer_orders"`
	DeliveryByChannelIsPickup  interface{} `json:"deliveryByChannelIsPickup" bson:"delivery_by_channel_is_pickup"`
	IgnoreUnknownOrderStatuses bool        `json:"ignoreUnknownOrderStatuses" bson:"ignore_unknown_order_statuses"`
	RecalcPriceOverrides       interface{} `json:"recalcPriceOverrides" bson:"recalc_price_overrides"`
	SeparateSameProducts       interface{} `json:"separateSameProducts" bson:"separate_same_products"`
	DontSendCancel             bool        `json:"dontSendCancel" bson:"dont_send_cancel"`
	Readonly                   bool        `json:"readonly" bson:"readonly"`
	LogOps                     bool        `json:"logOps" bson:"log_ops"`
	AveragePreparationTime     int         `json:"AveragePreparationTime" bson:"average_preparation_time"`
	OrderStatus                int         `json:"orderStatus" bson:"order_status"`
}
type PosSettings struct {
	Dummy Dummy `json:"dummy" bson:"dummy"`
}

type BrandDetails struct {
	BrandID   string `json:"brand_id" bson:"brand_id"`
	BrandName string `json:"brand_name" bson:"brand_name"`
}
type Subscriptions struct {
	BrandID string `json:"brandId" bson:"brand_id"`
}
type Address struct {
	Remarks         string      `json:"remarks" bson:"remarks"`
	PhoneNumber     string      `json:"phoneNumber" bson:"phone_number"`
	HouseNumber     string      `json:"houseNumber" bson:"house_number"`
	Street          string      `json:"street" bson:"street"`
	PostalCode      string      `json:"postalCode" bson:"postal_code"`
	City            string      `json:"city" bson:"city"`
	StateOrProvince string      `json:"stateOrProvince" bson:"state_or_province"`
	RestaurantName  string      `json:"restaurantName" bson:"restaurant_name"`
	Country         string      `json:"country"  bson:"country"`
	Coordinates     Coordinates `json:"coordinates"`
}
type Coordinates struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}
type Contact struct {
	PhoneNumber string `json:"phoneNumber" bson:"phone_number"`
	Email       string `json:"email" bson:"email"`
	FirstName   string `json:"firstName" bson:"first_name"`
	LastName    string `json:"lastName" bson:"last_name"`
}
type Geolocation struct {
	Country     string `json:"country" bson:"country"`
	Subdivision string `json:"subdivision" bson:"subdivision"`
	City        string `json:"city" bson:"city"`
}
type OpeningHours struct {
	DayOfWeek int    `json:"dayOfWeek" bson:"day_of_week"`
	StartTime string `json:"startTime" bson:"start_time"`
	EndTime   string `json:"endTime" bson:"end_time"`
}
type Account struct {
	Title string `json:"title" bson:"title"`
	Href  string `json:"href" bson:"href"`
}
type Items struct {
	ID                         string          `json:"_id" bson:"id"`
	PosSettings                PosSettings     `json:"posSettings" bson:"pos_settings"`
	Status                     string          `json:"status" bson:"status"`
	Subscriptions              []Subscriptions `json:"subscriptions" bson:"subscriptions"`
	BrandDetails               BrandDetails    `json:"brand_details" bson:"brand_details"`
	Comment                    string          `json:"comment" bson:"comment"`
	EnableOutOfStock           bool            `json:"enableOutOfStock" bson:"enable_out_of_stock"`
	SplitOrderItems            interface{}     `json:"splitOrderItems" bson:"split_order_items"`
	EnableWorkstations         bool            `json:"enableWorkstations" bson:"enable_workstations"`
	IgnorePOSOrderStatuses     bool            `json:"ignorePOSOrderStatuses" bson:"ignore_pos_order_statuses"`
	PosSystemID                int             `json:"posSystemId" bson:"pos_system_id"`
	Address                    Address         `json:"address" bson:"address"`
	Contact                    Contact         `json:"contact" bson:"contact"`
	Timezone                   string          `json:"timezone" bson:"timezone"`
	Name                       string          `json:"name" bson:"name"`
	TaxExcl                    interface{}     `json:"taxExcl" bson:"taxExcl"`
	EnablePickupScreen         bool            `json:"enablePickupScreen" bson:"enable_pickup_screen"`
	EnablePickupScreenQR       bool            `json:"enablePickupScreenQR" bson:"enable_pickup_screen_qr"`
	PickupScreenSubscriptionID string          `json:"pickupScreenSubscriptionId" bson:"pickup_screen_subscription_id"`
	Account                    string          `json:"account" bson:"account"`
	Updated                    time.Time       `json:"_updated" bson:"updated"`
	Created                    time.Time       `json:"_created" bson:"created"`
	Deleted                    bool            `json:"_deleted" bson:"deleted"`
	Geolocation                Geolocation     `json:"geolocation" bson:"geolocation"`
	Etag                       string          `json:"_etag" bson:"etag"`
	ChannelLinks               []string        `json:"channelLinks" bson:"channelLinks"`
	OpeningHours               []OpeningHours  `json:"openingHours" bson:"openingHours"`
	Region                     string          `json:"region" bson:"region"`
	IsCanary                   bool            `json:"isCanary" bson:"isCanary"`
}

func (s Location) IsEmpty() bool {
	return reflect.DeepEqual(s, Location{})
}

func (s Address) IsEmpty() bool {
	return reflect.DeepEqual(s, Address{})
}
func (s Coordinates) IsEmpty() bool {
	return reflect.DeepEqual(s, Coordinates{})
}

func (s Items) IsEmpty() bool {
	return reflect.DeepEqual(s, Items{})
}
