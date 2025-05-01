package handlers

import (
	"net/http"
	"strings"

	"github.com/Olt-Kondirolli91/healthcare-api/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func RegisterPatientRoutes(r *gin.Engine, db *gorm.DB) {
	group := r.Group("/patients")
	group.GET("", listPatients(db))
	group.POST("", createPatient(db))
	group.GET("/:id", getPatient(db))
	group.PUT("/:id", updatePatient(db))
	group.DELETE("/:id", deletePatient(db))
}

func listPatients(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var patients []models.Patient
		if err := db.Preload("Appointments").Find(&patients).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "query failed"})
			return
		}
		c.JSON(http.StatusOK, patients)
	}
}

func createPatient(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var p models.Patient
		if err := c.ShouldBindJSON(&p); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := db.Omit("id").Create(&p).Error; err != nil {
			if strings.Contains(err.Error(), "UNIQUE constraint failed") {
				c.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
				return
			}
			logrus.Errorf("createPatient: db.Create error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, p)
	}
}

func getPatient(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var p models.Patient
		if err := db.Preload("Appointments").
			First(&p, c.Param("id")).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "patient not found"})
			return
		}
		c.JSON(http.StatusOK, p)
	}
}

func updatePatient(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var p models.Patient
		if err := db.First(&p, c.Param("id")).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "patient not found"})
			return
		}
		if err := c.ShouldBindJSON(&p); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		db.Save(&p)
		c.JSON(http.StatusOK, p)
	}
}

func deletePatient(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := db.Delete(&models.Patient{}, c.Param("id")).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "delete failed"})
			return
		}
		c.Status(http.StatusNoContent)
	}
}
