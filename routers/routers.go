package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wangdayong228/cns-backend/controllers"
	"github.com/wangdayong228/cns-backend/utils/ginutils"
)

var (
	regOrderCtrl   = controllers.NewRegisterOrderCtrl()
	renewOrderCtrl = controllers.NewRenewOrderCtrl()
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

		reg := api.Group("registers")
		{
			reg.POST("/", nil)
			reg.GET("/", nil)
			regOrders := reg.Group("order")
			{
				regOrders.POST("/:commit_hash", regOrderCtrl.MakeOrder)
				regOrders.GET("/:commit_hash", regOrderCtrl.GetOrder)
				regOrders.PUT("/refresh-url/:commit_hash", regOrderCtrl.RefreshURL)
			}
		}

		renew := api.Group("renews")
		{
			renew.POST("/", nil)
			renewOrders := renew.Group("order")
			{
				renewOrders.POST("/", renewOrderCtrl.MakeOrder)
				renewOrders.GET("/:id", renewOrderCtrl.GetOrder)
				renewOrders.PUT("/refresh-url/:id", renewOrderCtrl.RefreshURL)
			}
		}

	}
}

func indexEndpoint(c *gin.Context) {
	c.JSON(http.StatusOK, ginutils.DataResponse("CNS_BACKEND"))
}
