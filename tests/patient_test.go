package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Olt-Kondirolli91/healthcare-api/internal/database"
	"github.com/Olt-Kondirolli91/healthcare-api/internal/handlers"
	"github.com/Olt-Kondirolli91/healthcare-api/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	db := database.Connect(":memory:")
	database.AutoMigrate(db)
	r := gin.New()
	handlers.RegisterPatientRoutes(r, db)
	return r
}

func TestCreatePatient(t *testing.T) {
	r := setupRouter()
	body, _ := json.Marshal(map[string]string{
		"first_name": "John", "last_name": "Doe", "email": "john@doe.com",
	})
	req, _ := http.NewRequest(http.MethodPost, "/patients", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, 201, resp.Code)
}

func TestListPatients(t *testing.T) {
	r := setupRouter()
	// First create a patient
	createPatientRequest(r)

	// Now list patients
	req, _ := http.NewRequest(http.MethodGet, "/patients", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	// Parse response
	var patients []models.Patient
	json.Unmarshal(resp.Body.Bytes(), &patients)
	assert.GreaterOrEqual(t, len(patients), 1)
}

func TestGetPatient(t *testing.T) {
	r := setupRouter()
	// First create a patient
	resp := createPatientRequest(r)

	var created models.Patient
	json.Unmarshal(resp.Body.Bytes(), &created)

	// Now get that patient
	req, _ := http.NewRequest(http.MethodGet, "/patients/"+fmt.Sprint(created.ID), nil)
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var patient models.Patient
	json.Unmarshal(resp.Body.Bytes(), &patient)
	assert.Equal(t, created.ID, patient.ID)
	assert.Equal(t, "John", patient.FirstName)
}

func TestUpdatePatient(t *testing.T) {
	r := setupRouter()
	// First create a patient
	resp := createPatientRequest(r)

	var created models.Patient
	json.Unmarshal(resp.Body.Bytes(), &created)

	// Now update the patient
	updateData, _ := json.Marshal(map[string]string{
		"first_name": "Jane",
		"last_name":  "Smith",
		"email":      "jane@example.com",
		"phone":      "555-9876",
	})

	req, _ := http.NewRequest(http.MethodPut, "/patients/"+fmt.Sprint(created.ID),
		bytes.NewBuffer(updateData))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	// Verify changes
	var response map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &response)
	patientMap := response["patient"].(map[string]interface{})
	assert.Equal(t, "Jane", patientMap["first_name"])
	assert.Equal(t, "Smith", patientMap["last_name"])
}

func TestDeletePatient(t *testing.T) {
	r := setupRouter()
	// First create a patient
	resp := createPatientRequest(r)

	var created models.Patient
	json.Unmarshal(resp.Body.Bytes(), &created)

	// Now delete the patient
	req, _ := http.NewRequest(http.MethodDelete, "/patients/"+fmt.Sprint(created.ID), nil)
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	// Verify patient was deleted
	req, _ = http.NewRequest(http.MethodGet, "/patients/"+fmt.Sprint(created.ID), nil)
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
}

// Helper function to create a test patient and return the response
func createPatientRequest(r *gin.Engine) *httptest.ResponseRecorder {
	body, _ := json.Marshal(map[string]string{
		"first_name": "John", "last_name": "Doe", "email": "john@doe.com", "phone": "555-1234",
	})
	req, _ := http.NewRequest(http.MethodPost, "/patients", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	return resp
}
