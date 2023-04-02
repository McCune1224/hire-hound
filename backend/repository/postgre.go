package repository

import (
	"hirehound/models"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Global variable to access database
var DB *gorm.DB

func Connect() {
	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.Logger.LogMode(logger.Info)

	// // Delete tables if exist
	// db.Migrator().DropTable(models.Application{})
	// db.Migrator().DropTable(models.Company{})

	// Automigrations
	db.AutoMigrate(models.User{})
	db.AutoMigrate(models.Application{})

	DB = db
}
