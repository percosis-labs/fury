package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	cdptypes "github.com/percosis-labs/fury/x/cdp/types"
	"github.com/percosis-labs/fury/x/incentive/types"
)

// AccumulateUSDFMintingRewards calculates new rewards to distribute this block and updates the global indexes to reflect this.
// The provided rewardPeriod must be valid to avoid panics in calculating time durations.
func (k Keeper) AccumulateUSDFMintingRewards(ctx sdk.Context, rewardPeriod types.RewardPeriod) {
	previousAccrualTime, found := k.GetPreviousUSDFMintingAccrualTime(ctx, rewardPeriod.CollateralType)
	if !found {
		previousAccrualTime = ctx.BlockTime()
	}

	factor, found := k.GetUSDFMintingRewardFactor(ctx, rewardPeriod.CollateralType)
	if !found {
		factor = sdk.ZeroDec()
	}
	// wrap in RewardIndexes for compatibility with Accumulator
	indexes := types.RewardIndexes{}.With(types.USDFMintingRewardDenom, factor)

	acc := types.NewAccumulator(previousAccrualTime, indexes)

	totalSource := k.getUSDFTotalSourceShares(ctx, rewardPeriod.CollateralType)

	acc.Accumulate(types.NewMultiRewardPeriodFromRewardPeriod(rewardPeriod), totalSource, ctx.BlockTime())

	k.SetPreviousUSDFMintingAccrualTime(ctx, rewardPeriod.CollateralType, acc.PreviousAccumulationTime)

	factor, found = acc.Indexes.Get(types.USDFMintingRewardDenom)
	if !found {
		panic("could not find factor that should never be missing when accumulating usdf rewards")
	}
	k.SetUSDFMintingRewardFactor(ctx, rewardPeriod.CollateralType, factor)
}

// getUSDFTotalSourceShares fetches the sum of all source shares for a usdf minting reward.
// In the case of usdf minting, this is the total debt from all cdps of a particular type, divided by the cdp interest factor.
// This gives the "pre interest" value of the total debt.
func (k Keeper) getUSDFTotalSourceShares(ctx sdk.Context, collateralType string) sdk.Dec {
	totalPrincipal := k.cdpKeeper.GetTotalPrincipal(ctx, collateralType, cdptypes.DefaultStableDenom)

	cdpFactor, found := k.cdpKeeper.GetInterestFactor(ctx, collateralType)
	if !found {
		// assume nothing has been borrowed so the factor starts at it's default value
		cdpFactor = sdk.OneDec()
	}
	// return debt/factor to get the "pre interest" value of the current total debt
	return sdk.NewDecFromInt(totalPrincipal).Quo(cdpFactor)
}

// InitializeUSDFMintingClaim creates or updates a claim such that no new rewards are accrued, but any existing rewards are not lost.
// this function should be called after a cdp is created. If a user previously had a cdp, then closed it, they shouldn't
// accrue rewards during the period the cdp was closed. By setting the reward factor to the current global reward factor,
// any unclaimed rewards are preserved, but no new rewards are added.
func (k Keeper) InitializeUSDFMintingClaim(ctx sdk.Context, cdp cdptypes.CDP) {
	claim, found := k.GetUSDFMintingClaim(ctx, cdp.Owner)
	if !found { // this is the owner's first usdf minting reward claim
		claim = types.NewUSDFMintingClaim(cdp.Owner, sdk.NewCoin(types.USDFMintingRewardDenom, sdk.ZeroInt()), types.RewardIndexes{})
	}

	globalRewardFactor, found := k.GetUSDFMintingRewardFactor(ctx, cdp.Type)
	if !found {
		globalRewardFactor = sdk.ZeroDec()
	}
	claim.RewardIndexes = claim.RewardIndexes.With(cdp.Type, globalRewardFactor)

	k.SetUSDFMintingClaim(ctx, claim)
}

// SynchronizeUSDFMintingReward updates the claim object by adding any accumulated rewards and updating the reward index value.
// this should be called before a cdp is modified.
func (k Keeper) SynchronizeUSDFMintingReward(ctx sdk.Context, cdp cdptypes.CDP) {
	claim, found := k.GetUSDFMintingClaim(ctx, cdp.Owner)
	if !found {
		return
	}

	sourceShares, err := cdp.GetNormalizedPrincipal()
	if err != nil {
		panic(fmt.Sprintf("during usdf reward sync, could not get normalized principal for %s: %s", cdp.Owner, err.Error()))
	}

	claim = k.synchronizeSingleUSDFMintingReward(ctx, claim, cdp.Type, sourceShares)

	k.SetUSDFMintingClaim(ctx, claim)
}

