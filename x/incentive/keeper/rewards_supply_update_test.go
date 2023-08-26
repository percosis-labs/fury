package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/percosis-labs/fury/x/incentive/types"
)

// UpdateJinxSupplyIndexDenomsTests runs unit tests for the keeper.UpdateJinxSupplyIndexDenoms method
type UpdateJinxSupplyIndexDenomsTests struct {
	unitTester
}

func TestUpdateJinxSupplyIndexDenoms(t *testing.T) {
	suite.Run(t, new(UpdateJinxSupplyIndexDenomsTests))
}

func (suite *UpdateJinxSupplyIndexDenomsTests) TestClaimIndexesAreRemovedForDenomsNoLongerSupplied() {
	claim := types.JinxLiquidityProviderClaim{
		BaseMultiClaim: types.BaseMultiClaim{
			Owner: arbitraryAddress(),
		},
		SupplyRewardIndexes: nonEmptyMultiRewardIndexes,
	}
	suite.storeJinxClaim(claim)
	suite.storeGlobalSupplyIndexes(claim.SupplyRewardIndexes)

	// remove one denom from the indexes already in the deposit
	expectedIndexes := claim.SupplyRewardIndexes[1:]
	deposit := NewJinxDepositBuilder(claim.Owner).
		WithArbitrarySourceShares(extractCollateralTypes(expectedIndexes)...).
		Build()

	suite.keeper.UpdateJinxSupplyIndexDenoms(suite.ctx, deposit)

	syncedClaim, _ := suite.keeper.GetJinxLiquidityProviderClaim(suite.ctx, claim.Owner)
	suite.Equal(expectedIndexes, syncedClaim.SupplyRewardIndexes)
}

func (suite *UpdateJinxSupplyIndexDenomsTests) TestClaimIndexesAreAddedForNewlySuppliedDenoms() {
	claim := types.JinxLiquidityProviderClaim{
		BaseMultiClaim: types.BaseMultiClaim{
			Owner: arbitraryAddress(),
		},
		SupplyRewardIndexes: nonEmptyMultiRewardIndexes,
	}
	suite.storeJinxClaim(claim)
	globalIndexes := appendUniqueMultiRewardIndex(claim.SupplyRewardIndexes)
	suite.storeGlobalSupplyIndexes(globalIndexes)

	deposit := NewJinxDepositBuilder(claim.Owner).
		WithArbitrarySourceShares(extractCollateralTypes(globalIndexes)...).
		Build()

	suite.keeper.UpdateJinxSupplyIndexDenoms(suite.ctx, deposit)

	syncedClaim, _ := suite.keeper.GetJinxLiquidityProviderClaim(suite.ctx, claim.Owner)
	suite.Equal(globalIndexes, syncedClaim.SupplyRewardIndexes)
}

func (suite *UpdateJinxSupplyIndexDenomsTests) TestClaimIndexesAreUnchangedWhenSuppliedDenomsUnchanged() {
	claim := types.JinxLiquidityProviderClaim{
		BaseMultiClaim: types.BaseMultiClaim{
			Owner: arbitraryAddress(),
		},
		SupplyRewardIndexes: nonEmptyMultiRewardIndexes,
	}
	suite.storeJinxClaim(claim)
	// Set global indexes with same denoms but different values.
	// UpdateJinxSupplyIndexDenoms should ignore the new values.
	suite.storeGlobalSupplyIndexes(increaseAllRewardFactors(claim.SupplyRewardIndexes))

	deposit := NewJinxDepositBuilder(claim.Owner).
		WithArbitrarySourceShares(extractCollateralTypes(claim.SupplyRewardIndexes)...).
		Build()

	suite.keeper.UpdateJinxSupplyIndexDenoms(suite.ctx, deposit)

	syncedClaim, _ := suite.keeper.GetJinxLiquidityProviderClaim(suite.ctx, claim.Owner)
	suite.Equal(claim.SupplyRewardIndexes, syncedClaim.SupplyRewardIndexes)
}

func (suite *UpdateJinxSupplyIndexDenomsTests) TestEmptyClaimIndexesAreAddedForNewlySuppliedButNotRewardedDenoms() {
	claim := types.JinxLiquidityProviderClaim{
		BaseMultiClaim: types.BaseMultiClaim{
			Owner: arbitraryAddress(),
		},
		SupplyRewardIndexes: nonEmptyMultiRewardIndexes,
	}
	suite.storeJinxClaim(claim)
	suite.storeGlobalSupplyIndexes(claim.SupplyRewardIndexes)

	// add a denom to the deposited amount that is not in the global or claim's indexes
	expectedIndexes := appendUniqueEmptyMultiRewardIndex(claim.SupplyRewardIndexes)
	depositedDenoms := extractCollateralTypes(expectedIndexes)
	deposit := NewJinxDepositBuilder(claim.Owner).
		WithArbitrarySourceShares(depositedDenoms...).
		Build()

	suite.keeper.UpdateJinxSupplyIndexDenoms(suite.ctx, deposit)

	syncedClaim, _ := suite.keeper.GetJinxLiquidityProviderClaim(suite.ctx, claim.Owner)
	suite.Equal(expectedIndexes, syncedClaim.SupplyRewardIndexes)
}
