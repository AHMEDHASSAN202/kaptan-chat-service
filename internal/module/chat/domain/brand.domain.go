package domain

import (
	"gorm.io/gorm"
)

type Brand struct {
	gorm.Model
	Name              string `gorm:"column:name"`
	Sort              string `gorm:"column:sort"`
	SupportedCarBrand string `gorm:"column:supported_car_brands"`
}
