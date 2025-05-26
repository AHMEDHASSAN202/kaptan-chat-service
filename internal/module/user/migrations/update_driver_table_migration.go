package migrations

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"kaptan/internal/module/user/domain"
)

func NewDriverMigration(db *gorm.DB) {
	// Add missing columns if they don't exist
	if !db.Migrator().HasColumn(&domain.Driver{}, "rating") {
		if err := db.Migrator().AddColumn(&domain.Driver{}, "rating"); err != nil {
			log.Error(err)
		}
	}

	if !db.Migrator().HasColumn(&domain.Driver{}, "sold_trips") {
		if err := db.Migrator().AddColumn(&domain.Driver{}, "sold_trips"); err != nil {
			log.Error(err)
		}
	}

	// Add indexes manually
	db.Exec("CREATE INDEX idx_drivers_rating ON drivers(rating)")
	db.Exec("CREATE INDEX idx_drivers_sold_trips ON drivers(sold_trips)")
}
