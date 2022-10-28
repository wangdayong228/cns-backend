package services

import (
	"context"
	"fmt"
	"math/big"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/wangdayong228/cns-backend/models"
	confluxpay "github.com/wangdayong228/conflux-pay-sdk-go"
	penums "github.com/wangdayong228/conflux-pay/models/enums"
	pservice "github.com/wangdayong228/conflux-pay/services"
)

type MakeRenewOrderReq struct {
	TradeProvider string           `json:"trade_provider" swaggertype:"string"`
	TradeType     penums.TradeType `json:"trade_type" binding:"required" swaggertype:"string"`
	Description   *string          `json:"description" binding:"required"`
	CnsName       string           `json:"cns_name" binding:"required"`
	Duration      int              `json:"duration" binding:"required"`
}

type MakeRenewOrderResp struct {
	pservice.MakeOrderResp
}

type RenewOrderService struct {
	modelOperator models.RenewOrderOperater
}

func (r *RenewOrderService) MakeOrder(req *MakeRenewOrderReq) (*models.RenewOrder, error) {
	// 获取价格
	price, err := web3RegController.RentPriceInFiat(nil, req.CnsName, big.NewInt(int64(req.Duration)))
	if err != nil {
		return nil, err
	}
	amount := new(big.Int).Add(price.Base, price.Premium)
	amount = amount.Div(amount, big.NewInt(1e6))
	fmt.Println("price", amount)
	if amount.Cmp(big.NewInt(1)) <= 0 {
		amount = big.NewInt(1)
	}

	// 2. call payservice.makeorder and save order
	provider, ok := penums.ParseTradeProviderByName(req.TradeProvider)
	if !ok {
		return nil, fmt.Errorf("invalid provider: %v", req.TradeProvider)
	}

	var payOrder *confluxpay.ModelsOrder

	switch *provider {
	case penums.TRADE_PROVIDER_WECHAT:
		wecahtOrdReq := *confluxpay.NewServicesMakeOrderReq(int32(amount.Int64()), *req.Description, int32(maxCommitmentAge.Int64()), req.TradeType.String())

		var resp *http.Response
		payOrder, resp, err = confluxPayClient.OrdersApi.MakeOrder(context.Background()).MakeOrdReq(wecahtOrdReq).Execute()
		if err != nil {
			logrus.WithError(err).WithField("order request", wecahtOrdReq).Info("failed to make order throught conflux-pay")
			return nil, err
		}
		fmt.Printf("make order resp %v\n", resp)
	default:
		return nil, fmt.Errorf("unspport")
	}

	RenewOrder, err := models.NewRenewOrderByPayResp(payOrder)
	if err != nil {
		return nil, err
	}

	if err := models.GetDB().Save(RenewOrder).Error; err != nil {
		return nil, err
	}

	return RenewOrder, nil
}

func (r *RenewOrderService) GetOrder(id int) (*models.RenewOrder, error) {
	o, err := r.modelOperator.FindOrderById(id)
	if err != nil {
		return nil, err
	}

	if o.IsStable() {
		return o, nil
	}

	resp, _, err := confluxPayClient.OrdersApi.QueryOrderSummary(context.Background(), o.TradeNo).Execute()
	if err != nil {
		return nil, err
	}

	r.modelOperator.UpdateOrderState(id, resp)
	return o, nil
}

func (r *RenewOrderService) RefreshURL(id int) (*models.RenewOrder, error) {
	o, err := r.GetOrder(id)
	if err != nil {
		return nil, err
	}

	if o.IsStable() {
		return nil, ErrOrderCompleted
	}

	resp, _, err := confluxPayClient.OrdersApi.RefreshPayUrl(context.Background(), o.TradeNo).Execute()
	if err != nil {
		return nil, err
	}

	o.CodeUrl = resp.CodeUrl
	o.H5Url = resp.H5Url
	return o, nil
}
