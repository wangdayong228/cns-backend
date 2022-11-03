package controllers

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wangdayong228/cns-backend/cns_errors"
	"github.com/wangdayong228/cns-backend/services"
	"github.com/wangdayong228/cns-backend/utils/ginutils"
	pservice "github.com/wangdayong228/conflux-pay/services"
)

type RenewCtrl struct {
	regOrderSev *services.RenewOrderService
}

func NewRenewCtrl() *RenewCtrl {
	return &RenewCtrl{&services.RenewOrderService{}}
}

// @Tags        Renews
// @ID          MakeRenewOrder
// @Summary     make renew order
// @Description make renew order
// @Produce     json
// @Param       make_renew_order_request body     services.MakeRenewOrderReq true "make renew order request"
// @Param       commit_hash              path     string                     true "commit hash"
// @Success     200                      {object} services.MakeRenewOrderResp
// @Failure     400                      {object} cns_errors.CnsErrorDetailInfo "Invalid request"
// @Failure     500                      {object} cns_errors.CnsErrorDetailInfo "Internal Server error"
// @Router      /renews/order [post]
func (r *RenewCtrl) MakeOrder(c *gin.Context) {
	var req services.MakeRenewOrderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ginutils.RenderRespError(c, err, cns_errors.ERR_INVALID_REQUEST_COMMON)
		return
	}

	order, err := r.regOrderSev.MakeOrder(&req)
	if err != nil {
		ginutils.RenderRespError(c, err)
		return
	}

	resp := services.MakeRenewOrderResp{
		ID: int(order.ID),
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

// @Tags        Renews
// @ID          GetRenewOrder
// @Summary     get renew order
// @Description get renew order
// @Produce     json
// @Param       id  path     number true "id"
// @Success     200 {object} models.RenewOrderCore
// @Failure     400 {object} cns_errors.CnsErrorDetailInfo "Invalid request"
// @Failure     500 {object} cns_errors.CnsErrorDetailInfo "Internal Server error"
// @Router      /renews/order/{id} [get]
func (r *RenewCtrl) GetOrder(c *gin.Context) {
	id, err := getParamId(c)
	if err != nil {
		ginutils.RenderRespError(c, err, cns_errors.ERR_INVALID_REQUEST_COMMON)
		return
	}

	order, err := r.regOrderSev.GetOrder(id)
	if err != nil {
		ginutils.RenderRespError(c, err)
		return
	}
	ginutils.RenderRespOK(c, order.RenewOrderCore)
}

// @Tags        Renews
// @ID          RefreshRenewOrderUrl
// @Summary     refresh renew order url
// @Description refresh renew order url
// @Produce     json
// @Param       id  path     number true "id"
// @Success     200 {object} services.MakeRenewOrderResp
// @Failure     400 {object} cns_errors.CnsErrorDetailInfo "Invalid request"
// @Failure     500 {object} cns_errors.CnsErrorDetailInfo "Internal Server error"
// @Router      /renews/order/refresh-url/{id} [put]
func (r *RenewCtrl) RefreshURL(c *gin.Context) {
	id, err := getParamId(c)
	if err != nil {
		ginutils.RenderRespError(c, err, cns_errors.ERR_INVALID_REQUEST_COMMON)
		return
	}

	order, err := r.regOrderSev.RefreshURL(id)
	if err != nil {
		ginutils.RenderRespError(c, err)
		return
	}

	resp := services.MakeRenewOrderResp{
		ID: int(order.ID),
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

func getParamId(c *gin.Context) (int, error) {
	idStr, ok := c.Params.Get("id")
	if !ok {
		return 0, fmt.Errorf("missing id")
	}

	return strconv.Atoi(idStr)
}