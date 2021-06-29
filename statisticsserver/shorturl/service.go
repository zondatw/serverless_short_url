package shorturl

import (
	"context"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
)

type Service interface {
	GetAll(context *gin.Context)
	Get(context *gin.Context)
}

type shorturlService struct {
	ctx    context.Context
	client *firestore.Client
	auth   *auth.Client
}

type GetAllParam struct {
	Start  string `form:"start"`
	Length int    `form:"length"`
}

type GetReportParam struct {
	Year  int `form:"year"`
	Month int `form:"month"`
}

//checkAuth Check ID token
func checkAuth(ctx context.Context, auth *auth.Client, idToken string) (*auth.Token, error) {
	token, err := auth.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (service *shorturlService) GetAll(context *gin.Context) {
	authEmail := ""
	if value := context.GetHeader("Authorization"); value != "" {
		token, err := checkAuth(service.ctx, service.auth, value)
		if err != nil {
			log.Printf("check auth error: %v", err)
			context.JSON(http.StatusUnauthorized, gin.H{"error": "Auth token error"})
			return
		}
		authEmail = token.Claims["email"].(string)
	}

	var param GetAllParam
	context.Bind(&param)
	log.Printf("Start %s, Length %d", param.Start, param.Length)

	if param.Length <= 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Length can't less than 1"})
	} else {
		context.JSON(http.StatusOK, getAllShortUrlList(service.ctx, service.client, authEmail, param.Start, param.Length))
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

func (service *shorturlService) GetDailyReport(context *gin.Context) {
	hash := context.Param("hash")

	var param GetReportParam
	context.Bind(&param)
	log.Printf("Year %d, Month %d", param.Year, param.Month)

	if param.Year <= 0 || param.Month < 1 || param.Month > 12 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Year must be greater than 1 or Month must be in range between 1 ~ 12"})
	} else {
		context.JSON(http.StatusOK, getShortUrlDailyReport(service.ctx, service.client, hash, param.Year, param.Month))
	}
}
