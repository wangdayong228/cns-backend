package cfx_tx_engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExport(t *testing.T) {
	am := NewPrivatekeyAccountManager(nil, 1)
	addr, err := am.ImportKey("7620ef675dc0df94d9081e4c4f64cc1d806927d52ca629284763d4e56d2c578b", "")
	assert.NoError(t, err)

	prv, err := am.Export(addr, "")
	assert.NoError(t, err)

	assert.Equal(t, prv, "0x7620ef675dc0df94d9081e4c4f64cc1d806927d52ca629284763d4e56d2c578b")
}
