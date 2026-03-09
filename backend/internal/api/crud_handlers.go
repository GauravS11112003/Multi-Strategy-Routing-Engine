package api

import (
	"net/http"
	"strconv"

	"shipt-route-optimizer/internal/database"
	"shipt-route-optimizer/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func paginationParams(c *gin.Context) (int, int) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	return limit, offset
}

func requireDB(c *gin.Context) bool {
	if database.Pool == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Database not connected"})
		return false
	}
	return true
}

// --- Order handlers ---

func CreateOrder(c *gin.Context) {
	if !requireDB(c) {
		return
	}

	var req struct {
		Lat            float64 `json:"lat" binding:"required"`
		Lng            float64 `json:"lng" binding:"required"`
		ItemCount      int     `json:"itemCount" binding:"required"`
		DeliveryWindow string  `json:"deliveryWindow"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	order, err := repository.CreateOrder(c.Request.Context(), req.Lat, req.Lng, req.ItemCount, req.DeliveryWindow)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

func GetOrder(c *gin.Context) {
	if !requireDB(c) {
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	order, err := repository.GetOrderByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if order == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, order)
}

func ListOrders(c *gin.Context) {
	if !requireDB(c) {
		return
	}

	limit, offset := paginationParams(c)
	status := c.Query("status")

	orders, total, err := repository.ListOrders(c.Request.Context(), status, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"orders": orders,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

func DeleteOrder(c *gin.Context) {
	if !requireDB(c) {
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	if err := repository.DeleteOrder(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"deleted": true})
}

// --- Shopper handlers ---

func CreateShopper(c *gin.Context) {
	if !requireDB(c) {
		return
	}

	var req struct {
		Name     string  `json:"name" binding:"required"`
		Lat      float64 `json:"lat" binding:"required"`
		Lng      float64 `json:"lng" binding:"required"`
		Capacity int     `json:"capacity" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	shopper, err := repository.CreateShopper(c.Request.Context(), req.Name, req.Lat, req.Lng, req.Capacity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, shopper)
}

func GetShopper(c *gin.Context) {
	if !requireDB(c) {
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid shopper ID"})
		return
	}

	shopper, err := repository.GetShopperByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if shopper == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Shopper not found"})
		return
	}

	c.JSON(http.StatusOK, shopper)
}

func ListShoppers(c *gin.Context) {
	if !requireDB(c) {
		return
	}

	limit, offset := paginationParams(c)
	status := c.Query("status")

	shoppers, total, err := repository.ListShoppers(c.Request.Context(), status, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"shoppers": shoppers,
		"total":    total,
		"limit":    limit,
		"offset":   offset,
	})
}

func DeleteShopper(c *gin.Context) {
	if !requireDB(c) {
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid shopper ID"})
		return
	}

	if err := repository.DeleteShopper(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"deleted": true})
}

// --- Optimization history handlers ---

func ListOptimizations(c *gin.Context) {
	if !requireDB(c) {
		return
	}

	limit, offset := paginationParams(c)

	runs, total, err := repository.ListOptimizationRuns(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"optimizations": runs,
		"total":         total,
		"limit":         limit,
		"offset":        offset,
	})
}

func GetOptimization(c *gin.Context) {
	if !requireDB(c) {
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid optimization ID"})
		return
	}

	run, err := repository.GetOptimizationRun(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if run == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Optimization run not found"})
		return
	}

	assignments, err := repository.GetRunAssignments(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"run":         run,
		"assignments": assignments,
	})
}
