package database

import (
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Connect opens (or creates) the SQLite file.
func Connect(dsn string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.Fatalf("cannot open database: %v", err)
	}
	return db
}
