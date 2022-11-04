package middlewares

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
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
		usersMap[u.ApiKeyHash] = u
	}

	return func(c *gin.Context) {
		apikey := c.GetHeader("X-Api-Key")
		apikeyHash := common.BytesToHash(crypto.Keccak256([]byte(apikey))).Hex()
		u := usersMap[apikeyHash]
		if u == nil {
			ginutils.RenderRespError(c, cns_errors.ERR_AUTHORIZATION_NO_PERMISSION)
			c.Abort()
			return
		}
		c.Set("user", u)
	}
}
