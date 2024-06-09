package dashboard

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type ListMenuGroupResponse struct {
	ID        primitive.ObjectID   `json:"id"`
	AccountId primitive.ObjectID   `json:"account_id"`
	Name      interface{}          `json:"name"`
	BranchIds []primitive.ObjectID `json:"branch_ids"`
	Status    string               `json:"status"`
	CreatedAt time.Time            `json:"created_at"`
	UpdateAt  time.Time            `json:"update_at"`
}
