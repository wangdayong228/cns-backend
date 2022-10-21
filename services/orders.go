package services

import (
	"context"
	"fmt"

	"github.com/wangdayong228/cns-backend/models"
	"github.com/wangdayong228/cns-backend/models/enums"
	confluxpay "github.com/wangdayong228/conflux-pay-sdk-go"
	penums "github.com/wangdayong228/conflux-pay/models/enums"
	pservice "github.com/wangdayong228/conflux-pay/services"
	"gorm.io/gorm"
)

type OrderReq struct {
	TradeProvider string `json:"trade_provider"`
	pservice.MakeOrderReq
}

type OrderResp struct {
	CommitHash string `json:"commit_hash"`
	pservice.MakeOrderResp
}

var (
	confluxPayClient *confluxpay.APIClient
)

func init() {
	configuration := confluxpay.NewConfiguration()
	configuration.Servers = confluxpay.ServerConfigurations{{
		URL:         "http://127.0.0.1:8080/v0",
		Description: "No description provided",
	}}
	confluxPayClient = confluxpay.NewAPIClient(configuration)
}

func MakeOrder(req *OrderReq, commitHash string) (*models.CnsOrder, error) {
	// TODO: verify
	// 1. check commitHash is valid by contract

	// 2. call payservice.makeorder and save order
	provider, ok := penums.ParseTradeProviderByName(req.TradeProvider)
	if !ok {
		return nil, fmt.Errorf("invalid provider: %v", req.TradeProvider)
	}

	var payOrder *confluxpay.ModelsOrder
	var err error

	switch *provider {
	case penums.TRADE_PROVIDER_WECHAT:
		wecahtOrdReq := *confluxpay.NewServicesMakeWechatOrderReq(int32(req.Amount), *req.Description, int32(req.TimeExpire), int32(req.TradeType))
		payOrder, _, err = confluxPayClient.OrdersApi.MakeOrder(context.Background()).WecahtOrdReq(wecahtOrdReq).Execute()
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unspport")
	}

	cnsOrder, err := models.NewOrderByPayResp(payOrder, commitHash)
	if err != nil {
		return nil, err
	}

	err = models.GetDB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(cnsOrder).Error; err != nil {
			return err
		}

		// 3. set commit order state
		commit, err := models.FindCommit(commitHash)
		if err != nil {
			return err
		}

		commit.OrderState = enums.ORDER_STATE_MADE
		if err := tx.Save(commit).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return cnsOrder, nil
}

func GetOrder(commitHash string) (*models.CnsOrder, error) {
	o, err := models.FindOrderByCommitHash(commitHash)
	if err != nil {
		return nil, err
	}

	if o.TradeState.IsStable() {
		return o, nil
	}

	resp, _, err := confluxPayClient.OrdersApi.QueryOrderSummary(context.Background(), o.TradeNo).Execute()
	if err != nil {
		return nil, err
	}

	if penums.TradeState(*resp.TradeState).IsStable() {
		o.TradeState = penums.TradeState(*resp.TradeState)
		models.GetDB().Save(o)
	}
	return o, nil
}
