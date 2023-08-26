package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	jinxtypes "github.com/percosis-labs/fury/x/jinx/types"
	"github.com/percosis-labs/fury/x/incentive/types"
)

// SynchronizeJinxSupplyRewardTests runs unit tests for the keeper.SynchronizeJinxSupplyReward method
type SynchronizeJinxSupplyRewardTests struct {
	unitTester
}

func TestSynchronizeJinxSupplyReward(t *testing.T) {
	suite.Run(t, new(SynchronizeJinxSupplyRewardTests))
}

func (suite *SynchronizeJinxSupplyRewardTests) TestClaimIndexesAreUpdatedWhenGlobalIndexesHaveIncreased() {
	// This is the normal case

	claim := types.JinxLiquidityProviderClaim{
		BaseMultiClaim: types.BaseMultiClaim{
			Owner: arbitraryAddress(),
		},
		SupplyRewardIndexes: nonEmptyMultiRewardIndexes,
	}
	suite.storeJinxClaim(claim)

	globalIndexes := increaseAllRewardFactors(nonEmptyMultiRewardIndexes)
	suite.storeGlobalSupplyIndexes(globalIndexes)
	deposit := NewJinxDepositBuilder(claim.Owner).
		WithArbitrarySourceShares(extractCollateralTypes(claim.SupplyRewardIndexes)...).
		Build()

	suite.keeper.SynchronizeJinxSupplyReward(suite.ctx, deposit)

	syncedClaim, _ := suite.keeper.GetJinxLiquidityProviderClaim(suite.ctx, claim.Owner)
	suite.Equal(globalIndexes, syncedClaim.SupplyRewardIndexes)
}

func (suite *SynchronizeJinxSupplyRewardTests) TestClaimIndexesAreUnchangedWhenGlobalIndexesUnchanged() {
	// It should be safe to call SynchronizeJinxSupplyReward multiple times

	unchangingIndexes := nonEmptyMultiRewardIndexes

	claim := types.JinxLiquidityProviderClaim{
		BaseMultiClaim: types.BaseMultiClaim{
			Owner: arbitraryAddress(),
		},
		SupplyRewardIndexes: unchangingIndexes,
	}
	suite.storeJinxClaim(claim)

	suite.storeGlobalSupplyIndexes(unchangingIndexes)

	deposit := NewJinxDepositBuilder(claim.Owner).
		WithArbitrarySourceShares(extractCollateralTypes(unchangingIndexes)...).
		Build()

	suite.keeper.SynchronizeJinxSupplyReward(suite.ctx, deposit)

	syncedClaim, _ := suite.keeper.GetJinxLiquidityProviderClaim(suite.ctx, claim.Owner)
	suite.Equal(unchangingIndexes, syncedClaim.SupplyRewardIndexes)
}

func (suite *SynchronizeJinxSupplyRewardTests) TestClaimIndexesAreUpdatedWhenNewRewardAdded() {
	// When a new reward is added (via gov) for a jinx deposit denom the user has already deposited, and the claim is synced;
	// Then the new reward's index should be added to the claim.

	claim := types.JinxLiquidityProviderClaim{
		BaseMultiClaim: types.BaseMultiClaim{
			Owner: arbitraryAddress(),
		},
		SupplyRewardIndexes: nonEmptyMultiRewardIndexes,
	}
	suite.storeJinxClaim(claim)

	globalIndexes := appendUniqueMultiRewardIndex(nonEmptyMultiRewardIndexes)
	suite.storeGlobalSupplyIndexes(globalIndexes)

	deposit := NewJinxDepositBuilder(claim.Owner).
		WithArbitrarySourceShares(extractCollateralTypes(globalIndexes)...).
		Build()

	suite.keeper.SynchronizeJinxSupplyReward(suite.ctx, deposit)

	syncedClaim, _ := suite.keeper.GetJinxLiquidityProviderClaim(suite.ctx, claim.Owner)
	suite.Equal(globalIndexes, syncedClaim.SupplyRewardIndexes)
}

