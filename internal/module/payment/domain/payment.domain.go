package domain

import (
	"github.com/kamva/mgm/v3"
)

type Payment struct {
	mgm.DefaultModel `bson:",inline"`
}
type PaymentUseCase interface {
}

type PaymentRepository interface {
}
