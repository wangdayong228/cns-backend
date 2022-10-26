package services

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/sirupsen/logrus"
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

var (
	ErrMakeCommithashFirst = errors.New("commitment not found, please make commit before make order")
	ErrCommitsUnexists     = errors.New("commitment invalid: not exist")
	ErrCommitsExpired      = errors.New("commitment invalid: expired")
	ErrOrderUnexists       = errors.New("order is exists, if need refresh url please invoke API 'RefreshUrl'")
	ErrOrderCompleted      = errors.New("order is completed")
)

func init() {
	configuration := confluxpay.NewConfiguration()
	configuration.Servers = confluxpay.ServerConfigurations{{
		URL:         "http://127.0.0.1:8080/v0",
		Description: "No description provided",
	}}
	confluxPayClient = confluxpay.NewAPIClient(configuration)
}

func MakeOrder(req *OrderReq, commitHash common.Hash) (*models.CnsOrder, error) {
	//  verify
	order, err := models.FindOrderByCommitHash(commitHash.Hex())
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// order 已存在，简单处理，直接返回错误，前端可以修改secret生成不同的commithash
	if order != nil {
		return nil, ErrOrderUnexists
	}

	// order 不存在
	// 1. check commitHash is valid by contract
	commit, err := models.FindCommit(commitHash.Hex())
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = ErrMakeCommithashFirst
		}
		return nil, err
	}

	// 看过期时间
	commitExpireTime, err := web3RegController.Commitments(nil, commitHash)
	if err != nil {
		return nil, err
	}

	if commitExpireTime.Cmp(big.NewInt(0)) == 0 {
		return nil, ErrCommitsUnexists
	}

	commitExpireTime = new(big.Int).Add(commitExpireTime, maxCommitmentAge)
	if commitExpireTime.Cmp(big.NewInt(time.Now().Unix())) < 0 {
		return nil, ErrCommitsExpired
	}

	// 获取价格
	price, err := web3RegController.RentPriceInFiat(nil, commit.Name, big.NewInt(int64(commit.Duration)))
	if err != nil {
		return nil, err
	}
	amount := new(big.Int).Add(price.Base, price.Premium)
	amount = amount.Div(amount, big.NewInt(1e6))
	if amount.Cmp(big.NewInt(0)) == 0 {
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
		wecahtOrdReq := *confluxpay.NewServicesMakeWechatOrderReq(int32(amount.Int64()), *req.Description, int32(commitExpireTime.Int64()), int32(req.TradeType))
		payOrder, _, err = confluxPayClient.OrdersApi.MakeOrder(context.Background()).WecahtOrdReq(wecahtOrdReq).Execute()
		if err != nil {
			logrus.WithError(err).WithField("order request", wecahtOrdReq).Info("failed to make order throught conflux-pay")
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unspport")
	}

	cnsOrder, err := models.NewOrderByPayResp(payOrder, commitHash.Hex())
	if err != nil {
		return nil, err
	}

	err = models.GetDB().Transaction(func(tx *gorm.DB) error {
		// 3. set commit order state
		commit, err := models.FindCommit(commitHash.Hex())
		if err != nil {
			return err
		}

		commit.OrderState = enums.ORDER_STATE_MADE
		if err := tx.Save(commit).Error; err != nil {
			return err
		}

		if err := tx.Save(cnsOrder).Error; err != nil {
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

func RefreshURL(commitHash string) (*models.CnsOrder, error) {
	order, err := GetOrder(commitHash)
	if err != nil {
		return nil, err
	}

	if order.TradeState.IsStable() {
		return nil, ErrOrderCompleted
	}

	resp, _, err := confluxPayClient.OrdersApi.RefreshPayUrl(context.Background(), order.TradeNo).Execute()
	if err != nil {
		return nil, err
	}

	order.CodeUrl = resp.CodeUrl
	order.H5Url = resp.H5Url
	return order, nil
}
