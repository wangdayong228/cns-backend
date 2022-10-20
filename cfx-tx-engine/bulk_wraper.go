package cfx_tx_engine

import (
	"fmt"
	"math/big"
	"time"

	sdk "github.com/Conflux-Chain/go-conflux-sdk"
	"github.com/Conflux-Chain/go-conflux-sdk/cfxclient/bulk"
	"github.com/Conflux-Chain/go-conflux-sdk/types"
	"github.com/Conflux-Chain/go-conflux-sdk/types/cfxaddress"
	"github.com/wangdayong228/cns-backend/models"

	"github.com/sirupsen/logrus"
)

type bulkHelper struct {
	client *sdk.Client
}

func (b *bulkHelper) createBulkSender(txWithRs []*TxWithResult) *bulk.BulkSender {
	bulkSender := bulk.NewBulkSender(*b.client)
	for _, txWithR := range txWithRs {
		bulkSender.AppendTransaction(txWithR.Tx)
	}
	return bulkSender
}

// populate and set error for all txWithResult
func (b *bulkHelper) populateTxs(txWithRs []*TxWithResult) {
	bulkSender := b.createBulkSender(txWithRs)
	populated, err := bulkSender.PopulateTransactions(false)
	if populated != nil {
		for i := range txWithRs {
			txWithRs[i].Tx = populated[i]
		}
	}

	if err == nil {
		return
	}

	if _estimateErrs, ok := err.(*bulk.ErrBulkEstimate); ok {
		for i, estimateErr := range *_estimateErrs {
			if estimateErr != nil {
				txWithRs[i].SetError(estimateErr)
				logrus.WithError(estimateErr).WithField("tx", txWithRs[i]).Info("populate txWithResult and estimate error")
			}
		}
		return
	}

	setAllErr(txWithRs, err)
	logrus.WithError(err).Info("populate txWithRs error (not estimate error)")
}

// sign&send and set error for all txWithResult
func (b *bulkHelper) signAndSend(txWithRs []*TxWithResult) {
	increaseSendCount(txWithRs)

	bulkSender := b.createBulkSender(txWithRs)
	hashes, errors, err := bulkSender.SignAndSend()

	if err != nil {
		setAllErr(txWithRs, err)
		logrus.WithError(err).Info("signAndSend txWithRs error (not rpc error)")
		return
	}

	for i, txWithR := range txWithRs {
		if errors[i] != nil {
			txWithR.SetError(errors[i])
			logrus.WithError(errors[i]).WithField("tx", txWithRs[i]).Info("signAndSend txWithRs error (rpc error)")
		}

		if hashes[i] != nil {
			txWithR.SetTxHash(*hashes[i])
		}
	}
}

func (b *bulkHelper) checkBalance(txWithRs []*TxWithResult) {
	needs := make(map[string]*big.Int)
	gots := make(map[string]*big.Int)
	balanceErrors := ErrNotEnoughCashes{}

	for _, txWithR := range txWithRs {
		tx := txWithR.Tx
		gas := new(big.Int).Mul(tx.Gas.ToInt(), tx.GasPrice.ToInt())
		storage := new(big.Int).Div(new(big.Int).SetUint64(uint64(*tx.StorageLimit)), big.NewInt(1024))
		storage = storage.Mul(storage, big.NewInt(1e18))
		need := new(big.Int).Add(gas, storage)
		need = need.Add(need, tx.Value.ToInt())
		needs[tx.From.String()] = new(big.Int).Add(needs[tx.From.String()], need)
	}

	for k, v := range needs {
		got, err := b.client.GetBalance(cfxaddress.MustNew(k))
		if err != nil {
			setAllErr(txWithRs, err)
			return
		}
		gots[k] = got.ToInt()
		if v.Cmp(gots[k]) < 0 {
			balanceErrors[k] = ErrNotEnoughCash{v, gots[k]}
		}
	}

	for _, txWithR := range txWithRs {
		if v, ok := balanceErrors[txWithR.Tx.From.String()]; ok {
			txWithR.SetError(v)
		}
	}
}

func (b *bulkHelper) waitTxsBeExecute(txWithRs []*TxWithResult, perTxSendLimit int) {
	for i := 0; i < 5; i++ {
		// wait last tx for every from; set status if executed, at most wait 25s
		txWithRs = filterPendings(txWithRs)
		if len(txWithRs) == 0 {
			return
		}

		bulkCaller := bulk.NewBulkCaller(b.client)
		receipts := make([]*types.TransactionReceipt, len(txWithRs))
		errors := make([]*error, len(txWithRs))
		for i, v := range txWithRs {
			receipts[i], errors[i] = bulkCaller.GetTransactionReceipt(*v.Hash)
		}

		// err means timeout
		if err := bulkCaller.Execute(); err != nil {
			continue
		}

		for i, err := range errors {
			if *err != nil {
				txWithRs[i].SetError(*err)
				logrus.WithError(*err).WithField("txWithRs", txWithRs[i]).Info("get tx receipt error")
				continue
			}

			if receipts[i].TransactionHash == "" {
				continue
			}

			_status := (uint64)(receipts[i].OutcomeStatus)
			txWithRs[i].SetTxStatus(_status, uint(*receipts[i].EpochNumber))
			txWithRs[i].SetTxHash(receipts[i].TransactionHash)
		}
		time.Sleep(time.Second * 5)
	}

	for _, v := range txWithRs {
		if v.State() == models.TX_STATE_PENDING && v.TrySendCnt >= perTxSendLimit {
			err := fmt.Errorf("sent %v times still pending", v.TrySendCnt)
			v.SetError(err)
			logrus.WithError(err).WithField("txWithRs", v).Info("wait tx receipt arrive limit")
		}
	}
}
