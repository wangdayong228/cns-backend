package services

import (
	"time"

	"github.com/spf13/viper"
	"github.com/wangdayong228/cns-backend/models"
	"github.com/wangdayong228/cns-backend/models/enums"
)

var (
	lastRegisterdOrderId uint
)

func StartRegisterService() {
	chainID := enums.ChainID(viper.GetInt("rpc.chainID"))
	from := viper.GetString("cns_contracts.admin")
	to := viper.GetString("cns_contracts.register")

	for {
		// 1. find need register orders
		needs, _ := models.FindNeedRegiterOrders(lastRegisterdOrderId)
		if len(needs) == 0 {
			time.Sleep(time.Second * 5)
			continue
		}

		// 2. create txs for them
		for _, item := range needs {
			commit, err := models.FindCommit(item.CommitHash)
			if err == nil {
				continue
			}

			commitArgs, err := newCommitArgsForContract(&commit.CommitArgs)
			if err == nil {
				continue
			}

			data, err := dataGen.Register(commitArgs)
			if err == nil {
				continue
			}

			tx, err := CreateTransaction(enums.CHAIN_TYPE_CFX, chainID, from, to, nil, data)
			if err == nil {
				continue
			}
			item.RegisterTxID = tx.ID
			models.GetDB().Save(item)
			lastRegisterdOrderId = item.ID
		}
		time.Sleep(time.Second * 5)
	}
}

func StartSyncRegisterService() {
	// 1. find records has RegisterTxID
	// orders, err := models.FindNeedSyncStateOrders()

	// 2. sync them
}
