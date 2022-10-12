package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/wangdayong228/cns-backend/logger"
	// "github.com/wangdayong228/cns-backend/middlewares"
	// "github.com/wangdayong228/cns-backend/models"
	// "github.com/wangdayong228/cns-backend/routers"
	// "github.com/wangdayong228/cns-backend/routers/assets"
	// "github.com/wangdayong228/cns-backend/services"
)

func initConfig() {
	viper.SetConfigName("config")             // name of config file (without extension)
	viper.SetConfigType("yaml")               // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/etc/rainbow_api/")  // path to look for the config file in
	viper.AddConfigPath("$HOME/.rainbow_api") // call multiple times to add many search paths
	viper.AddConfigPath(".")                  // optionally look for config in the working directory
	err := viper.ReadInConfig()               // Find and read the config file
	if err != nil {                           // Handle errors reading the config file
		log.Fatalln(fmt.Errorf("fatal error config file: %w", err))
	}
}

func initGin() *gin.Engine {
	engine := gin.New()
	engine.Use(gin.Logger())
	engine.Use(middlewares.Logger())
	// engine.Use(gin.Recovery())
	engine.Use(middlewares.Recovery())
	return engine
}

func init() {
	initConfig()
	logger.Init()
	services.InitCfxMainService()
	services.InitCfxTestService()
	middlewares.InitOpenJwtMiddleware()
	middlewares.InitDashboardJwtMiddleware()
	middlewares.InitAdminJwtMiddleware()
	middlewares.InitRateLimitMiddleware()
	logrus.Info("init done")
}

// @title       Rainbow-API
// @version     1.0
// @description The responses of the open api in swagger focus on the data field rather than the code and the message fields

// @license.name Apache 2.0
// @license.url  http://www.apache.org/licenses/LICENSE-2.0.html

// @host     api.nftrainbow.xyz
// @BasePath /v1
// @schemes  http https
func main() {
	models.ConnectDB()

	// go services.StartTXService()
	// go services.SyncNFTMintTaskStatus()
	// go services.SyncNFTMintBatchTaskStatus()
	// go services.SyncNFTTransferTaskStatus()
	// go services.SyncNFTTransferBatchTaskStatus()
	// go services.SyncNFTBurnTaskStatus()
	// go services.SyncContractDeployTaskStatus()
	// go services.RefreshLogOutput()
	// go services.StartWatcherService()

	app := initGin()
	app.Use(middlewares.RateLimitMiddleware)
	routers.SetupRoutes(app)
	admin.SetupRoutes(app)
	assets.SetupRoutes(app)

	port := viper.GetString("port")
	if port == "" {
		logrus.Panic("port must be specified")
	}

	address := fmt.Sprintf("0.0.0.0:%s", port)
	logrus.Info("Rainbow-API Start Listening and serving HTTP on ", address)
	err := app.Run(address)
	if err != nil {
		log.Panic(err)
	}
}
