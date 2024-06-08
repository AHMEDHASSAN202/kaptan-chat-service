package domain

import "github.com/kamva/mgm/v3"

type Location struct {
	mgm.DefaultModel `bson:",inline"`
}

type LocationUseCase interface {
}

type LocationRepository interface {
}
