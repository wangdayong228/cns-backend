package controllers

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/wangdayong228/cns-backend/cns_errors"
	"github.com/wangdayong228/cns-backend/services"
	"github.com/wangdayong228/cns-backend/utils/ginutils"
	pservice "github.com/wangdayong228/conflux-pay/services"
)

type RegisterOrderCtrl struct {
	regOrderSev *services.RegisterOrderService
}

func NewRegisterOrderCtrl() *RegisterOrderCtrl {
	return &RegisterOrderCtrl{&services.RegisterOrderService{}}
}

// CNS_BACKEND
func (r *RegisterOrderCtrl) MakeOrder(c *gin.Context) {
	var req services.MakeRegisterOrderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ginutils.RenderRespError(c, err, cns_errors.ERR_INVALID_REQUEST_COMMON)
		return
	}

	commitHash, ok := c.Params.Get("commit_hash")
	if !ok {
		ginutils.RenderRespError(c, fmt.Errorf("missing commit_hash"), cns_errors.ERR_INVALID_REQUEST_COMMON)
		return
	}

	order, err := r.regOrderSev.MakeOrder(&req, common.HexToHash(commitHash))
	if err != nil {
		ginutils.RenderRespError(c, err)
		return
	}

	resp := services.MakeRegisterOrderResp{
		CommitHash: commitHash,
		MakeOrderResp: pservice.MakeOrderResp{
			TradeProvider: order.Provider,
			TradeType:     order.TradeType,
			TradeNo:       order.TradeNo,
			CodeUrl:       order.CodeUrl,
			H5Url:         order.H5Url,
		},
	}

	ginutils.RenderRespOK(c, resp)
}

func (r *RegisterOrderCtrl) GetOrder(c *gin.Context) {
	commitHash, ok := c.Params.Get("commit_hash")
	if !ok {
		ginutils.RenderRespError(c, fmt.Errorf("missing commit_hash"), cns_errors.ERR_INVALID_REQUEST_COMMON)
		return
	}
	order, err := r.regOrderSev.GetOrder(commitHash)
	if err != nil {
		ginutils.RenderRespError(c, err)
		return
	}
	ginutils.RenderRespOK(c, order.RegisterOrderCore)
}

func (r *RegisterOrderCtrl) RefreshURL(c *gin.Context) {
	commitHash, ok := c.Params.Get("commit_hash")
	if !ok {
		ginutils.RenderRespError(c, fmt.Errorf("missing commit_hash"), cns_errors.ERR_INVALID_REQUEST_COMMON)
		return
	}
	order, err := r.regOrderSev.RefreshURL(commitHash)
	if err != nil {
		ginutils.RenderRespError(c, err)
		return
	}

	resp := services.MakeRegisterOrderResp{
		CommitHash: commitHash,
		MakeOrderResp: pservice.MakeOrderResp{
			TradeProvider: order.Provider,
			TradeType:     order.TradeType,
			TradeNo:       order.TradeNo,
			CodeUrl:       order.CodeUrl,
			H5Url:         order.H5Url,
		},
	}
	ginutils.RenderRespOK(c, resp)
}
