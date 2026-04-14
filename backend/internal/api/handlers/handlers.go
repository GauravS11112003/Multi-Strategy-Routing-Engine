package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"shipt-route-optimizer/internal/data"
	"shipt-route-optimizer/internal/database"
	"shipt-route-optimizer/internal/models"
	"shipt-route-optimizer/internal/optimizer"
	"shipt-route-optimizer/internal/optimizer/hybrid"
	"shipt-route-optimizer/internal/repository"
	"shipt-route-optimizer/internal/routing"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthCheck returns API health status
func HealthCheck(c *gin.Context) {
	apiKeySet := os.Getenv("OPENROUTE_API_KEY") != ""

	dbStatus := "disconnected"
	if database.Pool != nil {
		if err := database.HealthCheck(c.Request.Context()); err == nil {
			dbStatus = "connected"
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"service":   "shipt-route-optimizer",
		"apiKeySet": apiKeySet,
		"postgres":  dbStatus,
	})
}

// TestRouting tests the OpenRouteService API
func TestRouting(c *gin.Context) {
	// Test route from Birmingham coordinates
	segment, err := routing.GetRoute(33.5200, -86.8100, 33.5186, -86.8104)

	result := gin.H{
		"error":         nil,
		"pointCount":    0,
		"distance":      0.0,
		"duration":      0.0,
		"apiKeySet":     os.Getenv("OPENROUTE_API_KEY") != "",
		"apiKeyLength":  len(os.Getenv("OPENROUTE_API_KEY")),
		"usingFallback": false,
	}

	if err != nil {
		result["error"] = err.Error()
	}

	if segment != nil {
		result["pointCount"] = len(segment.Geometry)
		result["distance"] = segment.Distance
		result["duration"] = segment.Duration
		// If only 2 points, it's using fallback straight line
		result["usingFallback"] = len(segment.Geometry) == 2
		if len(segment.Geometry) <= 5 {
			result["geometry"] = segment.Geometry
		} else {
			result["geometrySample"] = gin.H{
				"first": segment.Geometry[0],
				"last":  segment.Geometry[len(segment.Geometry)-1],
				"total": len(segment.Geometry),
			}
		}
	}

	c.JSON(http.StatusOK, result)
}

// GetSampleData returns mock orders and shoppers
func GetSampleData(c *gin.Context) {
	sampleData := data.GenerateSampleData()
	c.JSON(http.StatusOK, sampleData)
}

// OptimizeRoutes assigns orders to shoppers and optimizes routes
func OptimizeRoutes(c *gin.Context) {
	var req models.OptimizeRequest

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

	// Run optimization algorithm
	assignments, totalBefore, totalAfter := optimizer.Optimize(req.Orders, req.Shoppers)
	optimizer.SortAssignmentsByShopper(assignments)

	response := models.OptimizeResponse{
		Assignments:         assignments,
		TotalDistanceBefore: totalBefore,
		TotalDistanceAfter:  totalAfter,
	}

	c.JSON(http.StatusOK, response)
}

// OptimizeWithAnalytics performs optimization and returns detailed analytics
func OptimizeWithAnalytics(c *gin.Context) {
	var req struct {
		Orders        []models.Order   `json:"orders"`
		Shoppers      []models.Shopper `json:"shoppers"`
		UseRealRoutes bool             `json:"useRealRoutes"`
		Algorithm     string           `json:"algorithm"` // "nearest-neighbor" or "astar"
		ApiKey        string           `json:"apiKey"`    // OpenRouteService API key from frontend
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

	start := time.Now()

	optimizeResponse, analyticsResponse := optimizer.OptimizeWithAnalytics(
		req.Orders,
		req.Shoppers,
		req.UseRealRoutes,
		req.Algorithm,
		req.ApiKey, // Pass API key to optimizer
	)

	elapsed := time.Since(start)

	response := gin.H{
		"optimization": optimizeResponse,
		"analytics":    analyticsResponse,
		"algorithm":    req.Algorithm,
	}

	go persistOptimizationResult(
		req.Algorithm,
		len(req.Orders), len(req.Shoppers),
		optimizeResponse.TotalDistanceBefore,
		optimizeResponse.TotalDistanceAfter,
		optimizeResponse.Assignments,
		elapsed,
	)

	c.JSON(http.StatusOK, response)
}

// HybridSolveStream runs the hybrid solver and streams progress events using NDJSON.
func HybridSolveStream(c *gin.Context) {
	var req models.HybridSolveRequest
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

	writer := c.Writer
	flusher, ok := writer.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Streaming unsupported"})
		return
	}

	writer.Header().Set("Content-Type", "application/x-ndjson")
	writer.Header().Set("Cache-Control", "no-cache")
	writer.Header().Set("Connection", "keep-alive")

	progressCh := make(chan models.HybridProgress, 32)
	resultCh := make(chan models.HybridSolveResponse, 1)
	errorCh := make(chan error, 1)

	ctx := c.Request.Context()

	go func() {
		defer close(progressCh)
		response, err := hybrid.Run(ctx, req.Orders, req.Shoppers, req.Options, func(progress models.HybridProgress) {
			select {
			case progressCh <- progress:
			case <-ctx.Done():
			}
		})
		if err != nil {
			errorCh <- err
			return
		}
		resultCh <- response
	}()

	encoder := json.NewEncoder(writer)

	for {
		select {
		case progress, ok := <-progressCh:
			if !ok {
				progressCh = nil
				continue
			}
			_ = encoder.Encode(gin.H{
				"type": "progress",
				"data": progress,
			})
			flusher.Flush()
		case result := <-resultCh:
			_ = encoder.Encode(gin.H{
				"type": "completed",
				"data": result,
			})
			flusher.Flush()
			return
		case err := <-errorCh:
			_ = encoder.Encode(gin.H{
				"type":  "error",
				"error": err.Error(),
			})
			flusher.Flush()
			return
		case <-ctx.Done():
			_ = encoder.Encode(gin.H{
				"type":  "error",
				"error": "request cancelled",
			})
			flusher.Flush()
			return
		}
	}
}

func persistOptimizationResult(algorithm string, totalOrders, totalShoppers int,
	distanceBefore, distanceAfter float64, assignments []models.Assignment, elapsed time.Duration) {

	if database.Pool == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	improvementPct := 0.0
	if distanceBefore > 0 {
		improvementPct = ((distanceBefore - distanceAfter) / distanceBefore) * 100
	}

	run, err := repository.SaveOptimizationRun(ctx, algorithm, totalOrders, totalShoppers,
		distanceBefore, distanceAfter, improvementPct, int(elapsed.Milliseconds()))
	if err != nil {
		log.Printf("Failed to persist optimization run: %v", err)
		return
	}

	for _, a := range assignments {
		for seq, orderID := range a.Route {
			if err := repository.SaveAssignment(ctx, run.ID, a.ShopperID, orderID, seq, a.TotalDistance); err != nil {
				log.Printf("Failed to persist assignment: %v", err)
			}
		}
	}
}
