package config

import (
	"fmt"
	"log"

	"github.com/Conflux-Chain/go-conflux-sdk/types/cfxaddress"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	RpcVal         *Rpc
	CnsContractVal *CnsContracts
	TxEngineVal    *TxEngine
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

	initConfigs()
	logrus.WithField("rpc", RpcVal).WithField("contract", CnsContractVal).WithField("tx_engine", TxEngineVal).Info("init config completed")
}

type Rpc struct {
	Url         string
	ChainID     uint32
	PrivateKeys []string
}

type CnsContracts struct {
	Register cfxaddress.Address
	Admin    cfxaddress.Address
}

type TxEngine struct {
	RetryLimit    int
	SendCountOnce int
}

func initConfigs() {
	if err := viper.UnmarshalKey("rpc", &RpcVal); err != nil {
		panic(err)
	}
	if err := viper.UnmarshalKey("cnsContracts", &CnsContractVal); err != nil {
		panic(err)
	}
	if err := viper.UnmarshalKey("txEngine", &TxEngineVal); err != nil {
		panic(err)
	}
}
