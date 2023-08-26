package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/percosis-labs/fury/x/incentive/types"
)

// UpdateJinxBorrowIndexDenomsTests runs unit tests for the keeper.UpdateJinxBorrowIndexDenoms method
type UpdateJinxBorrowIndexDenomsTests struct {
	unitTester
}

func TestUpdateJinxBorrowIndexDenoms(t *testing.T) {
	suite.Run(t, new(UpdateJinxBorrowIndexDenomsTests))
}

func (suite *UpdateJinxBorrowIndexDenomsTests) TestClaimIndexesAreRemovedForDenomsNoLongerBorrowed() {
	claim := types.JinxLiquidityProviderClaim{
		BaseMultiClaim: types.BaseMultiClaim{
			Owner: arbitraryAddress(),
		},
		BorrowRewardIndexes: nonEmptyMultiRewardIndexes,
	}
	suite.storeJinxClaim(claim)
	suite.storeGlobalBorrowIndexes(claim.BorrowRewardIndexes)

	// remove one denom from the indexes already in the borrow
	expectedIndexes := claim.BorrowRewardIndexes[1:]
	borrow := NewBorrowBuilder(claim.Owner).
		WithArbitrarySourceShares(extractCollateralTypes(expectedIndexes)...).
		Build()

	suite.keeper.UpdateJinxBorrowIndexDenoms(suite.ctx, borrow)

	syncedClaim, _ := suite.keeper.GetJinxLiquidityProviderClaim(suite.ctx, claim.Owner)
	suite.Equal(expectedIndexes, syncedClaim.BorrowRewardIndexes)
}

func (suite *UpdateJinxBorrowIndexDenomsTests) TestClaimIndexesAreAddedForNewlyBorrowedDenoms() {
	claim := types.JinxLiquidityProviderClaim{
		BaseMultiClaim: types.BaseMultiClaim{
			Owner: arbitraryAddress(),
		},
		BorrowRewardIndexes: nonEmptyMultiRewardIndexes,
	}
	suite.storeJinxClaim(claim)
	globalIndexes := appendUniqueMultiRewardIndex(claim.BorrowRewardIndexes)
	suite.storeGlobalBorrowIndexes(globalIndexes)

	borrow := NewBorrowBuilder(claim.Owner).
		WithArbitrarySourceShares(extractCollateralTypes(globalIndexes)...).
		Build()

	suite.keeper.UpdateJinxBorrowIndexDenoms(suite.ctx, borrow)

	syncedClaim, _ := suite.keeper.GetJinxLiquidityProviderClaim(suite.ctx, claim.Owner)
	suite.Equal(globalIndexes, syncedClaim.BorrowRewardIndexes)
}

func (suite *UpdateJinxBorrowIndexDenomsTests) TestClaimIndexesAreUnchangedWhenBorrowedDenomsUnchanged() {
	claim := types.JinxLiquidityProviderClaim{
		BaseMultiClaim: types.BaseMultiClaim{
			Owner: arbitraryAddress(),
		},
		BorrowRewardIndexes: nonEmptyMultiRewardIndexes,
	}
	suite.storeJinxClaim(claim)
	// Set global indexes with same denoms but different values.
	// UpdateJinxBorrowIndexDenoms should ignore the new values.
	suite.storeGlobalBorrowIndexes(increaseAllRewardFactors(claim.BorrowRewardIndexes))

	borrow := NewBorrowBuilder(claim.Owner).
		WithArbitrarySourceShares(extractCollateralTypes(claim.BorrowRewardIndexes)...).
		Build()

	suite.keeper.UpdateJinxBorrowIndexDenoms(suite.ctx, borrow)

	syncedClaim, _ := suite.keeper.GetJinxLiquidityProviderClaim(suite.ctx, claim.Owner)
	suite.Equal(claim.BorrowRewardIndexes, syncedClaim.BorrowRewardIndexes)
}

func (suite *UpdateJinxBorrowIndexDenomsTests) TestEmptyClaimIndexesAreAddedForNewlyBorrowedButNotRewardedDenoms() {
	claim := types.JinxLiquidityProviderClaim{
		BaseMultiClaim: types.BaseMultiClaim{
			Owner: arbitraryAddress(),
		},
		BorrowRewardIndexes: nonEmptyMultiRewardIndexes,
	}
	suite.storeJinxClaim(claim)
	suite.storeGlobalBorrowIndexes(claim.BorrowRewardIndexes)

	// add a denom to the borrowed amount that is not in the global or claim's indexes
	expectedIndexes := appendUniqueEmptyMultiRewardIndex(claim.BorrowRewardIndexes)
	borrowedDenoms := extractCollateralTypes(expectedIndexes)
	borrow := NewBorrowBuilder(claim.Owner).
		WithArbitrarySourceShares(borrowedDenoms...).
		Build()

	suite.keeper.UpdateJinxBorrowIndexDenoms(suite.ctx, borrow)

	syncedClaim, _ := suite.keeper.GetJinxLiquidityProviderClaim(suite.ctx, claim.Owner)
	suite.Equal(expectedIndexes, syncedClaim.BorrowRewardIndexes)
}
