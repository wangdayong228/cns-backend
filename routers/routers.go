package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wangdayong228/cns-backend/controllers"
	"github.com/wangdayong228/cns-backend/utils/ginutils"
)

var (
	regOrderCtrl = controllers.NewRegisterOrderCtrl()
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

		regOrders := api.Group("orders/register")
		{
			regOrders.POST("/:commit_hash", regOrderCtrl.MakeOrder)
			regOrders.GET("/:commit_hash", regOrderCtrl.GetOrder)
			regOrders.PUT("/refresh-url/:commit_hash", regOrderCtrl.RefreshURL)
		}

		renewOrders := api.Group("orders/renew")
		{
			renewOrders.POST("/", nil)
			renewOrders.GET("/:id", nil)
			renewOrders.PUT("/refresh-url/:id", nil)
		}

	}
}

func indexEndpoint(c *gin.Context) {
	c.JSON(http.StatusOK, ginutils.DataResponse("CNS_BACKEND"))
}
