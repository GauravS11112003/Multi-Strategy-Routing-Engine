package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"shipt-route-optimizer/internal/cache"
	"shipt-route-optimizer/internal/database"
	"shipt-route-optimizer/internal/messaging"
	"shipt-route-optimizer/internal/models"
	"shipt-route-optimizer/internal/optimizer"
	"shipt-route-optimizer/internal/repository"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func main() {
	log.Println("Starting optimization worker...")

	envLocations := []string{".env", "../.env", filepath.Join("backend", ".env")}
	for _, p := range envLocations {
		if err := godotenv.Load(p); err == nil {
			log.Printf("Loaded .env from: %s", p)
			break
		}
	}

	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		if err := database.Connect(dbURL); err != nil {
			log.Printf("WARNING: PostgreSQL unavailable: %v", err)
		} else {
			defer database.Close()
		}
	}

	if redisURL := os.Getenv("REDIS_URL"); redisURL != "" {
		if err := cache.Connect(redisURL); err != nil {
			log.Printf("WARNING: Redis unavailable: %v", err)
		} else {
			defer cache.Close()
		}
	}

	brokersStr := os.Getenv("KAFKA_BROKERS")
	if brokersStr == "" {
		brokersStr = "localhost:9092"
	}
	brokers := strings.Split(brokersStr, ",")

	if err := messaging.ConnectProducer(brokers); err != nil {
		log.Fatalf("Failed to connect Kafka producer: %v", err)
	}
	defer messaging.CloseProducer()

	consumerGroup, err := messaging.NewConsumerGroup(brokers, "optimization-workers")
	if err != nil {
		log.Fatalf("Failed to create consumer group: %v", err)
	}
	defer consumerGroup.Close()

	consumerGroup.RegisterHandler(messaging.TopicOptimizationRequests, handleOptimizationRequest)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigCh
		log.Println("Shutdown signal received")
		cancel()
	}()

	topics := []string{messaging.TopicOptimizationRequests}
	log.Printf("Worker listening on topics: %v", topics)

	if err := consumerGroup.Start(ctx, topics); err != nil {
		if ctx.Err() == nil {
			log.Fatalf("Consumer group error: %v", err)
		}
	}

	log.Println("Worker shutdown complete")
}

func handleOptimizationRequest(event messaging.Event) error {
	dataBytes, err := json.Marshal(event.Data)
	if err != nil {
		return err
	}

	var reqData messaging.OptimizationRequestData
	if err := json.Unmarshal(dataBytes, &reqData); err != nil {
		return err
	}

	log.Printf("Processing optimization job %s (algorithm=%s, orders=%d, shoppers=%d)",
		reqData.JobID, reqData.Algorithm, len(reqData.Orders), len(reqData.Shoppers))

	orders := make([]models.Order, len(reqData.Orders))
	for i, o := range reqData.Orders {
		orders[i] = models.Order{
			ID: o.ID, Lat: o.Lat, Lng: o.Lng,
			ItemCount: o.ItemCount, DeliveryWindow: o.DeliveryWindow,
		}
	}

	shoppers := make([]models.Shopper, len(reqData.Shoppers))
	for i, s := range reqData.Shoppers {
		shoppers[i] = models.Shopper{
			ID: s.ID, Lat: s.Lat, Lng: s.Lng, Capacity: s.Capacity,
		}
	}

	start := time.Now()

	var assignments []models.Assignment
	var distBefore, distAfter float64

	switch reqData.Algorithm {
	case "astar":
		assignments, distBefore, distAfter = optimizer.OptimizeAStar(orders, shoppers)
	default:
		assignments, distBefore, distAfter = optimizer.Optimize(orders, shoppers)
	}

	elapsed := time.Since(start)

	improvementPct := 0.0
	if distBefore > 0 {
		improvementPct = ((distBefore - distAfter) / distBefore) * 100
	}

	assignmentPayloads := make([]messaging.AssignmentPayload, len(assignments))
	for i, a := range assignments {
		assignmentPayloads[i] = messaging.AssignmentPayload{
			ShopperID: a.ShopperID, Route: a.Route, TotalDistance: a.TotalDistance,
		}
	}

	result := messaging.OptimizationResultData{
		JobID:          reqData.JobID,
		Algorithm:      reqData.Algorithm,
		Status:         "completed",
		DistanceBefore: distBefore,
		DistanceAfter:  distAfter,
		ImprovementPct: improvementPct,
		DurationMs:     elapsed.Milliseconds(),
		Assignments:    assignmentPayloads,
	}

	if err := messaging.PublishEvent(messaging.TopicOptimizationResults, messaging.EventOptimizationCompleted, result); err != nil {
		log.Printf("Failed to publish result for job %s: %v", reqData.JobID, err)
	}

	if database.Pool != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		run, err := repository.SaveOptimizationRun(ctx, reqData.Algorithm,
			len(orders), len(shoppers), distBefore, distAfter, improvementPct, int(elapsed.Milliseconds()))
		if err != nil {
			log.Printf("Failed to persist run: %v", err)
		} else {
			for _, a := range assignments {
				for seq, orderID := range a.Route {
					_ = repository.SaveAssignment(ctx, run.ID, a.ShopperID, orderID, seq, a.TotalDistance)
				}
			}
		}
	}

	log.Printf("Completed job %s: %.1f%% improvement in %dms", reqData.JobID, improvementPct, elapsed.Milliseconds())
	return nil
}

func init() {
	_ = uuid.New()
}
