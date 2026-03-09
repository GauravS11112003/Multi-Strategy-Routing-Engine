package optimizer

import (
	"math"
	"testing"

	"shipt-route-optimizer/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHaversineDistance(t *testing.T) {
	tests := []struct {
		name     string
		lat1     float64
		lng1     float64
		lat2     float64
		lng2     float64
		expected float64
		delta    float64
	}{
		{
			name: "same point",
			lat1: 33.52, lng1: -86.81,
			lat2: 33.52, lng2: -86.81,
			expected: 0,
			delta:    0.001,
		},
		{
			name: "Birmingham to Hoover AL",
			lat1: 33.5186, lng1: -86.8104,
			lat2: 33.4054, lng2: -86.8114,
			expected: 12.57,
			delta:    1.0,
		},
		{
			name: "short distance within city",
			lat1: 33.520, lng1: -86.810,
			lat2: 33.521, lng2: -86.811,
			expected: 0.14,
			delta:    0.1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dist := HaversineDistance(tt.lat1, tt.lng1, tt.lat2, tt.lng2)
			assert.InDelta(t, tt.expected, dist, tt.delta)
		})
	}
}

func TestHaversineDistanceSymmetric(t *testing.T) {
	d1 := HaversineDistance(33.52, -86.81, 33.40, -86.70)
	d2 := HaversineDistance(33.40, -86.70, 33.52, -86.81)
	assert.InDelta(t, d1, d2, 0.0001, "Haversine should be symmetric")
}

func TestOptimize_EmptyInputs(t *testing.T) {
	tests := []struct {
		name     string
		orders   []models.Order
		shoppers []models.Shopper
	}{
		{"no orders", nil, []models.Shopper{{ID: "s1", Lat: 33.5, Lng: -86.8, Capacity: 5}}},
		{"no shoppers", []models.Order{{ID: "o1", Lat: 33.5, Lng: -86.8, ItemCount: 3}}, nil},
		{"both empty", nil, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assignments, before, after := Optimize(tt.orders, tt.shoppers)
			assert.Empty(t, assignments)
			assert.Equal(t, 0.0, before)
			assert.Equal(t, 0.0, after)
		})
	}
}

func TestOptimize_SingleOrderSingleShopper(t *testing.T) {
	orders := []models.Order{
		{ID: "o1", Lat: 33.52, Lng: -86.81, ItemCount: 3, DeliveryWindow: "morning"},
	}
	shoppers := []models.Shopper{
		{ID: "s1", Lat: 33.50, Lng: -86.80, Capacity: 10},
	}

	assignments, _, after := Optimize(orders, shoppers)

	require.Len(t, assignments, 1)
	assert.Equal(t, "s1", assignments[0].ShopperID)
	assert.Equal(t, []string{"o1"}, assignments[0].Route)
	assert.Greater(t, after, 0.0)
}

func TestOptimize_AssignsToNearestShopper(t *testing.T) {
	orders := []models.Order{
		{ID: "o1", Lat: 33.52, Lng: -86.81, ItemCount: 2},
	}
	shoppers := []models.Shopper{
		{ID: "far", Lat: 34.00, Lng: -87.00, Capacity: 10},
		{ID: "near", Lat: 33.521, Lng: -86.811, Capacity: 10},
	}

	assignments, _, _ := Optimize(orders, shoppers)

	var nearAssignment *models.Assignment
	for _, a := range assignments {
		if a.ShopperID == "near" {
			nearAssignment = &a
			break
		}
	}

	require.NotNil(t, nearAssignment, "Order should be assigned to nearest shopper")
	assert.Contains(t, nearAssignment.Route, "o1")
}

func TestOptimize_RespectsCapacity(t *testing.T) {
	orders := []models.Order{
		{ID: "o1", Lat: 33.52, Lng: -86.81, ItemCount: 2},
		{ID: "o2", Lat: 33.521, Lng: -86.811, ItemCount: 2},
		{ID: "o3", Lat: 33.522, Lng: -86.812, ItemCount: 2},
	}
	shoppers := []models.Shopper{
		{ID: "s1", Lat: 33.52, Lng: -86.81, Capacity: 2},
		{ID: "s2", Lat: 33.55, Lng: -86.85, Capacity: 10},
	}

	assignments, _, _ := Optimize(orders, shoppers)

	totalAssigned := 0
	for _, a := range assignments {
		totalAssigned += len(a.Route)
	}
	assert.Equal(t, 3, totalAssigned, "All orders should be assigned")
}

func TestOptimize_MultipleShoppers(t *testing.T) {
	orders := []models.Order{
		{ID: "o1", Lat: 33.52, Lng: -86.81, ItemCount: 1},
		{ID: "o2", Lat: 33.53, Lng: -86.82, ItemCount: 1},
		{ID: "o3", Lat: 33.40, Lng: -86.70, ItemCount: 1},
		{ID: "o4", Lat: 33.41, Lng: -86.71, ItemCount: 1},
	}
	shoppers := []models.Shopper{
		{ID: "north", Lat: 33.55, Lng: -86.83, Capacity: 10},
		{ID: "south", Lat: 33.39, Lng: -86.69, Capacity: 10},
	}

	assignments, before, after := Optimize(orders, shoppers)

	assert.GreaterOrEqual(t, len(assignments), 1)
	assert.Greater(t, before, 0.0)
	assert.Greater(t, after, 0.0)

	totalAssigned := 0
	for _, a := range assignments {
		totalAssigned += len(a.Route)
	}
	assert.Equal(t, 4, totalAssigned)
}

func TestOptimize_DistanceImproves(t *testing.T) {
	orders := make([]models.Order, 10)
	for i := 0; i < 10; i++ {
		orders[i] = models.Order{
			ID:        "o" + string(rune('A'+i)),
			Lat:       33.5 + float64(i)*0.01,
			Lng:       -86.8 + float64(i)*0.01,
			ItemCount: 2,
		}
	}

	shoppers := []models.Shopper{
		{ID: "s1", Lat: 33.5, Lng: -86.8, Capacity: 5},
		{ID: "s2", Lat: 33.55, Lng: -86.75, Capacity: 5},
	}

	_, before, after := Optimize(orders, shoppers)
	assert.Less(t, after, before, "Optimized distance should be less than random")
}

func TestSortAssignmentsByShopper(t *testing.T) {
	assignments := []models.Assignment{
		{ShopperID: "c", Route: []string{"o3"}},
		{ShopperID: "a", Route: []string{"o1"}},
		{ShopperID: "b", Route: []string{"o2"}},
	}

	SortAssignmentsByShopper(assignments)

	assert.Equal(t, "a", assignments[0].ShopperID)
	assert.Equal(t, "b", assignments[1].ShopperID)
	assert.Equal(t, "c", assignments[2].ShopperID)
}

func TestCalculateRandomDistance(t *testing.T) {
	orders := []models.Order{
		{ID: "o1", Lat: 33.52, Lng: -86.81},
		{ID: "o2", Lat: 33.53, Lng: -86.82},
	}
	shoppers := []models.Shopper{
		{ID: "s1", Lat: 33.50, Lng: -86.80, Capacity: 5},
	}

	dist := calculateRandomDistance(orders, shoppers)
	assert.Greater(t, dist, 0.0)
	assert.False(t, math.IsNaN(dist))
	assert.False(t, math.IsInf(dist, 0))
}
