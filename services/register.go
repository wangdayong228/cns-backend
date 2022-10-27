package services

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/wangdayong228/cns-backend/config"
	"github.com/wangdayong228/cns-backend/models"
	"github.com/wangdayong228/cns-backend/models/enums"
	"gorm.io/gorm"
)

var (
	lastRegisterdOrderId uint
)

func RegisterService() {
	from := config.CnsContractVal.Admin
	to := config.CnsContractVal.Register

	for {
		// 1. find need register orders
		needs, _ := models.FindNeedRegiterOrders(lastRegisterdOrderId)
		if len(needs) == 0 {
			time.Sleep(time.Second * 5)
			continue
		}
		logrus.WithField("needs", needs).Info("find need register orders")

		// 2. create txs for them
		for _, item := range needs {
			logrus.WithField("order", item).Error("creat register tx for order")
			commit, err := models.FindCommit(item.CommitHash)
			if err != nil {
				logrus.WithField("commit hash", item.CommitHash).WithError(err).Error("failed find commit")
				continue
			}

			commitArgs, err := newCommitArgsForContract(&commit.CommitArgs)
			if err != nil {
				logrus.WithField("commit args", commit.CommitArgs).WithError(err).Error("failed convert commit args")
				continue
			}

			data, err := dataGen.Register(commitArgs)
			if err != nil {
				logrus.WithField("args", commit.CommitArgs).WithError(err).Error("failed gen register data")
				continue
			}

			tx, err := CreateTransaction(enums.CHAIN_TYPE_CFX, enums.ChainID(config.RpcVal.ChainID), from.String(), to.String(), nil, data)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"chainID": config.RpcVal.ChainID,
					"from":    from,
					"to":      to,
					"data":    data,
				}).WithError(err).Error("failed gen register data")
				continue
			}
			item.RegisterTxID = tx.ID
			models.GetDB().Save(item)
			lastRegisterdOrderId = item.ID
		}
		time.Sleep(time.Second * 5)
	}
}

// TODO: implement
func SyncRegisterStateService() {
	for {
		time.Sleep(time.Second * 5)
		// 1. find records has RegisterTxID and state is UnCompleted
		orders, err := models.FindNeedSyncStateOrders(500)
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				logrus.WithError(err).Error("failed find orders need by sync")
			}
			continue
		}

		if len(orders) == 0 {
			continue
		}

		logrus.WithField("orders", orders).Info("find orders need sync reigster state")

		// 2. sync them
		for _, o := range orders {
			tx, err := models.FindTransactionByID(o.RegisterTxID)
			if err != nil {
				logrus.WithField("tx_id", o.RegisterTxID).WithError(err).Error("failed find tx by id")
				continue
			}
			if tx.IsFinalized() {
				o.RegisterTxHash = tx.Hash
				o.RegisterTxState = models.TxState(tx.State)
				// refund
				_, _, err := confluxPayClient.OrdersApi.Refund(context.Background(), o.TradeNo).Execute()
				if err != nil {
					logrus.WithField("order", o).WithError(err).Error("failed refund order")
					continue
				}

				if err = models.GetDB().Save(o).Error; err != nil {
					logrus.WithField("order", o).WithError(err).Error("failed save order")
					continue
				}
			}
		}
	}
}
