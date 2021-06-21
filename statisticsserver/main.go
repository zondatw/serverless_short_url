package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
)

var router *gin.Engine

func main() {
	var (
		ip   = "0.0.0.0"
		port = "80"
	)

	// Use a service account
	ctx := context.Background()
	var app *firebase.App
	var auth *auth.Client
	var err error
	if projectID, ok := os.LookupEnv("projectID"); ok {
		fmt.Printf("On GCP: %v\n", projectID)
		conf := &firebase.Config{ProjectID: projectID}
		app, err = firebase.NewApp(ctx, conf)
	} else {
		fmt.Printf("On Local")
		sa := option.WithCredentialsFile("keys/key.json")
		app, err = firebase.NewApp(ctx, nil, sa)
	}
	if err != nil {
		log.Fatalln(err)
	}

	if auth, err = app.Auth(ctx); err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	urlSetting := fmt.Sprintf("%s:%s", ip, port)
	router = gin.Default()
	router.GET("/", health)
	router.GET("/health", health)
	initRoute(ctx, client, auth)
	router.Run(urlSetting)
}

func health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}
