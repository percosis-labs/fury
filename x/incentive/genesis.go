package incentive

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/percosis-labs/fury/x/incentive/keeper"
	"github.com/percosis-labs/fury/x/incentive/types"
)

// InitGenesis initializes the store state from a genesis state.
func InitGenesis(
	ctx sdk.Context,
	k keeper.Keeper,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	cdpKeeper types.CdpKeeper,
	gs types.GenesisState,
) {
	// check if the module account exists
	moduleAcc := accountKeeper.GetModuleAccount(ctx, types.IncentiveMacc)
	if moduleAcc == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.IncentiveMacc))
	}

	if err := gs.Validate(); err != nil {
		panic(fmt.Sprintf("failed to validate %s genesis state: %s", types.ModuleName, err))
	}

	for _, rp := range gs.Params.USDFMintingRewardPeriods {
		if _, found := cdpKeeper.GetCollateral(ctx, rp.CollateralType); !found {
			panic(fmt.Sprintf("incentive params contain collateral not found in cdp params: %s", rp.CollateralType))
		}
	}
	// TODO more param validation?

	k.SetParams(ctx, gs.Params)

	// USDF Minting
	for _, claim := range gs.USDFMintingClaims {
		k.SetUSDFMintingClaim(ctx, claim)
	}
	for _, gat := range gs.USDFRewardState.AccumulationTimes {
		if err := ValidateAccumulationTime(gat.PreviousAccumulationTime); err != nil {
			panic(err.Error())
		}
		k.SetPreviousUSDFMintingAccrualTime(ctx, gat.CollateralType, gat.PreviousAccumulationTime)
	}
	for _, mri := range gs.USDFRewardState.MultiRewardIndexes {
		factor, found := mri.RewardIndexes.Get(types.USDFMintingRewardDenom)
		if !found || len(mri.RewardIndexes) != 1 {
			panic(fmt.Sprintf("USDF Minting reward factors must only have denom %s", types.USDFMintingRewardDenom))
		}
		k.SetUSDFMintingRewardFactor(ctx, mri.CollateralType, factor)
	}

	// Jinx Supply / Borrow
	for _, claim := range gs.JinxLiquidityProviderClaims {
		k.SetJinxLiquidityProviderClaim(ctx, claim)
	}
	for _, gat := range gs.JinxSupplyRewardState.AccumulationTimes {
		if err := ValidateAccumulationTime(gat.PreviousAccumulationTime); err != nil {
			panic(err.Error())
		}
		k.SetPreviousJinxSupplyRewardAccrualTime(ctx, gat.CollateralType, gat.PreviousAccumulationTime)
	}
	for _, mri := range gs.JinxSupplyRewardState.MultiRewardIndexes {
		k.SetJinxSupplyRewardIndexes(ctx, mri.CollateralType, mri.RewardIndexes)
	}
	for _, gat := range gs.JinxBorrowRewardState.AccumulationTimes {
		if err := ValidateAccumulationTime(gat.PreviousAccumulationTime); err != nil {
			panic(err.Error())
		}
		k.SetPreviousJinxBorrowRewardAccrualTime(ctx, gat.CollateralType, gat.PreviousAccumulationTime)
	}
	for _, mri := range gs.JinxBorrowRewardState.MultiRewardIndexes {
		k.SetJinxBorrowRewardIndexes(ctx, mri.CollateralType, mri.RewardIndexes)
	}

	// Delegator
	for _, claim := range gs.DelegatorClaims {
		k.SetDelegatorClaim(ctx, claim)
	}
	for _, gat := range gs.DelegatorRewardState.AccumulationTimes {
		if err := ValidateAccumulationTime(gat.PreviousAccumulationTime); err != nil {
			panic(err.Error())
		}
		k.SetPreviousDelegatorRewardAccrualTime(ctx, gat.CollateralType, gat.PreviousAccumulationTime)
	}
	for _, mri := range gs.DelegatorRewardState.MultiRewardIndexes {
		k.SetDelegatorRewardIndexes(ctx, mri.CollateralType, mri.RewardIndexes)
	}

	// Swap
	for _, claim := range gs.SwapClaims {
		k.SetSwapClaim(ctx, claim)
	}
	for _, gat := range gs.SwapRewardState.AccumulationTimes {
		if err := ValidateAccumulationTime(gat.PreviousAccumulationTime); err != nil {
			panic(err.Error())
		}
		k.SetSwapRewardAccrualTime(ctx, gat.CollateralType, gat.PreviousAccumulationTime)
	}
	for _, mri := range gs.SwapRewardState.MultiRewardIndexes {
		k.SetSwapRewardIndexes(ctx, mri.CollateralType, mri.RewardIndexes)
	}

	// Savings
	for _, claim := range gs.SavingsClaims {
		k.SetSavingsClaim(ctx, claim)
	}
	for _, gat := range gs.SavingsRewardState.AccumulationTimes {
		if err := ValidateAccumulationTime(gat.PreviousAccumulationTime); err != nil {
			panic(err.Error())
		}
		k.SetSavingsRewardAccrualTime(ctx, gat.CollateralType, gat.PreviousAccumulationTime)
	}
	for _, mri := range gs.SavingsRewardState.MultiRewardIndexes {
		k.SetSavingsRewardIndexes(ctx, mri.CollateralType, mri.RewardIndexes)
	}

	// Earn
	for _, claim := range gs.EarnClaims {
		k.SetEarnClaim(ctx, claim)
	}
	for _, gat := range gs.EarnRewardState.AccumulationTimes {
		if err := ValidateAccumulationTime(gat.PreviousAccumulationTime); err != nil {
			panic(err.Error())
		}
		k.SetEarnRewardAccrualTime(ctx, gat.CollateralType, gat.PreviousAccumulationTime)
	}
	for _, mri := range gs.EarnRewardState.MultiRewardIndexes {
		k.SetEarnRewardIndexes(ctx, mri.CollateralType, mri.RewardIndexes)
	}
}

