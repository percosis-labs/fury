package keeper

import (
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/percosis-labs/fury/x/incentive/types"
)

// Keeper keeper for the incentive module
type Keeper struct {
	cdc           codec.Codec
	key           storetypes.StoreKey
	paramSubspace types.ParamSubspace
	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
	cdpKeeper     types.CdpKeeper
	jinxKeeper    types.JinxKeeper
	stakingKeeper types.StakingKeeper
	swapKeeper    types.SwapKeeper
	savingsKeeper types.SavingsKeeper
	liquidKeeper  types.LiquidKeeper
	earnKeeper    types.EarnKeeper

	// Keepers used for APY queries
	mintKeeper      types.MintKeeper
	distrKeeper     types.DistrKeeper
	pricefeedKeeper types.PricefeedKeeper
}

// NewKeeper creates a new keeper
func NewKeeper(
	cdc codec.Codec, key storetypes.StoreKey, paramstore types.ParamSubspace, bk types.BankKeeper,
	cdpk types.CdpKeeper, hk types.JinxKeeper, ak types.AccountKeeper, stk types.StakingKeeper,
	merk types.SwapKeeper, svk types.SavingsKeeper, lqk types.LiquidKeeper, ek types.EarnKeeper,
	mk types.MintKeeper, dk types.DistrKeeper, pfk types.PricefeedKeeper,
) Keeper {
	if !paramstore.HasKeyTable() {
		paramstore = paramstore.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		accountKeeper: ak,
		cdc:           cdc,
		key:           key,
		paramSubspace: paramstore,
		bankKeeper:    bk,
		cdpKeeper:     cdpk,
		jinxKeeper:    hk,
		stakingKeeper: stk,
		swapKeeper:    merk,
		savingsKeeper: svk,
		liquidKeeper:  lqk,
		earnKeeper:    ek,

		mintKeeper:      mk,
		distrKeeper:     dk,
		pricefeedKeeper: pfk,
	}
}

// GetUSDFMintingClaim returns the claim in the store corresponding the the input address collateral type and id and a boolean for if the claim was found
func (k Keeper) GetUSDFMintingClaim(ctx sdk.Context, addr sdk.AccAddress) (types.USDFMintingClaim, bool) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.USDFMintingClaimKeyPrefix)
	bz := store.Get(addr)
	if bz == nil {
		return types.USDFMintingClaim{}, false
	}
	var c types.USDFMintingClaim
	k.cdc.MustUnmarshal(bz, &c)
	return c, true
}

// SetUSDFMintingClaim sets the claim in the store corresponding to the input address, collateral type, and id
func (k Keeper) SetUSDFMintingClaim(ctx sdk.Context, c types.USDFMintingClaim) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.USDFMintingClaimKeyPrefix)
	bz := k.cdc.MustMarshal(&c)
	store.Set(c.Owner, bz)
}

// DeleteUSDFMintingClaim deletes the claim in the store corresponding to the input address, collateral type, and id
func (k Keeper) DeleteUSDFMintingClaim(ctx sdk.Context, owner sdk.AccAddress) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.USDFMintingClaimKeyPrefix)
	store.Delete(owner)
}

// IterateUSDFMintingClaims iterates over all claim  objects in the store and preforms a callback function
func (k Keeper) IterateUSDFMintingClaims(ctx sdk.Context, cb func(c types.USDFMintingClaim) (stop bool)) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.USDFMintingClaimKeyPrefix)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var c types.USDFMintingClaim
		k.cdc.MustUnmarshal(iterator.Value(), &c)
		if cb(c) {
			break
		}
	}
}

// GetAllUSDFMintingClaims returns all Claim objects in the store
func (k Keeper) GetAllUSDFMintingClaims(ctx sdk.Context) types.USDFMintingClaims {
	cs := types.USDFMintingClaims{}
	k.IterateUSDFMintingClaims(ctx, func(c types.USDFMintingClaim) (stop bool) {
		cs = append(cs, c)
		return false
	})
	return cs
}

// GetPreviousUSDFMintingAccrualTime returns the last time a collateral type accrued USDF minting rewards
func (k Keeper) GetPreviousUSDFMintingAccrualTime(ctx sdk.Context, ctype string) (blockTime time.Time, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.PreviousUSDFMintingRewardAccrualTimeKeyPrefix)
	b := store.Get([]byte(ctype))
	if b == nil {
		return time.Time{}, false
	}
	if err := blockTime.UnmarshalBinary(b); err != nil {
		panic(err)
	}
	return blockTime, true
}

