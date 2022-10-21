package contracts

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/wangdayong228/cns-backend/utils"
)

var (
	genCommitABIStr           = `[{"inputs":[{"internalType":"bytes32","name":"name","type":"bytes32"},{"internalType":"address","name":"owner","type":"address"},{"internalType":"uint256","name":"duration","type":"uint256"},{"internalType":"bytes32","name":"secret","type":"bytes32"},{"internalType":"address","name":"resolver","type":"address"},{"internalType":"bytes[]","name":"data","type":"bytes[]"},{"internalType":"bool","name":"reverseRecord","type":"bool"},{"internalType":"uint32","name":"fuses","type":"uint32"},{"internalType":"uint64","name":"wrapperExpiry","type":"uint64"}],"name":"genCommit","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"stateMutability":"pure","type":"function"}]`
	genCommitABI, _           = abi.JSON(strings.NewReader(genCommitABIStr))
	web3RegisterController, _ = abi.JSON(strings.NewReader(Web3RegisterControllerABI))
)

type DataGenerator struct {
}

type CommitArgs struct {
	Name          string
	Label         [32]byte
	Owner         common.Address
	Duration      *big.Int
	Secret        [32]byte
	Resolver      common.Address
	Data          [][]byte
	ReverseRecord bool
	Fuses         uint32
	WrapperExpiry uint64
}

func (*DataGenerator) MakeCommitment(args *CommitArgs) (string, error) {
	v, err := web3RegisterController.Pack("makeCommitment", args.Name, args.Owner, args.Duration, args.Secret, args.Resolver, args.Data, args.ReverseRecord, args.Fuses, args.WrapperExpiry)
	return utils.Bytes2Hex(v), err
}

func (*DataGenerator) GenCommit(args *CommitArgs) ([]byte, error) {
	return genCommitABI.Pack("genCommit", args.Label, args.Owner, args.Duration, args.Secret, args.Resolver, args.Data, args.ReverseRecord, args.Fuses, args.WrapperExpiry)
}

func (*DataGenerator) Register(args *CommitArgs) (string, error) {
	v, err := web3RegisterController.Pack("registerWithFiat", args.Name, args.Owner, args.Duration, args.Secret, args.Resolver, args.Data, args.ReverseRecord, args.Fuses, args.WrapperExpiry)
	return utils.Bytes2Hex(v), err
}
