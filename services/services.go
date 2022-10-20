package services

import (
	"os"

	sdk "github.com/Conflux-Chain/go-conflux-sdk"
	"github.com/spf13/viper"
	tx_engine "github.com/wangdayong228/cns-backend/cfx-tx-engine"
)

var (
	rpcClient *sdk.Client
)

type Pagination struct {
	Page     int `json:"page,omitempty"`
	PageSize int `json:"page_size,omitempty"`
}

func (p Pagination) CalcOffsetLimit() (offset int, limit int) {

	if p.Page <= 0 {
		p.Page = 1
	}

	switch {
	case p.PageSize > 100:
		p.PageSize = 100
	case p.PageSize <= 0:
		p.PageSize = 10
	}

	offset = (p.Page - 1) * p.PageSize
	limit = p.PageSize
	return
}

func Init() {
	rpcClient = sdk.MustNewClient(viper.GetString("rpc.url"), sdk.ClientOption{
		RetryCount: 3,
		Logger:     os.Stdout,
	})
	// TODO: add admin private key
	privKeys := viper.GetStringSlice("rpc.private_keys")
	chainID := viper.GetInt("rpc.chainID")
	rpcClient.SetAccountManager(tx_engine.NewPrivatekeyAccountManager(privKeys, uint32(chainID)))
}

func StartServices() {
	go StartTXService()
}
