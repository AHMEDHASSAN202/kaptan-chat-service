package responses

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Status struct {
	SnoozeTo string `json:"snooze_to"`
	Status   string `json:"status" bson:"status"`
	Meta     Meta   `json:"meta" bson:"meta"`
}
type Meta struct {
	NameEN string `json:"name_en" bson:"name_en"`
	NameAr string `json:"name_ar" bson:"name_ar"`
	Color  string `json:"color" bson:"color"`
}
type Name struct {
	Ar string `json:"ar" validate:"required,min=3"`
	En string `json:"en" validate:"required,min=3"`
}
type City struct {
	Id   primitive.ObjectID `json:"_id" bson:"id"`
	Name Name               `json:"name" bson:"name"`
}
type Coordinate struct {
	Type        string    `json:"type" bson:"type"`
	Coordinates []float64 `json:"coordinates" bson:"coordinates"`
}
type BrandDetails struct {
	Id       primitive.ObjectID `json:"_id" bson:"_id"`
	Name     Name               `json:"name" bson:"name"`
	Logo     string             `json:"logo" bson:"logo"`
	IsActive bool               `json:"is_active" bson:"is_active"`
	Cuisines []CuisineDetails   `json:"cuisines" bson:"cuisines"`
}
type WorkingHour struct {
	Day       string `json:"day" bson:"day"`
	From      string `json:"from" bson:"from"`
	IsFullDay bool   `json:"is_full_day" bson:"is_full_day"`
	To        string `json:"to" bson:"to"`
}
type Country struct {
	Id   string `json:"_id" bson:"_id"`
	Name struct {
		Ar string `json:"ar" bson:"ar"`
		En string `json:"en" bson:"en"`
	} `json:"name" bson:"name"`
	Timezone    string `json:"timezone" bson:"timezone"`
	Currency    string `json:"currency" bson:"currency"`
	PhonePrefix string `json:"phone_prefix" bson:"phone_prefix"`
}

type CuisineDetails struct {
	Id       string `json:"_id" bson:"_id"`
	Name     Name   `json:"name" bson:"name"`
	Logo     string `json:"logo" bson:"logo"`
	IsHidden bool   `json:"is_hidden" bson:"is_hidden"`
}
type LocationDetails struct {
	Id              string        `json:"_id" bson:"_id"`
	Name            Name          `json:"name" bson:"name"`
	City            City          `json:"city" bson:"city"`
	Street          Name          `json:"street" bson:"street"`
	CoverImage      string        `json:"cover_image" bson:"cover_image"`
	Logo            string        `json:"logo" bson:"logo"`
	SnoozeTo        *time.Time    `json:"snooze_to" bson:"snooze_to"`
	IsOpen          bool          `json:"is_open" bson:"is_open"`
	WorkingHour     []WorkingHour `json:"working_hour" bson:"working_hour"`
	Phone           string        `json:"phone" bson:"phone"`
	Coordinate      Coordinate    `json:"coordinate" bson:"coordinate"`
	BrandDetails    BrandDetails  `json:"brand_details" bson:"brand_details"`
	PreparationTime int           `json:"preparation_time" bson:"preparation_time"`
	Distance        float64       `json:"distance" bson:"distance"`
	Country         Country       `json:"country" bson:"country"`
	Status          Status        `json:"status" bson:"-"`
}