// ExportGenesis export genesis state for incentive module
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) types.GenesisState {
	params := k.GetParams(ctx)

	usdfClaims := k.GetAllUSDFMintingClaims(ctx)
	usdfRewardState := getUSDFMintingGenesisRewardState(ctx, k)

	jinxClaims := k.GetAllJinxLiquidityProviderClaims(ctx)
	jinxSupplyRewardState := getJinxSupplyGenesisRewardState(ctx, k)
	jinxBorrowRewardState := getJinxBorrowGenesisRewardState(ctx, k)

	delegatorClaims := k.GetAllDelegatorClaims(ctx)
	delegatorRewardState := getDelegatorGenesisRewardState(ctx, k)

	swapClaims := k.GetAllSwapClaims(ctx)
	swapRewardState := getSwapGenesisRewardState(ctx, k)

	savingsClaims := k.GetAllSavingsClaims(ctx)
	savingsRewardState := getSavingsGenesisRewardState(ctx, k)

	earnClaims := k.GetAllEarnClaims(ctx)
	earnRewardState := getEarnGenesisRewardState(ctx, k)

	return types.NewGenesisState(
		params,
		// Reward states
		usdfRewardState, jinxSupplyRewardState, jinxBorrowRewardState, delegatorRewardState, swapRewardState, savingsRewardState, earnRewardState,
		// Claims
		usdfClaims, jinxClaims, delegatorClaims, swapClaims, savingsClaims, earnClaims,
	)
}

func getUSDFMintingGenesisRewardState(ctx sdk.Context, keeper keeper.Keeper) types.GenesisRewardState {
	var ats types.AccumulationTimes
	keeper.IterateUSDFMintingAccrualTimes(ctx, func(ctype string, accTime time.Time) bool {
		ats = append(ats, types.NewAccumulationTime(ctype, accTime))
		return false
	})

	var mris types.MultiRewardIndexes
	keeper.IterateUSDFMintingRewardFactors(ctx, func(ctype string, factor sdk.Dec) bool {
		mris = append(
			mris,
			types.NewMultiRewardIndex(
				ctype,
				types.RewardIndexes{types.NewRewardIndex(types.USDFMintingRewardDenom, factor)},
			),
		)
		return false
	})

	return types.NewGenesisRewardState(ats, mris)
}

