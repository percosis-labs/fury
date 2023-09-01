package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/percosis-labs/fury/x/incentive/types"
	jinxtypes "github.com/percosis-labs/fury/x/jinx/types"
)

// AccumulateJinxBorrowRewards calculates new rewards to distribute this block and updates the global indexes to reflect this.
// The provided rewardPeriod must be valid to avoid panics in calculating time durations.
func (k Keeper) AccumulateJinxBorrowRewards(ctx sdk.Context, rewardPeriod types.MultiRewardPeriod) {
	previousAccrualTime, found := k.GetPreviousJinxBorrowRewardAccrualTime(ctx, rewardPeriod.CollateralType)
	if !found {
		previousAccrualTime = ctx.BlockTime()
	}

	indexes, found := k.GetJinxBorrowRewardIndexes(ctx, rewardPeriod.CollateralType)
	if !found {
		indexes = types.RewardIndexes{}
	}

	acc := types.NewAccumulator(previousAccrualTime, indexes)

	totalSource := k.getJinxBorrowTotalSourceShares(ctx, rewardPeriod.CollateralType)

	acc.Accumulate(rewardPeriod, totalSource, ctx.BlockTime())

	k.SetPreviousJinxBorrowRewardAccrualTime(ctx, rewardPeriod.CollateralType, acc.PreviousAccumulationTime)
	if len(acc.Indexes) > 0 {
		// the store panics when setting empty or nil indexes
		k.SetJinxBorrowRewardIndexes(ctx, rewardPeriod.CollateralType, acc.Indexes)
	}
}

// getJinxBorrowTotalSourceShares fetches the sum of all source shares for a borrow reward.
//
// In the case of jinx borrow, this is the total borrowed divided by the borrow interest factor (for a particular denom).
// This gives the "pre interest" or "normalized" value of the total borrowed. This is an amount, that if it was borrowed when
// the interest factor was zero (ie at time 0), the current value of it with interest would be equal to the current total borrowed.
//
// The normalized borrow is also used for each individual borrow's source shares amount. Normalized amounts do not change except through
// user input. This is essential as claims must be synced before any change to a source shares amount. The actual borrowed amounts cannot
// be used as they increase every block due to interest.
func (k Keeper) getJinxBorrowTotalSourceShares(ctx sdk.Context, denom string) sdk.Dec {
	totalBorrowedCoins, found := k.jinxKeeper.GetBorrowedCoins(ctx)
	if !found {
		// assume no coins have been borrowed
		totalBorrowedCoins = sdk.NewCoins()
	}
	totalBorrowed := totalBorrowedCoins.AmountOf(denom)

	interestFactor, found := k.jinxKeeper.GetBorrowInterestFactor(ctx, denom)
	if !found {
		// assume nothing has been borrowed so the factor starts at it's default value
		interestFactor = sdk.OneDec()
	}

	// return borrowed/factor to get the "pre interest" value of the current total borrowed
	return sdk.NewDecFromInt(totalBorrowed).Quo(interestFactor)
}

// InitializeJinxBorrowReward initializes the borrow-side of a jinx liquidity provider claim
// by creating the claim and setting the borrow reward factor index
func (k Keeper) InitializeJinxBorrowReward(ctx sdk.Context, borrow jinxtypes.Borrow) {
	claim, found := k.GetJinxLiquidityProviderClaim(ctx, borrow.Borrower)
	if !found {
		claim = types.NewJinxLiquidityProviderClaim(borrow.Borrower, sdk.Coins{}, nil, nil)
	}

	var borrowRewardIndexes types.MultiRewardIndexes
	for _, coin := range borrow.Amount {
		globalRewardIndexes, found := k.GetJinxBorrowRewardIndexes(ctx, coin.Denom)
		if !found {
			globalRewardIndexes = types.RewardIndexes{}
		}
		borrowRewardIndexes = borrowRewardIndexes.With(coin.Denom, globalRewardIndexes)
	}

	claim.BorrowRewardIndexes = borrowRewardIndexes
	k.SetJinxLiquidityProviderClaim(ctx, claim)
}

// SynchronizeJinxBorrowReward updates the claim object by adding any accumulated rewards
// and updating the reward index value
func (k Keeper) SynchronizeJinxBorrowReward(ctx sdk.Context, borrow jinxtypes.Borrow) {
	claim, found := k.GetJinxLiquidityProviderClaim(ctx, borrow.Borrower)
	if !found {
		return
	}

	// Source shares for jinx borrows is their normalized borrow amount
	normalizedBorrows, err := borrow.NormalizedBorrow()
	if err != nil {
		panic(fmt.Sprintf("during borrow reward sync, could not get normalized borrow for %s: %s", borrow.Borrower, err.Error()))
	}

	for _, normedBorrow := range normalizedBorrows {
		claim = k.synchronizeSingleJinxBorrowReward(ctx, claim, normedBorrow.Denom, normedBorrow.Amount)
	}
	k.SetJinxLiquidityProviderClaim(ctx, claim)
}

