package services

import (
	"context"
	"math/big"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/wangdayong228/cns-backend/config"
	"github.com/wangdayong228/cns-backend/models"
	"github.com/wangdayong228/cns-backend/models/enums"
	confluxpay "github.com/wangdayong228/conflux-pay-sdk-go"
	"gorm.io/gorm"
)

// create renew task which will create a tx, return task
func MakeRenewOrder(req *MakeRenewOrderReq) (*models.RenewOrder, error) {
	return nil, nil
}

// create renew tx
func createRenewTx(name string, duration *big.Int) (*models.Transaction, error) {
	data, err := dataGen.Renew(name, duration)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"chainID": config.RpcVal.ChainID,
			"data":    data,
		}).WithError(err).Error("failed gen Renew data")
		return nil, err
	}
	return CreateTransaction(enums.CHAIN_TYPE_CFX, enums.ChainID(config.RpcVal.ChainID),
		config.CnsContractVal.Admin.String(), config.CnsContractVal.Register.String(),
		big.NewInt(0), data)
}

var (
	lastRenewdOrderId  uint
	renewOrderOperator = models.RenewOrderOperater{}
)

func RenewService() {
	// from := config.CnsContractVal.Admin
	// to := config.CnsContractVal.Register

	for {
		// 1. find need Renew orders
		needs, _ := renewOrderOperator.FindNeedRnewOrders(lastRenewdOrderId)
		if len(needs) == 0 {
			time.Sleep(time.Second * 5)
			continue
		}
		logrus.WithField("needs", needs).Info("find need Renew orders")

		// 2. create txs for them
		for _, item := range needs {
			logrus.WithField("order", item).Error("creat Renew tx for order")
			// commit, err := renewOrderOperator.FindOrderById(item.CommitHash)
			// if err != nil {
			// 	logrus.WithField("commit hash", item.CommitHash).WithError(err).Error("failed find commit")
			// 	continue
			// }

			// commitArgs, err := newCommitArgsForContract(&commit.CommitArgs)
			// if err != nil {
			// 	logrus.WithField("commit args", commit.CommitArgs).WithError(err).Error("failed convert commit args")
			// 	continue
			// }

			// data, err := dataGen.Renew(commitArgs)
			// if err != nil {
			// 	logrus.WithField("args", commit.CommitArgs).WithError(err).Error("failed gen Renew data")
			// 	continue
			// }

			tx, err := createRenewTx(item.CnsName, big.NewInt(int64(item.Duration)))
			if err != nil {
				continue
			}
			item.TxID = tx.ID
			models.GetDB().Save(item)
			lastRenewdOrderId = item.ID
		}
		time.Sleep(time.Second * 5)
	}
}

// TODO: implement
func SyncRenewStateService() {
	for {
		time.Sleep(time.Second * 5)
		// 1. find records has RenewTxID and state is UnCompleted
		orders, err := renewOrderOperator.FindNeedSyncStateOrders(500)
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
			tx, err := models.FindTransactionByID(o.TxID)
			if err != nil {
				logrus.WithField("tx_id", o.TxID).WithError(err).Error("failed find tx by id")
				continue
			}
			if tx.IsFinalized() {
				o.TxHash = tx.Hash
				o.TxState = models.TxState(tx.State)

				if err = models.GetDB().Save(o).Error; err != nil {
					logrus.WithField("order", o).WithError(err).Error("failed save order")
					continue
				}

				if o.TxState.IsSuccess() {
					continue
				}

				// refund
				_, _, err := confluxPayClient.OrdersApi.Refund(context.Background(), o.TradeNo).
					RefundReq(confluxpay.ServicesRefundReq{Reason: "failed to Renew cns"}).Execute()
				logrus.WithField("order", o).WithError(err).Error("refund order")
			}
		}
	}
}
