package shorturl

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/api/iterator"
)

type ShortUrl struct {
	Hash      string    `json:"hash" form:"hash" example:"XXXXXX"`
	Target    string    `json:"target" form:"target" example:"http://localhost/"`
	Type      string    `json:"type" form:"type" example:"url"`
	CreatedAt time.Time `json:"createdAt" form:"createdAt" example:"2021-06-26T13:14:47.15739Z"`
}

type shortUrlReport struct {
	Year  int                   `json:"year" form:"year" example:"2021"`
	Month int                   `json:"month" form:"month" example:"06"`
	Dates []ShortUrlDailyReport `json:"dates" form:"dates"`
}

type ShortUrlDailyReport struct {
	Hash  string `json:"hash" form:"hash" example:"XXXXXX"`
	Count int64  `json:"count" form:"count" example:"30"`
	Date  string `json:"date" form:"date" example:"2021-6-30"`
}

type ShortUrlPaginate struct {
	Next   string     `json:"next" form:"next" example:"OOOOOO"`
	Data   []ShortUrl `json:"data" form:"data"`
	Start  string     `json:"start" form:"start" example:"0"`
	Length int        `json:"length" form:"length" example:"5"`
}

type ShortUrlDetail struct {
	Target    string    `json:"target" form:"target" example:"http://localhost/"`
	Count     int       `json:"count" form:"count" example:"30"`
	Type      string    `json:"type" form:"type" example:"url"`
	Owner     string    `json:"owner,omitempty" form:"owner,omitempty" example:"test@test.org"`
	CreatedAt time.Time `json:"createdAt" form:"createdAt" example:"2021-06-26T13:14:47.15739Z"`
}

func getAllShortUrlList(ctx context.Context, client *firestore.Client, authEmail string, start string, length int) ShortUrlPaginate {
	var data []ShortUrl = make([]ShortUrl, 0)
	collect := client.Collection("short-url-map").Where("owner", "==", authEmail)
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
			break
		}

		if detail, err := getShortUrlDetail(ctx, client, authEmail, doc.Ref.ID); err == nil {
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

	next := ""
	if len(data) > 0 {
		next = data[len(data)-1].Hash
	}

	ret := ShortUrlPaginate{
		Next:   next,
		Data:   data,
		Start:  start,
		Length: length,
	}
	return ret
}

func getShortUrlDetail(ctx context.Context, client *firestore.Client, authEmail string, shortUrlHash string) (ShortUrlDetail, error) {
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

	if shortUrlDetail.Owner != authEmail {
		log.Printf("Error: %v: url owner %s, current user: %s\n", "Owner not match", shortUrlDetail.Owner, authEmail)
		return shortUrlDetail, errors.New("Owner not match")
	}
	return shortUrlDetail, nil
}

func getShortUrlDailyReport(ctx context.Context, client *firestore.Client, shortUrlHash string, year int, month int) shortUrlReport {
	// TODO: check auth

	var dates []ShortUrlDailyReport = make([]ShortUrlDailyReport, 0)

	// Init dates
	firstDay := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	lastDay := firstDay.AddDate(0, 1, 0).Add(-time.Nanosecond)
	log.Printf("Day range: %v ~ %v", firstDay, lastDay)
	for day := firstDay.Day(); day <= lastDay.Day(); day++ {
		dates = append(dates, ShortUrlDailyReport{
			Hash:  shortUrlHash,
			Count: 0,
			Date:  fmt.Sprintf("%d-%d-%d", year, month, day),
		})
	}

	collect := client.Collection("daily-report")
	var iter *firestore.DocumentIterator
	iter = collect.
		Where("shortHash", "==", client.Doc(fmt.Sprintf("short-url-map/%s", shortUrlHash))).
		Where("datetime", ">=", firstDay).
		Where("datetime", "<=", lastDay).
		OrderBy("datetime", firestore.Asc).
		Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("Error to iterate: %v", err)
			break
		}
		data := doc.Data()
		index := data["datetime"].(time.Time).Day() - 1
		dates[index].Count = data["count"].(int64)
	}

	ret := shortUrlReport{
		Year:  year,
		Month: month,
		Dates: dates,
	}
	return ret
}
