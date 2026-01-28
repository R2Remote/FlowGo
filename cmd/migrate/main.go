package main

import (
	"log"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"FLOWGO/internal/domain/entity"
	"FLOWGO/internal/infrastructure/config"
)

func main() {
	// 1. Load Config
	if err := config.LoadConfig("config.yaml"); err != nil {
		log.Println("Could not load config.yaml, assuming environment variables or defaults")
	}

	// SQLite Destination
	sqliteFile := "flowgo.db"
	log.Printf("Connecting to SQLite destination: %s...", sqliteFile)
	sqliteDB, err := gorm.Open(sqlite.Open(sqliteFile), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to SQLite: %v", err)
	}
	log.Println("Connected to SQLite.")

	// 2. AutoMigrate Schema
	log.Println("Migrating schema...")
	err = sqliteDB.AutoMigrate(
		&entity.User{},
		&entity.Project{},
		&entity.Team{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate schema: %v", err)
	}

	log.Println("Migration complete!")
}
