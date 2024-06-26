package admin

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type FindAdminResponse struct {
	ID          primitive.ObjectID `json:"id"`
	Name        string             `json:"name"`
	Email       string             `json:"email"`
	Type        string             `json:"type"`
	Role        string             `json:"role"`
	Permissions []string           `json:"permissions"`
	CountryIds  []string           `json:"country_ids"`
	MetaData    MetaData           `json:"meta_data"`
	Status      string             `json:"status"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdateAt    time.Time          `json:"update_at"`
}
