package main

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/zondatw/serverless_short_url/statisticsserver/middleware"
	"github.com/zondatw/serverless_short_url/statisticsserver/shorturl"
)

func initRoute(ctx context.Context, client *firestore.Client) {
	router.Use(middleware.CORSMiddleware)

	api := router.Group("/api")
	shorturl.Route(ctx, api, client)
}
