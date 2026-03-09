package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"shipt-route-optimizer/internal/api"
	"shipt-route-optimizer/internal/cache"
	"shipt-route-optimizer/internal/database"
	"shipt-route-optimizer/internal/messaging"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	envLocations := []string{
		".env",
		"../.env",
		filepath.Join("backend", ".env"),
	}

	for _, envPath := range envLocations {
		if err := godotenv.Load(envPath); err == nil {
			log.Printf("Loaded .env file from: %s\n", envPath)
			break
		}
	}

	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		if err := database.Connect(dbURL); err != nil {
			log.Printf("WARNING: PostgreSQL unavailable: %v (running without persistence)", err)
		} else {
			defer database.Close()

			migrationPaths := []string{
				"internal/database/migrations",
				"../internal/database/migrations",
				filepath.Join("backend", "internal", "database", "migrations"),
			}
			for _, mp := range migrationPaths {
				if _, err := os.Stat(mp); err == nil {
					if err := database.RunMigrations(dbURL, mp); err != nil {
						log.Printf("WARNING: Migration failed: %v", err)
					}
					break
				}
			}
		}
	} else {
		log.Println("DATABASE_URL not set, running without persistence")
	}

	if redisURL := os.Getenv("REDIS_URL"); redisURL != "" {
		if err := cache.Connect(redisURL); err != nil {
			log.Printf("WARNING: Redis unavailable: %v (running without cache)", err)
		} else {
			defer cache.Close()
		}
	} else {
		log.Println("REDIS_URL not set, running without cache")
	}

	if brokersStr := os.Getenv("KAFKA_BROKERS"); brokersStr != "" {
		brokers := strings.Split(brokersStr, ",")
		if err := messaging.ConnectProducer(brokers); err != nil {
			log.Printf("WARNING: Kafka unavailable: %v (running without messaging)", err)
		} else {
			defer messaging.CloseProducer()
		}
	} else {
		log.Println("KAFKA_BROKERS not set, running without messaging")
	}

	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173", "http://localhost:5174", "http://localhost:3000", "http://localhost:80", "http://localhost"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept"}
	r.Use(cors.New(config))

	apiGroup := r.Group("/api")
	{
		apiGroup.GET("/health", api.HealthCheck)
		apiGroup.GET("/test-routing", api.TestRouting)
		apiGroup.GET("/sample-data", api.GetSampleData)
		apiGroup.POST("/optimize", api.OptimizeRoutes)
		apiGroup.POST("/optimize-analytics", api.OptimizeWithAnalytics)
		apiGroup.POST("/optimize-hybrid-stream", api.HybridSolveStream)

		apiGroup.GET("/orders", api.ListOrders)
		apiGroup.POST("/orders", api.CreateOrder)
		apiGroup.GET("/orders/:id", api.GetOrder)
		apiGroup.DELETE("/orders/:id", api.DeleteOrder)

		apiGroup.GET("/shoppers", api.ListShoppers)
		apiGroup.POST("/shoppers", api.CreateShopper)
		apiGroup.GET("/shoppers/:id", api.GetShopper)
		apiGroup.DELETE("/shoppers/:id", api.DeleteShopper)

		apiGroup.GET("/optimizations", api.ListOptimizations)
		apiGroup.GET("/optimizations/:id", api.GetOptimization)

		apiGroup.POST("/optimize-async", api.OptimizeAsync)
		apiGroup.GET("/optimize-async/:id", api.GetAsyncResult)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Multi-Strategy Routing Engine Backend starting on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
