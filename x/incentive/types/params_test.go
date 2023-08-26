package types_test

import (
	"fmt"
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stretchr/testify/suite"

	"github.com/percosis-labs/fury/x/incentive/types"
)

type ParamTestSuite struct {
	suite.Suite
}

func (suite *ParamTestSuite) SetupTest() {}

var rewardPeriodWithInvalidRewardsPerSecond = types.NewRewardPeriod(
	true,
	"bnb",
	time.Date(2020, 10, 15, 14, 0, 0, 0, time.UTC),
	time.Date(2024, 10, 15, 14, 0, 0, 0, time.UTC),
	sdk.Coin{Denom: "INVALID!@#😫", Amount: sdk.ZeroInt()},
)

var rewardPeriodWithZeroRewardsPerSecond = types.NewRewardPeriod(
	true,
	"bnb",
	time.Date(2020, 10, 15, 14, 0, 0, 0, time.UTC),
	time.Date(2024, 10, 15, 14, 0, 0, 0, time.UTC),
	sdk.Coin{Denom: "ufury", Amount: sdk.ZeroInt()},
)

var rewardMultiPeriodWithInvalidRewardsPerSecond = types.NewMultiRewardPeriod(
	true,
	"bnb",
	time.Date(2020, 10, 15, 14, 0, 0, 0, time.UTC),
	time.Date(2024, 10, 15, 14, 0, 0, 0, time.UTC),
	sdk.Coins{sdk.Coin{Denom: "INVALID!@#😫", Amount: sdk.ZeroInt()}},
)

var rewardMultiPeriodWithZeroRewardsPerSecond = types.NewMultiRewardPeriod(
	true,
	"bnb",
	time.Date(2020, 10, 15, 14, 0, 0, 0, time.UTC),
	time.Date(2024, 10, 15, 14, 0, 0, 0, time.UTC),
	sdk.Coins{sdk.Coin{Denom: "zero", Amount: sdk.ZeroInt()}},
)

var validMultiRewardPeriod = types.NewMultiRewardPeriod(
	true,
	"bnb",
	time.Date(2020, 10, 15, 14, 0, 0, 0, time.UTC),
	time.Date(2024, 10, 15, 14, 0, 0, 0, time.UTC),
	sdk.NewCoins(sdk.NewInt64Coin("swap", 1e9)),
)

var validRewardPeriod = types.NewRewardPeriod(
	true,
	"bnb-a",
	time.Date(2020, 10, 15, 14, 0, 0, 0, time.UTC),
	time.Date(2024, 10, 15, 14, 0, 0, 0, time.UTC),
	sdk.NewInt64Coin(types.USDFMintingRewardDenom, 1e9),
)

