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

// @Tags        Commits
// @ID          MakeCommits
// @Summary     make commit
// @Description make commit for record commit detials for using when register
// @Produce     json
// @Param       make_commit_req body     models.CommitCore true "make commit request"
// @Success     200             {object} services.MakeCommitResp
// @Failure     400             {object} cns_errors.CnsErrorDetailInfo "Invalid request"
// @Failure     500             {object} cns_errors.CnsErrorDetailInfo "Internal Server error"
// @Router      /commits [post]
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

// @Tags        Commits
// @ID          GetCommit
// @Summary     get commit
// @Description get commit details by commit hash
// @Produce     json
// @Param       commit_hash path     string true "commit hash"
// @Success     200         {object} models.CommitCore
// @Failure     400         {object} cns_errors.CnsErrorDetailInfo "Invalid request"
// @Failure     500         {object} cns_errors.CnsErrorDetailInfo "Internal Server error"
// @Router      /commits/{commit_hash} [get]
func GetCommit(c *gin.Context) {
	commitHash, ok := c.Params.Get("commit_hash")
	if !ok {
		ginutils.RenderRespError(c, ErrMissingCommitHash, cns_errors.ERR_INVALID_REQUEST_COMMON)
	}
	commit, err := services.GetCommit(commitHash)
	if err != nil {
		ginutils.RenderRespError(c, err)
		return
	}
	ginutils.RenderRespOK(c, commit.CommitCore)
}

// @Tags        Commits
// @ID          QueryCommit
// @Summary     query commit
// @Description query commit
// @Produce     json
// @Param       query_commit_request query    services.QueryCommitsReq true "query commit request"
// @Success     200                  {array}  models.CommitCore
// @Failure     400                  {object} cns_errors.CnsErrorDetailInfo "Invalid request"
// @Failure     500                  {object} cns_errors.CnsErrorDetailInfo "Internal Server error"
// @Router      /commits [get]
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
