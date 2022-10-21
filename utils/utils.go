package utils

import (
	"encoding/hex"
	"errors"

	"github.com/ethereum/go-ethereum/common"
)

var (
	ErrHashLength = errors.New("hash should be 32 bytes")
)

func Bytes2Hex(data []byte) string {
	return "0x" + common.Bytes2Hex(data)
}

func StrToHash(input string) (*common.Hash, error) {
	if input[:2] == "0x" {
		input = input[2:]
	}
	val, err := hex.DecodeString(input)
	if err != nil {
		return nil, err
	}
	if len(val) != common.HashLength {
		return nil, ErrHashLength
	}

	hash := common.BytesToHash(val)
	return &hash, nil
}
