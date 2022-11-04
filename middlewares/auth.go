package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/wangdayong228/cns-backend/cns_errors"
	"github.com/wangdayong228/cns-backend/models"
	"github.com/wangdayong228/cns-backend/utils/ginutils"
)

func Auth() gin.HandlerFunc {
	users, err := models.GetAllUsers()
	if err != nil {
		panic(err)
	}
	usersMap := make(map[string]*models.User)
	for _, u := range users {
		usersMap[u.ApiKey] = u
	}

	return func(c *gin.Context) {
		apikey := c.GetHeader("X-Api-Key")
		u := usersMap[apikey]
		if u == nil {
			ginutils.RenderRespError(c, cns_errors.ERR_AUTHORIZATION_NO_PERMISSION)
			c.Abort()
			return
		}
		c.Set("user_id", u.ID)
		c.Set("user_permission", uint(u.Permission))
	}
}
