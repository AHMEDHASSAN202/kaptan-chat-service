package role

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Name struct {
	Ar string `json:"ar"`
	En string `json:"en"`
}

type FindRoleResponse struct {
	ID          primitive.ObjectID `json:"id"`
	Name        Name               `json:"name"`
	Type        string             `json:"type"`
	Permissions []string           `json:"permissions"`
	CanDelete   bool               `json:"can_delete"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}
