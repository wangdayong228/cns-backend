package services

import (
	"math/big"
	"os"

	sdk "github.com/Conflux-Chain/go-conflux-sdk"
	tx_engine "github.com/wangdayong228/cns-backend/cfx-tx-engine"
	"github.com/wangdayong228/cns-backend/config"
	"github.com/wangdayong228/cns-backend/contracts"
	confluxpay "github.com/wangdayong228/conflux-pay-sdk-go"
)

var (
	rpcClient         *sdk.Client
	web3RegController *contracts.Web3RegisterController
	maxCommitmentAge  *big.Int
	confluxPayClient  *confluxpay.APIClient
)

func Init() {
	initRpcClient()
	initContracts()
	initConfluxPay()
}

func initRpcClient() {
	rpc := config.RpcVal
	rpcClient = sdk.MustNewClient(rpc.Url, sdk.ClientOption{
		RetryCount: 3,
		Logger:     os.Stdout,
	})
	rpcClient.SetAccountManager(tx_engine.NewPrivatekeyAccountManager(rpc.PrivateKeys, uint32(rpc.ChainID)))
}

func initContracts() {
	var err error
	web3RegController, err = contracts.NewWeb3RegisterController(config.CnsContractVal.Register, rpcClient)
	if err != nil {
		panic(err)
	}

	maxCommitmentAge, err = web3RegController.MaxCommitmentAge(nil)
	if err != nil {
		panic(err)
	}
}

func initConfluxPay() {
	configuration := confluxpay.NewConfiguration()
	configuration.Servers = confluxpay.ServerConfigurations{{
		URL:         "http://127.0.0.1:8080/v0",
		Description: "No description provided",
	}}
	confluxPayClient = confluxpay.NewAPIClient(configuration)
}

func StartServices() {
	go TxService()
	go RegisterService()
	go SyncRegisterStateService()
	go RenewService()
	go SyncRenewStateService()
}
