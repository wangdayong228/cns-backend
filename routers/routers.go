package routers

import "github.com/gin-gonic/gin"

func SetupRoutes(router *gin.Engine) {
	router.GET("/", nil)
	router.GET("/swagger/*any", nil)
}
