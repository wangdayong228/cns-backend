package contracts

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/wangdayong228/cns-backend/utils"
)

var (
	genCommitABIStr           = `[{"inputs":[{"internalType":"bytes32","name":"label","type":"bytes32"},{"internalType":"address","name":"owner","type":"address"},{"internalType":"uint256","name":"duration","type":"uint256"},{"internalType":"address","name":"resolver","type":"address"},{"internalType":"bytes[]","name":"data","type":"bytes[]"},{"internalType":"bytes32","name":"secret","type":"bytes32"},{"internalType":"bool","name":"reverseRecord","type":"bool"},{"internalType":"uint32","name":"fuses","type":"uint32"},{"internalType":"uint64","name":"wrapperExpiry","type":"uint64"}],"name":"genCommitData","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"stateMutability":"pure","type":"function"},{"inputs":[{"internalType":"string","name":"name","type":"string"},{"internalType":"address","name":"owner","type":"address"},{"internalType":"uint256","name":"duration","type":"uint256"},{"internalType":"bytes32","name":"secret","type":"bytes32"},{"internalType":"address","name":"resolver","type":"address"},{"internalType":"bytes[]","name":"data","type":"bytes[]"},{"internalType":"bool","name":"reverseRecord","type":"bool"},{"internalType":"uint32","name":"fuses","type":"uint32"},{"internalType":"uint64","name":"wrapperExpiry","type":"uint64"}],"name":"makeCommitEncodeData","outputs":[{"internalType":"bytes","name":"","type":"bytes"}],"stateMutability":"pure","type":"function"},{"inputs":[{"internalType":"string","name":"name","type":"string"},{"internalType":"address","name":"owner","type":"address"},{"internalType":"uint256","name":"duration","type":"uint256"},{"internalType":"bytes32","name":"secret","type":"bytes32"},{"internalType":"address","name":"resolver","type":"address"},{"internalType":"bytes[]","name":"data","type":"bytes[]"},{"internalType":"bool","name":"reverseRecord","type":"bool"},{"internalType":"uint32","name":"fuses","type":"uint32"},{"internalType":"uint64","name":"wrapperExpiry","type":"uint64"}],"name":"makeCommitment","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"stateMutability":"pure","type":"function"}]`
	genCommitABI, _           = abi.JSON(strings.NewReader(genCommitABIStr))
	web3RegisterController, _ = abi.JSON(strings.NewReader(Web3RegisterControllerABI))
)

type DataGenerator struct {
}

type CommitArgs struct {
	Name          string
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
	label := crypto.Keccak256Hash([]byte(args.Name))
	return genCommitABI.Methods["genCommitData"].Inputs.Pack(label, args.Owner, args.Duration, args.Resolver, args.Data, args.Secret, args.ReverseRecord, args.Fuses, args.WrapperExpiry)
}

func (*DataGenerator) Register(args *CommitArgs) (string, error) {
	v, err := web3RegisterController.Pack("registerWithFiat", args.Name, args.Owner, args.Duration, args.Secret, args.Resolver, args.Data, args.ReverseRecord, args.Fuses, args.WrapperExpiry)
	return utils.Bytes2Hex(v), err
}

func (*DataGenerator) Renew(name string, duration *big.Int) (string, error) {
	v, err := web3RegisterController.Pack("renew", name, duration)
	return utils.Bytes2Hex(v), err
}
