package keeper_test

import (
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/percosis-labs/fury/app"
	jinxtypes "github.com/percosis-labs/fury/x/jinx/types"
	"github.com/percosis-labs/fury/x/incentive/keeper"
	"github.com/percosis-labs/fury/x/incentive/types"
	"github.com/stretchr/testify/suite"
	tmprototypes "github.com/tendermint/tendermint/proto/tendermint/types"
)

const (
	oneYear time.Duration = 365 * 24 * time.Hour
)

type grpcQueryTestSuite struct {
	suite.Suite

	tApp        app.TestApp
	ctx         sdk.Context
	keeper      keeper.Keeper
	queryClient types.QueryClient
	addrs       []sdk.AccAddress

	genesisTime  time.Time
	genesisState types.GenesisState
}

func (suite *grpcQueryTestSuite) SetupTest() {
	suite.tApp = app.NewTestApp()
	cdc := suite.tApp.AppCodec()

	_, addrs := app.GeneratePrivKeyAddressPairs(5)
	suite.genesisTime = time.Date(1998, 1, 1, 0, 0, 0, 0, time.UTC)

	suite.addrs = addrs

	suite.ctx = suite.tApp.NewContext(true, tmprototypes.Header{}).
		WithBlockTime(time.Now().UTC())
	suite.keeper = suite.tApp.GetIncentiveKeeper()

	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, suite.tApp.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, keeper.NewQueryServerImpl(suite.keeper))

	suite.queryClient = types.NewQueryClient(queryHelper)

	loanToValue, _ := sdk.NewDecFromStr("0.6")
	borrowLimit := sdk.NewDec(1000000000000000)
	jinxGS := jinxtypes.NewGenesisState(
		jinxtypes.NewParams(
			jinxtypes.MoneyMarkets{
				jinxtypes.NewMoneyMarket("ufury", jinxtypes.NewBorrowLimit(false, borrowLimit, loanToValue), "fury:usd", sdkmath.NewInt(1000000), jinxtypes.NewInterestRateModel(sdk.MustNewDecFromStr("0.05"), sdk.MustNewDecFromStr("2"), sdk.MustNewDecFromStr("0.8"), sdk.MustNewDecFromStr("10")), sdk.MustNewDecFromStr("0.05"), sdk.ZeroDec()),
				jinxtypes.NewMoneyMarket("bnb", jinxtypes.NewBorrowLimit(false, borrowLimit, loanToValue), "bnb:usd", sdkmath.NewInt(1000000), jinxtypes.NewInterestRateModel(sdk.MustNewDecFromStr("0.05"), sdk.MustNewDecFromStr("2"), sdk.MustNewDecFromStr("0.8"), sdk.MustNewDecFromStr("10")), sdk.MustNewDecFromStr("0.05"), sdk.ZeroDec()),
			},
			sdk.NewDec(10),
		),
		jinxtypes.DefaultAccumulationTimes,
		jinxtypes.DefaultDeposits,
		jinxtypes.DefaultBorrows,
		jinxtypes.DefaultTotalSupplied,
		jinxtypes.DefaultTotalBorrowed,
		jinxtypes.DefaultTotalReserves,
	)

	suite.genesisState = types.NewGenesisState(
		types.NewParams(
			types.RewardPeriods{types.NewRewardPeriod(true, "bnb-a", suite.genesisTime.Add(-1*oneYear), suite.genesisTime.Add(oneYear), c("ufury", 122354))},
			types.MultiRewardPeriods{types.NewMultiRewardPeriod(true, "bnb", suite.genesisTime.Add(-1*oneYear), suite.genesisTime.Add(oneYear), cs(c("jinx", 122354)))},
			types.MultiRewardPeriods{types.NewMultiRewardPeriod(true, "bnb", suite.genesisTime.Add(-1*oneYear), suite.genesisTime.Add(oneYear), cs(c("jinx", 122354)))},
			types.MultiRewardPeriods{types.NewMultiRewardPeriod(true, "ufury", suite.genesisTime.Add(-1*oneYear), suite.genesisTime.Add(oneYear), cs(c("jinx", 122354)))},
			types.MultiRewardPeriods{types.NewMultiRewardPeriod(true, "btcb/usdf", suite.genesisTime.Add(-1*oneYear), suite.genesisTime.Add(oneYear), cs(c("mer", 122354)))},
			types.MultiRewardPeriods{types.NewMultiRewardPeriod(true, "ufury", suite.genesisTime.Add(-1*oneYear), suite.genesisTime.Add(oneYear), cs(c("jinx", 122354)))},
			types.MultiRewardPeriods{types.NewMultiRewardPeriod(true, "ufury", suite.genesisTime.Add(-1*oneYear), suite.genesisTime.Add(oneYear), cs(c("jinx", 122354)))},
			types.MultipliersPerDenoms{
				{
					Denom: "ufury",
					Multipliers: types.Multipliers{
						types.NewMultiplier("large", 12, d("1.0")),
					},
				},
				{
					Denom: "jinx",
					Multipliers: types.Multipliers{
						types.NewMultiplier("small", 1, d("0.25")),
						types.NewMultiplier("large", 12, d("1.0")),
					},
				},
				{
					Denom: "mer",
					Multipliers: types.Multipliers{
						types.NewMultiplier("small", 1, d("0.25")),
						types.NewMultiplier("medium", 6, d("0.8")),
					},
				},
			},
			suite.genesisTime.Add(5*oneYear),
		),
		types.NewGenesisRewardState(
			types.AccumulationTimes{
				types.NewAccumulationTime("bnb-a", suite.genesisTime),
			},
			types.MultiRewardIndexes{
				types.NewMultiRewardIndex("bnb-a", types.RewardIndexes{{CollateralType: "ufury", RewardFactor: d("0.3")}}),
			},
		),
		types.NewGenesisRewardState(
			types.AccumulationTimes{
				types.NewAccumulationTime("bnb", suite.genesisTime.Add(-1*time.Hour)),
			},
			types.MultiRewardIndexes{
				types.NewMultiRewardIndex("bnb", types.RewardIndexes{{CollateralType: "jinx", RewardFactor: d("0.1")}}),
			},
		),
		types.NewGenesisRewardState(
			types.AccumulationTimes{
				types.NewAccumulationTime("bnb", suite.genesisTime.Add(-2*time.Hour)),
			},
			types.MultiRewardIndexes{
				types.NewMultiRewardIndex("bnb", types.RewardIndexes{{CollateralType: "jinx", RewardFactor: d("0.05")}}),
			},
		),
		types.NewGenesisRewardState(
			types.AccumulationTimes{
				types.NewAccumulationTime("ufury", suite.genesisTime.Add(-3*time.Hour)),
			},
			types.MultiRewardIndexes{
				types.NewMultiRewardIndex("ufury", types.RewardIndexes{{CollateralType: "jinx", RewardFactor: d("0.2")}}),
			},
		),
		types.NewGenesisRewardState(
			types.AccumulationTimes{
				types.NewAccumulationTime("bctb/usdf", suite.genesisTime.Add(-4*time.Hour)),
			},
			types.MultiRewardIndexes{
				types.NewMultiRewardIndex("btcb/usdf", types.RewardIndexes{{CollateralType: "swap", RewardFactor: d("0.001")}}),
			},
		),
		types.NewGenesisRewardState(
			types.AccumulationTimes{
				types.NewAccumulationTime("ufury", suite.genesisTime.Add(-3*time.Hour)),
			},
			types.MultiRewardIndexes{
				types.NewMultiRewardIndex("ufury", types.RewardIndexes{{CollateralType: "ufury", RewardFactor: d("0.2")}}),
			},
		),
		types.NewGenesisRewardState(
			types.AccumulationTimes{
				types.NewAccumulationTime("usdf", suite.genesisTime.Add(-3*time.Hour)),
			},
			types.MultiRewardIndexes{
				types.NewMultiRewardIndex("usdf", types.RewardIndexes{{CollateralType: "usdf", RewardFactor: d("0.2")}}),
			},
		),
		types.USDFMintingClaims{
			types.NewUSDFMintingClaim(
				suite.addrs[0],
				c("ufury", 1e9),
				types.RewardIndexes{{CollateralType: "bnb-a", RewardFactor: d("0.3")}},
			),
			types.NewUSDFMintingClaim(
				suite.addrs[1],
				c("ufury", 1),
				types.RewardIndexes{{CollateralType: "bnb-a", RewardFactor: d("0.001")}},
			),
		},
		types.JinxLiquidityProviderClaims{
			types.NewJinxLiquidityProviderClaim(
				suite.addrs[0],
				cs(c("ufury", 1e9), c("jinx", 1e9)),
				types.MultiRewardIndexes{{CollateralType: "bnb", RewardIndexes: types.RewardIndexes{{CollateralType: "jinx", RewardFactor: d("0.01")}}}},
				types.MultiRewardIndexes{{CollateralType: "bnb", RewardIndexes: types.RewardIndexes{{CollateralType: "jinx", RewardFactor: d("0.0")}}}},
			),
			types.NewJinxLiquidityProviderClaim(
				suite.addrs[1],
				cs(c("jinx", 1)),
				types.MultiRewardIndexes{{CollateralType: "bnb", RewardIndexes: types.RewardIndexes{{CollateralType: "jinx", RewardFactor: d("0.1")}}}},
				types.MultiRewardIndexes{{CollateralType: "bnb", RewardIndexes: types.RewardIndexes{{CollateralType: "jinx", RewardFactor: d("0.0")}}}},
			),
		},
		types.DelegatorClaims{
			types.NewDelegatorClaim(
				suite.addrs[2],
				cs(c("jinx", 5)),
				types.MultiRewardIndexes{{CollateralType: "ufury", RewardIndexes: types.RewardIndexes{{CollateralType: "jinx", RewardFactor: d("0.2")}}}},
			),
		},
		types.SwapClaims{
			types.NewSwapClaim(
				suite.addrs[3],
				nil,
				types.MultiRewardIndexes{{CollateralType: "btcb/usdf", RewardIndexes: types.RewardIndexes{{CollateralType: "swap", RewardFactor: d("0.0")}}}},
			),
		},
		types.SavingsClaims{
			types.NewSavingsClaim(
				suite.addrs[3],
				nil,
				types.MultiRewardIndexes{{CollateralType: "ufury", RewardIndexes: types.RewardIndexes{{CollateralType: "ufury", RewardFactor: d("0.0")}}}},
			),
		},
		types.EarnClaims{
			types.NewEarnClaim(
				suite.addrs[3],
				nil,
				types.MultiRewardIndexes{{CollateralType: "usdf", RewardIndexes: types.RewardIndexes{{CollateralType: "usdf", RewardFactor: d("0.0")}}}},
			),
		},
	)

	err := suite.genesisState.Validate()
	suite.Require().NoError(err)

	suite.tApp = suite.tApp.InitializeFromGenesisStatesWithTime(
		suite.genesisTime,
		app.GenesisState{types.ModuleName: cdc.MustMarshalJSON(&suite.genesisState)},
		app.GenesisState{jinxtypes.ModuleName: cdc.MustMarshalJSON(&jinxGS)},
		NewCDPGenStateMulti(cdc),
		NewPricefeedGenStateMultiFromTime(cdc, suite.genesisTime),
	)

	suite.tApp.DeleteGenesisValidator(suite.T(), suite.ctx)
	claims := suite.keeper.GetAllDelegatorClaims(suite.ctx)
	for _, claim := range claims {
		// Delete the InitGenesis validator's claim
		if !claim.Owner.Equals(suite.addrs[2]) {
			suite.keeper.DeleteDelegatorClaim(suite.ctx, claim.Owner)
		}
	}
}