// SetPreviousUSDFMintingAccrualTime sets the last time a collateral type accrued USDF minting rewards
func (k Keeper) SetPreviousUSDFMintingAccrualTime(ctx sdk.Context, ctype string, blockTime time.Time) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.PreviousUSDFMintingRewardAccrualTimeKeyPrefix)
	bz, err := blockTime.MarshalBinary()
	if err != nil {
		panic(err)
	}
	store.Set([]byte(ctype), bz)
}

// IterateUSDFMintingAccrualTimes iterates over all previous USDF minting accrual times and preforms a callback function
func (k Keeper) IterateUSDFMintingAccrualTimes(ctx sdk.Context, cb func(string, time.Time) (stop bool)) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.PreviousUSDFMintingRewardAccrualTimeKeyPrefix)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var accrualTime time.Time
		if err := accrualTime.UnmarshalBinary(iterator.Value()); err != nil {
			panic(err)
		}
		denom := string(iterator.Key())
		if cb(denom, accrualTime) {
			break
		}
	}
}

// GetUSDFMintingRewardFactor returns the current reward factor for an individual collateral type
func (k Keeper) GetUSDFMintingRewardFactor(ctx sdk.Context, ctype string) (factor sdk.Dec, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.USDFMintingRewardFactorKeyPrefix)
	bz := store.Get([]byte(ctype))
	if bz == nil {
		return sdk.ZeroDec(), false
	}
	if err := factor.Unmarshal(bz); err != nil {
		panic(err)
	}
	return factor, true
}

// SetUSDFMintingRewardFactor sets the current reward factor for an individual collateral type
func (k Keeper) SetUSDFMintingRewardFactor(ctx sdk.Context, ctype string, factor sdk.Dec) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.USDFMintingRewardFactorKeyPrefix)
	bz, err := factor.Marshal()
	if err != nil {
		panic(err)
	}
	store.Set([]byte(ctype), bz)
}

// IterateUSDFMintingRewardFactors iterates over all USDF Minting reward factor objects in the store and preforms a callback function
func (k Keeper) IterateUSDFMintingRewardFactors(ctx sdk.Context, cb func(denom string, factor sdk.Dec) (stop bool)) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.USDFMintingRewardFactorKeyPrefix)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var factor sdk.Dec
		if err := factor.Unmarshal(iterator.Value()); err != nil {
			panic(err)
		}
		if cb(string(iterator.Key()), factor) {
			break
		}
	}
}

// GetJinxLiquidityProviderClaim returns the claim in the store corresponding the the input address collateral type and id and a boolean for if the claim was found
func (k Keeper) GetJinxLiquidityProviderClaim(ctx sdk.Context, addr sdk.AccAddress) (types.JinxLiquidityProviderClaim, bool) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.JinxLiquidityClaimKeyPrefix)
	bz := store.Get(addr)
	if bz == nil {
		return types.JinxLiquidityProviderClaim{}, false
	}
	var c types.JinxLiquidityProviderClaim
	k.cdc.MustUnmarshal(bz, &c)
	return c, true
}

// SetJinxLiquidityProviderClaim sets the claim in the store corresponding to the input address, collateral type, and id
func (k Keeper) SetJinxLiquidityProviderClaim(ctx sdk.Context, c types.JinxLiquidityProviderClaim) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.JinxLiquidityClaimKeyPrefix)
	bz := k.cdc.MustMarshal(&c)
	store.Set(c.Owner, bz)
}

// DeleteJinxLiquidityProviderClaim deletes the claim in the store corresponding to the input address, collateral type, and id
func (k Keeper) DeleteJinxLiquidityProviderClaim(ctx sdk.Context, owner sdk.AccAddress) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.JinxLiquidityClaimKeyPrefix)
	store.Delete(owner)
}

// IterateJinxLiquidityProviderClaims iterates over all claim  objects in the store and preforms a callback function
func (k Keeper) IterateJinxLiquidityProviderClaims(ctx sdk.Context, cb func(c types.JinxLiquidityProviderClaim) (stop bool)) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.JinxLiquidityClaimKeyPrefix)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var c types.JinxLiquidityProviderClaim
		k.cdc.MustUnmarshal(iterator.Value(), &c)
		if cb(c) {
			break
		}
	}
}

