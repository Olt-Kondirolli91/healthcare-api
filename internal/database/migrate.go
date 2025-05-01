package database

import (
	"github.com/Olt-Kondirolli91/healthcare-api/internal/models"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// AutoMigrate runs GORM migrations.
func AutoMigrate(db *gorm.DB) {
	if err := db.AutoMigrate(&models.Patient{}, &models.Appointment{}); err != nil {
		logrus.Fatalf("migration failed: %v", err)
	}
}
