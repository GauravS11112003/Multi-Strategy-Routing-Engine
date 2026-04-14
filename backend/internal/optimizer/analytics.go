package optimizer

import (
	"math"
	"shipt-route-optimizer/internal/models"
	"shipt-route-optimizer/internal/routing"
	"sort"
	"time"
)

// OptimizeWithAnalytics performs route optimization and calculates detailed analytics
func OptimizeWithAnalytics(orders []models.Order, shoppers []models.Shopper, useRealRoutes bool, algorithm string, apiKey string) (*models.OptimizeResponse, *models.AnalyticsResponse) {
	var assignments []models.Assignment
	var totalBefore, totalAfter float64

	// Choose algorithm
	switch algorithm {
	case "astar":
		assignments, totalBefore, totalAfter = OptimizeAStar(orders, shoppers)
	default: // "nearest-neighbor" or empty
		assignments, totalBefore, totalAfter = Optimize(orders, shoppers)
	}

	// Calculate analytics (pass API key)
	analytics := calculateAnalytics(orders, shoppers, assignments, useRealRoutes, apiKey)

	response := &models.OptimizeResponse{
		Assignments:         assignments,
		TotalDistanceBefore: totalBefore,
		TotalDistanceAfter:  totalAfter,
	}

	return response, analytics
}

// calculateAnalytics generates comprehensive analytics
func calculateAnalytics(orders []models.Order, shoppers []models.Shopper, assignments []models.Assignment, useRealRoutes bool, apiKey string) *models.AnalyticsResponse {
	shopperAnalytics := calculateShopperAnalytics(orders, shoppers, assignments, useRealRoutes)
	orderAnalytics := calculateOrderAnalytics(orders, assignments)
	systemAnalytics := calculateSystemAnalytics(shoppers, orders, assignments, shopperAnalytics)
	routeGeometries := calculateRouteGeometries(orders, shoppers, assignments, useRealRoutes, apiKey)

	return &models.AnalyticsResponse{
		System:          systemAnalytics,
		Shoppers:        shopperAnalytics,
		Orders:          orderAnalytics,
		RouteGeometries: routeGeometries,
	}
}

// AnalyticsFromAssignments is a helper to compute analytics for externally generated assignments.
func AnalyticsFromAssignments(orders []models.Order, shoppers []models.Shopper, assignments []models.Assignment, useRealRoutes bool, apiKey string) *models.AnalyticsResponse {
	if len(orders) == 0 || len(shoppers) == 0 {
		return &models.AnalyticsResponse{}
	}
	return calculateAnalytics(orders, shoppers, assignments, useRealRoutes, apiKey)
}

// calculateShopperAnalytics generates per-shopper metrics
func calculateShopperAnalytics(orders []models.Order, shoppers []models.Shopper, assignments []models.Assignment, useRealRoutes bool) []models.ShopperAnalytics {
	analytics := []models.ShopperAnalytics{}

	// Create order map for quick lookup
	orderMap := make(map[string]models.Order)
	for _, order := range orders {
		orderMap[order.ID] = order
	}

	// Create shopper map
	shopperMap := make(map[string]models.Shopper)
	for _, shopper := range shoppers {
		shopperMap[shopper.ID] = shopper
	}

	currentTime := time.Now()

	for _, assignment := range assignments {
		shopper := shopperMap[assignment.ShopperID]
		ordersAssigned := len(assignment.Route)

		// Calculate total distance and duration
		totalDistance := assignment.TotalDistance
		totalDuration := 0.0

		if useRealRoutes {
			// Estimate duration based on distance (avg 40 km/h)
			totalDuration = (totalDistance / 40.0) * 60.0
		} else {
			// Fallback estimation
			totalDuration = (totalDistance / 40.0) * 60.0
		}

		// Add time for each delivery (assume 10 min per delivery)
		totalDuration += float64(ordersAssigned) * 10.0

		// Calculate capacity utilization
		capacityUtil := 0.0
		if shopper.Capacity > 0 {
			capacityUtil = (float64(ordersAssigned) / float64(shopper.Capacity)) * 100.0
		}

		// Calculate average distance per order
		avgDistance := 0.0
		if ordersAssigned > 0 {
			avgDistance = totalDistance / float64(ordersAssigned)
		}

		// Calculate efficiency (orders per hour)
		efficiency := 0.0
		if totalDuration > 0 {
			efficiency = (float64(ordersAssigned) / totalDuration) * 60.0
		}

		// Time estimates
		startTime := currentTime.Add(15 * time.Minute) // 15 min prep time
		endTime := startTime.Add(time.Duration(totalDuration) * time.Minute)

		analytics = append(analytics, models.ShopperAnalytics{
			ShopperID:            assignment.ShopperID,
			OrdersAssigned:       ordersAssigned,
			TotalDistance:        totalDistance,
			TotalDuration:        math.Round(totalDuration*10) / 10,
			CapacityUtilization:  math.Round(capacityUtil*10) / 10,
			AverageOrderDistance: math.Round(avgDistance*100) / 100,
			EstimatedStartTime:   startTime.Format("3:04 PM"),
			EstimatedEndTime:     endTime.Format("3:04 PM"),
			Efficiency:           math.Round(efficiency*100) / 100,
		})
	}

	// Sort by shopper ID
	sort.Slice(analytics, func(i, j int) bool {
		return analytics[i].ShopperID < analytics[j].ShopperID
	})

	return analytics
}