// GetAllJinxLiquidityProviderClaims returns all Claim objects in the store
func (k Keeper) GetAllJinxLiquidityProviderClaims(ctx sdk.Context) types.JinxLiquidityProviderClaims {
	cs := types.JinxLiquidityProviderClaims{}
	k.IterateJinxLiquidityProviderClaims(ctx, func(c types.JinxLiquidityProviderClaim) (stop bool) {
		cs = append(cs, c)
		return false
	})
	return cs
}

// GetDelegatorClaim returns the claim in the store corresponding the the input address collateral type and id and a boolean for if the claim was found
func (k Keeper) GetDelegatorClaim(ctx sdk.Context, addr sdk.AccAddress) (types.DelegatorClaim, bool) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.DelegatorClaimKeyPrefix)
	bz := store.Get(addr)
	if bz == nil {
		return types.DelegatorClaim{}, false
	}
	var c types.DelegatorClaim
	k.cdc.MustUnmarshal(bz, &c)
	return c, true
}

// SetDelegatorClaim sets the claim in the store corresponding to the input address, collateral type, and id
func (k Keeper) SetDelegatorClaim(ctx sdk.Context, c types.DelegatorClaim) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.DelegatorClaimKeyPrefix)
	bz := k.cdc.MustMarshal(&c)
	store.Set(c.Owner, bz)
}

// DeleteDelegatorClaim deletes the claim in the store corresponding to the input address, collateral type, and id
func (k Keeper) DeleteDelegatorClaim(ctx sdk.Context, owner sdk.AccAddress) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.DelegatorClaimKeyPrefix)
	store.Delete(owner)
}

// IterateDelegatorClaims iterates over all claim  objects in the store and preforms a callback function
func (k Keeper) IterateDelegatorClaims(ctx sdk.Context, cb func(c types.DelegatorClaim) (stop bool)) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.DelegatorClaimKeyPrefix)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var c types.DelegatorClaim
		k.cdc.MustUnmarshal(iterator.Value(), &c)
		if cb(c) {
			break
		}
	}
}

// GetAllDelegatorClaims returns all DelegatorClaim objects in the store
func (k Keeper) GetAllDelegatorClaims(ctx sdk.Context) types.DelegatorClaims {
	cs := types.DelegatorClaims{}
	k.IterateDelegatorClaims(ctx, func(c types.DelegatorClaim) (stop bool) {
		cs = append(cs, c)
		return false
	})
	return cs
}

// GetSwapClaim returns the claim in the store corresponding the the input address.
func (k Keeper) GetSwapClaim(ctx sdk.Context, addr sdk.AccAddress) (types.SwapClaim, bool) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.SwapClaimKeyPrefix)
	bz := store.Get(addr)
	if bz == nil {
		return types.SwapClaim{}, false
	}
	var c types.SwapClaim
	k.cdc.MustUnmarshal(bz, &c)
	return c, true
}

// SetSwapClaim sets the claim in the store corresponding to the input address.
func (k Keeper) SetSwapClaim(ctx sdk.Context, c types.SwapClaim) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.SwapClaimKeyPrefix)
	bz := k.cdc.MustMarshal(&c)
	store.Set(c.Owner, bz)
}

// DeleteSwapClaim deletes the claim in the store corresponding to the input address.
func (k Keeper) DeleteSwapClaim(ctx sdk.Context, owner sdk.AccAddress) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.SwapClaimKeyPrefix)
	store.Delete(owner)
}

// IterateSwapClaims iterates over all claim  objects in the store and preforms a callback function
func (k Keeper) IterateSwapClaims(ctx sdk.Context, cb func(c types.SwapClaim) (stop bool)) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.SwapClaimKeyPrefix)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var c types.SwapClaim
		k.cdc.MustUnmarshal(iterator.Value(), &c)
		if cb(c) {
			break
		}
	}
}

// GetAllSwapClaims returns all Claim objects in the store
func (k Keeper) GetAllSwapClaims(ctx sdk.Context) types.SwapClaims {
	cs := types.SwapClaims{}
	k.IterateSwapClaims(ctx, func(c types.SwapClaim) (stop bool) {
		cs = append(cs, c)
		return false
	})
	return cs
}