func (suite *grpcQueryTestSuite) TestGrpcQueryParams() {
	res, err := suite.queryClient.Params(sdk.WrapSDKContext(suite.ctx), &types.QueryParamsRequest{})
	suite.Require().NoError(err)

	expected := suite.keeper.GetParams(suite.ctx)

	suite.Equal(expected, res.Params, "params should equal default params")
}

func (suite *grpcQueryTestSuite) TestGrpcQueryRewards() {
	res, err := suite.queryClient.Rewards(sdk.WrapSDKContext(suite.ctx), &types.QueryRewardsRequest{
		Unsynchronized: true,
	})
	suite.Require().NoError(err)

	suite.Equal(suite.genesisState.USDFMintingClaims, res.USDFMintingClaims)
	suite.Equal(suite.genesisState.JinxLiquidityProviderClaims, res.JinxLiquidityProviderClaims)
	suite.Equal(suite.genesisState.DelegatorClaims, res.DelegatorClaims)
	suite.Equal(suite.genesisState.SwapClaims, res.SwapClaims)
	suite.Equal(suite.genesisState.SavingsClaims, res.SavingsClaims)
	suite.Equal(suite.genesisState.EarnClaims, res.EarnClaims)
}

func (suite *grpcQueryTestSuite) TestGrpcQueryRewards_Owner() {
	res, err := suite.queryClient.Rewards(sdk.WrapSDKContext(suite.ctx), &types.QueryRewardsRequest{
		Owner: suite.addrs[0].String(),
	})
	suite.Require().NoError(err)

	suite.Len(res.USDFMintingClaims, 1)
	suite.Len(res.JinxLiquidityProviderClaims, 1)

	suite.Equal(suite.genesisState.USDFMintingClaims[0], res.USDFMintingClaims[0])
	suite.Equal(suite.genesisState.JinxLiquidityProviderClaims[0], res.JinxLiquidityProviderClaims[0])

	// No other claims - owner has none
	suite.Empty(res.DelegatorClaims)
	suite.Empty(res.SwapClaims)
	suite.Empty(res.SavingsClaims)
	suite.Empty(res.EarnClaims)
}

