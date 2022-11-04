package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wangdayong228/cns-backend/controllers"
	"github.com/wangdayong228/cns-backend/middlewares"
	"github.com/wangdayong228/cns-backend/utils/ginutils"
)

var (
	regCtrl   = controllers.NewRegisterCtrl()
	renewCtrl = controllers.NewRenewCtrl()
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
		reg.Use(middlewares.Auth())
		{
			reg.POST("/", regCtrl.RegisterByAdmin)
			reg.GET("/:commit_hash", regCtrl.GetRegister)
		}

		regOrders := api.Group("registers/order")
		{
			regOrders.POST("/:commit_hash", regCtrl.MakeOrder)
			regOrders.GET("/:commit_hash", regCtrl.GetOrder)
			regOrders.PUT("/refresh-url/:commit_hash", regCtrl.RefreshURL)
		}

		renew := api.Group("renews")
		renew.Use(middlewares.Auth())
		{
			renew.POST("/", renewCtrl.RenewByAdmin)
			renew.GET("/:commit_hash", renewCtrl.GetRenew)
		}

		renewOrders := api.Group("renews/order")
		{
			renewOrders.POST("/", renewCtrl.MakeOrder)
			renewOrders.GET("/:id", renewCtrl.GetOrder)
			renewOrders.PUT("/refresh-url/:id", renewCtrl.RefreshURL)
		}

	}
}

func indexEndpoint(c *gin.Context) {
	c.JSON(http.StatusOK, ginutils.DataResponse("CNS_BACKEND"))
}