func (suite *ParamTestSuite) TestParamValidation() {
	type errArgs struct {
		expectPass bool
		contains   string
	}
	type test struct {
		name    string
		params  types.Params
		errArgs errArgs
	}

	testCases := []test{
		{
			"default is valid",
			types.DefaultParams(),
			errArgs{
				expectPass: true,
			},
		},
		{
			"valid",
			types.Params{
				USDFMintingRewardPeriods: types.RewardPeriods{
					types.NewRewardPeriod(
						true,
						"bnb-a",
						time.Date(2020, 10, 15, 14, 0, 0, 0, time.UTC),
						time.Date(2024, 10, 15, 14, 0, 0, 0, time.UTC),
						sdk.NewCoin(types.USDFMintingRewardDenom, sdkmath.NewInt(122354)),
					),
				},
				JinxSupplyRewardPeriods: types.DefaultMultiRewardPeriods,
				JinxBorrowRewardPeriods: types.DefaultMultiRewardPeriods,
				DelegatorRewardPeriods:  types.DefaultMultiRewardPeriods,
				SwapRewardPeriods:       types.DefaultMultiRewardPeriods,
				SavingsRewardPeriods:    types.DefaultMultiRewardPeriods,
				ClaimMultipliers: types.MultipliersPerDenoms{
					{
						Denom: "jinx",
						Multipliers: types.Multipliers{
							types.NewMultiplier("small", 1, sdk.MustNewDecFromStr("0.25")),
							types.NewMultiplier("large", 12, sdk.MustNewDecFromStr("1.0")),
						},
					},
					{
						Denom: "ufury",
						Multipliers: types.Multipliers{
							types.NewMultiplier("small", 1, sdk.MustNewDecFromStr("0.2")),
							types.NewMultiplier("large", 12, sdk.MustNewDecFromStr("1.0")),
						},
					},
				},
				ClaimEnd: time.Date(2025, 10, 15, 14, 0, 0, 0, time.UTC),
			},
			errArgs{
				expectPass: true,
			},
		},
		{
			"invalid usdf minting period makes params invalid",
			types.Params{
				USDFMintingRewardPeriods: types.RewardPeriods{rewardPeriodWithInvalidRewardsPerSecond},
				JinxSupplyRewardPeriods:  types.DefaultMultiRewardPeriods,
				JinxBorrowRewardPeriods:  types.DefaultMultiRewardPeriods,
				DelegatorRewardPeriods:   types.DefaultMultiRewardPeriods,
				SwapRewardPeriods:        types.DefaultMultiRewardPeriods,
				SavingsRewardPeriods:     types.DefaultMultiRewardPeriods,
				ClaimMultipliers:         types.DefaultMultipliers,
				ClaimEnd:                 time.Date(2025, 10, 15, 14, 0, 0, 0, time.UTC),
			},
			errArgs{
				expectPass: false,
				contains:   fmt.Sprintf("reward denom must be %s", types.USDFMintingRewardDenom),
			},
		},
		{
			"invalid jinx supply periods makes params invalid",
			types.Params{
				USDFMintingRewardPeriods: types.DefaultRewardPeriods,
				JinxSupplyRewardPeriods:  types.MultiRewardPeriods{rewardMultiPeriodWithInvalidRewardsPerSecond},
				JinxBorrowRewardPeriods:  types.DefaultMultiRewardPeriods,
				DelegatorRewardPeriods:   types.DefaultMultiRewardPeriods,
				SwapRewardPeriods:        types.DefaultMultiRewardPeriods,
				SavingsRewardPeriods:     types.DefaultMultiRewardPeriods,
				ClaimMultipliers:         types.DefaultMultipliers,
				ClaimEnd:                 time.Date(2025, 10, 15, 14, 0, 0, 0, time.UTC),
			},
			errArgs{
				expectPass: false,
				contains:   "invalid reward amount",
			},
		},
		{
			"invalid jinx borrow periods makes params invalid",
			types.Params{
				USDFMintingRewardPeriods: types.DefaultRewardPeriods,
				JinxSupplyRewardPeriods:  types.DefaultMultiRewardPeriods,
				JinxBorrowRewardPeriods:  types.MultiRewardPeriods{rewardMultiPeriodWithInvalidRewardsPerSecond},
				DelegatorRewardPeriods:   types.DefaultMultiRewardPeriods,
				SwapRewardPeriods:        types.DefaultMultiRewardPeriods,
				SavingsRewardPeriods:     types.DefaultMultiRewardPeriods,
				ClaimMultipliers:         types.DefaultMultipliers,
				ClaimEnd:                 time.Date(2025, 10, 15, 14, 0, 0, 0, time.UTC),
			},
			errArgs{
				expectPass: false,
				contains:   "invalid reward amount",
			},
		},
		{
			"invalid delegator periods makes params invalid",
			types.Params{
				USDFMintingRewardPeriods: types.DefaultRewardPeriods,
				JinxSupplyRewardPeriods:  types.DefaultMultiRewardPeriods,
				JinxBorrowRewardPeriods:  types.DefaultMultiRewardPeriods,
				DelegatorRewardPeriods:   types.MultiRewardPeriods{rewardMultiPeriodWithInvalidRewardsPerSecond},
				SwapRewardPeriods:        types.DefaultMultiRewardPeriods,
				SavingsRewardPeriods:     types.DefaultMultiRewardPeriods,
				ClaimMultipliers:         types.DefaultMultipliers,
				ClaimEnd:                 time.Date(2025, 10, 15, 14, 0, 0, 0, time.UTC),
			},
			errArgs{
				expectPass: false,
				contains:   "invalid reward amount",
			},
		},
		{
			"invalid swap periods makes params invalid",
			types.Params{
				USDFMintingRewardPeriods: types.DefaultRewardPeriods,
				JinxSupplyRewardPeriods:  types.DefaultMultiRewardPeriods,
				JinxBorrowRewardPeriods:  types.DefaultMultiRewardPeriods,
				DelegatorRewardPeriods:   types.DefaultMultiRewardPeriods,
				SwapRewardPeriods:        types.MultiRewardPeriods{rewardMultiPeriodWithInvalidRewardsPerSecond},
				SavingsRewardPeriods:     types.DefaultMultiRewardPeriods,
				ClaimMultipliers:         types.DefaultMultipliers,
				ClaimEnd:                 time.Date(2025, 10, 15, 14, 0, 0, 0, time.UTC),
			},
			errArgs{
				expectPass: false,
				contains:   "invalid reward amount",
			},
		},
		{
			"invalid multipliers makes params invalid",
			types.Params{
				USDFMintingRewardPeriods: types.DefaultRewardPeriods,
				JinxSupplyRewardPeriods:  types.DefaultMultiRewardPeriods,
				JinxBorrowRewardPeriods:  types.DefaultMultiRewardPeriods,
				DelegatorRewardPeriods:   types.DefaultMultiRewardPeriods,
				SwapRewardPeriods:        types.DefaultMultiRewardPeriods,
				SavingsRewardPeriods:     types.DefaultMultiRewardPeriods,
				ClaimMultipliers: types.MultipliersPerDenoms{
					{
						Denom: "jinx",
						Multipliers: types.Multipliers{
							types.NewMultiplier("small", -9999, sdk.MustNewDecFromStr("0.25")),
						},
					},
				},
				ClaimEnd: time.Date(2025, 10, 15, 14, 0, 0, 0, time.UTC),
			},
			errArgs{
				expectPass: false,
				contains:   "expected non-negative lockup",
			},
		},
		{
			"invalid zero amount multi rewards per second",
			types.Params{
				USDFMintingRewardPeriods: types.DefaultRewardPeriods,
				JinxSupplyRewardPeriods: types.MultiRewardPeriods{
					rewardMultiPeriodWithZeroRewardsPerSecond,
				},
				JinxBorrowRewardPeriods: types.DefaultMultiRewardPeriods,
				DelegatorRewardPeriods:  types.DefaultMultiRewardPeriods,
				SwapRewardPeriods:       types.DefaultMultiRewardPeriods,
				SavingsRewardPeriods:    types.DefaultMultiRewardPeriods,
				ClaimMultipliers:        types.DefaultMultipliers,
				ClaimEnd:                time.Date(2025, 10, 15, 14, 0, 0, 0, time.UTC),
			},
			errArgs{
				expectPass: false,
				contains:   "invalid reward amount: 0zero",
			},
		},
		{
			"invalid zero amount single rewards per second",
			types.Params{
				USDFMintingRewardPeriods: types.RewardPeriods{
					rewardPeriodWithZeroRewardsPerSecond,
				},
				JinxSupplyRewardPeriods: types.DefaultMultiRewardPeriods,
				JinxBorrowRewardPeriods: types.DefaultMultiRewardPeriods,
				DelegatorRewardPeriods:  types.DefaultMultiRewardPeriods,
				SwapRewardPeriods:       types.DefaultMultiRewardPeriods,
				SavingsRewardPeriods:    types.DefaultMultiRewardPeriods,
				ClaimMultipliers:        types.DefaultMultipliers,
				ClaimEnd:                time.Date(2025, 10, 15, 14, 0, 0, 0, time.UTC),
			},
			errArgs{
				expectPass: false,
				contains:   "reward amount cannot be zero: 0ufury",
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			err := tc.params.Validate()

			if tc.errArgs.expectPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.errArgs.contains)
			}
		})
	}
}

