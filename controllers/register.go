package controllers

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/wangdayong228/cns-backend/cns_errors"
	"github.com/wangdayong228/cns-backend/models"
	"github.com/wangdayong228/cns-backend/services"
	"github.com/wangdayong228/cns-backend/utils/ginutils"
	pservice "github.com/wangdayong228/conflux-pay/services"
)

type RegisterCtrl struct {
	regSev *services.RegisterService
}

func NewRegisterCtrl() *RegisterCtrl {
	return &RegisterCtrl{&services.RegisterService{}}
}

func (r *RegisterCtrl) RegisterByAdmin(c *gin.Context) {
	var req models.CommitCore
	if err := c.ShouldBindJSON(&req); err != nil {
		ginutils.RenderRespError(c, err, cns_errors.ERR_INVALID_REQUEST_COMMON)
		return
	}
	user, ok := c.Get("user")
	if !ok {
		ginutils.RenderRespError(c, cns_errors.ERR_AUTHORIZATION_NO_PERMISSION)
		return
	}

	reg, err := r.regSev.RegisterByAdmin(&req, user.(*models.User))
	if err != nil {
		ginutils.RenderRespError(c, err)
		return
	}

	ginutils.RenderRespOK(c, services.NewRegisterByAdminRespByRaw(reg))
}

func (r *RegisterCtrl) GetRegister(c *gin.Context) {
	commitHash, ok := c.Params.Get("commit_hash")
	if !ok {
		ginutils.RenderRespError(c, fmt.Errorf("missing commit_hash"), cns_errors.ERR_INVALID_REQUEST_COMMON)
		return
	}
	reg, err := r.regSev.GetOrder(commitHash)
	if err != nil {
		ginutils.RenderRespError(c, err)
		return
	}
	ginutils.RenderRespOK(c, services.NewRegisterByAdminRespByRaw(reg))
}

// @Tags        Registers
// @ID          MakeRegisterOrder
// @Summary     make register order
// @Description make register order
// @Produce     json
// @Param       make_register_order_request body     services.MakeRegisterOrderReq true "make register order request"
// @Param       commit_hash                 path     string                        true "commit hash"
// @Success     200                         {object} services.MakeRegisterOrderResp
// @Failure     400                         {object} cns_errors.CnsErrorDetailInfo "Invalid request"
// @Failure     500                         {object} cns_errors.CnsErrorDetailInfo "Internal Server error"
// @Router      /registers/order/{commit_hash} [post]
func (r *RegisterCtrl) MakeOrder(c *gin.Context) {
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

	order, err := r.regSev.MakeOrder(&req, common.HexToHash(commitHash))
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

// @Tags        Registers
// @ID          GetRegisterOrder
// @Summary     get register order
// @Description get register order
// @Produce     json
// @Param       commit_hash path     string true "commit hash"
// @Success     200         {object} models.RegisterOrderCore
// @Failure     400         {object} cns_errors.CnsErrorDetailInfo "Invalid request"
// @Failure     500         {object} cns_errors.CnsErrorDetailInfo "Internal Server error"
// @Router      /registers/order/{commit_hash} [get]
func (r *RegisterCtrl) GetOrder(c *gin.Context) {
	commitHash, ok := c.Params.Get("commit_hash")
	if !ok {
		ginutils.RenderRespError(c, fmt.Errorf("missing commit_hash"), cns_errors.ERR_INVALID_REQUEST_COMMON)
		return
	}
	order, err := r.regSev.GetOrder(commitHash)
	if err != nil {
		ginutils.RenderRespError(c, err)
		return
	}
	ginutils.RenderRespOK(c, order.RegisterCore)
}

// @Tags        Registers
// @ID          RefreshRegisterOrderUrl
// @Summary     refresh register order url
// @Description refresh register order url
// @Produce     json
// @Param       commit_hash path     string true "commit hash"
// @Success     200         {object} services.MakeRegisterOrderResp
// @Failure     400         {object} cns_errors.CnsErrorDetailInfo "Invalid request"
// @Failure     500         {object} cns_errors.CnsErrorDetailInfo "Internal Server error"
// @Router      /registers/order/refresh-url/{commit_hash} [put]
func (r *RegisterCtrl) RefreshURL(c *gin.Context) {
	commitHash, ok := c.Params.Get("commit_hash")
	if !ok {
		ginutils.RenderRespError(c, fmt.Errorf("missing commit_hash"), cns_errors.ERR_INVALID_REQUEST_COMMON)
		return
	}
	order, err := r.regSev.RefreshURL(commitHash)
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
