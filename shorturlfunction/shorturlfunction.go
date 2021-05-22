package shorturlfunction

import (
	"encoding/json"
	"errors"
	"fmt"
	"hash/crc32"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

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
func convertToShort(originalUrl string) string {
	return fmt.Sprintf("%x%x", crc32.ChecksumIEEE([]byte(originalUrl+"AABBCCDD")), crc32.ChecksumIEEE([]byte(originalUrl+"ZZXXYYWW")))
}

// Register set new url on redis instance
// and return short url path
func Register(res http.ResponseWriter, req *http.Request) {
	shortUrlBase := os.Getenv("SHORTURLBASE")
	if shortUrlBase == "" {
		log.Printf("initializeEnvs: SHORTURLBASE must be set")
		http.Error(res, "Error initializing base url", http.StatusInternalServerError)
		return
	}

	// Initialize connection pool on first invocation
	if redisPool == nil {
		// Pre-declare err to avoid shadowing redisPool
		var err error
		redisPool, err = initializeRedis()
		if err != nil {
			log.Printf("initializeRedis: %v", err)
			http.Error(res, "Error initializing connection pool", http.StatusInternalServerError)
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
		http.Error(res, "Error parsing register request", http.StatusBadRequest)
		return
	}

	// Get short url
	var rgw registerResponseStruct
	shortHash := convertToShort(rg.Url)

	if u, err := url.Parse(shortUrlBase); err != nil {
		log.Printf("SHORTURLBASE value: %v", err)
		http.Error(res, "Error SHORTURLBASE", http.StatusInternalServerError)
		return
	} else {
		u.Path = path.Join(u.Path, shortHash)
		rgw.Url = u.String()
	}

	// Store to redis
	redisConn.Send("MULTI")
	redisConn.Send("SET", shortHash, rg.Url)
	redisConn.Send("EXPIRE", rg.Url, 30*24*60*60) // expire 30 days

	if _, err := redisConn.Do("EXEC"); err != nil {
		log.Printf("Store short url to redis: %v", err)
		http.Error(res, "Error storing data to redis", http.StatusInternalServerError)
		return
	}

	// Set response block
	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(rgw)
}

// Redirect direct user to other site with short url path
func Redirect(res http.ResponseWriter, req *http.Request) {
	// Initialize connection pool on first invocation
	if redisPool == nil {
		// Pre-declare err to avoid shadowing redisPool
		var err error
		redisPool, err = initializeRedis()
		if err != nil {
			log.Printf("initializeRedis: %v", err)
			http.Error(res, "Error initializing connection pool", http.StatusInternalServerError)
			return
		}
	}

	redisConn := redisPool.Get()
	defer redisConn.Close()

	// Search original url from path
	// if original url exist
	// 	return Redirect response
	// else
	// 	return 404
	ss := strings.Split(req.URL.Path, "/")
	shortHash := ss[len(ss)-1]
	if originalUrl, err := redis.String(redisConn.Do("GET", shortHash)); err != nil {
		log.Printf("Redirect key not exist: %v", err)
		http.Error(res, "Error redirect path", http.StatusNotFound)
		return
	} else {
		log.Printf("Redirect url: %v", originalUrl)
		http.Redirect(res, req, originalUrl, http.StatusSeeOther)
	}
}
