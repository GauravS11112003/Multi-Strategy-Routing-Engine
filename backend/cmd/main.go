package main

import (
	"log"
	"os"
	"path/filepath"

	"shipt-route-optimizer/internal/api"
	"shipt-route-optimizer/internal/database"

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

	r := api.SetupRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Multi-Strategy Routing Engine Backend starting on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