func (suite *grpcQueryTestSuite) TestGrpcQueryRewards_RewardType() {
	res, err := suite.queryClient.Rewards(sdk.WrapSDKContext(suite.ctx), &types.QueryRewardsRequest{
		RewardType:     keeper.RewardTypeJinx,
		Unsynchronized: true,
	})
	suite.Require().NoError(err)

	suite.Equal(suite.genesisState.JinxLiquidityProviderClaims, res.JinxLiquidityProviderClaims)

	// No other reward types when specifying rewardType
	suite.Empty(res.USDFMintingClaims)
	suite.Empty(res.DelegatorClaims)
	suite.Empty(res.SwapClaims)
	suite.Empty(res.SavingsClaims)
	suite.Empty(res.EarnClaims)
}

func (suite *grpcQueryTestSuite) TestGrpcQueryRewards_RewardType_and_Owner() {
	res, err := suite.queryClient.Rewards(sdk.WrapSDKContext(suite.ctx), &types.QueryRewardsRequest{
		Owner:      suite.addrs[0].String(),
		RewardType: keeper.RewardTypeJinx,
	})
	suite.Require().NoError(err)

	suite.Len(res.JinxLiquidityProviderClaims, 1)
	suite.Equal(suite.genesisState.JinxLiquidityProviderClaims[0], res.JinxLiquidityProviderClaims[0])

	suite.Empty(res.USDFMintingClaims)
	suite.Empty(res.DelegatorClaims)
	suite.Empty(res.SwapClaims)
	suite.Empty(res.SavingsClaims)
	suite.Empty(res.EarnClaims)
}

func (suite *grpcQueryTestSuite) TestGrpcQueryRewardFactors() {
	res, err := suite.queryClient.RewardFactors(sdk.WrapSDKContext(suite.ctx), &types.QueryRewardFactorsRequest{})
	suite.Require().NoError(err)

	suite.NotEmpty(res.UsdfMintingRewardFactors)
	suite.NotEmpty(res.JinxSupplyRewardFactors)
	suite.NotEmpty(res.JinxBorrowRewardFactors)
	suite.NotEmpty(res.DelegatorRewardFactors)
	suite.NotEmpty(res.SwapRewardFactors)
	suite.NotEmpty(res.SavingsRewardFactors)
	suite.NotEmpty(res.EarnRewardFactors)
}

func TestGrpcQueryTestSuite(t *testing.T) {
	suite.Run(t, new(grpcQueryTestSuite))
}
