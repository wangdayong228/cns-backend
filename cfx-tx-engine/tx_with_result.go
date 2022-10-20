package cfx_tx_engine

import (
	"encoding/json"
	"math/big"
	"strings"

	"github.com/Conflux-Chain/go-conflux-sdk/types"
	"github.com/Conflux-Chain/go-conflux-sdk/types/cfxaddress"
	"github.com/Conflux-Chain/go-conflux-sdk/utils"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/spf13/viper"
	"github.com/wangdayong228/cns-backend/models"
)

type TxWithResult struct {
	ID         int
	RawTx      *models.Transaction
	Tx         *types.UnsignedTransaction
	Hash       *types.Hash
	Status     *uint64
	Err        error
	TrySendCnt int

	stateChangeHandlers []*StateChangeHandler
}

// type TxState int

func NewTxWithResultFromTransaction(tx *models.Transaction) *TxWithResult {
	utx := types.UnsignedTransaction{}
	_from := cfxaddress.MustNew(tx.From, uint32(tx.ChainId))
	utx.From = &_from
	if tx.To != "" {
		_to := cfxaddress.MustNew(tx.To, uint32(tx.ChainId))
		utx.To = &_to
	}

	if tx.Data != "" {
		utx.Data = hexutil.MustDecode(tx.Data)
	}

	utx.Value = (*hexutil.Big)(tx.Value.BigInt())

	return &TxWithResult{int(tx.ID), tx, &utx, nil, nil, nil, 0, nil}
}

func (t *TxWithResult) registerStateChangeEvents(handlers ...*StateChangeHandler) {
	t.stateChangeHandlers = append(t.stateChangeHandlers, handlers...)
}

func (t *TxWithResult) FillBackTransaction() {

	tx := t.RawTx
	tx.State = int(t.State())

	if t.Err != nil {
		tx.Error = t.Err.Error()
	}
	// pending and exceed retry limit
	if t.Hash != nil {
		tx.Hash = t.Hash.String()
	}
	if t.Tx.Nonce != nil {
		tx.Nonce = uint(t.Tx.Nonce.ToInt().Uint64())
	}
	models.GetDB().Save(tx)
}

func (t *TxWithResult) SetError(err error) {
	t.doAndEmitStateChange(func() {
		t.Err = err
		t.FillBackTransaction()
	})
}

func (t *TxWithResult) SetTxStatus(txStatus uint64, epochNumber uint) {
	t.doAndEmitStateChange(func() {
		t.Status = &txStatus
		t.FillBackTransaction()
	})
}

func (t *TxWithResult) SetTxHash(txhash types.Hash) {
	t.doAndEmitStateChange(func() {
		t.Hash = &txhash
		t.FillBackTransaction()
	})
}

func (t *TxWithResult) doAndEmitStateChange(doFn func()) {
	oldstate := t.State()

	doFn()

	newstate := t.State()
	t.emitStateChangeEvent(t.RawTx, oldstate, newstate)
}

func (t *TxWithResult) emitStateChangeEvent(tx *models.Transaction, oldstate models.TxState, newstate models.TxState) {
	if oldstate != newstate {
		for _, h := range t.stateChangeHandlers {
			if h == nil {
				continue
			}
			(*h)(tx, oldstate, newstate)
		}
	}
}

func (t TxWithResult) MarshalJSON() ([]byte, error) {
	type alias TxWithResult
	type composite struct {
		alias
		Err string
	}
	tmp := composite{}
	tmp.alias = alias(t)
	if t.Err != nil {
		tmp.Err = t.Err.Error()
	}

	return json.Marshal(tmp)
}

func (t *TxWithResult) State() models.TxState {
	if t.Hash == nil && t.Status == nil && t.Err == nil {
		if t.Tx.Gas != nil && t.Tx.GasPrice != nil && t.Tx.StorageLimit != nil && t.Tx.ChainID != nil {
			return models.TX_STATE_POPULATED
		}
		return models.TX_STATE_INIT
	}

	if t.Status != nil {
		switch *t.Status {
		case 0:
			return models.TX_STATE_EXECUTED
		}
		return models.TX_STATE_EXECUTE_FAILED
	}

	if t.Err != nil {
		retry, upperGas := needRetry(t)
		if !retry {
			return models.TX_STATE_SEND_FAILED
		}
		if upperGas {
			return models.TX_STATE_SEND_FAILED_RETRY_UPPER_GAS
		}
		return models.TX_STATE_SEND_FAILED_RETRY
	}

	return models.TX_STATE_PENDING
}

//  executed 不重发
//  rpc error 且 非 tx_pool full 不重发
//  out of balance 不需要重发
//  其余重发
func (t *TxWithResult) IsFinalized() bool {
	state := t.State()

	switch state {
	case models.TX_STATE_EXECUTED:
		return true
	case models.TX_STATE_EXECUTE_FAILED:
		return true
	case models.TX_STATE_SEND_FAILED:
		return true
	case models.TX_STATE_CONFIRMED:
		return true
	}

	return false
}

func (t *TxWithResult) reset() {
	t.Tx.GasPrice = nil
	t.Tx.EpochHeight = nil
	t.Tx.Gas = nil
	t.Tx.StorageLimit = nil
	t.Tx.Nonce = nil

	t.RawTx.Hash = ""
	t.RawTx.Error = ""
	t.RawTx.Nonce = 0

	t.Err = nil
	t.Hash = nil
}

func (t *TxWithResult) improveGasPrice(delta *big.Int) {
	// if tx already exists, set old gasprice to 2e9
	if t.Tx.GasPrice == nil {
		t.Tx.GasPrice = types.NewBigInt(1e9)
	}
	oldGasPrice := t.Tx.GasPrice.ToInt()
	_newGasPrice := new(big.Int).Add(oldGasPrice, delta)
	t.Tx.GasPrice = types.NewBigIntByRaw(_newGasPrice)
}

func needRetry(t *TxWithResult) (_needRetry bool, needUpperGas bool) {
	if t.TrySendCnt > viper.GetInt("tx_engine.retryLimit") {
		return false, false
	}

	if ok, rpcErr := IsRpcError(t.Err); ok {
		if rpcErr == TX_ERR_RPC_TXPOOL_FULL ||
			rpcErr == TX_ERR_NORMAL_TOO_MANY_REQUEST {
			return true, false
		}
		if rpcErr == TX_ERR_NORMAL_ALREADY_EXIST {
			return true, true
		}
	}

	if ok, normalErr := IsNormalError(t.Err); ok &&
		(normalErr == TX_ERR_NORMAL_TIMEOUT) {
		return true, false
	}

	return false, false
}

func IsRpcError(err error) (bool, TxRpcError) {
	v, err := utils.ToRpcError(err)
	if err != nil {
		return false, TX_ERR_RPC_OTHER
	}
	return true, getRpcErrorType(v)
}

func IsNormalError(err error) (bool, TxNormalError) {
	if utils.IsNil(err) {
		return false, TX_ERR_NORMAL_OTHER
	}

	if ok, _ := IsRpcError(err); ok {
		return false, TX_ERR_NORMAL_OTHER
	}
	if strings.Contains(err.Error(), "timeout") {
		return true, TX_ERR_NORMAL_TIMEOUT
	}
	if strings.Contains(err.Error(), "still pending") {
		return true, TX_ERR_NORMAL_PENDING_LIMIT
	}

	return true, TX_ERR_NORMAL_OTHER
}

func getIdsOfTxWithResults(ts []*TxWithResult) []int {
	results := []int{}
	for _, v := range ts {
		results = append(results, v.ID)
	}
	return results
}