// GetSavingsClaim returns the claim in the store corresponding the the input address.
func (k Keeper) GetSavingsClaim(ctx sdk.Context, addr sdk.AccAddress) (types.SavingsClaim, bool) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.SavingsClaimKeyPrefix)
	bz := store.Get(addr)
	if bz == nil {
		return types.SavingsClaim{}, false
	}
	var c types.SavingsClaim
	k.cdc.MustUnmarshal(bz, &c)
	return c, true
}

// SetSavingsClaim sets the claim in the store corresponding to the input address.
func (k Keeper) SetSavingsClaim(ctx sdk.Context, c types.SavingsClaim) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.SavingsClaimKeyPrefix)
	bz := k.cdc.MustMarshal(&c)
	store.Set(c.Owner, bz)
}

// DeleteSavingsClaim deletes the claim in the store corresponding to the input address.
func (k Keeper) DeleteSavingsClaim(ctx sdk.Context, owner sdk.AccAddress) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.SavingsClaimKeyPrefix)
	store.Delete(owner)
}

// IterateSavingsClaims iterates over all savings claim objects in the store and preforms a callback function
func (k Keeper) IterateSavingsClaims(ctx sdk.Context, cb func(c types.SavingsClaim) (stop bool)) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.SavingsClaimKeyPrefix)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var c types.SavingsClaim
		k.cdc.MustUnmarshal(iterator.Value(), &c)
		if cb(c) {
			break
		}
	}
}

// GetAllSavingsClaims returns all savings claim objects in the store
func (k Keeper) GetAllSavingsClaims(ctx sdk.Context) types.SavingsClaims {
	cs := types.SavingsClaims{}
	k.IterateSavingsClaims(ctx, func(c types.SavingsClaim) (stop bool) {
		cs = append(cs, c)
		return false
	})
	return cs
}

// GetEarnClaim returns the claim in the store corresponding the the input address.
func (k Keeper) GetEarnClaim(ctx sdk.Context, addr sdk.AccAddress) (types.EarnClaim, bool) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.EarnClaimKeyPrefix)
	bz := store.Get(addr)
	if bz == nil {
		return types.EarnClaim{}, false
	}
	var c types.EarnClaim
	k.cdc.MustUnmarshal(bz, &c)
	return c, true
}

// SetEarnClaim sets the claim in the store corresponding to the input address.
func (k Keeper) SetEarnClaim(ctx sdk.Context, c types.EarnClaim) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.EarnClaimKeyPrefix)
	bz := k.cdc.MustMarshal(&c)
	store.Set(c.Owner, bz)
}

// DeleteEarnClaim deletes the claim in the store corresponding to the input address.
func (k Keeper) DeleteEarnClaim(ctx sdk.Context, owner sdk.AccAddress) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.EarnClaimKeyPrefix)
	store.Delete(owner)
}

// IterateEarnClaims iterates over all claim  objects in the store and preforms a callback function
func (k Keeper) IterateEarnClaims(ctx sdk.Context, cb func(c types.EarnClaim) (stop bool)) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.EarnClaimKeyPrefix)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var c types.EarnClaim
		k.cdc.MustUnmarshal(iterator.Value(), &c)
		if cb(c) {
			break
		}
	}
}

// GetAllEarnClaims returns all Claim objects in the store
func (k Keeper) GetAllEarnClaims(ctx sdk.Context) types.EarnClaims {
	cs := types.EarnClaims{}
	k.IterateEarnClaims(ctx, func(c types.EarnClaim) (stop bool) {
		cs = append(cs, c)
		return false
	})
	return cs
}

// SetJinxSupplyRewardIndexes sets the current reward indexes for an individual denom
func (k Keeper) SetJinxSupplyRewardIndexes(ctx sdk.Context, denom string, indexes types.RewardIndexes) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.JinxSupplyRewardIndexesKeyPrefix)
	bz := k.cdc.MustMarshal(&types.RewardIndexesProto{
		RewardIndexes: indexes,
	})
	store.Set([]byte(denom), bz)
}

