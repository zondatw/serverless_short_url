package shorturlsourcefunction

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/linkedin/goavro"
)

// PubSubMessage is the payload of a Pub/Sub event. Please refer to the docs for
// additional information regarding Pub/Sub events.
type PubSubMessage struct {
	Data []byte `json:"data"`
}

// initializeFireBase initializes firebase client
func initializeFireBase(ctx context.Context, projectID string) (*firestore.Client, error) {
	conf := &firebase.Config{ProjectID: projectID}
	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		return nil, err
	}

	return app.Firestore(ctx)
}

// getAccessId return sha1(data)
func getAccessId(data []byte) []byte {
	h := sha1.New()
	h.Write(data)
	bs := h.Sum(nil)
	return bs
}

// HelloPubSub consumes a Pub/Sub message.
func ShortUrlSource(ctx context.Context, m PubSubMessage) error {
	projectID := os.Getenv("PROJECTID")
	if projectID == "" {
		log.Printf("ShortUrlSource: PROJECTID must be set")
		return errors.New(fmt.Sprintf("ShortUrlSource: PROJECTID must be set"))
	}

	codec, err := goavro.NewCodec(AVRO_SOURCE)
	if err != nil {
		log.Printf("ShortUrlSource: goavro.NewCodec err: %v", err)
		return errors.New(fmt.Sprintf("ShortUrlSource: goavro.NewCodec err: %v", err))
	}

	native, _, err := codec.NativeFromTextual(m.Data)
	if err != nil {
		log.Printf("ShortUrlSource: goavro.NativeFromTextual err: %v", err)
		return errors.New(fmt.Sprintf("ShortUrlSource: goavro.NativeFromTextual err: %v", err))
	}
	log.Printf("ShortUrlSource: native: %v", native)
	record := native.(map[string]interface{})
	log.Printf("ShortUrlSource: record.Datetime: %v", record["Datetime"])
	log.Printf("ShortUrlSource: record.SourceIp: %v", record["SourceIp"])
	log.Printf("ShortUrlSource: record.Agent: %v", record["Agent"])
	log.Printf("ShortUrlSource: record.ShortHash: %v", record["ShortHash"])

	client, err := initializeFireBase(ctx, projectID)
	if err != nil {
		log.Printf("initializeFireBase: %v", err)
		return errors.New(fmt.Sprintf("initializeFireBase: %v", err))
	}

	datetime, err := time.Parse(time.RFC3339, record["Datetime"].(string))
	if err != nil {
		log.Printf("ShortUrlSource: convert string to datetime error: %v", err)
		return errors.New(fmt.Sprintf("ShortUrlSource: convert string to datetime error: %v", err))
	}

	accessID := fmt.Sprintf("%x", getAccessId(m.Data))
	access := map[string]interface{}{
		"datetime":   datetime,
		"source-ip":  record["SourceIp"],
		"agent":      record["Agent"],
		"short-hash": client.Doc(fmt.Sprintf("short-url-map/%s", record["ShortHash"])),
	}
	log.Printf("ShortUrlSource: access id: %v, data: %v", accessID, access)
	if _, err := client.Collection("access").Doc(accessID).Set(ctx, access); err != nil {
		log.Printf("ShortUrlSource: add access to firebase error: %v", err)
		return errors.New(fmt.Sprintf("ShortUrlSource: add access to firebase error: %v", err))
	}
	return nil
}
