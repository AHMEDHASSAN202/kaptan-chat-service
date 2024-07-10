package mongodb

import (
	"context"
	. "github.com/gobeam/mongo-go-pagination"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"samm/internal/module/payment/domain"
	"samm/internal/module/payment/dto/card"
	"samm/pkg/utils"
	"time"
)

type CardRepository struct {
	cardCollection *mgm.Collection
}

func (c CardRepository) StoreCard(ctx context.Context, card *domain.Card) (err error) {
	err = mgm.Coll(card).CreateWithCtx(ctx, card)
	if err != nil {
		return err
	}
	return nil
}

func (c CardRepository) FindCard(ctx context.Context, Id primitive.ObjectID, userId primitive.ObjectID) (card *domain.Card, err error) {
	domainData := domain.Card{}
	filter := bson.M{"deleted_at": nil, "_id": Id, "user_id": userId}
	err = c.cardCollection.FirstWithCtx(ctx, filter, &domainData)

	return &domainData, err
}

func (l CardRepository) UpdateCard(ctx context.Context, card *domain.Card) (err error) {
	update := bson.M{"$set": card}
	_, err = mgm.Coll(card).UpdateByID(ctx, card.ID, update)
	return
}
func (c CardRepository) DeleteCard(ctx context.Context, Id primitive.ObjectID, userId primitive.ObjectID) (err error) {
	userData, err := c.FindCard(ctx, Id, userId)
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	userData.DeletedAt = &now
	userData.UpdatedAt = now
	return c.UpdateCard(ctx, userData)
}

func (c CardRepository) ListCard(ctx context.Context, payload *card.ListCardDto) (cards []domain.Card, paginationResult PaginationData, err error) {
	models := make([]domain.Card, 0)

	matching := bson.M{"$match": bson.M{"$and": []interface{}{
		bson.M{"deleted_at": nil},
		bson.M{"user_id": utils.ConvertStringIdToObjectId(payload.UserId)},
	}}}

	data, err := New(c.cardCollection.Collection).Context(ctx).Limit(payload.Limit).Page(payload.Page).Sort("created_at", -1).Aggregate(matching)
	if data == nil || data.Data == nil {
		return models, paginationResult, err
	}

	for _, raw := range data.Data {
		model := domain.Card{}
		errUnmarshal := bson.Unmarshal(raw, &model)
		if errUnmarshal != nil {
			break
		}
		models = append(models, model)
	}
	return models, data.Pagination, err
}

func NewCardMongoRepository(dbs *mongo.Database) domain.CardRepository {
	cardDbCollection := mgm.Coll(&domain.Card{})

	return &CardRepository{
		cardCollection: cardDbCollection,
	}
}
