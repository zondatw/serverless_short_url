package shorturl

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
)

func Route(ctx context.Context, api *gin.RouterGroup, client *firestore.Client) {
	sS := shorturlService{ctx: ctx, client: client}

	shorturlRoute := api.Group("/shorturl")
	{
		shorturlRoute.GET("/", sS.GetAll)
		shorturlRoute.GET("/:hash", sS.Get)
	}
}
