package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/percosis-labs/fury/x/incentive/types"
)

// InitializeJinxSupplyRewardTests runs unit tests for the keeper.InitializeJinxSupplyReward method
type InitializeJinxSupplyRewardTests struct {
	unitTester
}

func TestInitializeJinxSupplyReward(t *testing.T) {
	suite.Run(t, new(InitializeJinxSupplyRewardTests))
}

func (suite *InitializeJinxSupplyRewardTests) TestClaimIndexesAreSetWhenClaimExists() {
	claim := types.JinxLiquidityProviderClaim{
		BaseMultiClaim: types.BaseMultiClaim{
			Owner: arbitraryAddress(),
		},
		// Indexes should always be empty when initialize is called.
		// If initialize is called then the user must have repaid their deposit positions,
		// which means UpdateJinxSupplyIndexDenoms was called and should have remove indexes.
		SupplyRewardIndexes: types.MultiRewardIndexes{},
	}
	suite.storeJinxClaim(claim)

	globalIndexes := nonEmptyMultiRewardIndexes
	suite.storeGlobalSupplyIndexes(globalIndexes)

	deposit := NewJinxDepositBuilder(claim.Owner).
		WithArbitrarySourceShares(extractCollateralTypes(globalIndexes)...).
		Build()

	suite.keeper.InitializeJinxSupplyReward(suite.ctx, deposit)

	syncedClaim, _ := suite.keeper.GetJinxLiquidityProviderClaim(suite.ctx, claim.Owner)
	suite.Equal(globalIndexes, syncedClaim.SupplyRewardIndexes)
}

func (suite *InitializeJinxSupplyRewardTests) TestClaimIndexesAreSetWhenClaimDoesNotExist() {
	globalIndexes := nonEmptyMultiRewardIndexes
	suite.storeGlobalSupplyIndexes(globalIndexes)

	owner := arbitraryAddress()
	deposit := NewJinxDepositBuilder(owner).
		WithArbitrarySourceShares(extractCollateralTypes(globalIndexes)...).
		Build()

	suite.keeper.InitializeJinxSupplyReward(suite.ctx, deposit)

	syncedClaim, found := suite.keeper.GetJinxLiquidityProviderClaim(suite.ctx, owner)
	suite.True(found)
	suite.Equal(globalIndexes, syncedClaim.SupplyRewardIndexes)
}

func (suite *InitializeJinxSupplyRewardTests) TestClaimIndexesAreSetEmptyForMissingIndexes() {
	globalIndexes := nonEmptyMultiRewardIndexes
	suite.storeGlobalSupplyIndexes(globalIndexes)

	owner := arbitraryAddress()
	// Supply a denom that is not in the global indexes.
	// This happens when a deposit denom has no rewards associated with it.
	expectedIndexes := appendUniqueEmptyMultiRewardIndex(globalIndexes)
	depositedDenoms := extractCollateralTypes(expectedIndexes)
	deposit := NewJinxDepositBuilder(owner).
		WithArbitrarySourceShares(depositedDenoms...).
		Build()

	suite.keeper.InitializeJinxSupplyReward(suite.ctx, deposit)

	syncedClaim, _ := suite.keeper.GetJinxLiquidityProviderClaim(suite.ctx, owner)
	suite.Equal(expectedIndexes, syncedClaim.SupplyRewardIndexes)
}
