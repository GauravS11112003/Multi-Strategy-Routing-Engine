package cache

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var Client *redis.Client

func Connect(redisURL string) error {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return fmt.Errorf("parse redis URL: %w", err)
	}

	opts.PoolSize = 20
	opts.MinIdleConns = 5

	Client = redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := Client.Ping(ctx).Err(); err != nil {
		Client = nil
		return fmt.Errorf("ping redis: %w", err)
	}

	log.Println("Connected to Redis")
	return nil
}

func Close() {
	if Client != nil {
		Client.Close()
		log.Println("Redis connection closed")
	}
}

func HealthCheck(ctx context.Context) error {
	if Client == nil {
		return fmt.Errorf("redis not connected")
	}
	return Client.Ping(ctx).Err()
}

// --- Distance cache ---

func distanceKey(lat1, lng1, lat2, lng2 float64) string {
	return fmt.Sprintf("dist:%.6f,%.6f:%.6f,%.6f", lat1, lng1, lat2, lng2)
}

func GetCachedDistance(ctx context.Context, lat1, lng1, lat2, lng2 float64) (float64, bool) {
	if Client == nil {
		return 0, false
	}

	val, err := Client.Get(ctx, distanceKey(lat1, lng1, lat2, lng2)).Float64()
	if err != nil {
		return 0, false
	}
	return val, true
}

func SetCachedDistance(ctx context.Context, lat1, lng1, lat2, lng2, distance float64) {
	if Client == nil {
		return
	}

	Client.Set(ctx, distanceKey(lat1, lng1, lat2, lng2), distance, 24*time.Hour)
}

// --- Optimization result cache ---

func optimizationCacheKey(input interface{}) string {
	data, _ := json.Marshal(input)
	hash := sha256.Sum256(data)
	return fmt.Sprintf("opt:%x", hash[:16])
}

func GetCachedOptimization(ctx context.Context, input interface{}) ([]byte, bool) {
	if Client == nil {
		return nil, false
	}

	key := optimizationCacheKey(input)
	val, err := Client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, false
	}
	return val, true
}

func SetCachedOptimization(ctx context.Context, input interface{}, result interface{}) {
	if Client == nil {
		return
	}

	data, err := json.Marshal(result)
	if err != nil {
		return
	}

	key := optimizationCacheKey(input)
	Client.Set(ctx, key, data, 10*time.Minute)
}

// --- Rate limiting ---

func CheckRateLimit(ctx context.Context, identifier string, maxRequests int, window time.Duration) (bool, int, error) {
	if Client == nil {
		return true, maxRequests, nil
	}

	key := fmt.Sprintf("ratelimit:%s", identifier)

	pipe := Client.Pipeline()
	incr := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, window)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return true, maxRequests, err
	}

	count := int(incr.Val())
	remaining := maxRequests - count
	if remaining < 0 {
		remaining = 0
	}

	return count <= maxRequests, remaining, nil
}
