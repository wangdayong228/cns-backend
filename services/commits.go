package services

import (
	"encoding/hex"
	"errors"
	"math/big"

	"github.com/Conflux-Chain/go-conflux-sdk/types/cfxaddress"
	sdkutils "github.com/Conflux-Chain/go-conflux-sdk/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/sirupsen/logrus"
	"github.com/wangdayong228/cns-backend/contracts"
	"github.com/wangdayong228/cns-backend/models"
	"github.com/wangdayong228/cns-backend/models/enums"
	"github.com/wangdayong228/cns-backend/utils"
)

var (
	ErrCommitAlreadyExists   = errors.New("commit already exists")
	ErrCommitHashNotMatch    = errors.New("commit hash not match args")
	ErrOrderStateUnrecognize = errors.New("unrecognized order state")
	dataGen                  = contracts.DataGenerator{}
)

type MakeCommitResp struct {
	CommitHash string `json:"commit_hash"`
}

type QueryCommitsReq struct {
	utils.Pagination
	OrderState *string `json:"order_state,omitempty"`
	Owner      string  `json:"owner,omitempty"`
}

func MakeCommits(c *models.CommitCore) (*models.Commit, error) {
	if v, _ := models.FindCommit(c.CommitHash); v != nil {
		return nil, ErrCommitAlreadyExists
	}

	// 1. verify commitHash is clac right
	targeHash, err := calcCommitHash(&c.CommitArgs)
	if err != nil {
		return nil, err
	}

	sourceHash := common.HexToHash(c.CommitHash)
	if sourceHash != targeHash {
		logrus.WithField("correct", targeHash).Info("commit hash not match")
		return nil, ErrCommitHashNotMatch
	}

	// 2. save
	commit := &models.Commit{CommitCore: *c}
	commit.OrderState = enums.ORDER_STATE_INIT
	if err = models.GetDB().Save(commit).Error; err != nil {
		return nil, err
	}
	return commit, nil
}

func QueryCommits(req *QueryCommitsReq) ([]*models.Commit, error) {
	cond := &models.Commit{}
	if req.OrderState != nil {
		orderState, ok := enums.ParseOrderState(*req.OrderState)
		if !ok {
			return nil, ErrOrderStateUnrecognize
		}
		cond.OrderState = *orderState
	}
	cond.Owner = req.Owner
	offset, limit := req.Pagination.CalcOffsetLimit()
	return models.FindCommits(cond, offset, limit)
}

func GetCommit(commitHash string) (*models.Commit, error) {
	return models.FindCommit(commitHash)
}

// ==================== utils ======================

func newCommitArgsForContract(c *models.CommitArgs) (*contracts.CommitArgs, error) {

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

	return &contracts.CommitArgs{
		Name:          c.Name,
		Owner:         owner.MustGetCommonAddress(),
		Duration:      duration,
		Secret:        *secretBytes,
		Resolver:      resolver.MustGetCommonAddress(),
		Data:          data,
		ReverseRecord: c.ReverseRecord,
		Fuses:         uint32(c.Fuses),
		WrapperExpiry: uint64(c.WrapperExpiry),
	}, nil
}

func genCommitData(c *models.CommitArgs) ([]byte, error) {
	arg, err := newCommitArgsForContract(c)
	if err != nil {
		return nil, err
	}
	return dataGen.GenCommit(arg)
}

func calcCommitHash(c *models.CommitArgs) (common.Hash, error) {
	data, err := genCommitData(c)
	logrus.WithField("data", hex.EncodeToString(data)).Info("gen commit data")
	if err != nil {
		return common.Hash{}, err
	}
	return crypto.Keccak256Hash(data), nil
}
