package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

func Init() {
	viper.SetConfigName("config")             // name of config file (without extension)
	viper.SetConfigType("yaml")               // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/etc/cns_backend/")  // path to look for the config file in
	viper.AddConfigPath("$HOME/.cns_backend") // call multiple times to add many search paths
	viper.AddConfigPath(".")                  // optionally look for config in the working directory
	viper.AddConfigPath("..")                 // optionally look for config in the working directory
	err := viper.ReadInConfig()               // Find and read the config file
	if err != nil {                           // Handle errors reading the config file
		log.Fatalln(fmt.Errorf("fatal error config file: %w", err))
	}
}
