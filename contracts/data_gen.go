package contracts

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/wangdayong228/cns-backend/utils"
)

var (
	genCommitABIStr           = `{"inputs":[{"internalType":"bytes32","name":"name","type":"bytes32"},{"internalType":"address","name":"owner","type":"address"},{"internalType":"uint256","name":"duration","type":"uint256"},{"internalType":"bytes32","name":"secret","type":"bytes32"},{"internalType":"address","name":"resolver","type":"address"},{"internalType":"bytes[]","name":"data","type":"bytes[]"},{"internalType":"bool","name":"reverseRecord","type":"bool"},{"internalType":"uint32","name":"fuses","type":"uint32"},{"internalType":"uint64","name":"wrapperExpiry","type":"uint64"}],"name":"genCommit","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"stateMutability":"pure","type":"function"}`
	genCommitABI, _           = abi.JSON(strings.NewReader(genCommitABIStr))
	web3RegisterController, _ = abi.JSON(strings.NewReader(Web3RegisterControllerABI))
)

type DataGenerator struct {
}

func (*DataGenerator) MakeCommitment(name string, owner common.Address, duration *big.Int, secret [32]byte, resolver common.Address, data [][]byte, reverseRecord bool, fuses uint32, wrapperExpiry uint64) (string, error) {
	v, err := web3RegisterController.Pack("makeCommitment", name, owner, duration, secret, resolver, data, reverseRecord, fuses, wrapperExpiry)
	return utils.Bytes2Hex(v), err
}

func (*DataGenerator) GenCommit(label [32]byte, owner common.Address, duration *big.Int, secret [32]byte, resolver common.Address, data [][]byte, reverseRecord bool, fuses uint32, wrapperExpiry uint64) ([]byte, error) {
	return genCommitABI.Pack("genCommitHash", label, owner, duration, secret, resolver, data, reverseRecord, fuses, wrapperExpiry)
}
