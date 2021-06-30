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
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/zondatw/serverless_short_url/statisticsserver/docs"
	"google.golang.org/api/option"
)

var router *gin.Engine

// @title Swagger API
// @version 1.0
// @description Gin swagger.

// @contact.name API Support
// @contact.url https://github.com/zondatw/serverless_short_url

// @host localhost
// schemes http
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

	if mode := gin.Mode(); mode == gin.DebugMode {
		url := ginSwagger.URL(fmt.Sprintf("http://localhost:%s/swagger/doc.json", port))
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	}
	router.Run(urlSetting)
}

// @Summary Health
// @Tags Base
// @version 1.0
// @produce text/plain
// @Success 200 {string} json "{"status": "OK"}"
// @Router /health [get]
func health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}
