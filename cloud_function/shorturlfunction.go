package shorturlfunction

import (
	"encoding/json"
	"errors"
	"fmt"
	"hash/crc32"
	"log"
	"net/http"
	"os"

	"github.com/gomodule/redigo/redis"
)

var redisPool *redis.Pool

type registerRequestStruct struct {
	Url string `json:"url"`
}

type registerResponseStruct struct {
	Url string `json:"url"`
}

// initializeRedis initializes and returns a connection pool
func initializeRedis() (*redis.Pool, error) {
	redisHost := os.Getenv("REDISHOST")
	if redisHost == "" {
		return nil, errors.New("REDISHOST must be set")
	}
	redisPort := os.Getenv("REDISPORT")
	if redisPort == "" {
		return nil, errors.New("REDISPORT must be set")
	}
	redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)

	const maxConnections = 10
	return &redis.Pool{
		MaxIdle: maxConnections,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", redisAddr)
			if err != nil {
				return nil, fmt.Errorf("redis.Dial: %v", err)
			}
			return c, err
		},
	}, nil
}

// conertToShort use crc32 to get short hash
func convertToShort(original_url string) string {
	return fmt.Sprintf("%x%x", crc32.ChecksumIEEE([]byte(original_url+"AABBCCDD")), crc32.ChecksumIEEE([]byte(original_url+"ZZXXYYWW")))
}

// Register set new url on redis instance
// and return short url path
func Register(rw http.ResponseWriter, req *http.Request) {
	shortUrlBase := os.Getenv("SHORTURLBASE")
	if shortUrlBase == "" {
		log.Printf("initializeEnvs: SHORTURLBASE must be set")
		http.Error(rw, "Error initializing base url", http.StatusInternalServerError)
		return
	}

	// Initialize connection pool on first invocation
	if redisPool == nil {
		// Pre-declare err to avoid shadowing redisPool
		var err error
		redisPool, err = initializeRedis()
		if err != nil {
			log.Printf("initializeRedis: %v", err)
			http.Error(rw, "Error initializing connection pool", http.StatusInternalServerError)
			return
		}
	}

	redisConn := redisPool.Get()
	defer redisConn.Close()

	// Parse body
	decoder := json.NewDecoder(req.Body)
	var rg registerRequestStruct
	if err := decoder.Decode(&rg); err != nil {
		log.Printf("Parse Register request error: %v", err)
		http.Error(rw, "Error parsing register request", http.StatusBadRequest)
		return
	}

	// Get short url
	var rgw registerResponseStruct
	shortHash := convertToShort(rg.Url)
	rgw.Url = shortUrlBase + shortHash

	// Store to redis
	redisConn.Send("MULTI")
	redisConn.Send("SET", rg.Url, rgw.Url)
	redisConn.Send("EXPIRE", rg.Url, 30*24*60*60) // expire 30 days

	if _, err := redisConn.Do("EXEC"); err != nil {
		log.Printf("Store short url to redis: %v", err)
		http.Error(rw, "Error storing data to redis", http.StatusInternalServerError)
		return
	}

	// Set response block
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(rgw)
}