// calculateOrderAnalytics generates order-level insights
func calculateOrderAnalytics(orders []models.Order, assignments []models.Assignment) models.OrderAnalytics {
	totalItems := 0
	timeWindowCount := make(map[string]int)

	// Create set of assigned orders
	assignedSet := make(map[string]bool)
	for _, assignment := range assignments {
		for _, orderID := range assignment.Route {
			assignedSet[orderID] = true
		}
	}

	for _, order := range orders {
		totalItems += order.ItemCount
		timeWindowCount[order.DeliveryWindow]++
	}

	avgItemCount := 0.0
	if len(orders) > 0 {
		avgItemCount = float64(totalItems) / float64(len(orders))
	}

	// Calculate average distance between orders
	avgDistance := 0.0
	if len(orders) > 1 {
		totalDist := 0.0
		count := 0
		for i := 0; i < len(orders)-1; i++ {
			for j := i + 1; j < len(orders); j++ {
				dist := HaversineDistance(
					orders[i].Lat, orders[i].Lng,
					orders[j].Lat, orders[j].Lng,
				)
				totalDist += dist
				count++
			}
		}
		if count > 0 {
			avgDistance = totalDist / float64(count)
		}
	}

	return models.OrderAnalytics{
		TotalOrders:         len(orders),
		AverageItemCount:    math.Round(avgItemCount*10) / 10,
		TotalItems:          totalItems,
		OrderDensity:        calculateOrderDensity(orders),
		AverageDistance:     math.Round(avgDistance*100) / 100,
		UnassignedOrders:    len(orders) - len(assignedSet),
		TimeWindowBreakdown: timeWindowCount,
	}
}

// calculateSystemAnalytics generates system-wide metrics
func calculateSystemAnalytics(shoppers []models.Shopper, orders []models.Order, assignments []models.Assignment, shopperAnalytics []models.ShopperAnalytics) models.SystemAnalytics {
	totalDistance := 0.0
	totalDuration := 0.0
	totalEfficiency := 0.0
	activeShoppers := len(assignments)
	assignedOrders := 0

	for _, sa := range shopperAnalytics {
		totalDistance += sa.TotalDistance
		totalDuration += sa.TotalDuration
		totalEfficiency += sa.Efficiency
		assignedOrders += sa.OrdersAssigned
	}

	avgEfficiency := 0.0
	if len(shopperAnalytics) > 0 {
		avgEfficiency = totalEfficiency / float64(len(shopperAnalytics))
	}

	// Calculate optimization score (based on capacity utilization and distance efficiency)
	optimizationScore := calculateOptimizationScore(shoppers, assignments, shopperAnalytics)

	// Estimate fuel cost (assume $0.15 per km)
	fuelCost := totalDistance * 0.15

	// Estimate CO2 saved compared to unoptimized routing (assume 30% savings, 0.2 kg CO2 per km)
	co2Saved := totalDistance * 0.3 * 0.2

	return models.SystemAnalytics{
		TotalShoppers:     len(shoppers),
		ActiveShoppers:    activeShoppers,
		TotalOrders:       len(orders),
		AssignedOrders:    assignedOrders,
		TotalDistance:     math.Round(totalDistance*100) / 100,
		TotalDuration:     math.Round(totalDuration*10) / 10,
		AverageEfficiency: math.Round(avgEfficiency*100) / 100,
		OptimizationScore: math.Round(optimizationScore*10) / 10,
		EstimatedFuelCost: math.Round(fuelCost*100) / 100,
		CO2Saved:          math.Round(co2Saved*100) / 100,
	}
}

