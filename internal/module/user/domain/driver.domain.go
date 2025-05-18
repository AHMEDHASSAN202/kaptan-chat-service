package domain

import (
	"context"
)

type Driver struct {
	ID      uint   `gorm:"primarykey"`
	Name    string `gorm:"column:name"`
	Phone   string `gorm:"column:phone"`
	Address string `gorm:"column:name"`
}

type DriverRepository interface {
	Find(ctx *context.Context, id uint) (domainData *Driver, err error)
}
