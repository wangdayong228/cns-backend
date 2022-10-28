package services

import (
	"math/big"

	"github.com/wangdayong228/cns-backend/config"
	"github.com/wangdayong228/cns-backend/models"
	"github.com/wangdayong228/cns-backend/models/enums"
)

type MakeRenewOrderReq struct {
}

// create renew task which will create a tx, return task
func MakeRenewOrder(req *MakeRenewOrderReq) (*models.RenewOrder, error) {
	return nil, nil
}

// create renew tx
func createRenewTx(name string, duration *big.Int) (*models.Transaction, error) {
	data, err := dataGen.Renew(name, duration)
	if err != nil {
		return nil, err
	}
	return CreateTransaction(enums.CHAIN_TYPE_CFX, enums.ChainID(config.RpcVal.ChainID),
		config.CnsContractVal.Admin.String(), config.CnsContractVal.Register.String(),
		big.NewInt(0), data)
}
