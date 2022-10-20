package utils

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/wangdayong228/cns-backend/models"
)

func CalcCommitHash(input *models.CommitArgs) common.Hash {

	labelHash := crypto.Keccak256([]byte(input.Name))

	arg := abi.Arguments{}

	arg.Pack()

	// bytes32 label = keccak256(bytes(name));
	// if (data.length > 0 && resolver == address(0)) {
	// 	revert ResolverRequiredWhenDataSupplied();
	// }
	// return
	// 	keccak256(
	// 		abi.encode(
	// 			label,
	// 			owner,
	// 			duration,
	// 			resolver,
	// 			data,
	// 			secret,
	// 			reverseRecord,
	// 			fuses,
	// 			wrapperExpiry
	// 		)
	// 	);
}
