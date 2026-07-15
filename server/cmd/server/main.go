package main

import (
	"log"

	_ "github.com/1024XEngineer/xinfra/server/docs"
	"github.com/1024XEngineer/xinfra/server/internal/config"
	"github.com/1024XEngineer/xinfra/server/internal/database"
	"github.com/1024XEngineer/xinfra/server/internal/router"
)

func main() {
	cfg := config.Load()

	db, err := database.Open(cfg.MySQLDSN)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}

	if cfg.AutoMigrate {
		if err := database.AutoMigrate(db); err != nil {
			log.Fatalf("auto migrate: %v", err)
		}
	}
	if err := database.SeedDefaults(db, cfg); err != nil {
		log.Fatalf("seed defaults: %v", err)
	}

	r := router.New(router.Dependencies{
		Config: cfg,
		DB:     db,
	})

	log.Printf("authserver listening on %s", cfg.HTTPAddr)
	if err := r.Run(cfg.HTTPAddr); err != nil {
		log.Fatalf("run server: %v", err)
	}
}
