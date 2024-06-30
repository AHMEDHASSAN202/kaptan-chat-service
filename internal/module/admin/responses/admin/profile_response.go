package admin

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AdminProfileResponse struct {
	ID          primitive.ObjectID `json:"id"`
	Name        string             `json:"name"`
	Email       string             `json:"email"`
	Type        string             `json:"type"`
	Role        string             `json:"role"`
	Permissions []string           `json:"permissions"`
	CountryIds  []string           `json:"country_ids"`
	AccountId   string             `json:"account_id"`
}
