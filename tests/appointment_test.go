package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func TestListAppointments(t *testing.T) {
	r, _ := setupAppRouter()
	// First create an appointment
	createAppointmentRequest(r)

	// Now list appointments
	req, _ := http.NewRequest(http.MethodGet, "/appointments", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	// Parse response
	var appointments []models.Appointment
	json.Unmarshal(resp.Body.Bytes(), &appointments)
	assert.GreaterOrEqual(t, len(appointments), 1)
}

func TestGetAppointment(t *testing.T) {
	r, _ := setupAppRouter()
	// First create an appointment
	resp := createAppointmentRequest(r)

	var created models.Appointment
	json.Unmarshal(resp.Body.Bytes(), &created)

	// Now get that appointment
	req, _ := http.NewRequest(http.MethodGet, "/appointments/"+fmt.Sprint(created.ID), nil)
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var appointment models.Appointment
	json.Unmarshal(resp.Body.Bytes(), &appointment)
	assert.Equal(t, created.ID, appointment.ID)
	assert.Equal(t, uint(1), appointment.PatientID)
}

func TestUpdateAppointment(t *testing.T) {
	r, _ := setupAppRouter()
	// First create an appointment
	resp := createAppointmentRequest(r)

	var created models.Appointment
	json.Unmarshal(resp.Body.Bytes(), &created)

	// Now update the appointment
	tomorrow := time.Now().Add(24 * time.Hour).Format(time.RFC3339)
	updateData, _ := json.Marshal(map[string]interface{}{
		"patient_id": 1,
		"date_time":  tomorrow,
		"notes":      "rescheduled checkup",
	})

	req, _ := http.NewRequest(http.MethodPut, "/appointments/"+fmt.Sprint(created.ID),
		bytes.NewBuffer(updateData))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	// Verify changes
	var response map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &response)
	apptMap := response["appointment"].(map[string]interface{})
	assert.Equal(t, "rescheduled checkup", apptMap["notes"])
}

func TestDeleteAppointment(t *testing.T) {
	r, _ := setupAppRouter()
	// First create an appointment
	resp := createAppointmentRequest(r)

	var created models.Appointment
	json.Unmarshal(resp.Body.Bytes(), &created)

	// Now delete the appointment
	req, _ := http.NewRequest(http.MethodDelete, "/appointments/"+fmt.Sprint(created.ID), nil)
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	// Verify appointment was deleted
	req, _ = http.NewRequest(http.MethodGet, "/appointments/"+fmt.Sprint(created.ID), nil)
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
}

// Helper function to create a test appointment and return the response
func createAppointmentRequest(r *gin.Engine) *httptest.ResponseRecorder {
	payload, _ := json.Marshal(map[string]any{
		"patient_id": 1,
		"date_time":  time.Now().Format(time.RFC3339),
		"notes":      "annual checkup",
	})
	req, _ := http.NewRequest(http.MethodPost, "/appointments", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	return resp
}