func (suite *ParamTestSuite) TestRewardPeriods() {
	suite.Run("Validate", func() {
		type err struct {
			pass     bool
			contains string
		}
		testCases := []struct {
			name    string
			periods types.RewardPeriods
			expect  err
		}{
			{
				name: "single period is valid",
				periods: types.RewardPeriods{
					validRewardPeriod,
				},
				expect: err{
					pass: true,
				},
			},
			{
				name: "duplicated reward period is invalid",
				periods: types.RewardPeriods{
					validRewardPeriod,
					validRewardPeriod,
				},
				expect: err{
					contains: "duplicated reward period",
				},
			},
			{
				name: "invalid reward denom is invalid",
				periods: types.RewardPeriods{
					types.NewRewardPeriod(
						true,
						"bnb-a",
						time.Date(2020, 10, 15, 14, 0, 0, 0, time.UTC),
						time.Date(2024, 10, 15, 14, 0, 0, 0, time.UTC),
						sdk.NewInt64Coin("jinx", 1e9),
					),
				},
				expect: err{
					contains: fmt.Sprintf("reward denom must be %s", types.USDFMintingRewardDenom),
				},
			},
		}
		for _, tc := range testCases {

			err := tc.periods.Validate()

			if tc.expect.pass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
				suite.Contains(err.Error(), tc.expect.contains)
			}
		}
	})
}

func (suite *ParamTestSuite) TestMultiRewardPeriods() {
	suite.Run("Validate", func() {
		type err struct {
			pass     bool
			contains string
		}
		testCases := []struct {
			name    string
			periods types.MultiRewardPeriods
			expect  err
		}{
			{
				name: "single period is valid",
				periods: types.MultiRewardPeriods{
					validMultiRewardPeriod,
				},
				expect: err{
					pass: true,
				},
			},
			{
				name: "duplicated reward period is invalid",
				periods: types.MultiRewardPeriods{
					validMultiRewardPeriod,
					validMultiRewardPeriod,
				},
				expect: err{
					contains: "duplicated reward period",
				},
			},
			{
				name: "invalid reward period is invalid",
				periods: types.MultiRewardPeriods{
					rewardMultiPeriodWithInvalidRewardsPerSecond,
				},
				expect: err{
					contains: "invalid reward amount",
				},
			},
		}
		for _, tc := range testCases {

			err := tc.periods.Validate()

			if tc.expect.pass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
				suite.Contains(err.Error(), tc.expect.contains)
			}
		}
	})
}

func TestParamTestSuite(t *testing.T) {
	suite.Run(t, new(ParamTestSuite))
}
