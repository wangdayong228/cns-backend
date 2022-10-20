package services

import (
	"os"

	sdk "github.com/Conflux-Chain/go-conflux-sdk"
	"github.com/spf13/viper"
)

var (
	rpcClient *sdk.Client
)

func Init() {
	rpcClient = sdk.MustNewClient(viper.GetString("rpc.url"), sdk.ClientOption{
		RetryCount: 3,
		Logger:     os.Stdout,
	})
	// TODO: add admin private key
}
