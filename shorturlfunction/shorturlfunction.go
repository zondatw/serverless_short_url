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
	"time"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/pubsub"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/gomodule/redigo/redis"
	"github.com/linkedin/goavro"
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

// getFirebaseApp get firbase app
func getFirebaseApp(ctx context.Context, projectID string) (*firebase.App, error) {
	conf := &firebase.Config{ProjectID: projectID}
	return firebase.NewApp(ctx, conf)
}

// initializeFireBase initializes firebase client
func initializeFireBase(ctx context.Context, projectID string) (*firestore.Client, error) {
	app, err := getFirebaseApp(ctx, projectID)
	if err != nil {
		return nil, err
	}

	return app.Firestore(ctx)
}

// initializeAuth initializes Auth client
func initializeAuth(ctx context.Context, projectID string) (*auth.Client, error) {
	app, err := getFirebaseApp(ctx, projectID)
	if err != nil {
		return nil, err
	}

	return app.Auth(ctx)
}

//checkAuth Check ID token
func checkAuth(ctx context.Context, projectID string, idToken string) (*auth.Token, error) {
	auth, err := initializeAuth(ctx, projectID)
	if err != nil {
		return nil, err
	}

	token, err := auth.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, err
	}
	return token, nil
}

// conertToShort use crc32 to get short hash
func convertToShort(url string, extra string) string {
	return fmt.Sprintf("%x%x", crc32.ChecksumIEEE([]byte(url+extra+"AABBCCDD")), crc32.ChecksumIEEE([]byte(url+extra+"ZZXXYYWW")))
}

// storeShortUrlToRedis
func storeShortUrlToRedis(redisConn redis.Conn, shortHash string, url string) error {
	redisConn.Send("MULTI")
	redisConn.Send("SET", shortHash, url)
	redisConn.Send("EXPIRE", shortHash, 30*24*60*60) // expire 30 days
	_, err := redisConn.Do("EXEC")
	return err
}

// sendClientSourceToPub will send source ip and agent to publisher
func sendClientSourceToPub(ctx context.Context, projectID string, shortHash string, sourceIP string, agent string) {
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Printf("sendClientSourceToPub: %v", err)
		return
	}
	defer client.Close()

	topicID := os.Getenv("TOPICID")
	if topicID == "" {
		log.Printf("sendClientSourceToPub: TOPICID must be set")
		return
	}

	record := map[string]interface{}{
		"Datetime":  time.Now().Format(time.RFC3339),
		"SourceIp":  sourceIP,
		"Agent":     agent,
		"ShortHash": shortHash,
	}
	log.Printf("sendClientSourceToPub: original: %v", record)

	codec, err := goavro.NewCodec(AVRO_SOURCE)
	if err != nil {
		log.Printf("sendClientSourceToPub: goavro.NewCodec err: %v", err)
		return
	}

	topic := client.Topic(topicID)

	msg, err := codec.TextualFromNative(nil, record)
	if err != nil {
		log.Printf("sendClientSourceToPub: codec.TextualFromNative err: %v", err)
		return
	}

	result := topic.Publish(ctx, &pubsub.Message{
		Data: msg,
	})
	id, err := result.Get(ctx)
	if err != nil {
		log.Printf("sendClientSourceToPub err: original get: %v", err)
		return
	}
	log.Printf("sendClientSourceToPub: Published message with custom attributes; msg ID: %v\n", id)
}

// Register set new url on redis instance when sign in
// and return short url path
func RegisterWithAuth(res http.ResponseWriter, req *http.Request) {
	authEmail := ""
	// This function only execute on gcp
	if value, exists := os.LookupEnv("ISONGCP"); exists && value == "True" {
		projectID := os.Getenv("PROJECTID")
		if projectID == "" {
			log.Printf("initializeEnvs: PROJECTID must be set")
			http.Error(res, "Error initializing project id", http.StatusInternalServerError)
			return
		}
		ctx := context.Background()
		log.Printf("Auth Token: %v", req.Header.Get("Authorization"))
		token, err := checkAuth(ctx, projectID, req.Header.Get("Authorization"))
		if err != nil {
			log.Printf("check auth error: %v", err)
			http.Error(res, "Auth token error", http.StatusBadRequest)
			return
		}
		authEmail = token.Claims["email"].(string)
	}
	RegisterBase(res, req, true, authEmail)
}

// Register set new url on redis instance
// and return short url path
func Register(res http.ResponseWriter, req *http.Request) {
	RegisterBase(res, req, false, "")
}

// RegisterBase set new url on redis instance
// and when fromAuth is true, it will store auth's email to url info
// and return short url path
func RegisterBase(res http.ResponseWriter, req *http.Request, fromAuth bool, authEmail string) {
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
	shortHash := convertToShort(rg.Url, authEmail)

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
		projectID := os.Getenv("PROJECTID")
		if projectID == "" {
			log.Printf("initializeEnvs: PROJECTID must be set")
			http.Error(res, "Error initializing project id", http.StatusInternalServerError)
			return
		}

		shortUrlData := map[string]interface{}{
			"createdAt": time.Now(),
			"target":    rg.Url,
			"type":      "url",
		}

		if fromAuth {
			shortUrlData["owner"] = authEmail
		}

		ctx := context.Background()
		client, err := initializeFireBase(ctx, projectID)
		if err != nil {
			log.Printf("initializeFireBase: %v", err)
			http.Error(res, "Error initializing firebase", http.StatusInternalServerError)
			return
		}
		if _, err := client.Collection("short-url-map").Doc(shortHash).Set(ctx, shortUrlData); err != nil {
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
	var ctx context.Context
	var projectID string
	var isOnGCP bool = false
	// This function only execute on gcp
	if value, exists := os.LookupEnv("ISONGCP"); exists && value == "True" {
		isOnGCP = true
		ctx = context.Background()
		projectID = os.Getenv("PROJECTID")
		if projectID == "" {
			log.Printf("initializeEnvs: PROJECTID must be set")
			http.Error(res, "Error initializing project id", http.StatusInternalServerError)
			return
		}
	}

	targetUrl, err = redis.String(redisConn.Do("GET", shortHash))
	if err != nil {
		// This function only execute on gcp
		if isOnGCP {
			client, err := initializeFireBase(ctx, projectID)
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
	// This function only execute on gcp
	if isOnGCP {
		// Sending client source to publisher
		ipAdr := req.Header.Get("X-Real-Ip")
		if ipAdr == "" {
			ipAdr = req.Header.Get("X-Forwarded-For")
		}
		if ipAdr == "" {
			ipAdr = req.RemoteAddr
		}
		sendClientSourceToPub(ctx, projectID, shortHash, ipAdr, req.UserAgent())
	}

	log.Printf("Redirect url: %v", targetUrl)
	http.Redirect(res, req, targetUrl, http.StatusSeeOther)
}
