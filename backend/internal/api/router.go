package api

import (
	"shipt-route-optimizer/internal/api/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRouter configures the gin engine and registers all routes
func SetupRouter() *gin.Engine {
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173", "http://localhost:5174", "http://localhost:3000", "http://localhost:80", "http://localhost"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept"}
	r.Use(cors.New(config))

	apiGroup := r.Group("/api")
	{
		apiGroup.GET("/health", handlers.HealthCheck)
		apiGroup.GET("/test-routing", handlers.TestRouting)
		apiGroup.GET("/sample-data", handlers.GetSampleData)
		apiGroup.POST("/optimize", handlers.OptimizeRoutes)
		apiGroup.POST("/optimize-analytics", handlers.OptimizeWithAnalytics)
		apiGroup.POST("/optimize-hybrid-stream", handlers.HybridSolveStream)

		apiGroup.GET("/orders", handlers.ListOrders)
		apiGroup.POST("/orders", handlers.CreateOrder)
		apiGroup.GET("/orders/:id", handlers.GetOrder)
		apiGroup.DELETE("/orders/:id", handlers.DeleteOrder)

		apiGroup.GET("/shoppers", handlers.ListShoppers)
		apiGroup.POST("/shoppers", handlers.CreateShopper)
		apiGroup.GET("/shoppers/:id", handlers.GetShopper)
		apiGroup.DELETE("/shoppers/:id", handlers.DeleteShopper)

		apiGroup.GET("/optimizations", handlers.ListOptimizations)
		apiGroup.GET("/optimizations/:id", handlers.GetOptimization)
	}

	return r
}
