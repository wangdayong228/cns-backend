package controllers

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/wangdayong228/cns-backend/cns_errors"
	"github.com/wangdayong228/cns-backend/models"
	"github.com/wangdayong228/cns-backend/services"
	"github.com/wangdayong228/cns-backend/utils/ginutils"
)

var (
	ErrMissingCommitHash = errors.New("missing commit hash")
)

func MakeCommits(c *gin.Context) {
	var commitCore models.CommitCore
	if err := c.ShouldBindJSON(&commitCore); err != nil {
		ginutils.RenderRespError(c, err, cns_errors.ERR_INVALID_REQUEST_COMMON)
		return
	}
	commits, err := services.MakeCommits(&commitCore)
	if err != nil {
		ginutils.RenderRespError(c, err)
		return
	}
	ginutils.RenderRespOK(c, services.MakeCommitResp{commits.CommitHash})
}

func GetCommit(c *gin.Context) {
	commitHash, ok := c.Params.Get("commit_hash")
	if !ok {
		ginutils.RenderRespError(c, ErrMissingCommitHash, cns_errors.ERR_INVALID_REQUEST_COMMON)
	}
	commit, err := services.GetCommit(commitHash)
	ginutils.RenderResp(c, commit, err)
}

func QueryCommits(c *gin.Context) {
	commitReq := &services.QueryCommitsReq{}
	if err := c.ShouldBindQuery(&commitReq); err != nil {
		ginutils.RenderRespError(c, err, cns_errors.ERR_INVALID_REQUEST_COMMON)
		return
	}

	commits, err := services.QueryCommits(commitReq)

	resp := []models.CommitCore{}
	for _, v := range commits {
		resp = append(resp, v.CommitCore)
	}

	ginutils.RenderResp(c, resp, err)
}
