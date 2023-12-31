package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/percosis-labs/fury/x/evmutil/keeper"
	"github.com/percosis-labs/fury/x/evmutil/testutil"
	"github.com/percosis-labs/fury/x/evmutil/types"
)

type ParamsTestSuite struct {
	testutil.Suite
}

func TestParamsSuite(t *testing.T) {
	suite.Run(t, new(ParamsTestSuite))
}

func (suite *ParamsTestSuite) TestEnabledConversionPair() {
	pairAddr := testutil.MustNewInternalEVMAddressFromString("0x8Cd469d13d804f44Ea6bB2F7d7bF8919Ae45400D")
	expPair := types.ConversionPair{
		FuryERC20Address: pairAddr.Bytes(),
		Denom:            "weth",
	}
	params := types.DefaultParams()
	params.EnabledConversionPairs = []types.ConversionPair{expPair}
	suite.Keeper.SetParams(suite.Ctx, params)

	actualPair, err := suite.Keeper.GetEnabledConversionPairFromERC20Address(
		suite.Ctx,
		pairAddr,
	)
	suite.Require().NoError(err)
	suite.Require().Equal(expPair, actualPair)
}

func (suite *ParamsTestSuite) TestHistoricParamsQuery() {
	// setup a params store that lacks allowed_cosmos_denoms param (as was the case in v1)
	oldParamStore := suite.App.GetParamsKeeper().Subspace("test_subspace_for_evmutil")
	oldParamStore.WithKeyTable(types.ParamKeyTable())
	oldParamStore.Set(suite.Ctx, types.KeyEnabledConversionPairs, types.ConversionPairs{})

	suite.True(oldParamStore.Has(suite.Ctx, types.KeyEnabledConversionPairs))
	suite.False(oldParamStore.Has(suite.Ctx, types.KeyAllowedCosmosDenoms))

	oldStateKeeper := keeper.NewKeeper(
		suite.App.AppCodec(),
		sdk.NewKVStoreKey(types.StoreKey),
		oldParamStore,
		suite.App.GetBankKeeper(),
		suite.App.GetAccountKeeper(),
	)

	// prior to making GetParams() use GetParamSetIfExists, this would panic.
	suite.NotPanics(func() {
		_ = oldStateKeeper.GetParams(suite.Ctx)
	})
}

func (suite *keeperTestSuite) TestGetAllowedTokenMetadata() {
	suite.SetupTest()

	atom := types.NewAllowedCosmosCoinERC20Token(
		"ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2",
		"Fury EVM ATOM", "ATOM", 6,
	)
	jinx := types.NewAllowedCosmosCoinERC20Token("jinx", "Fury EVM Jinx", "JINX", 6)

	// init state with some allowed tokens
	params := suite.Keeper.GetParams(suite.Ctx)
	params.AllowedCosmosDenoms = types.NewAllowedCosmosCoinERC20Tokens(atom, jinx)
	suite.Keeper.SetParams(suite.Ctx, params)

	// finds allowed tokens by denom
	storedAtom, allowed := suite.Keeper.GetAllowedTokenMetadata(suite.Ctx, atom.CosmosDenom)
	suite.True(allowed)
	suite.Equal(atom, storedAtom)
	storedJinx, allowed := suite.Keeper.GetAllowedTokenMetadata(suite.Ctx, jinx.CosmosDenom)
	suite.True(allowed)
	suite.Equal(jinx, storedJinx)

	// returns not-allowed when token not allowed
	_, allowed = suite.Keeper.GetAllowedTokenMetadata(suite.Ctx, "not-in-list")
	suite.False(allowed)
}
