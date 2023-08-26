package incentive

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/percosis-labs/fury/x/incentive/keeper"
	"github.com/percosis-labs/fury/x/incentive/types"
)

// BeginBlocker runs at the start of every block
func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)

	params := k.GetParams(ctx)

	for _, rp := range params.USDFMintingRewardPeriods {
		k.AccumulateUSDFMintingRewards(ctx, rp)
	}
	for _, rp := range params.JinxSupplyRewardPeriods {
		k.AccumulateJinxSupplyRewards(ctx, rp)
	}
	for _, rp := range params.JinxBorrowRewardPeriods {
		k.AccumulateJinxBorrowRewards(ctx, rp)
	}
	for _, rp := range params.DelegatorRewardPeriods {
		k.AccumulateDelegatorRewards(ctx, rp)
	}
	for _, rp := range params.SwapRewardPeriods {
		k.AccumulateSwapRewards(ctx, rp)
	}
	for _, rp := range params.SavingsRewardPeriods {
		k.AccumulateSavingsRewards(ctx, rp)
	}
	for _, rp := range params.EarnRewardPeriods {
		if err := k.AccumulateEarnRewards(ctx, rp); err != nil {
			panic(fmt.Sprintf("failed to accumulate earn rewards: %s", err))
		}
	}
}
