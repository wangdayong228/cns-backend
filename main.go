package main

import (
	"fmt"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/wangdayong228/cns-backend/config"
	"github.com/wangdayong228/cns-backend/logger"
	"github.com/wangdayong228/cns-backend/middlewares"
	"github.com/wangdayong228/cns-backend/models"
	"github.com/wangdayong228/cns-backend/routers"
	"github.com/wangdayong228/cns-backend/services"
	// "gorm.io/gorm/logger"
	// "github.com/wangdayong228/cns-backend/routers/assets"
	// "github.com/wangdayong228/cns-backend/services"
)

func initGin() *gin.Engine {
	engine := gin.New()
	engine.Use(gin.Logger())
	engine.Use(middlewares.Logger())
	// engine.Use(gin.Recovery())
	engine.Use(middlewares.Recovery())
	return engine
}

// func init() {
// logger.Init()
// middlewares.InitOpenJwtMiddleware()
// middlewares.InitRateLimitMiddleware()
// logrus.Info("init done")
// }

// @title       CNS-BACKEND
// @version     1.0
// @description The responses of the open api in swagger focus on the data field rather than the code and the message fields

// @license.name Apache 2.0
// @license.url  http://www.apache.org/licenses/LICENSE-2.0.html

// @host     127.0.0.1:8081
// @BasePath /v1
// @schemes  http https
func main() {
	config.Init()
	logger.Init()
	models.ConnectDB()

	services.StartServices()

	app := initGin()
	// app.Use(middlewares.RateLimitMiddleware)
	app.Use(cors.Default())
	routers.SetupRoutes(app)

	port := viper.GetString("port")
	if port == "" {
		logrus.Panic("port must be specified")
	}

	address := fmt.Sprintf("0.0.0.0:%s", port)
	logrus.Info("Cns-Backend Start Listening and serving HTTP on ", address)
	err := app.Run(address)
	if err != nil {
		log.Panic(err)
	}
}
