package shorturl

import (
	"context"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
)

func Route(ctx context.Context, api *gin.RouterGroup, client *firestore.Client, auth *auth.Client) {
	sS := shorturlService{ctx: ctx, client: client, auth: auth}

	shorturlreportRoute := api.Group("/shorturlreport")
	{
		shorturlreportRoute.GET("/daily/:hash", sS.GetDailyReport)
	}
	shorturlRoute := api.Group("/shorturl")
	{
		shorturlRoute.GET("/", sS.GetAll)
		shorturlRoute.GET("/:hash", sS.Get)
	}
}
