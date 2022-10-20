package cfx_tx_engine

import (
	"math/big"
	"time"

	sdk "github.com/Conflux-Chain/go-conflux-sdk"
	"github.com/sirupsen/logrus"
	"github.com/wangdayong228/cns-backend/models"
)

type StateChangeHandler func(tx *models.Transaction, oldstate models.TxState, newstate models.TxState)
type TransactionSender struct {
	client     *sdk.Client
	bulkHelper bulkHelper
	// per tx max send count, if still not execute after try to send the count, set the tx to error
	perTxSendLimit int

	// state change event
	stateChangeHandlers []*StateChangeHandler
}

func NewTransactionSender(client *sdk.Client, perTxSendLimit int) *TransactionSender {
	t := &TransactionSender{client, bulkHelper{client}, perTxSendLimit, nil}
	return t
}

func (t *TransactionSender) BulkSend(txs []*models.Transaction) []*models.Transaction {
	logrus.WithField("tx ids", models.GetTxIds(txs)).Info("tx engine bulk send start")
	var txWithResults []*TxWithResult
	for _, tx := range txs {
		tWithR := NewTxWithResultFromTransaction(tx)
		tWithR.registerStateChangeEvents(t.stateChangeHandlers...)
		txWithResults = append(txWithResults, tWithR)
	}

	t.send(txWithResults)

	for i := range txWithResults {
		txWithResults[i].FillBackTransaction()
	}

	logrus.WithField("tx ids", models.GetTxIds(txs)).Info("tx engine bulk send completed")
	return txs
}

func (t *TransactionSender) send(tWithRs []*TxWithResult) {
	if tWithRs = filterNeedSends(tWithRs); tWithRs == nil {
		return
	}
	logrus.WithField("len", len(tWithRs)).WithField("ids", getIdsOfTxWithResults(tWithRs)).Info("filtered need sends")

	cleanTxs(tWithRs)

	t.bulkHelper.populateTxs(tWithRs)
	logrus.WithField("len", len(tWithRs)).WithField("failed ids", getIdsOfTxWithResults(filterErrorTxs(tWithRs))).Info("populate txWithRs completed")

	populated := filterByState(tWithRs, models.TX_STATE_POPULATED)
	t.bulkHelper.signAndSend(populated)
	logrus.WithField("len", len(populated)).WithField("failed ids", getIdsOfTxWithResults(filterErrorTxs(populated))).Info("signAndSend populated tWithRs completed")

	pending := filterByState(tWithRs, models.TX_STATE_PENDING)
	t.bulkHelper.waitTxsBeExecute(pending, t.perTxSendLimit)
	logrus.WithField("len", len(pending)).WithField("failed ids", getIdsOfTxWithResults(filterErrorTxs(pending))).Info("waitTxsBeExecute pending txWithRs completed")

	time.Sleep(time.Second)
	t.send(tWithRs)
}

func (t *TransactionSender) RegisterStateChangeEvent(h *StateChangeHandler) {
	t.stateChangeHandlers = append(t.stateChangeHandlers, h)
}

func (t *TransactionSender) UnregisterStateChangeEvent(h *StateChangeHandler) {
	for i, sh := range t.stateChangeHandlers {
		if h == sh {
			t.stateChangeHandlers = append(t.stateChangeHandlers[:i], t.stateChangeHandlers[i+1:]...)
			break
		}
	}
}

func Verify(txs []*models.Transaction) {
	for _, v := range txs {
		if err := v.Verify(); err != nil {
			v.Error = err.Error()
		}
	}
}

func filterNeedSends(all []*TxWithResult) []*TxWithResult {
	var needSends []*TxWithResult
	for _, v := range all {
		if !v.IsFinalized() {
			needSends = append(needSends, v)
		}
	}
	return needSends
}

func filterPendings(all []*TxWithResult) []*TxWithResult {
	return filterByState(all, models.TX_STATE_PENDING)
}

func filterByState(all []*TxWithResult, state models.TxState) []*TxWithResult {
	var results []*TxWithResult
	for _, v := range all {
		if v.State() == state {
			results = append(results, v)
		}
	}
	return results
}

func filterErrorTxs(all []*TxWithResult) []*TxWithResult {
	var results []*TxWithResult
	for _, v := range all {
		if v.State() == models.TX_STATE_SEND_FAILED ||
			v.State() == models.TX_STATE_SEND_FAILED_RETRY ||
			v.State() == models.TX_STATE_SEND_FAILED_RETRY_UPPER_GAS {
			results = append(results, v)
		}
	}
	return results
}

func increaseSendCount(tWithRs []*TxWithResult) {
	for _, v := range tWithRs {
		v.TrySendCnt++
	}
}

func cleanTxs(tWithRs []*TxWithResult) {
	for _, v := range tWithRs {
		// improve gas for pending, otherwise reset
		state := v.State()
		if state == models.TX_STATE_PENDING {
			v.improveGasPrice(big.NewInt(1e9))
			continue
		}
		if state == models.TX_STATE_SEND_FAILED_RETRY_UPPER_GAS {
			_gasPrice := v.Tx.GasPrice
			v.reset()
			v.Tx.GasPrice = _gasPrice
			v.improveGasPrice(big.NewInt(1e9))
			continue
		}
		v.reset()
	}
}

func setAllErr(txWithRs []*TxWithResult, err error) {
	for _, v := range txWithRs {
		v.SetError(err)
	}
}
