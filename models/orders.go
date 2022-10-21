package models

import (
	"time"

	confluxpay "github.com/wangdayong228/conflux-pay-sdk-go"
	pmodels "github.com/wangdayong228/conflux-pay/models"
	"github.com/wangdayong228/conflux-pay/models/enums"
	penums "github.com/wangdayong228/conflux-pay/models/enums"
)

type CnsOrder struct {
	BaseModel
	CnsOrderCore
}

type CnsOrderCore struct {
	pmodels.OrderCore
	CommitHash      string  `gorm:"type:varchar(255);uniqueIndex" json:"commit_hash"`
	RegisterTxID    uint    `json:"-"`
	RegisterTxHash  string  `gorm:"type:varchar(255)" json:"register_tx_hash"`
	RegisterTxState TxState `gorm:"type:varchar(255)" json:"register_tx_state"`
}

func FindOrderByCommitHash(commitHash string) (*CnsOrder, error) {
	o := CnsOrder{}
	o.CommitHash = commitHash
	return &o, GetDB().Where(&o).First(&o).Error
}

func FindNeedRegiterOrders(startID uint) ([]*CnsOrder, error) {
	o := CnsOrder{}
	o.TradeState = penums.TRADE_STATE_SUCCESSS
	o.RegisterTxID = 0

	var orders []*CnsOrder
	return orders, GetDB().Where("id > ?", startID).Where(&o).Find(&orders).Error
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
func FindNeedSyncStateOrders(count int) ([]*CnsOrder, error) {
	o := CnsOrder{}
	o.RegisterTxID = 0

	var orders []*CnsOrder
	return orders, GetDB().Not(&o).
		Where("register_tx_state = ?", TX_STATE_INIT).
		Or("register_tx_state = ?", TX_STATE_SEND_FAILED_RETRY).
		Or("register_tx_state = ?", TX_STATE_SEND_FAILED_RETRY_UPPER_GAS).
		Or("register_tx_state = ?", TX_STATE_POPULATED).
		Or("register_tx_state = ?", TX_STATE_PENDING).
		Find(&orders).
		Limit(count).Error
}

func NewOrderByPayResp(payResp *confluxpay.ModelsOrder, commitHash string) (*CnsOrder, error) {
	tv, err := time.Parse("2006-01-02T15:04:05+08:00", *payResp.TimeExpire)
	if err != nil {
		return nil, err
	}

	o := CnsOrder{}
	o.Amount = uint(*payResp.Amount)
	o.AppName = *payResp.AppName
	o.CodeUrl = payResp.CodeUrl
	o.CommitHash = commitHash
	o.Description = payResp.Description
	o.H5Url = payResp.H5Url
	o.Provider = enums.TradeProvider(*payResp.TradeProvider)
	o.TimeExpire = &tv
	o.TradeNo = *payResp.TradeNo
	o.TradeState = enums.TradeState(*payResp.TradeState)
	o.TradeType = enums.TradeType(*payResp.TradeType)
	o.RegisterTxState = TX_STATE_INIT
	return &o, nil
}
