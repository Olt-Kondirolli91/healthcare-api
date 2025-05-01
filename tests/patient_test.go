package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Olt-Kondirolli91/healthcare-api/internal/database"
	"github.com/Olt-Kondirolli91/healthcare-api/internal/handlers"
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
