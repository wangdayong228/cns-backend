package services

import (
	"errors"
	"math/big"

	"github.com/Conflux-Chain/go-conflux-sdk/types/cfxaddress"
	sdkutils "github.com/Conflux-Chain/go-conflux-sdk/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/wangdayong228/cns-backend/contracts"
	"github.com/wangdayong228/cns-backend/models"
	"github.com/wangdayong228/cns-backend/utils"
)

var (
	ErrCommitHashWrong = errors.New("commit hash is wrong")
	dataGen            = contracts.DataGenerator{}
)

type QueryCommitsReq struct {
	Skip        int    `json:"skip,omitempty"`
	Limit       int    `json:"limit,omitempty"`
	IsOrderMade *bool  `json:"is_order_made,omitempty"`
	Owner       string `json:"owner,omitempty"`
}

func MakeCommits(c *models.CommitCore) (*models.Commit, error) {
	// 1. verify commitHash is clac right
	targeHash, err := calcCommitHash(&c.CommitArgs)
	if err != nil {
		return nil, err
	}

	sourceHash, _ := utils.StrToHash(c.CommitHash)
	if *sourceHash != targeHash {
		return nil, ErrCommitHashWrong
	}

	// 2. save
	commit := &models.Commit{CommitCore: *c}
	models.GetDB().Save(commit)
	return commit, nil
}

func QueryCommits(req *QueryCommitsReq) ([]*models.Commit, error) {
	cond := &models.Commit{}
	if req.IsOrderMade != nil {
		cond.IsOrderMade = *req.IsOrderMade
	}
	cond.Owner = req.Owner
	return models.FindCommits(cond, req.Skip, req.Limit)
}

func GetCommit(commitHash string) (*models.Commit, error) {
	return models.FindCommit(commitHash)
}

// ==================== utils ======================

func genCommitData(c *models.CommitArgs) ([]byte, error) {

	label := crypto.Keccak256Hash([]byte(c.Name))

	owner, err := cfxaddress.New(c.Owner)
	if err != nil {
		return nil, err
	}

	resolver, err := cfxaddress.New(c.Resolver)
	if err != nil {
		return nil, err
	}

	duration := big.NewInt(int64(c.Duration))

	secretBytes, err := utils.StrToHash(c.Secret)
	if err != nil {
		return nil, err
	}

	data := [][]byte{}
	for _, d := range c.Data {
		bytes, err := sdkutils.HexStringToBytes(d)
		if err != nil {
			return nil, err
		}
		data = append(data, bytes)
	}

	return dataGen.GenCommit(label, owner.MustGetCommonAddress(), duration, *secretBytes, resolver.MustGetCommonAddress(), data, c.ReverseRecord, uint32(c.Fuses), uint64(c.WrapperExpiry))
}

func calcCommitHash(c *models.CommitArgs) (common.Hash, error) {
	data, err := genCommitData(c)
	if err != nil {
		return common.Hash{}, err
	}
	return crypto.Keccak256Hash(data), nil
}
