package models

import (
	"errors"
	"time"

	confluxpay "github.com/wangdayong228/conflux-pay-sdk-go"
	pmodels "github.com/wangdayong228/conflux-pay/models"
	penums "github.com/wangdayong228/conflux-pay/models/enums"
)

type RegisterOrder struct {
	BaseModel
	RegisterOrderCore
}

type RegisterOrderCore struct {
	pmodels.OrderCore
	CommitHash      string  `gorm:"type:varchar(255);uniqueIndex" json:"commit_hash"`
	RegisterTxID    uint    `json:"-"`
	RegisterTxHash  string  `gorm:"type:varchar(255)" json:"register_tx_hash"`
	RegisterTxState TxState `gorm:"type:varchar(255)" json:"register_tx_state"`
}

func NewOrderByPayResp(payResp *confluxpay.ModelsOrder, commitHash string) (*RegisterOrder, error) {
	tv, err := time.Parse("2006-01-02T15:04:05+08:00", *payResp.TimeExpire)
	if err != nil {
		return nil, err
	}

	provider, ok1 := penums.ParseTradeProviderByName(*payResp.TradeProvider)
	tradeState, ok2 := penums.ParseTradeState(*payResp.TradeState)
	tradeType, ok3 := penums.ParseTradeType(*payResp.TradeType)
	refundState, ok4 := penums.ParserefundState(*payResp.RefundState)

	if !ok1 || !ok2 || !ok3 || !ok4 {
		return nil, errors.New("unkown trade type or trade provider or trade state or refund state")
	}

	o := RegisterOrder{}
	o.Amount = uint(*payResp.Amount)
	o.AppName = *payResp.AppName
	o.CodeUrl = payResp.CodeUrl
	o.CommitHash = commitHash
	o.Description = payResp.Description
	o.H5Url = payResp.H5Url
	o.Provider = *provider
	o.TimeExpire = &tv
	o.TradeNo = *payResp.TradeNo
	o.TradeState = *tradeState
	o.RefundState = *refundState
	o.TradeType = *tradeType
	o.RegisterTxState = TX_STATE_INIT
	return &o, nil
}

func FindOrderByCommitHash(commitHash string) (*RegisterOrder, error) {
	o := RegisterOrder{}
	o.CommitHash = commitHash
	if err := GetDB().Where(&o).First(&o).Error; err != nil {
		return nil, err
	}
	return &o, nil
}

func FindNeedRegiterOrders(startID uint) ([]*RegisterOrder, error) {
	o := RegisterOrder{}
	o.TradeState = penums.TRADE_STATE_SUCCESSS
	o.RegisterTxID = 0

	var orders []*RegisterOrder
	return orders, GetDB().Where("id > ? and register_tx_id = ?", startID, 0).Where(&o).Find(&orders).Error
}

// TX_STATE_SEND_FAILED_RETRY_UPPER_GAS TxState = iota - 4 // -4
// TX_STATE_SEND_FAILED_RETRY                              // -3
// TX_STATE_EXECUTE_FAILED                                 // -2
// TX_STATE_SEND_FAILED                                    // -1
// TX_STATE_INIT                                           // 0
// TX_STATE_POPULATED                                      // 1
// TX_STATE_PENDING                                        // 2
// TX_STATE_EXECUTED                                       // 3
// TX_STATE_CONFIRMED
func FindNeedSyncStateOrders(count int) ([]*RegisterOrder, error) {
	var orders []*RegisterOrder
	return orders, GetDB().
		Not("register_tx_id = ?", 0).
		Where("register_tx_state = ?", TX_STATE_INIT).
		// Or("register_tx_state = ?", TX_STATE_SEND_FAILED_RETRY).
		// Or("register_tx_state = ?", TX_STATE_SEND_FAILED_RETRY_UPPER_GAS).
		// Or("register_tx_state = ?", TX_STATE_POPULATED).
		// Or("register_tx_state = ?", TX_STATE_PENDING).
		Find(&orders).
		Limit(count).Error
}

func (o *RegisterOrderCore) IsStable() bool {
	if o.OrderCore.IsStable() {
		return true
	}
	return o.TradeState == penums.TRADE_STATE_SUCCESSS && o.RegisterTxState.IsSuccess()
}

func UpdateOrderState(commitHash string, raw *confluxpay.ModelsOrder) error {
	o, err := FindOrderByCommitHash(commitHash)
	if err != nil {
		return err
	}

	tradeState, _ := penums.ParseTradeState(*raw.TradeState)
	o.TradeState = *tradeState

	refundState, _ := penums.ParserefundState(*raw.RefundState)
	o.RefundState = *refundState

	return GetDB().Save(o).Error
}
