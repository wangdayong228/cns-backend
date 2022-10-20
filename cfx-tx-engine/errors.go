package cfx_tx_engine

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/Conflux-Chain/go-conflux-sdk/utils"
)

type ErrNotEnoughCashes map[string]ErrNotEnoughCash

func (e ErrNotEnoughCashes) Error() string {
	msgs := []string{}
	for user, err := range e {
		msgs = append(msgs, fmt.Sprintf("%v %v", user, err.Error()))
	}
	return strings.Join(msgs, "\n")
}

type ErrNotEnoughCash struct {
	Need *big.Int
	Got  *big.Int
}

func (e ErrNotEnoughCash) Error() string {
	return fmt.Sprintf("out of balance, need %v, got %v", e.Need, e.Got)
}

type TxRpcError int

const (
	TX_ERR_RPC_OUT_OF_BALANCE TxRpcError = iota
	TX_ERR_RPC_TXPOOL_FULL
	TX_ERR_RPC_OTHER
	//TODO: need retry
	TX_ERR_NORMAL_TOO_MANY_REQUEST
	//TODO: need retry
	TX_ERR_NORMAL_ALREADY_EXIST
)

func getRpcErrorType(err *utils.RpcError) TxRpcError {
	if strings.Contains(err.Message, "txpool is full") {
		return TX_ERR_RPC_TXPOOL_FULL
	}
	if strings.Contains(err.Message, "discarded due to out of balance") {
		return TX_ERR_RPC_OUT_OF_BALANCE
	}
	if strings.Contains(err.Message, "too many requests") {
		return TX_ERR_NORMAL_TOO_MANY_REQUEST
	}
	if strings.Contains(err.Data.(string), "already inserted") {
		return TX_ERR_NORMAL_ALREADY_EXIST
	}

	return TX_ERR_RPC_OTHER
}

type TxNormalError int

const (
	TX_ERR_NORMAL_OTHER TxNormalError = iota
	TX_ERR_NORMAL_TIMEOUT
	TX_ERR_NORMAL_PENDING_LIMIT
)
