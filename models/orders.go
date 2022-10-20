package models

import (
	"time"

	confluxpay "github.com/wangdayong228/conflux-pay-sdk-go"
	pmodels "github.com/wangdayong228/conflux-pay/models"
	"github.com/wangdayong228/conflux-pay/models/enums"
)

type CnsOrder struct {
	BaseModel
	CnsOrderCore
}

type CnsOrderCore struct {
	pmodels.OrderCore
	CommitHash      string `gorm:"type:varchar(255);uniqueIndex" json:"commit_hash"`
	RegisterTxHash  string `gorm:"type:varchar(255)" json:"register_tx_hash"`
	RegisterTxState uint   `gorm:"type:varchar(255)" json:"register_tx_state"`
}

func FindOrderByCommitHash(commitHash string) (*CnsOrder, error) {
	o := CnsOrder{}
	o.CommitHash = commitHash
	return &o, GetDB().Where(&o).First(&o).Error
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
	return &o, nil
}
