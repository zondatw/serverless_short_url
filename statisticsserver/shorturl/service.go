package shorturl

import (
	"context"
	"errors"
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

type ErrorMsg struct {
	Error string `json:"error" form:"v" example:"Error Message"`
}

type GetAllParam struct {
	Start  string `form:"start" example:"0"`  // start of page
	Length int    `form:"length" example:"5"` // length per page
}

type GetReportParam struct {
	Year  int `form:"year" example:"2021"` // year of search
	Month int `form:"month" example:"06"`  // month of search
}

func checkAuth(ctx context.Context, auth *auth.Client, idToken string) (*auth.Token, error) {
	token, err := auth.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func getAuthEmail(context *gin.Context, serviceContext context.Context, auth *auth.Client) (string, error) {
	authEmail := ""
	if value := context.GetHeader("Authorization"); value != "" {
		token, err := checkAuth(serviceContext, auth, value)
		if err != nil {
			log.Printf("check auth error: %v", err)
			return "", errors.New("Auth token error")
		}
		authEmail = token.Claims["email"].(string)
	}
	return authEmail, nil
}

// @Summary Get all short url
// @Tags shorturl
// @version 1.0
// @Description Get all short url
// @produce application/json
// @param object query GetAllParam true "query list"
// @Success 200 {array} ShortUrlPaginate
// @Success 400 {object} ErrorMsg
// @Router /api/shorturl/ [get]
func (service *shorturlService) GetAll(context *gin.Context) {
	authEmail := ""
	if value, err := getAuthEmail(context, service.ctx, service.auth); err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	} else {
		authEmail = value
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

// @Summary Get short url detail
// @Tags shorturl
// @version 1.0
// @Description Get short url detail
// @produce application/json
// @param hash path string true "hash"
// @Success 200 {object} ShortUrlDetail
// @Success 400 {object} ErrorMsg
// @Router /api/shorturl/{hash} [get]
func (service *shorturlService) Get(context *gin.Context) {
	authEmail := ""
	if value, err := getAuthEmail(context, service.ctx, service.auth); err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	} else {
		authEmail = value
	}

	hash := context.Param("hash")
	shortUrlDetail, err := getShortUrlDetail(service.ctx, service.client, authEmail, hash)
	if err == nil {
		context.JSON(http.StatusOK, shortUrlDetail)
	} else {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

// @Summary Get daily report
// @Tags shorturlreport
// @version 1.0
// @Description Get daily report
// @accept application/json
// @produce application/json
// @param hash path string true "hash"
// @param object query GetReportParam true "query list"
// @Success 200 {object} shortUrlReport
// @Success 400 {object} ErrorMsg
// @Router /api/shorturlreport/daily/{hash} [get]
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
