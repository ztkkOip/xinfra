package database

import (
	"github.com/1024XEngineer/xinfra/server/internal/config"

	"gorm.io/gorm"
)

func SeedDefaults(_ *gorm.DB, _ config.Config) error {
	return nil
}
