package services

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
	"time"

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
	models.RenewOrderArgs
}

type MakeRenewOrderResp struct {
	ID uint `json:"id"`
	pservice.MakeOrderResp
}

type RnewByAdminResp struct {
	ID uint `json:"id"`
	models.RenewOrderArgs
	models.TxSummary
}

func NewRenewByAdminRespByRaw(reg *models.Renew) *RnewByAdminResp {
	return &RnewByAdminResp{
		reg.ID, reg.RenewOrderArgs, reg.TxSummary,
	}
}

type RenewService struct {
	modelOperator models.RenewOrderOperater
}

// 双镜等有注册权限的用户才可以调用
// make commit
func (r *RenewService) RenewByAdmin(req *models.RenewOrderArgs, user *models.User) (*models.Renew, error) {
	renew := models.Renew{}
	renew.RenewOrderArgs = *req
	renew.UserID = user.ID
	renew.UserPermission = user.Permission

	if err := renew.Save(models.GetDB()); err != nil {
		return nil, err
	}

	return &renew, nil
}

func (r *RenewService) MakeOrder(req *MakeRenewOrderReq) (*models.Renew, error) {
	// 过期时间
	expireTime := maxCommitmentAge.Int64() + time.Now().Unix()

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
		wecahtOrdReq := *confluxpay.NewServicesMakeOrderReq(int32(amount.Int64()), *req.Description, int32(expireTime), req.TradeType.String())

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

	RenewOrder, err := models.NewRenewOrderByPayResp(payOrder, &req.RenewOrderArgs)
	if err != nil {
		return nil, err
	}

	if err := RenewOrder.Save(models.GetDB()); err != nil {
		return nil, err
	}

	return RenewOrder, nil
}

func (r *RenewService) GetOrder(id int) (*models.Renew, error) {
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

func (r *RenewService) RefreshURL(id int) (*models.Renew, error) {
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

	if err := o.Save(models.GetDB()); err != nil {
		return nil, err
	}

	return o, nil
}
