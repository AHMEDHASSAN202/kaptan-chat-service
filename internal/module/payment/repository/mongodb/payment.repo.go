package mongodb

import (
	"context"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"samm/internal/module/payment/domain"
	"samm/pkg/utils"
)

type PaymentRepository struct {
	paymentCollection *mgm.Collection
}

func NewPaymentMongoRepository(dbs *mongo.Database) domain.PaymentRepository {
	paymentDbCollection := mgm.Coll(&domain.Payment{})

	return &PaymentRepository{
		paymentCollection: paymentDbCollection,
	}
}

func (p PaymentRepository) CreateTransaction(ctx context.Context, document *domain.Payment) (response *domain.Payment, err error) {
	document.ID = primitive.NewObjectID()
	err = mgm.Coll(document).CreateWithCtx(ctx, document)
	if err != nil {
		return nil, err
	}
	return document, nil
}

func (p PaymentRepository) UpdateTransaction(ctx context.Context, document *domain.Payment) (err error) {
	update := bson.M{"$set": document}
	_, err = mgm.Coll(document).UpdateByID(ctx, document.ID, update)
	return
}

func (p PaymentRepository) FindPaymentTransaction(ctx context.Context, id string, transactionId string, transactionType string) (payment *domain.Payment, err error) {
	domainData := domain.Payment{}
	var filter bson.M
	if id != "" {
		filter = bson.M{"_id": utils.ConvertStringIdToObjectId(id)}
	} else {
		filter = bson.M{"transaction_id": utils.ConvertStringIdToObjectId(transactionId), "transaction_type": transactionType}
	}
	err = p.paymentCollection.FirstWithCtx(ctx, filter, &domainData)

	return &domainData, err
}