func (suite *SynchronizeJinxSupplyRewardTests) TestClaimIndexesAreUpdatedWhenNewRewardDenomAdded() {
	// When a new reward coin is added (via gov) to an already rewarded deposit denom (that the user has already deposited), and the claim is synced;
	// Then the new reward coin's index should be added to the claim.

	claim := types.JinxLiquidityProviderClaim{
		BaseMultiClaim: types.BaseMultiClaim{
			Owner: arbitraryAddress(),
		},
		SupplyRewardIndexes: nonEmptyMultiRewardIndexes,
	}
	suite.storeJinxClaim(claim)

	globalIndexes := appendUniqueRewardIndexToFirstItem(nonEmptyMultiRewardIndexes)
	suite.storeGlobalSupplyIndexes(globalIndexes)

	deposit := NewJinxDepositBuilder(claim.Owner).
		WithArbitrarySourceShares(extractCollateralTypes(globalIndexes)...).
		Build()

	suite.keeper.SynchronizeJinxSupplyReward(suite.ctx, deposit)

	syncedClaim, _ := suite.keeper.GetJinxLiquidityProviderClaim(suite.ctx, claim.Owner)
	suite.Equal(globalIndexes, syncedClaim.SupplyRewardIndexes)
}

func (suite *SynchronizeJinxSupplyRewardTests) TestRewardIsIncrementedWhenGlobalIndexesHaveIncreased() {
	// This is the normal case
	// Given some time has passed (meaning the global indexes have increased)
	// When the claim is synced
	// The user earns rewards for the time passed

	originalReward := arbitraryCoins()

	claim := types.JinxLiquidityProviderClaim{
		BaseMultiClaim: types.BaseMultiClaim{
			Owner:  arbitraryAddress(),
			Reward: originalReward,
		},
		SupplyRewardIndexes: types.MultiRewardIndexes{
			{
				CollateralType: "depositdenom",
				RewardIndexes: types.RewardIndexes{
					{
						CollateralType: "rewarddenom",
						RewardFactor:   d("1000.001"),
					},
				},
			},
		},
	}
	suite.storeJinxClaim(claim)

	suite.storeGlobalSupplyIndexes(types.MultiRewardIndexes{
		{
			CollateralType: "depositdenom",
			RewardIndexes: types.RewardIndexes{
				{
					CollateralType: "rewarddenom",
					RewardFactor:   d("2000.002"),
				},
			},
		},
	})

	deposit := NewJinxDepositBuilder(claim.Owner).
		WithSourceShares("depositdenom", 1e9).
		Build()

	suite.keeper.SynchronizeJinxSupplyReward(suite.ctx, deposit)

	// new reward is (new index - old index) * deposit amount
	syncedClaim, _ := suite.keeper.GetJinxLiquidityProviderClaim(suite.ctx, claim.Owner)
	suite.Equal(
		cs(c("rewarddenom", 1_000_001_000_000)).Add(originalReward...),
		syncedClaim.Reward,
	)
}

func (suite *SynchronizeJinxSupplyRewardTests) TestRewardIsIncrementedWhenNewRewardAdded() {
	// When a new reward is added (via gov) for a jinx deposit denom the user has already deposited, and the claim is synced
	// Then the user earns rewards for the time since the reward was added

	originalReward := arbitraryCoins()
	claim := types.JinxLiquidityProviderClaim{
		BaseMultiClaim: types.BaseMultiClaim{
			Owner:  arbitraryAddress(),
			Reward: originalReward,
		},
		SupplyRewardIndexes: types.MultiRewardIndexes{
			{
				CollateralType: "rewarded",
				RewardIndexes: types.RewardIndexes{
					{
						CollateralType: "reward",
						RewardFactor:   d("1000.001"),
					},
				},
			},
		},
	}
	suite.storeJinxClaim(claim)

	globalIndexes := types.MultiRewardIndexes{
		{
			CollateralType: "rewarded",
			RewardIndexes: types.RewardIndexes{
				{
					CollateralType: "reward",
					RewardFactor:   d("2000.002"),
				},
			},
		},
		{
			CollateralType: "newlyrewarded",
			RewardIndexes: types.RewardIndexes{
				{
					CollateralType: "otherreward",
					// Indexes start at 0 when the reward is added by gov,
					// so this represents the syncing happening some time later.
					RewardFactor: d("1000.001"),
				},
			},
		},
	}
	suite.storeGlobalSupplyIndexes(globalIndexes)

	deposit := NewJinxDepositBuilder(claim.Owner).
		WithSourceShares("rewarded", 1e9).
		WithSourceShares("newlyrewarded", 1e9).
		Build()

	suite.keeper.SynchronizeJinxSupplyReward(suite.ctx, deposit)

	// new reward is (new index - old index) * deposit amount for each deposited denom
	// The old index for `newlyrewarded` isn't in the claim, so it's added starting at 0 for calculating the reward.
	syncedClaim, _ := suite.keeper.GetJinxLiquidityProviderClaim(suite.ctx, claim.Owner)
	suite.Equal(
		cs(c("otherreward", 1_000_001_000_000), c("reward", 1_000_001_000_000)).Add(originalReward...),
		syncedClaim.Reward,
	)
}

