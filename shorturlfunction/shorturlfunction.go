package shorturlfunction

import (
	"context"
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

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
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

// initializeFireBase initializes firebase client
func initializeFireBase(ctx context.Context) (*firestore.Client, error) {
	projectID := os.Getenv("PROJECTID")
	if projectID == "" {
		return nil, errors.New("PROJECTID must be set")
	}

	conf := &firebase.Config{ProjectID: projectID}
	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		return nil, err
	}

	return app.Firestore(ctx)
}

// conertToShort use crc32 to get short hash
func convertToShort(url string) string {
	return fmt.Sprintf("%x%x", crc32.ChecksumIEEE([]byte(url+"AABBCCDD")), crc32.ChecksumIEEE([]byte(url+"ZZXXYYWW")))
}

// storeShortUrlToRedis
func storeShortUrlToRedis(redisConn redis.Conn, shortHash string, url string) error {
	redisConn.Send("MULTI")
	redisConn.Send("SET", shortHash, url)
	redisConn.Send("EXPIRE", shortHash, 30*24*60*60) // expire 30 days
	_, err := redisConn.Do("EXEC")
	return err
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

	// Store to firebase
	// This function only execute on gcp
	if value, exists := os.LookupEnv("ISONGCP"); exists && value == "True" {
		ctx := context.Background()
		client, err := initializeFireBase(ctx)
		if err != nil {
			log.Printf("initializeFireBase: %v", err)
			http.Error(res, "Error initializing firebase", http.StatusInternalServerError)
			return
		}
		if _, err := client.Collection("short-url-map").Doc(shortHash).Set(ctx, map[string]interface{}{
			"target": rg.Url,
			"type":   "url",
		}); err != nil {
			log.Printf("Add short url to firebase error: %v", err)
			http.Error(res, "Error store data to firebase", http.StatusInternalServerError)
			return
		}
	}

	// Store to redis
	if err := storeShortUrlToRedis(redisConn, shortHash, rg.Url); err != nil {
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
	var targetUrl string = ""
	var err error = nil
	targetUrl, err = redis.String(redisConn.Do("GET", shortHash))
	if err != nil {
		// This function only execute on gcp
		if value, exists := os.LookupEnv("ISONGCP"); exists && value == "True" {
			ctx := context.Background()
			client, err := initializeFireBase(ctx)
			if err != nil {
				log.Printf("initializeFireBase: %v", err)
				http.Error(res, "Error initializing firebase", http.StatusInternalServerError)
				return
			}
			if result, err := client.Collection("short-url-map").Doc(shortHash).Get(ctx); err == nil {
				targetUrl = result.Data()["target"].(string)
				// Store to redis
				if err := storeShortUrlToRedis(redisConn, shortHash, targetUrl); err != nil {
					log.Printf("Store short url to redis: %v", err)
					http.Error(res, "Error storing data to redis", http.StatusInternalServerError)
					return
				}
			}
		}
		if targetUrl == "" {
			log.Printf("Redirect key not exist: %v", err)
			http.Error(res, "Error redirect path", http.StatusNotFound)
			return
		}
	}

	log.Printf("Redirect url: %v", targetUrl)
	http.Redirect(res, req, targetUrl, http.StatusSeeOther)
}