// calculateRouteGeometries generates actual road paths for each route
func calculateRouteGeometries(orders []models.Order, shoppers []models.Shopper, assignments []models.Assignment, useRealRoutes bool, apiKey string) []models.RouteGeometry {
	geometries := []models.RouteGeometry{}

	// Create order map
	orderMap := make(map[string]models.Order)
	for _, order := range orders {
		orderMap[order.ID] = order
	}

	// Create shopper map
	shopperMap := make(map[string]models.Shopper)
	for _, shopper := range shoppers {
		shopperMap[shopper.ID] = shopper
	}

	for _, assignment := range assignments {
		shopper := shopperMap[assignment.ShopperID]
		points := [][]float64{}

		// Build waypoints
		waypoints := []routing.RoutePoint{{Lat: shopper.Lat, Lng: shopper.Lng}}
		for _, orderID := range assignment.Route {
			order := orderMap[orderID]
			waypoints = append(waypoints, routing.RoutePoint{Lat: order.Lat, Lng: order.Lng})
		}

		if useRealRoutes && len(waypoints) > 1 {
			// Get real route (with fallback to straight lines)
			for i := 0; i < len(waypoints)-1; i++ {
				segment, err := routing.GetRouteWithKey(
					waypoints[i].Lat, waypoints[i].Lng,
					waypoints[i+1].Lat, waypoints[i+1].Lng,
					apiKey,
				)

				if err == nil && segment != nil {
					// Add route geometry points
					for _, pt := range segment.Geometry {
						points = append(points, []float64{pt.Lat, pt.Lng})
					}
				} else {
					// Fallback to straight line
					points = append(points, []float64{waypoints[i].Lat, waypoints[i].Lng})
					points = append(points, []float64{waypoints[i+1].Lat, waypoints[i+1].Lng})
				}
			}
		} else {
			// Simple straight lines
			for _, wp := range waypoints {
				points = append(points, []float64{wp.Lat, wp.Lng})
			}
		}

		geometries = append(geometries, models.RouteGeometry{
			ShopperID: assignment.ShopperID,
			Points:    points,
		})
	}

	return geometries
}

// calculateOrderDensity calculates orders per square kilometer
func calculateOrderDensity(orders []models.Order) float64 {
	if len(orders) < 2 {
		return 0
	}

	// Find bounding box
	minLat, maxLat := orders[0].Lat, orders[0].Lat
	minLng, maxLng := orders[0].Lng, orders[0].Lng

	for _, order := range orders {
		if order.Lat < minLat {
			minLat = order.Lat
		}
		if order.Lat > maxLat {
			maxLat = order.Lat
		}
		if order.Lng < minLng {
			minLng = order.Lng
		}
		if order.Lng > maxLng {
			maxLng = order.Lng
		}
	}

	// Approximate area in square kilometers
	latDiff := maxLat - minLat
	lngDiff := maxLng - minLng

	// 1 degree ≈ 111 km
	heightKm := latDiff * 111.0
	widthKm := lngDiff * 111.0
	areaKm2 := heightKm * widthKm

	if areaKm2 == 0 {
		return 0
	}

	density := float64(len(orders)) / areaKm2
	return math.Round(density*100) / 100
}

// calculateOptimizationScore calculates a score from 0-100 based on efficiency
func calculateOptimizationScore(shoppers []models.Shopper, assignments []models.Assignment, analytics []models.ShopperAnalytics) float64 {
	if len(analytics) == 0 {
		return 0
	}

	// Factors: capacity utilization, distance efficiency, even distribution
	totalCapUtil := 0.0
	for _, sa := range analytics {
		totalCapUtil += sa.CapacityUtilization
	}
	avgCapUtil := totalCapUtil / float64(len(analytics))

	// Normalize to 0-100 scale
	capacityScore := math.Min(avgCapUtil, 100.0)

	// Check distribution evenness
	distributionScore := 100.0
	if len(analytics) > 1 {
		avgOrders := float64(0)
		for _, sa := range analytics {
			avgOrders += float64(sa.OrdersAssigned)
		}
		avgOrders /= float64(len(analytics))

		variance := 0.0
		for _, sa := range analytics {
			diff := float64(sa.OrdersAssigned) - avgOrders
			variance += diff * diff
		}
		variance /= float64(len(analytics))

		// Lower variance = better distribution
		distributionScore = math.Max(0, 100.0-variance*10.0)
	}

	// Combined score
	score := (capacityScore*0.6 + distributionScore*0.4)
	return math.Min(100.0, score)
}
