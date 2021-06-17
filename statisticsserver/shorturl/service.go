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

type GetAllParam struct {
	Start  int `form:"start"`
	Length int `form:"length"`
}

func (service *shorturlService) GetAll(context *gin.Context) {
	var param GetAllParam
	context.Bind(&param)
	if param.Start < 0 || param.Length < 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Values can't less than 0"})
	} else {
		context.JSON(http.StatusOK, getAllShortUrlList(service.ctx, service.client, param.Start, param.Length))
	}
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
