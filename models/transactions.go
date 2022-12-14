package models

import (
	"encoding/hex"
	"encoding/json"

	// "github.com/nft-rainbow/rainbow-api/utils"

	"github.com/Conflux-Chain/go-conflux-sdk/types/cfxaddress"
	"github.com/shopspring/decimal"
	"github.com/wangdayong228/cns-backend/models/enums"
	// "github.com/shopspring/decimal"
)

type TxState int

const (
	TX_STATE_SEND_FAILED_RETRY_UPPER_GAS TxState = iota - 4 // -4
	TX_STATE_SEND_FAILED_RETRY                              // -3
	TX_STATE_EXECUTE_FAILED                                 // -2
	TX_STATE_SEND_FAILED                                    // -1
	TX_STATE_INIT                                           // 0
	TX_STATE_POPULATED                                      // 1
	TX_STATE_PENDING                                        // 2
	TX_STATE_EXECUTED                                       // 3
	TX_STATE_CONFIRMED                                      // 4
)

type TxStateDesc struct {
	Code        string
	Description string
}

var (
	tradeProviderValue2StrMap  map[TxState]TxStateDesc
	tradeProviderName2ValueMap map[string]TxState
	tradeProviderCode2ValueMap map[string]TxState
)

func init() {
	tradeProviderValue2StrMap = map[TxState]TxStateDesc{
		TX_STATE_SEND_FAILED_RETRY_UPPER_GAS: {"SEND_FAILED_RETRY_UPPER_GAS", "Failed and will retry with improve gas price"},
		TX_STATE_SEND_FAILED_RETRY:           {"SEND_FAILED_RETRY", "Failed and will retry"},
		TX_STATE_EXECUTE_FAILED:              {"EXECUTE_FAILED", "Exectued and failed"},
		TX_STATE_SEND_FAILED:                 {"SEND_FAILED", "Send failed"},
		TX_STATE_INIT:                        {"INIT", "Ready to send"},
		TX_STATE_POPULATED:                   {"POPULATED", "Populated"},
		TX_STATE_PENDING:                     {"PENDING", "Pending"},
		TX_STATE_EXECUTED:                    {"EXECUTED_SUCCESS", "Excuted and success"},
		TX_STATE_CONFIRMED:                   {"CONFIRMED", "Confirmed"},
	}

	tradeProviderName2ValueMap = make(map[string]TxState)
	tradeProviderCode2ValueMap = make(map[string]TxState)
	for k, v := range tradeProviderValue2StrMap {
		tradeProviderName2ValueMap[v.Description] = k
		tradeProviderCode2ValueMap[v.Code] = k
	}
}

func (t TxState) String() string {
	v, ok := tradeProviderValue2StrMap[t]
	if ok {
		return v.Description
	}
	return "unknown"
}

func (t TxState) Code() string {
	v, ok := tradeProviderValue2StrMap[t]
	if ok {
		return v.Code
	}
	return "unknown"
}

func (t TxState) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Code())
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

func FindTransactionByIDs(ids []uint) ([]*Transaction, error) {
	txs := []*Transaction{}
	err := GetDB().Where(ids).Find(txs).Error
	return txs, err
}
