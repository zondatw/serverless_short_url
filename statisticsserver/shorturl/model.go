package shorturl

import (
	"context"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/api/iterator"
)

type ShortUrl struct {
	Hash      string    `json:"hash" form:"hash"`
	Target    string    `json:"target" form:"target"`
	Type      string    `json:"type" form:"type"`
	CreatedAt time.Time `json:"createdAt" form:"createdAt"`
}

type ShortUrlPaginate struct {
	Next   string     `json:"next" form:"next"`
	Data   []ShortUrl `json:"data" form:"data"`
	Start  string     `json:"start" form:"start"`
	Length int        `json:"length" form:"length"`
}

type ShortUrlDetail struct {
	Target    string    `json:"target" form:"target"`
	Type      string    `json:"type" form:"type"`
	Owner     string    `json:"owner,omitempty" form:"owner,omitempty"`
	CreatedAt time.Time `json:"createdAt" form:"createdAt"`
}

func getAllShortUrlList(ctx context.Context, client *firestore.Client, start string, length int) ShortUrlPaginate {
	var data []ShortUrl = make([]ShortUrl, 0)
	collect := client.Collection("short-url-map")
	var iter *firestore.DocumentIterator
	if start == "" {
		iter = collect.OrderBy(firestore.DocumentID, firestore.Asc).Limit(length).Documents(ctx)
	} else {
		iter = collect.OrderBy(firestore.DocumentID, firestore.Asc).StartAfter(start).Limit(length).Documents(ctx)
	}
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("Error to iterate: %v", err)
		}

		if detail, err := getShortUrlDetail(ctx, client, doc.Ref.ID); err == nil {
			data = append(data, ShortUrl{
				Hash:      doc.Ref.ID,
				Target:    detail.Target,
				Type:      detail.Type,
				CreatedAt: detail.CreatedAt,
			})
		} else {
			break
		}
	}

	ret := ShortUrlPaginate{
		Next:   data[len(data)-1].Hash,
		Data:   data,
		Start:  start,
		Length: length,
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
