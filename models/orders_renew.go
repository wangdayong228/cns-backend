package models

import (
	confluxpay "github.com/wangdayong228/conflux-pay-sdk-go"
	penums "github.com/wangdayong228/conflux-pay/models/enums"
)

type RenewOrder struct {
	BaseModel
	OrderWithTx
	CnsName  string `json:"cns_name" binding:"required"`
	Duration int    `gorm:"type:varchar(255)" json:"duration" binding:"required"`
}

func NewRenewOrderByPayResp(payResp *confluxpay.ModelsOrder) (*RenewOrder, error) {
	o, err := NewOrderWithTxByPayResp(payResp)
	if err != nil {
		return nil, err
	}
	result := RenewOrder{}
	result.OrderWithTx = *o
	return &result, nil
}

type RenewOrderOperater struct {
}

func (*RenewOrderOperater) FindOrderById(id int) (*RenewOrder, error) {
	o := RenewOrder{}
	if err := GetDB().Where(id).First(&o).Error; err != nil {
		return nil, err
	}
	return &o, nil
}

func (*RenewOrderOperater) FindNeedRnewOrders(startID uint) ([]*RenewOrder, error) {
	o := RenewOrder{}
	o.TradeState = penums.TRADE_STATE_SUCCESSS
	o.TxID = 0

	var orders []*RenewOrder
	return orders, GetDB().Where("id > ? and tx_id = ?", startID, 0).Where(&o).Find(&orders).Error
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
func (*RenewOrderOperater) FindNeedSyncStateOrders(count int) ([]*RenewOrder, error) {
	var orders []*RenewOrder
	return orders, GetDB().
		Not("tx_id = ?", 0).
		Where("tx_state = ?", TX_STATE_INIT).
		Find(&orders).
		Limit(count).Error
}

func (r *RenewOrderOperater) UpdateOrderState(id int, raw *confluxpay.ModelsOrder) error {
	o, err := r.FindOrderById(id)
	if err != nil {
		return err
	}

	tradeState, _ := penums.ParseTradeState(*raw.TradeState)
	o.TradeState = *tradeState

	refundState, _ := penums.ParserefundState(*raw.RefundState)
	o.RefundState = *refundState

	return GetDB().Save(o).Error
}
