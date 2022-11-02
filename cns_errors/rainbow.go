package cns_errors

import (
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

type CnsError int

type CnsErrorInfo struct {
	Message        string
	HttpStatusCode int
}

type CnsErrorDetailInfo struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var cnsErrorInfos = make(map[CnsError]CnsErrorInfo)

func (r CnsError) HttpStatusCode() int {
	return cnsErrorInfos[r].HttpStatusCode
}

func (r CnsError) Error() string {
	return cnsErrorInfos[r].Message
}

func (r CnsError) RenderJSON(c *gin.Context) {
	httpStatusCode := cnsErrorInfos[r].HttpStatusCode
	c.JSON(httpStatusCode, r.ErrorResponse())
}

func (r CnsError) AbortWithRenderJSON(c *gin.Context) {
	debug.PrintStack()
	c.Abort()
	r.RenderJSON(c)
}

func (r CnsError) ErrorResponse() *CnsErrorDetailInfo {
	return &CnsErrorDetailInfo{
		Code:    int(r),
		Message: r.Error(),
	}
}
