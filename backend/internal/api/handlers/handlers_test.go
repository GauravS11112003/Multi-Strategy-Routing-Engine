package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"shipt-route-optimizer/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	api := r.Group("/api")
	{
		api.GET("/health", HealthCheck)
		api.GET("/sample-data", GetSampleData)
		api.POST("/optimize", OptimizeRoutes)
		api.POST("/optimize-analytics", OptimizeWithAnalytics)
	}
	return r
}

func TestHealthCheck(t *testing.T) {
	router := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/health", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "ok", resp["status"])
	assert.Equal(t, "shipt-route-optimizer", resp["service"])
}

func TestGetSampleData(t *testing.T) {
	router := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/sample-data", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp models.SampleDataResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.NotEmpty(t, resp.Orders)
	assert.NotEmpty(t, resp.Shoppers)
}

func TestOptimizeRoutes_ValidRequest(t *testing.T) {
	router := setupRouter()

	body := models.OptimizeRequest{
		Orders: []models.Order{
			{ID: "o1", Lat: 33.52, Lng: -86.81, ItemCount: 3, DeliveryWindow: "morning"},
			{ID: "o2", Lat: 33.53, Lng: -86.82, ItemCount: 2, DeliveryWindow: "afternoon"},
		},
		Shoppers: []models.Shopper{
			{ID: "s1", Lat: 33.50, Lng: -86.80, Capacity: 10},
		},
	}

	bodyBytes, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/optimize", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp models.OptimizeResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.NotEmpty(t, resp.Assignments)
	assert.Greater(t, resp.TotalDistanceAfter, 0.0)
}

func TestOptimizeRoutes_NoOrders(t *testing.T) {
	router := setupRouter()

	body := models.OptimizeRequest{
		Orders:   []models.Order{},
		Shoppers: []models.Shopper{{ID: "s1", Lat: 33.5, Lng: -86.8, Capacity: 10}},
	}

	bodyBytes, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/optimize", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestOptimizeRoutes_NoShoppers(t *testing.T) {
	router := setupRouter()

	body := models.OptimizeRequest{
		Orders:   []models.Order{{ID: "o1", Lat: 33.5, Lng: -86.8, ItemCount: 3}},
		Shoppers: []models.Shopper{},
	}

	bodyBytes, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/optimize", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestOptimizeRoutes_InvalidJSON(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/optimize", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestOptimizeWithAnalytics_ValidRequest(t *testing.T) {
	router := setupRouter()

	body := map[string]interface{}{
		"orders": []models.Order{
			{ID: "o1", Lat: 33.52, Lng: -86.81, ItemCount: 3, DeliveryWindow: "morning"},
			{ID: "o2", Lat: 33.53, Lng: -86.82, ItemCount: 2, DeliveryWindow: "afternoon"},
		},
		"shoppers": []models.Shopper{
			{ID: "s1", Lat: 33.50, Lng: -86.80, Capacity: 10},
		},
		"algorithm":     "nearest-neighbor",
		"useRealRoutes": false,
	}

	bodyBytes, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/optimize-analytics", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Contains(t, resp, "optimization")
	assert.Contains(t, resp, "analytics")
	assert.Equal(t, "nearest-neighbor", resp["algorithm"])
}
