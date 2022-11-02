package services

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
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

type MakeRegisterOrderReq struct {
	TradeProvider string           `json:"trade_provider" swaggertype:"string"`
	TradeType     penums.TradeType `json:"trade_type" binding:"required" swaggertype:"string"`
	Description   *string          `json:"description" binding:"required"`
}

type MakeRegisterOrderResp struct {
	CommitHash string `json:"commit_hash"`
	pservice.MakeOrderResp
}

type RegisterOrderService struct {
	modelOperator models.RegisterOrderOperater
}

func (r *RegisterOrderService) MakeOrder(req *MakeRegisterOrderReq, commitHash common.Hash) (*models.RegisterOrder, error) {
	//  verify
	order, err := r.modelOperator.FindRegOrderByCommitHash(commitHash.Hex())
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
	commitSubmitTime, err := web3RegController.Commitments(nil, commitHash)
	if err != nil {
		return nil, err
	}

	if commitSubmitTime.Cmp(big.NewInt(0)) == 0 {
		return nil, ErrCommitsUnsubmitOnContract
	}

	commitExpireTime := new(big.Int).Add(commitSubmitTime, maxCommitmentAge)
	if commitExpireTime.Cmp(big.NewInt(time.Now().Unix())) < 0 {
		return nil, ErrCommitsExpired
	}
	fmt.Printf("commit times %v %v %v\n", commitSubmitTime, maxCommitmentAge, commitExpireTime)

	// 获取价格
	price, err := web3RegController.RentPriceInFiat(nil, commit.Name, big.NewInt(int64(commit.Duration)))
	if err != nil {
		return nil, err
	}
	amount := new(big.Int).Add(price.Base, price.Premium)
	// fmt.Println("price 1", amount)
	amount = amount.Div(amount, big.NewInt(1e6))
	// fmt.Println("price 2", amount)
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
		wecahtOrdReq := *confluxpay.NewServicesMakeOrderReq(int32(amount.Int64()), *req.Description, int32(commitExpireTime.Int64()), req.TradeType.String())

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

	RegisterOrder, err := models.NewRegOrderByPayResp(payOrder, commitHash.Hex())
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

		if err := RegisterOrder.Save(tx); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return RegisterOrder, nil
}

func (r *RegisterOrderService) GetOrder(commitHash string) (*models.RegisterOrder, error) {
	o, err := r.modelOperator.FindRegOrderByCommitHash(commitHash)
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

	r.modelOperator.UpdateRegOrderState(commitHash, resp)
	return o, nil
}

func (r *RegisterOrderService) RefreshURL(commitHash string) (*models.RegisterOrder, error) {
	o, err := r.GetOrder(commitHash)
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
