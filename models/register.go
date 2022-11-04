package models

import (
	"errors"
	"time"

	"github.com/wangdayong228/cns-backend/models/enums"
	confluxpay "github.com/wangdayong228/conflux-pay-sdk-go"

	pmodels "github.com/wangdayong228/conflux-pay/models"
	penums "github.com/wangdayong228/conflux-pay/models/enums"
	"gorm.io/gorm"
)

type Register struct {
	BaseModel
	RegisterCore
}

func (o *Register) Save(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if o.Order != nil {
			err := tx.Save(&o.Order).Error
			if err != nil {
				return err
			}
			o.OrderID = &o.Order.ID
			o.OrderTradeState = &o.Order.TradeState
		}
		return tx.Save(o).Error
	})
}

type RegisterCore struct {
	ProcessInfo
	CommitHash string `gorm:"type:varchar(255);uniqueIndex" json:"commit_hash"`
}

type ProcessInfo struct {
	TxSummary
	*pmodels.Order  `gorm:"-"`
	OrderID         *uint                `json:"-"`
	OrderTradeState *penums.TradeState   `json:"-"`
	UserID          uint                 `json:"-"`
	UserPermission  enums.UserPermission `json:"-"`
}

type TxSummary struct {
	TxID    uint    `json:"-"`
	TxHash  string  `gorm:"type:varchar(255)" json:"tx_hash"`
	TxState TxState `gorm:"uint" json:"tx_state"`
	TxError string  `gorm:"type:varchar(255)" json:"tx_error"`
}

func NewTxSummaryByRaw(tx *Transaction) *TxSummary {
	return &TxSummary{
		TxID:    tx.ID,
		TxHash:  tx.Hash,
		TxState: TxState(tx.State),
		TxError: tx.Error,
	}
}

func (o *ProcessInfo) IsStable() bool {
	if o == nil || o.Order == nil {
		return true
	}
	if o.OrderCore.IsStable() {
		return true
	}
	return o.TradeState == penums.TRADE_STATE_SUCCESSS && o.TxState.IsSuccess()
}

func NewOrderWithTxByPayResp(payResp *confluxpay.ModelsOrder) (*ProcessInfo, error) {
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

	o := ProcessInfo{}
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

func NewRegOrderByPayResp(payResp *confluxpay.ModelsOrder, commitHash string) (*Register, error) {
	o, err := NewOrderWithTxByPayResp(payResp)
	if err != nil {
		return nil, err
	}
	regOrder := Register{}
	regOrder.ProcessInfo = *o
	regOrder.CommitHash = commitHash
	return &regOrder, nil
}

type RegisterOrderOperater struct {
}

func (*RegisterOrderOperater) FindRegOrderByCommitHash(commitHash string) (*Register, error) {
	o := Register{}
	o.CommitHash = commitHash
	if err := GetDB().Where(&o).First(&o).Error; err != nil {
		return nil, err
	}

	if err := CompleteRegisterOrders([]*Register{&o}); err != nil {
		return nil, err
	}

	return &o, nil
}

func (*RegisterOrderOperater) FindNeedRegiterOrders(startID uint) ([]*Register, error) {
	o := Register{}
	// o.TradeState = penums.TRADE_STATE_SUCCESSS

	var orders []*Register
	if err := GetDB().Debug().
		Where(" tx_id = ? and ( order_trade_state = ? or user_permission > ?)", 0, penums.TRADE_STATE_SUCCESSS, 0).
		Where(&o).
		Find(&orders).Error; err != nil {
		return nil, err
	}

	if err := CompleteRegisterOrders(orders); err != nil {
		return nil, err
	}

	// filter trade successs
	// orders = FilterRegisterOrdersByTxState(orders, penums.TRADE_STATE_SUCCESSS)

	return orders, nil
}

func (*RegisterOrderOperater) FindNeedSyncStateRegOrders(count int) ([]*Register, error) {
	var orders []*Register
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
func CompleteRegisterOrders(regOrders []*Register) error {
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

// func FilterRegisterOrdersByTxState(regOrders []*Register, tradeState penums.TradeState) []*Register {
// 	var result []*Register
// 	for _, o := range regOrders {
// 		if o.Order != nil && o.Order.TradeState == tradeState {
// 			result = append(result, o)
// 		}
// 	}
// 	return result
// }
