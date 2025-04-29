package mysql

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"kaptan/pkg/config"
	"os"
	"time"
)

func NewClient(config *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.Mysql.USERNAME, config.Mysql.PASSWORD, config.Mysql.HOST, config.Mysql.PORT, config.Mysql.DATABASE)
	var enableLogs logger.LogLevel
	if os.Getenv("LOG_DB_QUERIES") == "TRUE" {
		enableLogs = logger.Info
	} else {
		enableLogs = logger.Silent
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(enableLogs)})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)
	return db, err
}
