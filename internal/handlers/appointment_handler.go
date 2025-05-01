package handlers

import (
	"net/http"
	"time"

	"github.com/Olt-Kondirolli91/healthcare-api/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterAppointmentRoutes(r *gin.Engine, db *gorm.DB) {
	group := r.Group("/appointments")
	group.GET("", listAppointments(db))
	group.POST("", createAppointment(db))
	group.GET("/:id", getAppointment(db))
	group.PUT("/:id", updateAppointment(db))
	group.DELETE("/:id", deleteAppointment(db))
}

func listAppointments(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var a []models.Appointment
		if err := db.Find(&a).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "query failed"})
			return
		}
		c.JSON(http.StatusOK, a)
	}
}

func createAppointment(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var a models.Appointment
		if err := c.ShouldBindJSON(&a); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// parse ISO8601 if client sends string
		if c.Request.Header.Get("Content-Type") == "application/json" {
			if dt, err := time.Parse(time.RFC3339, c.PostForm("date_time")); err == nil {
				a.DateTime = dt
			}
		}

		// transaction to enforce patient existence
		err := db.Transaction(func(tx *gorm.DB) error {
			var p models.Patient
			if err := tx.First(&p, a.PatientID).Error; err != nil {
				return err
			}
			return tx.Create(&a).Error
		})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid patient or insert failed"})
			return
		}
		c.JSON(http.StatusCreated, a)
	}
}

func getAppointment(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var a models.Appointment
		if err := db.First(&a, c.Param("id")).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "appointment not found"})
			return
		}
		c.JSON(http.StatusOK, a)
	}
}

func updateAppointment(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var a models.Appointment
		if err := db.First(&a, c.Param("id")).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "appointment not found"})
			return
		}
		if err := c.ShouldBindJSON(&a); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		db.Save(&a)
		c.JSON(http.StatusOK, a)
	}
}

func deleteAppointment(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := db.Delete(&models.Appointment{}, c.Param("id")).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "delete failed"})
			return
		}
		c.Status(http.StatusNoContent)
	}
}