// synchronizeSingleJinxBorrowReward synchronizes a single rewarded borrow denom in a jinx claim.
// It returns the claim without setting in the store.
// The public methods for accessing and modifying claims are preferred over this one. Direct modification of claims is easy to get wrong.
func (k Keeper) synchronizeSingleJinxBorrowReward(ctx sdk.Context, claim types.JinxLiquidityProviderClaim, denom string, sourceShares sdk.Dec) types.JinxLiquidityProviderClaim {
	globalRewardIndexes, found := k.GetJinxBorrowRewardIndexes(ctx, denom)
	if !found {
		// The global factor is only not found if
		// - the borrowed denom has not started accumulating rewards yet (either there is no reward specified in params, or the reward start time hasn't been hit)
		// - OR it was wrongly deleted from state (factors should never be removed while unsynced claims exist)
		// If not found we could either skip this sync, or assume the global factor is zero.
		// Skipping will avoid storing unnecessary factors in the claim for non rewarded denoms.
		// And in the event a global factor is wrongly deleted, it will avoid this function panicking when calculating rewards.
		return claim
	}

	userRewardIndexes, found := claim.BorrowRewardIndexes.Get(denom)
	if !found {
		// Normally the reward indexes should always be found.
		// But if a denom was not rewarded then becomes rewarded (ie a reward period is added to params), then the indexes will be missing from claims for that borrowed denom.
		// So given the reward period was just added, assume the starting value for any global reward indexes, which is an empty slice.
		userRewardIndexes = types.RewardIndexes{}
	}

	newRewards, err := k.CalculateRewards(userRewardIndexes, globalRewardIndexes, sourceShares)
	if err != nil {
		// Global reward factors should never decrease, as it would lead to a negative update to claim.Rewards.
		// This panics if a global reward factor decreases or disappears between the old and new indexes.
		panic(fmt.Sprintf("corrupted global reward indexes found: %v", err))
	}

	claim.Reward = claim.Reward.Add(newRewards...)
	claim.BorrowRewardIndexes = claim.BorrowRewardIndexes.With(denom, globalRewardIndexes)

	return claim
}

// UpdateJinxBorrowIndexDenoms adds or removes reward indexes from a claim to match the denoms in the borrow.
func (k Keeper) UpdateJinxBorrowIndexDenoms(ctx sdk.Context, borrow jinxtypes.Borrow) {
	claim, found := k.GetJinxLiquidityProviderClaim(ctx, borrow.Borrower)
	if !found {
		claim = types.NewJinxLiquidityProviderClaim(borrow.Borrower, sdk.Coins{}, nil, nil)
	}

	borrowDenoms := getDenoms(borrow.Amount)
	borrowRewardIndexDenoms := claim.BorrowRewardIndexes.GetCollateralTypes()

	borrowRewardIndexes := claim.BorrowRewardIndexes

	// Create a new multi-reward index in the claim for every new borrow denom
	uniqueBorrowDenoms := setDifference(borrowDenoms, borrowRewardIndexDenoms)

	for _, denom := range uniqueBorrowDenoms {
		globalBorrowRewardIndexes, found := k.GetJinxBorrowRewardIndexes(ctx, denom)
		if !found {
			globalBorrowRewardIndexes = types.RewardIndexes{}
		}
		borrowRewardIndexes = borrowRewardIndexes.With(denom, globalBorrowRewardIndexes)
	}

	// Delete multi-reward index from claim if the collateral type is no longer borrowed
	uniqueBorrowRewardDenoms := setDifference(borrowRewardIndexDenoms, borrowDenoms)

	for _, denom := range uniqueBorrowRewardDenoms {
		borrowRewardIndexes = borrowRewardIndexes.RemoveRewardIndex(denom)
	}

	claim.BorrowRewardIndexes = borrowRewardIndexes
	k.SetJinxLiquidityProviderClaim(ctx, claim)
}

// CalculateRewards computes how much rewards should have accrued to a reward source (eg a user's jinx borrowed btc amount)
// between two index values.
//
// oldIndex is normally the index stored on a claim, newIndex the current global value, and sourceShares a jinx borrowed/supplied amount.
//
// It returns an error if newIndexes does not contain all CollateralTypes from oldIndexes, or if any value of oldIndex.RewardFactor > newIndex.RewardFactor.
// This should never happen, as it would mean that a global reward index has decreased in value, or that a global reward index has been deleted from state.
func (k Keeper) CalculateRewards(oldIndexes, newIndexes types.RewardIndexes, sourceShares sdk.Dec) (sdk.Coins, error) {
	// check for missing CollateralType's
	for _, oldIndex := range oldIndexes {
		if newIndex, found := newIndexes.Get(oldIndex.CollateralType); !found {
			return nil, errorsmod.Wrapf(types.ErrDecreasingRewardFactor, "old: %v, new: %v", oldIndex, newIndex)
		}
	}
	var reward sdk.Coins
	for _, newIndex := range newIndexes {
		oldFactor, found := oldIndexes.Get(newIndex.CollateralType)
		if !found {
			oldFactor = sdk.ZeroDec()
		}

		rewardAmount, err := k.CalculateSingleReward(oldFactor, newIndex.RewardFactor, sourceShares)
		if err != nil {
			return nil, err
		}

		reward = reward.Add(
			sdk.NewCoin(newIndex.CollateralType, rewardAmount),
		)
	}
	return reward, nil
}

// CalculateSingleReward computes how much rewards should have accrued to a reward source (eg a user's btcb-a cdp principal)
// between two index values.
//
// oldIndex is normally the index stored on a claim, newIndex the current global value, and sourceShares a cdp principal amount.
//
// Returns an error if oldIndex > newIndex. This should never happen, as it would mean that a global reward index has decreased in value,
// or that a global reward index has been deleted from state.
func (k Keeper) CalculateSingleReward(oldIndex, newIndex, sourceShares sdk.Dec) (sdkmath.Int, error) {
	increase := newIndex.Sub(oldIndex)
	if increase.IsNegative() {
		return sdkmath.Int{}, errorsmod.Wrapf(types.ErrDecreasingRewardFactor, "old: %v, new: %v", oldIndex, newIndex)
	}
	reward := increase.Mul(sourceShares).RoundInt()
	return reward, nil
}
