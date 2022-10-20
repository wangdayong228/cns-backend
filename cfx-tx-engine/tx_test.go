package cfx_tx_engine

import (
	"fmt"
	"testing"

	sdk "github.com/Conflux-Chain/go-conflux-sdk"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/wangdayong228/cns-backend/models"
	// "github.com/stretchr/testify/assert"
)

func TestBatchSend(t *testing.T) {
	client, err := sdk.NewClient("https://test.confluxrpc.com")
	assert.NoError(t, err)

	client.SetAccountManager(NewPrivatekeyAccountManager(nil, 1))

	client.AccountManager.ImportKey("7620ef675dc0df94d9081e4c4f64cc1d806927d52ca629284763d4e56d2c578b", "")
	// client.AccountManager.Unlock(cfxaddress.MustNew("cfxtest:aatk708nbb7573bkwumsu00h0r1rtkcdz2chwhttzk"), "")

	txs := []*models.Transaction{
		// PENDING -> EXECUTED
		{BaseModel: models.BaseModel{ID: 0}, From: "cfxtest:aatk708nbb7573bkwumsu00h0r1rtkcdz2chwhttzk", To: "cfxtest:aatk708nbb7573bkwumsu00h0r1rtkcdz2chwhttzk", Value: decimal.NewFromInt(0)},
		// PENDING -> EXECUTED FAIL
		{BaseModel: models.BaseModel{ID: 1}, From: "cfxtest:aatk708nbb7573bkwumsu00h0r1rtkcdz2chwhttzk", Value: decimal.NewFromInt(0)},
		// SEND FAILED
		{BaseModel: models.BaseModel{ID: 2}, From: "cfxtest:aatk708nbb7573bkwumsu00h0r1rtkcdz2chwhttzk", To: "cfxtest:aatk708nbb7573bkwumsu00h0r1rtkcdz2chwhttzk", Value: decimal.NewFromInt(1e18).Mul(decimal.NewFromInt(2000))},
	}
	bSender := NewTransactionSender(client, 20)
	h := StateChangeHandler(func(tx *models.Transaction, oldstate models.TxState, newstate models.TxState) {
		fmt.Printf("%v: oldstate = %v, newstate = %v\n", tx.ID, oldstate, newstate)
	})
	bSender.RegisterStateChangeEvent(&h)

	bSender.BulkSend(txs)
}
