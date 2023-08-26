package types_test

import (
	"testing"

	"github.com/percosis-labs/fury/x/swap/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestKeys(t *testing.T) {
	key := types.PoolKey(types.PoolID("ufury", "usdf"))
	assert.Equal(t, types.PoolID("ufury", "usdf"), string(key))

	key = types.DepositorPoolSharesKey(sdk.AccAddress("testaddress1"), types.PoolID("ufury", "usdf"))
	assert.Equal(t, string(sdk.AccAddress("testaddress1"))+"|"+types.PoolID("ufury", "usdf"), string(key))
}
