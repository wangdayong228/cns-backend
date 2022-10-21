package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wangdayong228/cns-backend/controllers"
	"github.com/wangdayong228/cns-backend/utils/ginutils"
)

func SetupRoutes(router *gin.Engine) {
	router.GET("/", indexEndpoint)

	api := router.Group("v0")
	{
		commit := api.Group("commits")
		{
			commit.POST("/", controllers.MakeCommits)
			commit.GET("/:commit_hash", controllers.GetCommit)
			commit.GET("/", controllers.QueryCommits)
		}

		order := api.Group("orders")
		{
			order.POST("/:commit_hash", controllers.MakeOrder)
			order.GET("/:commit_hash", controllers.GetOrder)
		}
	}
}

func indexEndpoint(c *gin.Context) {
	c.JSON(http.StatusOK, ginutils.DataResponse("CNS_BACKEND"))
}