// GetJinxSupplyRewardIndexes gets the current reward indexes for an individual denom
func (k Keeper) GetJinxSupplyRewardIndexes(ctx sdk.Context, denom string) (types.RewardIndexes, bool) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.JinxSupplyRewardIndexesKeyPrefix)
	bz := store.Get([]byte(denom))
	if bz == nil {
		return types.RewardIndexes{}, false
	}
	var proto types.RewardIndexesProto
	k.cdc.MustUnmarshal(bz, &proto)

	return proto.RewardIndexes, true
}

// IterateJinxSupplyRewardIndexes iterates over all Jinx supply reward index objects in the store and preforms a callback function
func (k Keeper) IterateJinxSupplyRewardIndexes(ctx sdk.Context, cb func(denom string, indexes types.RewardIndexes) (stop bool)) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.JinxSupplyRewardIndexesKeyPrefix)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var proto types.RewardIndexesProto
		k.cdc.MustUnmarshal(iterator.Value(), &proto)
		if cb(string(iterator.Key()), proto.RewardIndexes) {
			break
		}
	}
}

func (k Keeper) IterateJinxSupplyRewardAccrualTimes(ctx sdk.Context, cb func(string, time.Time) (stop bool)) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.PreviousJinxSupplyRewardAccrualTimeKeyPrefix)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var accrualTime time.Time
		if err := accrualTime.UnmarshalBinary(iterator.Value()); err != nil {
			panic(err)
		}
		denom := string(iterator.Key())
		if cb(denom, accrualTime) {
			break
		}
	}
}

// SetJinxBorrowRewardIndexes sets the current reward indexes for an individual denom
func (k Keeper) SetJinxBorrowRewardIndexes(ctx sdk.Context, denom string, indexes types.RewardIndexes) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.JinxBorrowRewardIndexesKeyPrefix)
	bz := k.cdc.MustMarshal(&types.RewardIndexesProto{
		RewardIndexes: indexes,
	})
	store.Set([]byte(denom), bz)
}

// GetJinxBorrowRewardIndexes gets the current reward indexes for an individual denom
func (k Keeper) GetJinxBorrowRewardIndexes(ctx sdk.Context, denom string) (types.RewardIndexes, bool) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.JinxBorrowRewardIndexesKeyPrefix)
	bz := store.Get([]byte(denom))
	if bz == nil {
		return types.RewardIndexes{}, false
	}
	var proto types.RewardIndexesProto
	k.cdc.MustUnmarshal(bz, &proto)

	return proto.RewardIndexes, true
}

// IterateJinxBorrowRewardIndexes iterates over all Jinx borrow reward index objects in the store and preforms a callback function
func (k Keeper) IterateJinxBorrowRewardIndexes(ctx sdk.Context, cb func(denom string, indexes types.RewardIndexes) (stop bool)) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.JinxBorrowRewardIndexesKeyPrefix)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var proto types.RewardIndexesProto
		k.cdc.MustUnmarshal(iterator.Value(), &proto)
		if cb(string(iterator.Key()), proto.RewardIndexes) {
			break
		}
	}
}

func (k Keeper) IterateJinxBorrowRewardAccrualTimes(ctx sdk.Context, cb func(string, time.Time) (stop bool)) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.PreviousJinxBorrowRewardAccrualTimeKeyPrefix)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		denom := string(iterator.Key())
		var accrualTime time.Time
		if err := accrualTime.UnmarshalBinary(iterator.Value()); err != nil {
			panic(err)
		}
		if cb(denom, accrualTime) {
			break
		}
	}
}

// GetDelegatorRewardIndexes gets the current reward indexes for an individual denom
func (k Keeper) GetDelegatorRewardIndexes(ctx sdk.Context, denom string) (types.RewardIndexes, bool) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.DelegatorRewardIndexesKeyPrefix)
	bz := store.Get([]byte(denom))
	if bz == nil {
		return types.RewardIndexes{}, false
	}
	var proto types.RewardIndexesProto
	k.cdc.MustUnmarshal(bz, &proto)

	return proto.RewardIndexes, true
}

// SetDelegatorRewardIndexes sets the current reward indexes for an individual denom
func (k Keeper) SetDelegatorRewardIndexes(ctx sdk.Context, denom string, indexes types.RewardIndexes) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.DelegatorRewardIndexesKeyPrefix)
	bz := k.cdc.MustMarshal(&types.RewardIndexesProto{
		RewardIndexes: indexes,
	})
	store.Set([]byte(denom), bz)
}

