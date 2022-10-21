package services

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/wangdayong228/cns-backend/models"
)

// TODO: implement
func TestCalcCommitHash(t *testing.T) {
	commitArgs := models.CommitArgs{
		Name:          "hehe",
		Owner:         "0x1e3b4ab1d17a5e233949f64bfb15d2dc27337f8d",
		Duration:      1666003350,
		Secret:        "0x00a28927897fddc247757c1d693760b8bc88017abf7293fd36d08cc3c5a57171",
		Resolver:      "0x1e3b4ab1d17a5e233949f64bfb15d2dc27337f8d",
		Data:          nil,
		ReverseRecord: true,
		Fuses:         1,
		WrapperExpiry: 1,
	}
	target := common.HexToHash("0xc3c4622e04c642dd1726d2fa7f217b27998b42dae5d1c40e0a3c4d2e0d810046")
	actual, err := calcCommitHash(&commitArgs)
	assert.NoError(t, err)
	assert.Equal(t, target, actual)
}
