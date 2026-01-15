package cache

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	redisURL   string
	redisToken string
	httpClient *http.Client
)

// InitRedis initializes connection using Upstash REST API
func InitRedis() error {
	redisURL = os.Getenv("UPSTASH_REDIS_REST_URL")
	redisToken = os.Getenv("UPSTASH_REDIS_REST_TOKEN")

	if redisURL == "" || redisToken == "" {
		log.Println("Warning: Redis credentials not found, caching disabled")
		return nil
	}

	// Clean up URL for REST API - ensure https://
	if !strings.HasPrefix(redisURL, "https://") {
		redisURL = strings.TrimPrefix(redisURL, "rediss://")
		redisURL = strings.TrimPrefix(redisURL, "redis://")
		if idx := strings.Index(redisURL, ":"); idx != -1 {
			redisURL = redisURL[:idx]
		}
		redisURL = "https://" + redisURL
	}

	log.Printf("Using Upstash REST API at: %s", redisURL)

	httpClient = &http.Client{Timeout: 10 * time.Second}

	// Test connection with PING
	if err := ping(); err != nil {
		return err
	}

	log.Println("Redis connected successfully (via REST API)")
	return nil
}

// Execute Redis command via REST API
func command(args ...string) (interface{}, error) {
	if httpClient == nil {
		return nil, nil
	}

	body, err := json.Marshal(args)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", redisURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+redisToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("redis error: %s", string(respBody))
	}

	var result struct {
		Result interface{} `json:"result"`
		Error  string      `json:"error,omitempty"`
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, err
	}

	if result.Error != "" {
		return nil, fmt.Errorf("redis error: %s", result.Error)
	}

	return result.Result, nil
}

func ping() error {
	result, err := command("PING")
	if err != nil {
		return err
	}
	if result != "PONG" {
		return fmt.Errorf("unexpected PING response: %v", result)
	}
	return nil
}

// IsEnabled returns true if Redis is configured and connected
func IsEnabled() bool {
	return httpClient != nil
}

// Get retrieves a string value from cache
func Get(key string) (string, error) {
	if httpClient == nil {
		return "", nil
	}

	result, err := command("GET", key)
	if err != nil {
		return "", err
	}

	if result == nil {
		return "", nil // Cache miss
	}

	if str, ok := result.(string); ok {
		return str, nil
	}

	return "", nil
}

// Set stores a string value in cache with TTL in seconds
func Set(key string, value string, ttlSeconds int) error {
	if httpClient == nil {
		return nil
	}

	if ttlSeconds > 0 {
		_, err := command("SET", key, value, "EX", strconv.Itoa(ttlSeconds))
		return err
	}

	_, err := command("SET", key, value)
	return err
}

// Delete removes a key from cache
func Delete(key string) error {
	if httpClient == nil {
		return nil
	}

	_, err := command("DEL", key)
	return err
}

// GetOrFetch tries cache first, otherwise calls fetch function and caches result
func GetOrFetch(key string, ttlSeconds int, fetch func() (string, error)) (string, error) {
	// Try cache first
	cached, err := Get(key)
	if err != nil {
		log.Printf("Cache get error for %s: %v", key, err)
	}
	if cached != "" {
		log.Printf("Cache hit: %s", key)
		return cached, nil
	}

	// Cache miss - fetch fresh data
	log.Printf("Cache miss: %s", key)
	data, err := fetch()
	if err != nil {
		return "", err
	}

	// Cache the result
	if err := Set(key, data, ttlSeconds); err != nil {
		log.Printf("Cache set error for %s: %v", key, err)
		// Don't fail - we still have the data
	}

	return data, nil
}

// Cache keys for NASCAR pool
const (
	KeyStandings = "nascar:standings"
	KeyRaces     = "nascar:races"
	KeyDrivers   = "nascar:drivers"
)

// Cache TTLs in seconds
const (
	TTLStandings = 60    // 1 minute - changes after results
	TTLRaces     = 300   // 5 minutes
	TTLDrivers   = 86400 // 24 hours - rarely changes
)
