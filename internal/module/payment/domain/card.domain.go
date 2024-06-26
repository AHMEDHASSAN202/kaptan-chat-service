package domain

import (
	"context"
	. "github.com/gobeam/mongo-go-pagination"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"samm/internal/module/payment/dto/card"
	"samm/pkg/validators"
	"time"
)

type Card struct {
	mgm.DefaultModel `bson:",inline"`
	Type             string             `json:"type" bson:"type"`
	Number           string             `json:"number" bson:"number"`
	MFToken          string             `json:"mf_token" bson:"mf_token"`
	UserId           primitive.ObjectID `json:"user_id" bson:"user_id"`
	DeletedAt        *time.Time         `json:"deleted_at" bson:"deleted_at"`
}
type CardUseCase interface {
	StoreCard(ctx context.Context, payload *card.CreateCardDto) (err validators.ErrorResponse)
	FindCard(ctx context.Context, Id string, userID string) (card Card, err validators.ErrorResponse)
	DeleteCard(ctx context.Context, Id string, userID string) (err validators.ErrorResponse)
	ListCard(ctx context.Context, payload *card.ListCardDto) (cards []Card, paginationResult PaginationData, err validators.ErrorResponse)
}

type CardRepository interface {
	StoreCard(ctx context.Context, user *Card) (err error)
	UpdateCard(ctx context.Context, card *Card) (err error)
	FindCard(ctx context.Context, Id primitive.ObjectID, userID primitive.ObjectID) (card *Card, err error)
	DeleteCard(ctx context.Context, Id primitive.ObjectID, userId primitive.ObjectID) (err error)
	ListCard(ctx context.Context, payload *card.ListCardDto) (cards []Card, paginationResult PaginationData, err error)
}
