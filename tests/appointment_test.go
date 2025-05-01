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
)

func setupAppRouter() (*gin.Engine, *database.DB) {
	gin.SetMode(gin.TestMode)
	db := database.Connect(":memory:")
	database.AutoMigrate(db)
	// seed a patient
	db.Create(&models.Patient{FirstName: "A", LastName: "B"})
	r := gin.New()
	handlers.RegisterAppointmentRoutes(r, db)
	return r, db
}

func TestCreateAppointment(t *testing.T) {
	r, _ := setupAppRouter()
	body, _ := json.Marshal(map[string]any{
		"patient_id": 1,
		"date_time":  time.Now().Format(time.RFC3339),
		"notes":      "annual checkup",
	})
	req, _ := http.NewRequest(http.MethodPost, "/appointments", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, 201, resp.Code)
}
