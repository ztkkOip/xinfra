package database

import (
	"github.com/1024XEngineer/xinfra/server/internal/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Open(dsn string) (*gorm.DB, error) {
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.User{},
		&model.WayenCredential{},
		&model.BusinessLine{},
		&model.BusinessLineUser{},
		&model.BusinessLineWayneNamespace{},
		&model.AccessToken{},
		&model.AuditLog{},
	)
}
