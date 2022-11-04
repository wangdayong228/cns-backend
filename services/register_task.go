package services

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/wangdayong228/cns-backend/config"
	"github.com/wangdayong228/cns-backend/models"
	"github.com/wangdayong228/cns-backend/models/enums"
	confluxpay "github.com/wangdayong228/conflux-pay-sdk-go"
	"gorm.io/gorm"
)

var (
	lastRegisterdOrderId uint
	regOrderOperator     = models.RegisterOrderOperater{}
)

func LoopSendRegisterTx() {
	from := config.CnsContractVal.Admin
	to := config.CnsContractVal.Register

	for {
		// 1. find need register orders
		needs, _ := regOrderOperator.FindNeedRegiterOrders(lastRegisterdOrderId)
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

			data, err := dataGen.RegisterWithFiat(commitArgs)
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
			item.TxID = tx.ID
			err = item.Save(models.GetDB())
			if err != nil {
				logrus.WithField("value", item).WithError(err).Error("failed save register data")
				continue
			}

			// lastRegisterdOrderId = item.ID
		}
		time.Sleep(time.Second * 5)
	}
}

// TODO: implement
func LoopSyncRegisterState() {
	for {
		time.Sleep(time.Second * 5)
		// 1. find records has RegisterTxID and state is UnCompleted
		orders, err := regOrderOperator.FindNeedSyncStateRegOrders(500)
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
				o.TxSummary = *models.NewTxSummaryByRaw(tx)

				if err = o.Save(models.GetDB()); err != nil {
					logrus.WithField("order", o).WithError(err).Error("failed save order")
					continue
				}

				if o.TxState.IsSuccess() {
					continue
				}

				// refund
				_, _, err := confluxPayClient.OrdersApi.Refund(context.Background(), o.TradeNo).
					RefundReq(confluxpay.ServicesRefundReq{Reason: "failed to register cns"}).Execute()
				logrus.WithField("order", o).WithError(err).Error("refund order")
			}
		}
	}
}
