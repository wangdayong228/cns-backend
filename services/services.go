package services

import (
	"os"

	sdk "github.com/Conflux-Chain/go-conflux-sdk"
	"github.com/spf13/viper"
	tx_engine "github.com/wangdayong228/cns-backend/cfx-tx-engine"
	"github.com/wangdayong228/cns-backend/config"
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
	rpc := config.RpcVal
	rpcClient = sdk.MustNewClient(viper.GetString(rpc.Url), sdk.ClientOption{
		RetryCount: 3,
		Logger:     os.Stdout,
	})
	// TODO: add admin private key
	rpcClient.SetAccountManager(tx_engine.NewPrivatekeyAccountManager(rpc.PrivateKeys, uint32(rpc.ChainID)))
}

func StartServices() {
	go TxService()
	go RegisterService()
	go SyncRegisterStateService()
}