// IterateDelegatorRewardIndexes iterates over all delegator reward index objects in the store and preforms a callback function
func (k Keeper) IterateDelegatorRewardIndexes(ctx sdk.Context, cb func(denom string, indexes types.RewardIndexes) (stop bool)) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.DelegatorRewardIndexesKeyPrefix)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var proto types.RewardIndexesProto
		k.cdc.MustUnmarshal(iterator.Value(), &proto)
		if cb(string(iterator.Key()), proto.RewardIndexes) {
			break
		}
	}
}

func (k Keeper) IterateDelegatorRewardAccrualTimes(ctx sdk.Context, cb func(string, time.Time) (stop bool)) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.PreviousDelegatorRewardAccrualTimeKeyPrefix)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		denom := string(iterator.Key())
		var accrualTime time.Time
		if err := accrualTime.UnmarshalBinary(iterator.Value()); err != nil {
			panic(err)
		}
		if cb(denom, accrualTime) {
			break
		}
	}
}

// GetPreviousJinxSupplyRewardAccrualTime returns the last time a denom accrued Jinx protocol supply-side rewards
func (k Keeper) GetPreviousJinxSupplyRewardAccrualTime(ctx sdk.Context, denom string) (blockTime time.Time, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.PreviousJinxSupplyRewardAccrualTimeKeyPrefix)
	bz := store.Get([]byte(denom))
	if bz == nil {
		return time.Time{}, false
	}
	if err := blockTime.UnmarshalBinary(bz); err != nil {
		panic(err)
	}
	return blockTime, true
}

// SetPreviousJinxSupplyRewardAccrualTime sets the last time a denom accrued Jinx protocol supply-side rewards
func (k Keeper) SetPreviousJinxSupplyRewardAccrualTime(ctx sdk.Context, denom string, blockTime time.Time) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.PreviousJinxSupplyRewardAccrualTimeKeyPrefix)
	bz, err := blockTime.MarshalBinary()
	if err != nil {
		panic(err)
	}
	store.Set([]byte(denom), bz)
}

// GetPreviousJinxBorrowRewardAccrualTime returns the last time a denom accrued Jinx protocol borrow-side rewards
func (k Keeper) GetPreviousJinxBorrowRewardAccrualTime(ctx sdk.Context, denom string) (blockTime time.Time, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.PreviousJinxBorrowRewardAccrualTimeKeyPrefix)
	b := store.Get([]byte(denom))
	if b == nil {
		return time.Time{}, false
	}
	if err := blockTime.UnmarshalBinary(b); err != nil {
		panic(err)
	}
	return blockTime, true
}

// SetPreviousJinxBorrowRewardAccrualTime sets the last time a denom accrued Jinx protocol borrow-side rewards
func (k Keeper) SetPreviousJinxBorrowRewardAccrualTime(ctx sdk.Context, denom string, blockTime time.Time) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.PreviousJinxBorrowRewardAccrualTimeKeyPrefix)
	bz, err := blockTime.MarshalBinary()
	if err != nil {
		panic(err)
	}
	store.Set([]byte(denom), bz)
}

// GetPreviousDelegatorRewardAccrualTime returns the last time a denom accrued protocol delegator rewards
func (k Keeper) GetPreviousDelegatorRewardAccrualTime(ctx sdk.Context, denom string) (blockTime time.Time, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.PreviousDelegatorRewardAccrualTimeKeyPrefix)
	bz := store.Get([]byte(denom))
	if bz == nil {
		return time.Time{}, false
	}
	if err := blockTime.UnmarshalBinary(bz); err != nil {
		panic(err)
	}
	return blockTime, true
}

// SetPreviousDelegatorRewardAccrualTime sets the last time a denom accrued protocol delegator rewards
func (k Keeper) SetPreviousDelegatorRewardAccrualTime(ctx sdk.Context, denom string, blockTime time.Time) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.PreviousDelegatorRewardAccrualTimeKeyPrefix)
	bz, err := blockTime.MarshalBinary()
	if err != nil {
		panic(err)
	}
	store.Set([]byte(denom), bz)
}

