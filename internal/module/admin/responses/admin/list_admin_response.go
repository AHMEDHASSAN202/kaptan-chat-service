package admin

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"samm/internal/module/admin/responses/role"
	"time"
)

type MetaData struct {
	AccountId string `json:"account_id" bson:"account_id"`
}

type ListAdminResponse struct {
	ID         primitive.ObjectID    `json:"id"`
	Name       string                `json:"name"`
	Email      string                `json:"email"`
	Type       string                `json:"type"`
	Role       role.FindRoleResponse `json:"role"`
	CountryIds []string              `json:"country_ids"`
	MetaData   MetaData              `json:"meta_data"`
	Status     string                `json:"status"`
	CreatedAt  time.Time             `json:"created_at"`
	UpdateAt   time.Time             `json:"update_at"`
}
