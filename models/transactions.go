package models

import (
	"encoding/hex"

	// "github.com/nft-rainbow/rainbow-api/utils"

	"github.com/Conflux-Chain/go-conflux-sdk/types/cfxaddress"
	"github.com/shopspring/decimal"
	"github.com/wangdayong228/cns-backend/models/enums"
	// "github.com/shopspring/decimal"
)

type TxState int

const (
	TX_STATE_SEND_FAILED_RETRY_UPPER_GAS TxState = iota - 4 // -3
	TX_STATE_SEND_FAILED_RETRY                              // -3
	TX_STATE_EXECUTE_FAILED                                 // -2
	TX_STATE_SEND_FAILED                                    // -1
	TX_STATE_INIT                                           // 0
	TX_STATE_POPULATED                                      // 1
	TX_STATE_PENDING                                        // 2
	TX_STATE_EXECUTED                                       // 3
	TX_STATE_CONFIRMED                                      // 4
)

func (t TxState) String() string {
	switch t {
	case TX_STATE_SEND_FAILED_RETRY_UPPER_GAS:
		return "Failed and will retry with improve gas price"
	case TX_STATE_SEND_FAILED_RETRY:
		return "Failed and will retry"
	case TX_STATE_EXECUTE_FAILED:
		return "Exectued and failed"
	case TX_STATE_SEND_FAILED:
		return "Send failed"
	case TX_STATE_INIT:
		return "Ready to send"
	case TX_STATE_POPULATED:
		return "Populated"
	case TX_STATE_PENDING:
		return "Pending"
	case TX_STATE_EXECUTED:
		return "Excuted and success"
	case TX_STATE_CONFIRMED:
		return "Confirmed"
	}
	return "Unrecgonize"
}

func (t TxState) IsSuccess() bool {
	return t == TX_STATE_EXECUTED || t == TX_STATE_CONFIRMED
}

func (t TxState) IsFailed() bool {
	return t == TX_STATE_SEND_FAILED || t == TX_STATE_EXECUTE_FAILED
}

func (t TxState) IsFinalized() bool {
	return t.IsSuccess() || t.IsFailed()
}

type Transaction struct {
	BaseModel
	ChainType   uint            `gorm:"type:int"`
	ChainId     uint            `gorm:"type:int"`
	From        string          `gorm:"type:varchar(256);index"`
	To          string          `gorm:"type:varchar(256);index"`
	Nonce       uint            `gorm:"index"`
	Value       decimal.Decimal `gorm:"type:decimal(65,0)"`
	Data        string          `gorm:"type:text"`
	Hash        string          `gorm:"type:varchar(256);index"`
	State       int             `gorm:"default:0"`
	EpochNumber uint            `gorm:"type:uint" json:"epoch_number"`
	Error       string          `gorm:"type:text" json:"error"`
}

// gas, gasPrice, storageLimit, nonce, epochHeight

func (tx *Transaction) Verify() error {
	if tx.ChainType == uint(enums.CHAIN_TYPE_CFX) && (tx.ChainId == uint(enums.CHAINID_CFX_TESTNET) || tx.ChainId == uint(enums.CHAINID_CFX_MAINNET)) {
		if _, err := cfxaddress.New(tx.From, uint32(tx.ChainId)); err != nil {
			return err
		}
		if tx.To != "" {
			if _, err := cfxaddress.New(tx.To, uint32(tx.ChainId)); err != nil {
				return err
			}
		}
		if tx.Data != "" {
			if _, err := hex.DecodeString(tx.Data); err != nil {
				return err
			}
		}
	}
	return nil
}

func (tx *Transaction) IsSuccess() bool {
	return TxState(tx.State).IsSuccess()
}

func (tx *Transaction) IsFailed() bool {
	return TxState(tx.State).IsFailed()
}

func (tx *Transaction) IsFinalized() bool {
	return TxState(tx.State).IsFinalized()
}

func GetTxIds(txs []*Transaction) []uint {
	ids := []uint{}
	for _, tx := range txs {
		ids = append(ids, tx.ID)
	}
	return ids
}

func FindTransactions(option *Transaction) ([]Transaction, error) {
	txs := []Transaction{}
	err := GetDB().Where(option).Find(&txs).Error
	return txs, err
}

func FindTransactionByID(id uint) (*Transaction, error) {
	tx := Transaction{}
	err := GetDB().First(&tx, id).Error
	return &tx, err
}
