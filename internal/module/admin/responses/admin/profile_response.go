package admin

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"samm/internal/module/admin/responses/role"
)

type AdminProfileResponse struct {
	ID         primitive.ObjectID    `json:"id"`
	Name       string                `json:"name"`
	Email      string                `json:"email"`
	Type       string                `json:"type"`
	Role       role.FindRoleResponse `json:"role"`
	CountryIds []string              `json:"country_ids"`
	AccountId  string                `json:"account_id"`
}
