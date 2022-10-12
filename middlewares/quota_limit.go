package middlewares

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/nft-rainbow/rainbow-api/models"
	rainbow_errors "github.com/nft-rainbow/rainbow-api/rainbow-errors"
	"github.com/nft-rainbow/rainbow-api/utils/ginutils"
)

func QuotaLimitMint() gin.HandlerFunc {
	return quotaLimit(checkMintLimit)
}
func QuotaLimitDeploy() gin.HandlerFunc {
	return quotaLimit(checkDeployLimit)
}
func QuotaLimitFile() gin.HandlerFunc {
	return quotaLimit(checkFileLimit)
}

func quotaLimit(checkFn func(c *gin.Context, userId uint)) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		kycType := uint(claims[KYCTypeKey].(float64))
		userId := uint(claims[AppUserIdKey].(float64))
		if c.Request.Method == "POST" && kycType == models.USER_TYPE_NORMAL {
			checkFn(c, userId)
		}
	}
}

func checkMintLimit(c *gin.Context, userId uint) {
	mintCount, err := models.UserMonthMintCount(userId)
	if err == nil && mintCount >= 100 {
		ginutils.RenderRespError(c, rainbow_errors.ERR_MINT_LIMIT_EXCEEDED)
		c.Abort()
		return
	}
}

func checkDeployLimit(c *gin.Context, userId uint) {
	deployCount, err := models.UserMonthDeployContractCount(userId)
	if err == nil && deployCount >= 10 {
		ginutils.RenderRespError(c, rainbow_errors.ERR_DEPLOY_LIMIT_EXCEEDED)
		c.Abort()
		return
	}
}

func checkFileLimit(c *gin.Context, userId uint) {
	uploadFileCount, err := models.UserUploadFileCount(userId)
	if err == nil && uploadFileCount >= 10 {
		ginutils.RenderRespError(c, rainbow_errors.ERR_UPLOADE_FILE_LIMIT_EXCEEDED)
		c.Abort()
		return
	}
}
