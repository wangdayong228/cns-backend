package models

import (
	confluxpay "github.com/wangdayong228/conflux-pay-sdk-go"
	pmodels "github.com/wangdayong228/conflux-pay/models"
	penums "github.com/wangdayong228/conflux-pay/models/enums"
	"gorm.io/gorm"
)

type RenewOrder struct {
	BaseModel
	RenewOrderCore
}

func (o *RenewOrder) Save(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if o.Order != nil {
			err := tx.Save(&o.Order).Error
			if err != nil {
				return err
			}
			o.OrderID = &o.Order.ID
		}
		return tx.Save(o).Error
	})
}

type RenewOrderCore struct {
	ProcessInfo
	RenewOrderArgs
}

type RenewOrderArgs struct {
	CnsName       string `json:"cns_name" binding:"required"`
	Duration      uint   ` json:"duration" binding:"required"`
	Fuses         uint32 `json:"fuses"`
	WrapperExpiry uint64 `json:"wrapper_expiry" binding:"required"`
}

func NewRenewOrderByPayResp(payResp *confluxpay.ModelsOrder, renewArgs *RenewOrderArgs) (*RenewOrder, error) {
	o, err := NewOrderWithTxByPayResp(payResp)
	if err != nil {
		return nil, err
	}
	result := RenewOrder{}
	result.ProcessInfo = *o
	result.RenewOrderArgs = *renewArgs
	return &result, nil
}

type RenewOrderOperater struct {
}

func (*RenewOrderOperater) FindOrderById(id int) (*RenewOrder, error) {
	o := RenewOrder{}
	if err := GetDB().Where(id).First(&o).Error; err != nil {
		return nil, err
	}

	if err := CompleteRenewOrders([]*RenewOrder{&o}); err != nil {
		return nil, err
	}

	return &o, nil
}

func (*RenewOrderOperater) FindNeedRnewOrders(startID uint) ([]*RenewOrder, error) {
	o := RenewOrder{}
	// o.TradeState = penums.TRADE_STATE_SUCCESSS
	// o.TxID = 0

	var orders []*RenewOrder
	if err := GetDB().
		Where(" tx_id = ? and ( order_trade_state = ? or user_permission > ?)", 0, penums.TRADE_STATE_SUCCESSS, 0).
		Where(&o).Find(&orders).Error; err != nil {
		return nil, err
	}

	if err := CompleteRenewOrders(orders); err != nil {
		return nil, err
	}

	// orders = FilterRenewOrdersByTxState(orders, penums.TRADE_STATE_SUCCESSS)

	return orders, nil
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
	if err := GetDB().
		Not("tx_id = ?", 0).
		Where("tx_state = ?", TX_STATE_INIT).
		Find(&orders).
		Limit(count).Error; err != nil {
		return nil, err
	}

	if err := CompleteRenewOrders(orders); err != nil {
		return nil, err
	}

	return orders, nil
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

	return o.Save(GetDB())
}

// TODO: 用泛型代替
func CompleteRenewOrders(regOrders []*RenewOrder) error {
	// find orders
	var ids []uint
	for _, o := range regOrders {
		if o.OrderID != nil {
			ids = append(ids, *o.OrderID)
		}
	}

	var orders []*pmodels.Order
	if err := GetDB().Find(&orders, ids).Error; err != nil {
		return nil
	}

	// map orders
	var ordersCache map[uint]*pmodels.Order = make(map[uint]*pmodels.Order)
	for _, o := range orders {
		ordersCache[o.ID] = o
	}

	for _, ro := range regOrders {
		if ro.OrderID != nil {
			ro.Order = ordersCache[*ro.OrderID]
		}
	}

	return nil
}

// func FilterRenewOrdersByTxState(renewOrders []*RenewOrder, tradeState penums.TradeState) []*RenewOrder {
// 	var result []*RenewOrder
// 	for _, o := range renewOrders {
// 		if o.Order != nil && o.Order.TradeState == tradeState {
// 			result = append(result, o)
// 		}
// 	}
// 	return result
// }
