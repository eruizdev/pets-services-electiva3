package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	return r
}

func TestCreatePet_InvalidJSON(t *testing.T) {
	r := setupRouter()
	r.POST("/pets/", CreatePet)

	req := httptest.NewRequest(http.MethodPost, "/pets/", strings.NewReader("not-json"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-Id", "1")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Esperaba 400, got %d", w.Code)
	}
}

func TestGetPetByID_NoDB(t *testing.T) {
	r := setupRouter()
	r.GET("/pets/:id", GetPetByID)

	req := httptest.NewRequest(http.MethodGet, "/pets/999", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	// Sin DB esperamos 404 o 500
	if w.Code == http.StatusOK {
		t.Errorf("No debería devolver 200 sin DB")
	}
}

func TestUpdatePet_InvalidJSON(t *testing.T) {
	r := setupRouter()
	r.PUT("/pets/:id", UpdatePet)

	req := httptest.NewRequest(http.MethodPut, "/pets/1", strings.NewReader("bad-json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Esperaba 400, got %d", w.Code)
	}
}
