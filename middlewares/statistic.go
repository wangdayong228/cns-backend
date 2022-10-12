package middlewares

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/nft-rainbow/rainbow-api/models"
)

func Statistic() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		claims := jwt.ExtractClaims(c)
		userId := uint(claims[AppUserIdKey].(float64))
		models.IncreaseStatistic(userId, c.Request.Method, c.FullPath())
	}
}
