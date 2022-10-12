package middlewares

import (
	"bytes"

	"github.com/gin-gonic/gin"
	rainbow_errors "github.com/nft-rainbow/rainbow-api/rainbow-errors"
	"github.com/nft-rainbow/rainbow-api/utils/ginutils"
	"github.com/sirupsen/logrus"
)

func Recovery() gin.HandlerFunc {
	var buf bytes.Buffer
	return gin.CustomRecoveryWithWriter(&buf, gin.RecoveryFunc(func(c *gin.Context, err interface{}) {
		defer func() {
			logrus.WithField("recovered", buf.String()).Error("panic and recovery")
			buf.Reset()
		}()
		ginutils.RenderRespError(c, rainbow_errors.ERR_INTERNAL_SERVER_COMMON)
		c.Abort()
	}))
}
