package migrations

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"kaptan/internal/module/transfer/domain"
)

func NewTransferMigration(db *gorm.DB) {
	if !db.Migrator().HasColumn(&domain.Transfer{}, "seller_id") {
		if err := db.Migrator().AddColumn(&domain.Transfer{}, "seller_id"); err != nil {
			log.Error(err)
		}
		db.Exec("CREATE INDEX idx_seller_transfers ON transfers(seller_id)")
	}
}
