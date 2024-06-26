package mongodb

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo"
	"samm/internal/module/payment/domain"
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
