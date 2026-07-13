package database

import (
	"authserver/internal/config"

	"gorm.io/gorm"
)

func SeedDefaults(_ *gorm.DB, _ config.Config) error {
	return nil
}
