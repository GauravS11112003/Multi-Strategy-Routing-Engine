package messaging

import "time"

const (
	TopicOrderEvents         = "order.events"
	TopicOptimizationRequests = "optimization.requests"
	TopicOptimizationResults  = "optimization.results"
)

type EventType string

const (
	EventOrderCreated          EventType = "order.created"
	EventOrderAssigned         EventType = "order.assigned"
	EventOrderDelivered        EventType = "order.delivered"
	EventOptimizationRequested EventType = "optimization.requested"
	EventOptimizationCompleted EventType = "optimization.completed"
	EventOptimizationFailed    EventType = "optimization.failed"
)

type Event struct {
	ID        string      `json:"id"`
	Type      EventType   `json:"type"`
	Timestamp time.Time   `json:"timestamp"`
	Data      interface{} `json:"data"`
}

type OrderEventData struct {
	OrderID        string  `json:"orderId"`
	Lat            float64 `json:"lat"`
	Lng            float64 `json:"lng"`
	ItemCount      int     `json:"itemCount"`
	DeliveryWindow string  `json:"deliveryWindow"`
	Status         string  `json:"status"`
}

type OptimizationRequestData struct {
	JobID     string                 `json:"jobId"`
	Algorithm string                 `json:"algorithm"`
	Orders    []OrderPayload         `json:"orders"`
	Shoppers  []ShopperPayload       `json:"shoppers"`
}

type OrderPayload struct {
	ID             string  `json:"id"`
	Lat            float64 `json:"lat"`
	Lng            float64 `json:"lng"`
	ItemCount      int     `json:"itemCount"`
	DeliveryWindow string  `json:"deliveryWindow"`
}

type ShopperPayload struct {
	ID       string  `json:"id"`
	Lat      float64 `json:"lat"`
	Lng      float64 `json:"lng"`
	Capacity int     `json:"capacity"`
}

type OptimizationResultData struct {
	JobID           string               `json:"jobId"`
	Algorithm       string               `json:"algorithm"`
	Status          string               `json:"status"`
	DistanceBefore  float64              `json:"distanceBefore"`
	DistanceAfter   float64              `json:"distanceAfter"`
	ImprovementPct  float64              `json:"improvementPct"`
	DurationMs      int64                `json:"durationMs"`
	Assignments     []AssignmentPayload  `json:"assignments"`
	Error           string               `json:"error,omitempty"`
}

type AssignmentPayload struct {
	ShopperID     string   `json:"shopperId"`
	Route         []string `json:"route"`
	TotalDistance  float64  `json:"totalDistance"`
}
