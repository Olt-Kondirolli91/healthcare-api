package main

import (
	"net/http"

	"github.com/Olt-Kondirolli91/healthcare-api/internal/config"
	"github.com/Olt-Kondirolli91/healthcare-api/internal/database"
	"github.com/Olt-Kondirolli91/healthcare-api/internal/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg := config.MustLoad()

	db := database.Connect(cfg.DB)
	database.AutoMigrate(db)

	// Seed initial data for testing
	database.SeedData(db)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	router.Use(cors.Default())
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome to Healthcare Appointment System")
	})

	handlers.RegisterPatientRoutes(router, db)
	handlers.RegisterAppointmentRoutes(router, db)

	logrus.Infof("server listening on :%s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		logrus.Fatalf("server error: %v", err)
	}
}
