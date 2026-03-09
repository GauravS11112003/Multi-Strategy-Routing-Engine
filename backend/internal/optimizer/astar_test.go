package optimizer

import (
	"testing"

	"shipt-route-optimizer/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOptimizeAStar_EmptyInputs(t *testing.T) {
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
			assignments, before, after := OptimizeAStar(tt.orders, tt.shoppers)
			assert.Empty(t, assignments)
			assert.Equal(t, 0.0, before)
			assert.Equal(t, 0.0, after)
		})
	}
}

func TestOptimizeAStar_SingleOrder(t *testing.T) {
	orders := []models.Order{
		{ID: "o1", Lat: 33.52, Lng: -86.81, ItemCount: 3},
	}
	shoppers := []models.Shopper{
		{ID: "s1", Lat: 33.50, Lng: -86.80, Capacity: 10},
	}

	assignments, _, after := OptimizeAStar(orders, shoppers)

	require.Len(t, assignments, 1)
	assert.Equal(t, "s1", assignments[0].ShopperID)
	assert.Equal(t, []string{"o1"}, assignments[0].Route)
	assert.Greater(t, after, 0.0)
}

func TestOptimizeAStar_MultipleOrders(t *testing.T) {
	orders := []models.Order{
		{ID: "o1", Lat: 33.52, Lng: -86.81, ItemCount: 1},
		{ID: "o2", Lat: 33.53, Lng: -86.82, ItemCount: 1},
		{ID: "o3", Lat: 33.54, Lng: -86.83, ItemCount: 1},
	}
	shoppers := []models.Shopper{
		{ID: "s1", Lat: 33.50, Lng: -86.80, Capacity: 10},
	}

	assignments, before, after := OptimizeAStar(orders, shoppers)

	require.Len(t, assignments, 1)
	assert.Len(t, assignments[0].Route, 3)
	assert.Greater(t, before, 0.0)
	assert.Greater(t, after, 0.0)
}

func TestOptimizeRouteAStar_OrdersRouteIsOptimal(t *testing.T) {
	shopper := models.Shopper{ID: "s1", Lat: 33.50, Lng: -86.80, Capacity: 10}

	orders := []models.Order{
		{ID: "far", Lat: 33.55, Lng: -86.85, ItemCount: 1},
		{ID: "near", Lat: 33.501, Lng: -86.801, ItemCount: 1},
		{ID: "mid", Lat: 33.52, Lng: -86.82, ItemCount: 1},
	}

	route := OptimizeRouteAStar(shopper, orders)

	require.Len(t, route, 3)
	assert.Equal(t, "near", route[0].ID, "Nearest order should come first")
}

func TestOptimizeAStar_TwoShoppers(t *testing.T) {
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

	assignments, _, _ := OptimizeAStar(orders, shoppers)

	totalAssigned := 0
	for _, a := range assignments {
		totalAssigned += len(a.Route)
	}
	assert.Equal(t, 4, totalAssigned, "All orders should be assigned")
}

func TestOptimizeAStar_ProducesBetterOrEqualRoute(t *testing.T) {
	orders := []models.Order{
		{ID: "o1", Lat: 33.52, Lng: -86.81, ItemCount: 1},
		{ID: "o2", Lat: 33.53, Lng: -86.82, ItemCount: 1},
		{ID: "o3", Lat: 33.54, Lng: -86.83, ItemCount: 1},
		{ID: "o4", Lat: 33.55, Lng: -86.84, ItemCount: 1},
		{ID: "o5", Lat: 33.56, Lng: -86.85, ItemCount: 1},
	}
	shoppers := []models.Shopper{
		{ID: "s1", Lat: 33.50, Lng: -86.80, Capacity: 10},
	}

	_, _, astarDist := OptimizeAStar(orders, shoppers)
	_, _, greedyDist := Optimize(orders, shoppers)

	assert.LessOrEqual(t, astarDist, greedyDist+0.01,
		"A* should produce routes equal to or better than greedy")
}
