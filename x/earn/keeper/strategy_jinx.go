package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/percosis-labs/fury/x/earn/types"
)

// JinxStrategy defines the strategy that deposits assets to Jinx
type JinxStrategy Keeper

var _ Strategy = (*JinxStrategy)(nil)

// GetStrategyType returns the strategy type
func (s *JinxStrategy) GetStrategyType() types.StrategyType {
	return types.STRATEGY_TYPE_JINX
}

// GetEstimatedTotalAssets returns the current value of all assets deposited
// in jinx.
func (s *JinxStrategy) GetEstimatedTotalAssets(ctx sdk.Context, denom string) (sdk.Coin, error) {
	macc := s.accountKeeper.GetModuleAccount(ctx, types.ModuleName)
	deposit, found := s.jinxKeeper.GetSyncedDeposit(ctx, macc.GetAddress())
	if !found {
		// Return 0 if no deposit exists for module account
		return sdk.NewCoin(denom, sdk.ZeroInt()), nil
	}

	// Only return the deposit for the vault denom.
	for _, coin := range deposit.Amount {
		if coin.Denom == denom {
			return coin, nil
		}
	}

	// Return 0 if no deposit exists for the vault denom
	return sdk.NewCoin(denom, sdk.ZeroInt()), nil
}

// Deposit deposits the specified amount of coins into jinx.
func (s *JinxStrategy) Deposit(ctx sdk.Context, amount sdk.Coin) error {
	macc := s.accountKeeper.GetModuleAccount(ctx, types.ModuleName)
	return s.jinxKeeper.Deposit(ctx, macc.GetAddress(), sdk.NewCoins(amount))
}

// Withdraw withdraws the specified amount of coins from jinx.
func (s *JinxStrategy) Withdraw(ctx sdk.Context, amount sdk.Coin) error {
	macc := s.accountKeeper.GetModuleAccount(ctx, types.ModuleName)
	return s.jinxKeeper.Withdraw(ctx, macc.GetAddress(), sdk.NewCoins(amount))
}
