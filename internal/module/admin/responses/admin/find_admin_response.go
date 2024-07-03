package admin

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"samm/internal/module/admin/responses/role"
	"time"
)

type Name struct {
	Ar string `json:"ar" bson:"ar"`
	En string `json:"en" bson:"en"`
}

type AccountResp struct {
	Id   primitive.ObjectID `json:"id" bson:"_id"`
	Name Name               `json:"name" bson:"name"`
}

type FindAdminResponse struct {
	ID         primitive.ObjectID    `json:"id"`
	Name       string                `json:"name"`
	Email      string                `json:"email"`
	Type       string                `json:"type"`
	Role       role.FindRoleResponse `json:"role"`
	CountryIds []string              `json:"country_ids"`
	Account    *AccountResp          `json:"account"`
	Status     string                `json:"status"`
	CreatedAt  time.Time             `json:"created_at"`
	UpdateAt   time.Time             `json:"update_at"`
}
