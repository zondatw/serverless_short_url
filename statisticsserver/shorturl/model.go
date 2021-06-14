package shorturl

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/api/iterator"
)

type ShortUrl struct {
	Hash   string `json:"hash" form:"hash"`
	Target string `json:"target" form:"target"`
	Type   string `json:"type" form:"type"`
}

type ShortUrlDetail struct {
	Target string `json:"target" form:"target"`
	Type   string `json:"type" form:"type"`
	Owner  string `json:"owner,omitempty" form:"owner,omitempty"`
}

func getAllShortUrlList(ctx context.Context, client *firestore.Client) []ShortUrl {
	var ret []ShortUrl = make([]ShortUrl, 0)
	iter := client.Collection("short-url-map").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("Error to iterate: %v", err)
		}

		if detail, err := getShortUrlDetail(ctx, client, doc.Ref.ID); err == nil {
			ret = append(ret, ShortUrl{
				Hash:   doc.Ref.ID,
				Target: detail.Target,
				Type:   detail.Type,
			})
		} else {
			break
		}
	}
	return ret
}

func getShortUrlDetail(ctx context.Context, client *firestore.Client, shortUrlHash string) (ShortUrlDetail, error) {
	var shortUrlDetail ShortUrlDetail
	if result, err := client.Collection("short-url-map").Doc(shortUrlHash).Get(ctx); err == nil {
		if err := mapstructure.Decode(result.Data(), &shortUrlDetail); err != nil {
			log.Printf("Error: %v\n", err)
			return shortUrlDetail, err
		}
	} else {
		log.Printf("Error: %v\n", err)
		return shortUrlDetail, err
	}
	return shortUrlDetail, nil
}
