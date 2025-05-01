package database

import (
	"time"

	"github.com/Olt-Kondirolli91/healthcare-api/internal/models"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// SeedData populates the database with initial data for testing
func SeedData(db *gorm.DB) {
	// Check if there's existing data
	var count int64
	db.Model(&models.Patient{}).Count(&count)
	if count > 0 {
		logrus.Info("Database already contains data, skipping seed")
		return
	}

	logrus.Info("Seeding database with initial data")

	// Create patients
	patients := []models.Patient{
		{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
			Phone:     "555-123-4567",
		},
		{
			FirstName: "Jane",
			LastName:  "Smith",
			Email:     "jane.smith@example.com",
			Phone:     "555-765-4321",
		},
		{
			FirstName: "Robert",
			LastName:  "Johnson",
			Email:     "robert.johnson@example.com",
			Phone:     "555-987-6543",
		},
	}

	for i := range patients {
		if err := db.Create(&patients[i]).Error; err != nil {
			logrus.Warnf("Error seeding patient: %v", err)
		}
	}

	// Create appointments
	now := time.Now()
	appointments := []models.Appointment{
		{
			PatientID: 1,
			DateTime:  now.Add(24 * time.Hour),
			Notes:     "Annual physical examination",
		},
		{
			PatientID: 2,
			DateTime:  now.Add(48 * time.Hour),
			Notes:     "Follow-up appointment",
		},
		{
			PatientID: 1,
			DateTime:  now.Add(7 * 24 * time.Hour),
			Notes:     "Vaccination",
		},
		{
			PatientID: 3,
			DateTime:  now.Add(2 * 24 * time.Hour),
			Notes:     "Initial consultation",
		},
	}

	for i := range appointments {
		if err := db.Create(&appointments[i]).Error; err != nil {
			logrus.Warnf("Error seeding appointment: %v", err)
		}
	}

	logrus.Info("Database seeding completed")
}
