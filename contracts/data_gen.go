package contracts

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/wangdayong228/cns-backend/utils"
)

var (
	web3RegisterController, _ = abi.JSON(strings.NewReader(Web3RegisterControllerABI))
)

func GenMakeCommitDataInStr(name string, owner common.Address, duration *big.Int, secret [32]byte, resolver common.Address, data [][]byte, reverseRecord bool, fuses uint32, wrapperExpiry uint64) (string, error) {
	v, err := GenMakeCommitDataInBytes(name, owner, duration, secret, resolver, data, reverseRecord, fuses, wrapperExpiry)
	return utils.Bytes2Hex(v), err
}

func GenMakeCommitDataInBytes(name string, owner common.Address, duration *big.Int, secret [32]byte, resolver common.Address, data [][]byte, reverseRecord bool, fuses uint32, wrapperExpiry uint64) ([]byte, error) {
	return web3RegisterController.Pack("makeCommitment", name, owner, duration, secret, resolver, data, reverseRecord, fuses, wrapperExpiry)
}
