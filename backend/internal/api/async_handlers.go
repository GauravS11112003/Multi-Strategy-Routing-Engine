package api

import (
	"net/http"

	"shipt-route-optimizer/internal/messaging"
	"shipt-route-optimizer/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func OptimizeAsync(c *gin.Context) {
	if messaging.Producer == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Kafka not connected, async optimization unavailable"})
		return
	}

	var req struct {
		Orders    []models.Order   `json:"orders"`
		Shoppers  []models.Shopper `json:"shoppers"`
		Algorithm string           `json:"algorithm"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if len(req.Orders) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No orders provided"})
		return
	}
	if len(req.Shoppers) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No shoppers provided"})
		return
	}

	if req.Algorithm == "" {
		req.Algorithm = "nearest-neighbor"
	}

	jobID := uuid.New().String()

	orders := make([]messaging.OrderPayload, len(req.Orders))
	for i, o := range req.Orders {
		orders[i] = messaging.OrderPayload{
			ID: o.ID, Lat: o.Lat, Lng: o.Lng,
			ItemCount: o.ItemCount, DeliveryWindow: o.DeliveryWindow,
		}
	}

	shoppers := make([]messaging.ShopperPayload, len(req.Shoppers))
	for i, s := range req.Shoppers {
		shoppers[i] = messaging.ShopperPayload{
			ID: s.ID, Lat: s.Lat, Lng: s.Lng, Capacity: s.Capacity,
		}
	}

	reqData := messaging.OptimizationRequestData{
		JobID:    jobID,
		Algorithm: req.Algorithm,
		Orders:   orders,
		Shoppers: shoppers,
	}

	if err := messaging.PublishEvent(messaging.TopicOptimizationRequests, messaging.EventOptimizationRequested, reqData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit optimization job"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"jobId":   jobID,
		"status":  "queued",
		"message": "Optimization job submitted. Poll GET /api/optimize-async/:id for results.",
	})
}

func GetAsyncResult(c *gin.Context) {
	jobID := c.Param("id")
	if _, err := uuid.Parse(jobID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job ID"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"jobId":   jobID,
		"status":  "processing",
		"message": "Job is being processed. Check /api/optimizations for completed results.",
	})
}
