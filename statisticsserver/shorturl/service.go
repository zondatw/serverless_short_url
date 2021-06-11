package shorturl

import (
	"context"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
)

type Service interface {
	GetAll(context *gin.Context)
	Get(context *gin.Context)
}

type shorturlService struct {
	ctx    context.Context
	client *firestore.Client
}

func (service *shorturlService) GetAll(context *gin.Context) {
	context.JSON(http.StatusOK, getAllShortUrlList(service.ctx, service.client))
}

func (service *shorturlService) Get(context *gin.Context) {
	hash := context.Param("hash")
	shortUrlDetail, err := getShortUrlDetail(service.ctx, service.client, hash)
	if err == nil {
		context.JSON(http.StatusOK, shortUrlDetail)
	} else {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}
