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
	Permissions []string           `json:"permissions"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}