// synchronizeSingleUSDFMintingReward synchronizes a single rewarded cdp collateral type in a usdf minting claim.
// It returns the claim without setting in the store.
// The public methods for accessing and modifying claims are preferred over this one. Direct modification of claims is easy to get wrong.
func (k Keeper) synchronizeSingleUSDFMintingReward(ctx sdk.Context, claim types.USDFMintingClaim, ctype string, sourceShares sdk.Dec) types.USDFMintingClaim {
	globalRewardFactor, found := k.GetUSDFMintingRewardFactor(ctx, ctype)
	if !found {
		// The global factor is only not found if
		// - the cdp collateral type has not started accumulating rewards yet (either there is no reward specified in params, or the reward start time hasn't been hit)
		// - OR it was wrongly deleted from state (factors should never be removed while unsynced claims exist)
		// If not found we could either skip this sync, or assume the global factor is zero.
		// Skipping will avoid storing unnecessary factors in the claim for non rewarded denoms.
		// And in the event a global factor is wrongly deleted, it will avoid this function panicking when calculating rewards.
		return claim
	}

	userRewardFactor, found := claim.RewardIndexes.Get(ctype)
	if !found {
		// Normally the factor should always be found, as it is added when the cdp is created in InitializeUSDFMintingClaim.
		// However if a cdp type is not rewarded then becomes rewarded (ie a reward period is added to params), existing cdps will not have the factor in their claims.
		// So assume the factor is the starting value for any global factor: 0.
		userRewardFactor = sdk.ZeroDec()
	}

	newRewardsAmount, err := k.CalculateSingleReward(userRewardFactor, globalRewardFactor, sourceShares)
	if err != nil {
		// Global reward factors should never decrease, as it would lead to a negative update to claim.Rewards.
		// This panics if a global reward factor decreases or disappears between the old and new indexes.
		panic(fmt.Sprintf("corrupted global reward indexes found: %v", err))
	}
	newRewardsCoin := sdk.NewCoin(types.USDFMintingRewardDenom, newRewardsAmount)

	claim.Reward = claim.Reward.Add(newRewardsCoin)
	claim.RewardIndexes = claim.RewardIndexes.With(ctype, globalRewardFactor)

	return claim
}

// SimulateUSDFMintingSynchronization calculates a user's outstanding USDF minting rewards by simulating reward synchronization
func (k Keeper) SimulateUSDFMintingSynchronization(ctx sdk.Context, claim types.USDFMintingClaim) types.USDFMintingClaim {
	for _, ri := range claim.RewardIndexes {
		_, found := k.GetUSDFMintingRewardPeriod(ctx, ri.CollateralType)
		if !found {
			continue
		}

		globalRewardFactor, found := k.GetUSDFMintingRewardFactor(ctx, ri.CollateralType)
		if !found {
			globalRewardFactor = sdk.ZeroDec()
		}

		// the owner has an existing usdf minting reward claim
		index, hasRewardIndex := claim.HasRewardIndex(ri.CollateralType)
		if !hasRewardIndex { // this is the owner's first usdf minting reward for this collateral type
			claim.RewardIndexes = append(claim.RewardIndexes, types.NewRewardIndex(ri.CollateralType, globalRewardFactor))
		}
		userRewardFactor := claim.RewardIndexes[index].RewardFactor
		rewardsAccumulatedFactor := globalRewardFactor.Sub(userRewardFactor)
		if rewardsAccumulatedFactor.IsZero() {
			continue
		}

		claim.RewardIndexes[index].RewardFactor = globalRewardFactor

		cdp, found := k.cdpKeeper.GetCdpByOwnerAndCollateralType(ctx, claim.GetOwner(), ri.CollateralType)
		if !found {
			continue
		}
		newRewardsAmount := rewardsAccumulatedFactor.Mul(sdk.NewDecFromInt(cdp.GetTotalPrincipal().Amount)).RoundInt()
		if newRewardsAmount.IsZero() {
			continue
		}
		newRewardsCoin := sdk.NewCoin(types.USDFMintingRewardDenom, newRewardsAmount)
		claim.Reward = claim.Reward.Add(newRewardsCoin)
	}

	return claim
}

// SynchronizeUSDFMintingClaim updates the claim object by adding any rewards that have accumulated.
// Returns the updated claim object
func (k Keeper) SynchronizeUSDFMintingClaim(ctx sdk.Context, claim types.USDFMintingClaim) (types.USDFMintingClaim, error) {
	for _, ri := range claim.RewardIndexes {
		cdp, found := k.cdpKeeper.GetCdpByOwnerAndCollateralType(ctx, claim.Owner, ri.CollateralType)
		if !found {
			// if the cdp for this collateral type has been closed, no updates are needed
			continue
		}
		claim = k.synchronizeRewardAndReturnClaim(ctx, cdp)
	}
	return claim, nil
}

// this function assumes a claim already exists, so don't call it if that's not the case
func (k Keeper) synchronizeRewardAndReturnClaim(ctx sdk.Context, cdp cdptypes.CDP) types.USDFMintingClaim {
	k.SynchronizeUSDFMintingReward(ctx, cdp)
	claim, _ := k.GetUSDFMintingClaim(ctx, cdp.Owner)
	return claim
}

// ZeroUSDFMintingClaim zeroes out the claim object's rewards and returns the updated claim object
func (k Keeper) ZeroUSDFMintingClaim(ctx sdk.Context, claim types.USDFMintingClaim) types.USDFMintingClaim {
	claim.Reward = sdk.NewCoin(claim.Reward.Denom, sdk.ZeroInt())
	k.SetUSDFMintingClaim(ctx, claim)
	return claim
}
