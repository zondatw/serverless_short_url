package shorturl

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/api/iterator"
)

type ShortUrl struct {
	Target string `json:"target" form:"target"`
	Type   string `json:"type" form:"type"`
	Owner  string `json:"owner,omitempty" form:"owner,omitempty"`
}

func getAllShortUrlList(ctx context.Context, client *firestore.Client) []string {
	var ret []string = make([]string, 0)
	iter := client.Collection("short-url-map").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("Error to iterate: %v", err)
		}
		ret = append(ret, doc.Ref.ID)
	}
	return ret
}

func getShortUrlDetail(ctx context.Context, client *firestore.Client, shortUrlHash string) (ShortUrl, error) {
	var shortUrl ShortUrl
	if result, err := client.Collection("short-url-map").Doc(shortUrlHash).Get(ctx); err == nil {
		if err := mapstructure.Decode(result.Data(), &shortUrl); err != nil {
			log.Printf("Error: %v\n", err)
			return shortUrl, err
		}
	} else {
		log.Printf("Error: %v\n", err)
		return shortUrl, err
	}
	return shortUrl, nil
}
