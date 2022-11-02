package models

import (
	"errors"
	"time"

	confluxpay "github.com/wangdayong228/conflux-pay-sdk-go"

	pmodels "github.com/wangdayong228/conflux-pay/models"
	penums "github.com/wangdayong228/conflux-pay/models/enums"
	"gorm.io/gorm"
)

type RegisterOrder struct {
	BaseModel
	RegisterOrderCore
}

func (o *RegisterOrder) Save(db *gorm.DB) error {
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

type RegisterOrderCore struct {
	OrderWithTx
	CommitHash string `gorm:"type:varchar(255);uniqueIndex" json:"commit_hash"`
}

type OrderWithTx struct {
	*pmodels.Order `gorm:"-"`
	OrderID        *uint
	TxID           uint    `json:"-"`
	TxHash         string  `gorm:"type:varchar(255)" json:"tx_hash"`
	TxState        TxState `gorm:"type:varchar(255)" json:"tx_state"`
}

func (o *OrderWithTx) IsStable() bool {
	if o.OrderCore.IsStable() {
		return true
	}
	return o.TradeState == penums.TRADE_STATE_SUCCESSS && o.TxState.IsSuccess()
}

func NewOrderWithTxByPayResp(payResp *confluxpay.ModelsOrder) (*OrderWithTx, error) {
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

	o := OrderWithTx{}
	o.Order = &pmodels.Order{}
	o.Amount = uint(*payResp.Amount)
	o.AppName = *payResp.AppName
	o.CodeUrl = payResp.CodeUrl
	o.Description = payResp.Description
	o.H5Url = payResp.H5Url
	o.Provider = *provider
	o.TimeExpire = &tv
	o.TradeNo = *payResp.TradeNo
	o.TradeState = *tradeState
	o.RefundState = *refundState
	o.TradeType = *tradeType
	o.TxState = TX_STATE_INIT
	return &o, nil
}

func NewRegOrderByPayResp(payResp *confluxpay.ModelsOrder, commitHash string) (*RegisterOrder, error) {
	o, err := NewOrderWithTxByPayResp(payResp)
	if err != nil {
		return nil, err
	}
	regOrder := RegisterOrder{}
	regOrder.OrderWithTx = *o
	regOrder.CommitHash = commitHash
	return &regOrder, nil
}

type RegisterOrderOperater struct {
}

func (*RegisterOrderOperater) FindRegOrderByCommitHash(commitHash string) (*RegisterOrder, error) {
	o := RegisterOrder{}
	o.CommitHash = commitHash
	if err := GetDB().Where(&o).First(&o).Error; err != nil {
		return nil, err
	}

	if err := CompleteRegisterOrders([]*RegisterOrder{&o}); err != nil {
		return nil, err
	}

	return &o, nil
}

func (*RegisterOrderOperater) FindNeedRegiterOrders(startID uint) ([]*RegisterOrder, error) {
	o := RegisterOrder{}
	// o.TradeState = penums.TRADE_STATE_SUCCESSS

	var orders []*RegisterOrder
	if err := GetDB().Where("id > ? and tx_id = ?", startID, 0).Where(&o).Find(&orders).Error; err != nil {
		return nil, err
	}

	if err := CompleteRegisterOrders(orders); err != nil {
		return nil, err
	}

	// filter trade successs
	orders = FilterRegisterOrdersByTxState(orders, penums.TRADE_STATE_SUCCESSS)

	return orders, nil
}

func (*RegisterOrderOperater) FindNeedSyncStateRegOrders(count int) ([]*RegisterOrder, error) {
	var orders []*RegisterOrder
	if err := GetDB().
		Not("tx_id = ?", 0).
		Where("tx_state = ?", TX_STATE_INIT).
		// Or("tx_state = ?", TX_STATE_PENDING).
		Find(&orders).
		Limit(count).Error; err != nil {
		return nil, err
	}

	if err := CompleteRegisterOrders(orders); err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *RegisterOrderOperater) UpdateRegOrderState(commitHash string, raw *confluxpay.ModelsOrder) error {
	o, err := r.FindRegOrderByCommitHash(commitHash)
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
func CompleteRegisterOrders(regOrders []*RegisterOrder) error {
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

func FilterRegisterOrdersByTxState(regOrders []*RegisterOrder, tradeState penums.TradeState) []*RegisterOrder {
	var result []*RegisterOrder
	for _, o := range regOrders {
		if o.Order != nil && o.Order.TradeState == tradeState {
			result = append(result, o)
		}
	}
	return result
}