// SetSwapRewardIndexes stores the global reward indexes that track total rewards to a swap pool.
func (k Keeper) SetSwapRewardIndexes(ctx sdk.Context, poolID string, indexes types.RewardIndexes) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.SwapRewardIndexesKeyPrefix)
	bz := k.cdc.MustMarshal(&types.RewardIndexesProto{
		RewardIndexes: indexes,
	})
	store.Set([]byte(poolID), bz)
}

// GetSwapRewardIndexes fetches the global reward indexes that track total rewards to a swap pool.
func (k Keeper) GetSwapRewardIndexes(ctx sdk.Context, poolID string) (types.RewardIndexes, bool) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.SwapRewardIndexesKeyPrefix)
	bz := store.Get([]byte(poolID))
	if bz == nil {
		return types.RewardIndexes{}, false
	}
	var proto types.RewardIndexesProto
	k.cdc.MustUnmarshal(bz, &proto)
	return proto.RewardIndexes, true
}

// IterateSwapRewardIndexes iterates over all swap reward index objects in the store and preforms a callback function
func (k Keeper) IterateSwapRewardIndexes(ctx sdk.Context, cb func(poolID string, indexes types.RewardIndexes) (stop bool)) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.SwapRewardIndexesKeyPrefix)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var proto types.RewardIndexesProto
		k.cdc.MustUnmarshal(iterator.Value(), &proto)
		if cb(string(iterator.Key()), proto.RewardIndexes) {
			break
		}
	}
}

// GetSwapRewardAccrualTime fetches the last time rewards were accrued for a swap pool.
func (k Keeper) GetSwapRewardAccrualTime(ctx sdk.Context, poolID string) (blockTime time.Time, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.PreviousSwapRewardAccrualTimeKeyPrefix)
	b := store.Get([]byte(poolID))
	if b == nil {
		return time.Time{}, false
	}
	if err := blockTime.UnmarshalBinary(b); err != nil {
		panic(err)
	}
	return blockTime, true
}

// SetSwapRewardAccrualTime stores the last time rewards were accrued for a swap pool.
func (k Keeper) SetSwapRewardAccrualTime(ctx sdk.Context, poolID string, blockTime time.Time) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.PreviousSwapRewardAccrualTimeKeyPrefix)
	bz, err := blockTime.MarshalBinary()
	if err != nil {
		panic(err)
	}
	store.Set([]byte(poolID), bz)
}

func (k Keeper) IterateSwapRewardAccrualTimes(ctx sdk.Context, cb func(string, time.Time) (stop bool)) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.PreviousSwapRewardAccrualTimeKeyPrefix)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		poolID := string(iterator.Key())
		var accrualTime time.Time
		if err := accrualTime.UnmarshalBinary(iterator.Value()); err != nil {
			panic(err)
		}
		if cb(poolID, accrualTime) {
			break
		}
	}
}

// SetSavingsRewardIndexes stores the global reward indexes that rewards for an individual denom type
func (k Keeper) SetSavingsRewardIndexes(ctx sdk.Context, denom string, indexes types.RewardIndexes) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.SavingsRewardIndexesKeyPrefix)
	bz := k.cdc.MustMarshal(&types.RewardIndexesProto{
		RewardIndexes: indexes,
	})
	store.Set([]byte(denom), bz)
}

// GetSavingsRewardIndexes fetches the global reward indexes that track rewards for an individual denom type
func (k Keeper) GetSavingsRewardIndexes(ctx sdk.Context, denom string) (types.RewardIndexes, bool) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.SavingsRewardIndexesKeyPrefix)
	bz := store.Get([]byte(denom))
	if bz == nil {
		return types.RewardIndexes{}, false
	}
	var proto types.RewardIndexesProto
	k.cdc.MustUnmarshal(bz, &proto)
	return proto.RewardIndexes, true
}

// IterateSavingsRewardIndexes iterates over all savings reward index objects in the store and preforms a callback function
func (k Keeper) IterateSavingsRewardIndexes(ctx sdk.Context, cb func(poolID string, indexes types.RewardIndexes) (stop bool)) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.SavingsRewardIndexesKeyPrefix)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var proto types.RewardIndexesProto
		k.cdc.MustUnmarshal(iterator.Value(), &proto)
		if cb(string(iterator.Key()), proto.RewardIndexes) {
			break
		}
	}
}