func getJinxSupplyGenesisRewardState(ctx sdk.Context, keeper keeper.Keeper) types.GenesisRewardState {
	var ats types.AccumulationTimes
	keeper.IterateJinxSupplyRewardAccrualTimes(ctx, func(ctype string, accTime time.Time) bool {
		ats = append(ats, types.NewAccumulationTime(ctype, accTime))
		return false
	})

	var mris types.MultiRewardIndexes
	keeper.IterateJinxSupplyRewardIndexes(ctx, func(ctype string, indexes types.RewardIndexes) bool {
		mris = append(mris, types.NewMultiRewardIndex(ctype, indexes))
		return false
	})

	return types.NewGenesisRewardState(ats, mris)
}

func getJinxBorrowGenesisRewardState(ctx sdk.Context, keeper keeper.Keeper) types.GenesisRewardState {
	var ats types.AccumulationTimes
	keeper.IterateJinxBorrowRewardAccrualTimes(ctx, func(ctype string, accTime time.Time) bool {
		ats = append(ats, types.NewAccumulationTime(ctype, accTime))
		return false
	})

	var mris types.MultiRewardIndexes
	keeper.IterateJinxBorrowRewardIndexes(ctx, func(ctype string, indexes types.RewardIndexes) bool {
		mris = append(mris, types.NewMultiRewardIndex(ctype, indexes))
		return false
	})

	return types.NewGenesisRewardState(ats, mris)
}

func getDelegatorGenesisRewardState(ctx sdk.Context, keeper keeper.Keeper) types.GenesisRewardState {
	var ats types.AccumulationTimes
	keeper.IterateDelegatorRewardAccrualTimes(ctx, func(ctype string, accTime time.Time) bool {
		ats = append(ats, types.NewAccumulationTime(ctype, accTime))
		return false
	})

	var mris types.MultiRewardIndexes
	keeper.IterateDelegatorRewardIndexes(ctx, func(ctype string, indexes types.RewardIndexes) bool {
		mris = append(mris, types.NewMultiRewardIndex(ctype, indexes))
		return false
	})

	return types.NewGenesisRewardState(ats, mris)
}

func getSwapGenesisRewardState(ctx sdk.Context, keeper keeper.Keeper) types.GenesisRewardState {
	var ats types.AccumulationTimes
	keeper.IterateSwapRewardAccrualTimes(ctx, func(ctype string, accTime time.Time) bool {
		ats = append(ats, types.NewAccumulationTime(ctype, accTime))
		return false
	})

	var mris types.MultiRewardIndexes
	keeper.IterateSwapRewardIndexes(ctx, func(ctype string, indexes types.RewardIndexes) bool {
		mris = append(mris, types.NewMultiRewardIndex(ctype, indexes))
		return false
	})

	return types.NewGenesisRewardState(ats, mris)
}

func getSavingsGenesisRewardState(ctx sdk.Context, keeper keeper.Keeper) types.GenesisRewardState {
	var ats types.AccumulationTimes
	keeper.IterateSavingsRewardAccrualTimes(ctx, func(ctype string, accTime time.Time) bool {
		ats = append(ats, types.NewAccumulationTime(ctype, accTime))
		return false
	})

	var mris types.MultiRewardIndexes
	keeper.IterateSavingsRewardIndexes(ctx, func(ctype string, indexes types.RewardIndexes) bool {
		mris = append(mris, types.NewMultiRewardIndex(ctype, indexes))
		return false
	})

	return types.NewGenesisRewardState(ats, mris)
}

func getEarnGenesisRewardState(ctx sdk.Context, keeper keeper.Keeper) types.GenesisRewardState {
	var ats types.AccumulationTimes
	keeper.IterateEarnRewardAccrualTimes(ctx, func(ctype string, accTime time.Time) bool {
		ats = append(ats, types.NewAccumulationTime(ctype, accTime))
		return false
	})

	var mris types.MultiRewardIndexes
	keeper.IterateEarnRewardIndexes(ctx, func(ctype string, indexes types.RewardIndexes) bool {
		mris = append(mris, types.NewMultiRewardIndex(ctype, indexes))
		return false
	})

	return types.NewGenesisRewardState(ats, mris)
}

func ValidateAccumulationTime(previousAccumulationTime time.Time) error {
	if previousAccumulationTime.Equal(time.Time{}) {
		return fmt.Errorf("accumulation time is not set")
	}
	return nil
}