func (suite *SynchronizeJinxSupplyRewardTests) TestRewardIsIncrementedWhenNewRewardDenomAdded() {
	// When a new reward coin is added (via gov) to an already rewarded deposit denom (that the user has already deposited), and the claim is synced;
	// Then the user earns rewards for the time since the reward was added

	originalReward := arbitraryCoins()
	claim := types.JinxLiquidityProviderClaim{
		BaseMultiClaim: types.BaseMultiClaim{
			Owner:  arbitraryAddress(),
			Reward: originalReward,
		},
		SupplyRewardIndexes: types.MultiRewardIndexes{
			{
				CollateralType: "deposited",
				RewardIndexes: types.RewardIndexes{
					{
						CollateralType: "reward",
						RewardFactor:   d("1000.001"),
					},
				},
			},
		},
	}
	suite.storeJinxClaim(claim)

	globalIndexes := types.MultiRewardIndexes{
		{
			CollateralType: "deposited",
			RewardIndexes: types.RewardIndexes{
				{
					CollateralType: "reward",
					RewardFactor:   d("2000.002"),
				},
				{
					CollateralType: "otherreward",
					// Indexes start at 0 when the reward is added by gov,
					// so this represents the syncing happening some time later.
					RewardFactor: d("1000.001"),
				},
			},
		},
	}
	suite.storeGlobalSupplyIndexes(globalIndexes)

	deposit := NewJinxDepositBuilder(claim.Owner).
		WithSourceShares("deposited", 1e9).
		Build()

	suite.keeper.SynchronizeJinxSupplyReward(suite.ctx, deposit)

	// new reward is (new index - old index) * deposit amount for each deposited denom
	// The old index for `otherreward` isn't in the claim, so it's added starting at 0 for calculating the reward.
	syncedClaim, _ := suite.keeper.GetJinxLiquidityProviderClaim(suite.ctx, claim.Owner)
	suite.Equal(
		cs(c("reward", 1_000_001_000_000), c("otherreward", 1_000_001_000_000)).Add(originalReward...),
		syncedClaim.Reward,
	)
}

// JinxDepositBuilder is a tool for creating a jinx deposit in tests.
// The builder inherits from jinx.Deposit, so fields can be accessed directly if a helper method doesn't exist.
type JinxDepositBuilder struct {
	jinxtypes.Deposit
}

// NewJinxDepositBuilder creates a JinxDepositBuilder containing an empty deposit.
func NewJinxDepositBuilder(depositor sdk.AccAddress) JinxDepositBuilder {
	return JinxDepositBuilder{
		Deposit: jinxtypes.Deposit{
			Depositor: depositor,
		},
	}
}

// Build assembles and returns the final deposit.
func (builder JinxDepositBuilder) Build() jinxtypes.Deposit { return builder.Deposit }

// WithSourceShares adds a deposit amount and factor such that the source shares for this deposit is equal to specified.
// With a factor of 1, the deposit amount is the source shares. This picks an arbitrary factor to ensure factors are accounted for in production code.
func (builder JinxDepositBuilder) WithSourceShares(denom string, shares int64) JinxDepositBuilder {
	if !builder.Amount.AmountOf(denom).Equal(sdk.ZeroInt()) {
		panic("adding to amount with existing denom not implemented")
	}
	if _, f := builder.Index.GetInterestFactor(denom); f {
		panic("adding to indexes with existing denom not implemented")
	}

	// pick arbitrary factor
	factor := sdk.MustNewDecFromStr("2")

	// Calculate deposit amount that would equal the requested source shares given the above factor.
	amt := sdkmath.NewInt(shares).Mul(factor.RoundInt())

	builder.Amount = builder.Amount.Add(sdk.NewCoin(denom, amt))
	builder.Index = builder.Index.SetInterestFactor(denom, factor)
	return builder
}

// WithArbitrarySourceShares adds arbitrary deposit amounts and indexes for each specified denom.
func (builder JinxDepositBuilder) WithArbitrarySourceShares(denoms ...string) JinxDepositBuilder {
	const arbitraryShares = 1e9
	for _, denom := range denoms {
		builder = builder.WithSourceShares(denom, arbitraryShares)
	}
	return builder
}
