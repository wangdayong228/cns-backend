package services

import (
	"math/big"
	"time"

	tx_engine "github.com/wangdayong228/cns-backend/cfx-tx-engine"
	"github.com/wangdayong228/cns-backend/models"
	"github.com/wangdayong228/cns-backend/models/enums"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type txProcesser func()

func CreateTransaction(chainType enums.ChainType, chainId enums.ChainID, from string, to string, value *big.Int, data string) (*models.Transaction, error) {
	valueInDecimal := decimal.NewFromBigInt(value, 0)
	tx := models.Transaction{
		ChainType: uint(chainType),
		ChainId:   uint(chainId),
		From:      from,
		To:        to,
		Value:     valueInDecimal,
		Data:      data,
		State:     int(models.TX_STATE_INIT),
	}
	res := models.GetDB().Create(&tx)
	return &tx, res.Error
}

func FindTxNeedToBeProcessed(chainType enums.ChainType, chainId enums.ChainID, count int) ([]*models.Transaction, error) {
	var txs []*models.Transaction
	findRes := models.GetDB().Where("chain_type = ? and chain_id = ? and (state = ? or state = ?)", chainType, chainId, models.STATUS_INIT, models.TX_STATE_PENDING).Order("id asc").Limit(count).Find(&txs)
	return txs, findRes.Error
}

func FindTransactionByID(id uint) (*models.Transaction, error) {
	return models.FindTransactionByID(id)
}

func StartTXService() {
	logrus.Info("start task for sending transactions")
	sendInterval := time.Second * 5

	txProcessers := []txProcesser{}
	txProcessers = append(txProcessers, MustGetTxProcesser(enums.CHAIN_TYPE_CFX, enums.CHAINID_CFX_MAINNET))
	txProcessers = append(txProcessers, MustGetTxProcesser(enums.CHAIN_TYPE_CFX, enums.CHAINID_CFX_TESTNET))

	for {
		for _, p := range txProcessers {
			p()
		}
		time.Sleep(sendInterval)
	}
}

func MustGetTxProcesser(chainType enums.ChainType, chainId enums.ChainID) txProcesser {
	logEntry := logrus.WithFields(logrus.Fields{
		"chain type": chainType,
		"chain ID":   chainId,
	})
	if chainType == enums.CHAIN_TYPE_CFX {
		// client := MustGetCfxChainEnv(chainId).Client
		retryLimit := viper.GetInt("tx_engine.retryLimit")
		sendCountOnce := viper.GetInt("tx_engine.sendCountOnce")
		logEntry.WithField("retryLimit", retryLimit).WithField("sendCountOnce", sendCountOnce).Info("created tx processer")
		return func() {
			txs, err := FindTxNeedToBeProcessed(chainType, chainId, sendCountOnce)
			if len(txs) == 0 {
				return
			}

			if err != nil {
				logEntry.WithError(err).Info("find txs need to be processed error")
				return
			}

			logEntry.WithField("len", len(txs)).Info("find txs need to be processed")
			bSender := tx_engine.NewTransactionSender(rpcClient, retryLimit)
			txWithSendInfo := bSender.BulkSend(txs)
			// update tx
			for _, tx := range txWithSendInfo {
				models.GetDB().Save(&tx)
			}
			logEntry.WithField("len", len(txs)).Info("txs process completed")
		}
	}
	logEntry.Panic("unsupported")
	return nil
}
