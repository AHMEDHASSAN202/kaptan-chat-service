package driver

import (
	"context"
	"gorm.io/gorm"
	"kaptan/internal/module/user/domain"
	"kaptan/pkg/logger"
)

type Repository struct {
	logger logger.ILogger
	db     *gorm.DB
}

func NewDriverRepository(log logger.ILogger, db *gorm.DB) domain.DriverRepository {
	return &Repository{
		logger: log,
		db:     db,
	}
}

func (r *Repository) Find(ctx *context.Context, id uint) (*domain.Driver, error) {
	driver := domain.Driver{ID: id}
	r.db.First(&driver)
	return &driver, nil
}

func (r *Repository) FindWithMedia(ctx *context.Context, id uint) (*domain.Driver, error) {
	var driver domain.Driver
	err := r.db.
		Preload("Media", func(db *gorm.DB) *gorm.DB {
			return db.Where("model_type = ?", "App\\Models\\Driver").
				Order("order_column ASC")
		}).
		First(&driver, id).Error
	if err != nil {
		return nil, err
	}
	return &driver, nil
}

func (r *Repository) FindByAccessTokenId(ctx *context.Context, id uint) (*domain.Driver, error) {
	var tokenableID uint
	err := r.db.Table("personal_access_tokens").
		Select("tokenable_id").
		Where("id = ? AND tokenable_type = ?", id, "App\\Models\\Driver").
		Scan(&tokenableID).Error
	if err != nil {
		return nil, err
	}
	driver := domain.Driver{ID: tokenableID}
	r.db.First(&driver)
	return &driver, nil
}

// IncrementSoldTripsByValue increments sold_trips by a specific value
func (r *Repository) IncrementSoldTripsByValue(ctx *context.Context, id uint, value int) error {
	if value <= 0 {
		return nil
	}

	result := r.db.Model(&domain.Driver{}).
		Where("id = ?", id).
		Update("sold_trips", gorm.Expr("sold_trips + ?", value))

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
