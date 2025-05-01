package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Olt-Kondirolli91/healthcare-api/internal/database"
	"github.com/Olt-Kondirolli91/healthcare-api/internal/handlers"
	"github.com/Olt-Kondirolli91/healthcare-api/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupAppRouter() (*gin.Engine, *gorm.DB) {
	gin.SetMode(gin.TestMode)
	db := database.Connect(":memory:")
	database.AutoMigrate(db)
	// seed a patient so appointments can reference it
	db.Create(&models.Patient{FirstName: "A", LastName: "B", Email: "a@b.com"})
	r := gin.New()
	handlers.RegisterAppointmentRoutes(r, db)
	return r, db
}

func TestCreateAppointment(t *testing.T) {
	r, _ := setupAppRouter()
	payload, _ := json.Marshal(map[string]any{
		"patient_id": 1,
		"date_time":  time.Now().Format(time.RFC3339),
		"notes":      "annual checkup",
	})
	req, _ := http.NewRequest(http.MethodPost, "/appointments", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
}
