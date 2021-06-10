package main

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/zondatw/serverless_short_url/statisticsserver/shorturl"
)

func initRoute(ctx context.Context, client *firestore.Client) {
	api := router.Group("/api")
	shorturl.Route(ctx, api, client)
}
