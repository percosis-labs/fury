package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/percosis-labs/fury/x/incentive/types"
)

// InitializeJinxBorrowRewardTests runs unit tests for the keeper.InitializeJinxBorrowReward method
type InitializeJinxBorrowRewardTests struct {
	unitTester
}

func TestInitializeJinxBorrowReward(t *testing.T) {
	suite.Run(t, new(InitializeJinxBorrowRewardTests))
}

func (suite *InitializeJinxBorrowRewardTests) TestClaimIndexesAreSetWhenClaimExists() {
	claim := types.JinxLiquidityProviderClaim{
		BaseMultiClaim: types.BaseMultiClaim{
			Owner: arbitraryAddress(),
		},
		// Indexes should always be empty when initialize is called.
		// If initialize is called then the user must have repaid their borrow positions,
		// which means UpdateJinxBorrowIndexDenoms was called and should have remove indexes.
		BorrowRewardIndexes: types.MultiRewardIndexes{},
	}
	suite.storeJinxClaim(claim)

	globalIndexes := nonEmptyMultiRewardIndexes
	suite.storeGlobalBorrowIndexes(globalIndexes)

	borrow := NewBorrowBuilder(claim.Owner).
		WithArbitrarySourceShares(extractCollateralTypes(globalIndexes)...).
		Build()

	suite.keeper.InitializeJinxBorrowReward(suite.ctx, borrow)

	syncedClaim, _ := suite.keeper.GetJinxLiquidityProviderClaim(suite.ctx, claim.Owner)
	suite.Equal(globalIndexes, syncedClaim.BorrowRewardIndexes)
}

func (suite *InitializeJinxBorrowRewardTests) TestClaimIndexesAreSetWhenClaimDoesNotExist() {
	globalIndexes := nonEmptyMultiRewardIndexes
	suite.storeGlobalBorrowIndexes(globalIndexes)

	owner := arbitraryAddress()
	borrow := NewBorrowBuilder(owner).
		WithArbitrarySourceShares(extractCollateralTypes(globalIndexes)...).
		Build()

	suite.keeper.InitializeJinxBorrowReward(suite.ctx, borrow)

	syncedClaim, found := suite.keeper.GetJinxLiquidityProviderClaim(suite.ctx, owner)
	suite.True(found)
	suite.Equal(globalIndexes, syncedClaim.BorrowRewardIndexes)
}

func (suite *InitializeJinxBorrowRewardTests) TestClaimIndexesAreSetEmptyForMissingIndexes() {
	globalIndexes := nonEmptyMultiRewardIndexes
	suite.storeGlobalBorrowIndexes(globalIndexes)

	owner := arbitraryAddress()
	// Borrow a denom that is not in the global indexes.
	// This happens when a borrow denom has no rewards associated with it.
	expectedIndexes := appendUniqueEmptyMultiRewardIndex(globalIndexes)
	borrowedDenoms := extractCollateralTypes(expectedIndexes)
	borrow := NewBorrowBuilder(owner).
		WithArbitrarySourceShares(borrowedDenoms...).
		Build()

	suite.keeper.InitializeJinxBorrowReward(suite.ctx, borrow)

	syncedClaim, _ := suite.keeper.GetJinxLiquidityProviderClaim(suite.ctx, owner)
	suite.Equal(expectedIndexes, syncedClaim.BorrowRewardIndexes)
}
