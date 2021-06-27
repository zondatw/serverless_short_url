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

// getDailyReportId return sha1(data)
func getDailyReportId(data []byte) []byte {
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

	// convert []byte to json data
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

	// store access
	var shortHash string = record["ShortHash"].(string)
	accessID := fmt.Sprintf("%x", getAccessId(m.Data))
	access := map[string]interface{}{
		"datetime":  datetime,
		"sourceIp":  record["SourceIp"],
		"agent":     record["Agent"],
		"shortHash": client.Doc(fmt.Sprintf("short-url-map/%s", shortHash)),
	}
	log.Printf("ShortUrlSource: access id: %v, data: %v", accessID, access)
	if _, err := client.Collection("access").Doc(accessID).Set(ctx, access); err != nil {
		log.Printf("ShortUrlSource: add access to firebase error: %v", err)
		return errors.New(fmt.Sprintf("ShortUrlSource: add access to firebase error: %v", err))
	}

	// create daily report doc when it not exist
	dateStr := fmt.Sprintf("%s/%s/%s", datetime.Year(), datetime.Month(), datetime.Day())
	dailyReportId := fmt.Sprintf("%x", getDailyReportId([]byte(dateStr+shortHash)))
	if _, err := client.Collection("daily-report").Doc(dailyReportId).Get(ctx); err != nil {
		dailyReport := map[string]interface{}{
			"shortHash": client.Doc(fmt.Sprintf("short-url-map/%s", shortHash)),
			"datetime":  time.Date(datetime.Year(), datetime.Month(), datetime.Day(), 0, 0, 0, datetime.Nanosecond(), datetime.Location()),
		}
		if _, err := client.Collection("daily-report").Doc(dailyReportId).Set(ctx, dailyReport); err != nil {
			log.Printf("ShortUrlSource: init daily report to firebase error: %v", err)
			return errors.New(fmt.Sprintf("ShortUrlSource: init daily report  to firebase error: %v", err))
		}
	}

	// store daily report
	countDailyReport := []firestore.Update{
		{Path: "count", Value: firestore.Increment(1)},
	}
	log.Printf("ShortUrlSource: datetime: %v, report: %v", datetime, countDailyReport)
	if _, err := client.Collection("daily-report").Doc(dailyReportId).Update(ctx, countDailyReport); err != nil {
		log.Printf("ShortUrlSource: update daily report from firebase error: %v", err)
		return errors.New(fmt.Sprintf("ShortUrlSource: update daily report from firebase error: %v", err))
	}

	// update short url count
	countShortUrl := []firestore.Update{
		{Path: "count", Value: firestore.Increment(1)},
	}
	log.Printf("ShortUrlSource: shortHASH: %v update count", shortHash)
	if _, err := client.Collection("short-url-map").Doc(shortHash).Update(ctx, countShortUrl); err != nil {
		log.Printf("ShortUrlSource: update update short url count from firebase error: %v", err)
		return errors.New(fmt.Sprintf("ShortUrlSource: update short url count from firebase error: %v", err))
	}
	return nil
}
