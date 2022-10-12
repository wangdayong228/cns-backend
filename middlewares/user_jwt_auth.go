package middlewares

import (
	"runtime/debug"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	"github.com/nft-rainbow/rainbow-api/models"
	rainbow_errors "github.com/nft-rainbow/rainbow-api/rainbow-errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	JwtAuthMiddleware *jwt.GinJWTMiddleware
)

type login struct {
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

var JwtIdentityKey = "id"

type User struct {
	Id    uint
	Email string
	Name  string
}

func InitDashboardJwtMiddleware() {
	var err error
	JwtAuthMiddleware, err = jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "rainbow-api-jwt",
		Key:         []byte(viper.GetString("jwtKeys.dashboard")),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour * 24,
		IdentityKey: JwtIdentityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*User); ok {
				return jwt.MapClaims{
					JwtIdentityKey: v.Id,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			id := claims[JwtIdentityKey]
			return uint(id.(float64))
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			email := loginVals.Email
			password := loginVals.Password
			pwdHash := crypto.Keccak256Hash([]byte(password)).String()

			user, err := models.FindUserByEmail(email)
			if err == nil && user.Password == pwdHash {
				return &User{
					Id:    user.ID,
					Name:  user.Name,
					Email: user.Email,
				}, nil
			}

			return nil, jwt.ErrFailedAuthentication
		},
		LoginResponse: func(c *gin.Context, code int, message string, time time.Time) {
			c.JSON(code, gin.H{
				"token":  message,
				"expire": time,
			})
		},
		RefreshResponse: func(c *gin.Context, code int, message string, time time.Time) {
			c.JSON(code, gin.H{
				"token":  message,
				"expire": time,
			})
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    rainbow_errors.ERR_AUTHORIZATION_JWT,
				"message": message,
			})
		},
		TokenLookup:   "header: Authorization", // cookie: jwt, query: token
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})

	if err != nil {
		logrus.WithError(err).WithField("stack", string(debug.Stack())).Fatal("init DashboardJWT middleware error")
		return
	}

	logrus.Info("init dashboard jwt middleware done")

}