// GetSavingsRewardAccrualTime fetches the last time rewards were accrued for an individual denom type
func (k Keeper) GetSavingsRewardAccrualTime(ctx sdk.Context, poolID string) (blockTime time.Time, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.PreviousSavingsRewardAccrualTimeKeyPrefix)
	b := store.Get([]byte(poolID))
	if b == nil {
		return time.Time{}, false
	}
	if err := blockTime.UnmarshalBinary(b); err != nil {
		panic(err)
	}
	return blockTime, true
}

// SetSavingsRewardAccrualTime stores the last time rewards were accrued for a savings deposit denom type
func (k Keeper) SetSavingsRewardAccrualTime(ctx sdk.Context, poolID string, blockTime time.Time) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.PreviousSavingsRewardAccrualTimeKeyPrefix)
	bz, err := blockTime.MarshalBinary()
	if err != nil {
		panic(err)
	}
	store.Set([]byte(poolID), bz)
}

// IterateSavingsRewardAccrualTimesiterates over all the previous savings reward accrual times in the store
func (k Keeper) IterateSavingsRewardAccrualTimes(ctx sdk.Context, cb func(string, time.Time) (stop bool)) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.PreviousSavingsRewardAccrualTimeKeyPrefix)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		poolID := string(iterator.Key())
		var accrualTime time.Time
		if err := accrualTime.UnmarshalBinary(iterator.Value()); err != nil {
			panic(err)
		}
		if cb(poolID, accrualTime) {
			break
		}
	}
}

// SetEarnRewardIndexes stores the global reward indexes that track total rewards to a earn vault.
func (k Keeper) SetEarnRewardIndexes(ctx sdk.Context, vaultDenom string, indexes types.RewardIndexes) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.EarnRewardIndexesKeyPrefix)
	bz := k.cdc.MustMarshal(&types.RewardIndexesProto{
		RewardIndexes: indexes,
	})
	store.Set([]byte(vaultDenom), bz)
}

// GetEarnRewardIndexes fetches the global reward indexes that track total rewards to a earn vault.
func (k Keeper) GetEarnRewardIndexes(ctx sdk.Context, vaultDenom string) (types.RewardIndexes, bool) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.EarnRewardIndexesKeyPrefix)
	bz := store.Get([]byte(vaultDenom))
	if bz == nil {
		return types.RewardIndexes{}, false
	}
	var proto types.RewardIndexesProto
	k.cdc.MustUnmarshal(bz, &proto)
	return proto.RewardIndexes, true
}

// IterateEarnRewardIndexes iterates over all earn reward index objects in the store and preforms a callback function
func (k Keeper) IterateEarnRewardIndexes(ctx sdk.Context, cb func(vaultDenom string, indexes types.RewardIndexes) (stop bool)) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.EarnRewardIndexesKeyPrefix)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var proto types.RewardIndexesProto
		k.cdc.MustUnmarshal(iterator.Value(), &proto)
		if cb(string(iterator.Key()), proto.RewardIndexes) {
			break
		}
	}
}

// GetEarnRewardAccrualTime fetches the last time rewards were accrued for an earn vault.
func (k Keeper) GetEarnRewardAccrualTime(ctx sdk.Context, vaultDenom string) (blockTime time.Time, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.PreviousEarnRewardAccrualTimeKeyPrefix)
	b := store.Get([]byte(vaultDenom))
	if b == nil {
		return time.Time{}, false
	}
	if err := blockTime.UnmarshalBinary(b); err != nil {
		panic(err)
	}
	return blockTime, true
}

// SetEarnRewardAccrualTime stores the last time rewards were accrued for a earn vault.
func (k Keeper) SetEarnRewardAccrualTime(ctx sdk.Context, vaultDenom string, blockTime time.Time) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.PreviousEarnRewardAccrualTimeKeyPrefix)
	bz, err := blockTime.MarshalBinary()
	if err != nil {
		panic(err)
	}
	store.Set([]byte(vaultDenom), bz)
}

func (k Keeper) IterateEarnRewardAccrualTimes(ctx sdk.Context, cb func(string, time.Time) (stop bool)) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.PreviousEarnRewardAccrualTimeKeyPrefix)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		poolID := string(iterator.Key())
		var accrualTime time.Time
		if err := accrualTime.UnmarshalBinary(iterator.Value()); err != nil {
			panic(err)
		}
		if cb(poolID, accrualTime) {
			break
		}
	}
}
