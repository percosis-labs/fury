<!--
order: 6
-->

# Hooks

This module implements the `Hooks` interface for the following modules:

- cdp
- jinx
- swap
- staking (defined in cosmos-sdk)

CDP module hooks manage the creation and synchronization of USDF minting incentives.

```go
// ------------------- Cdp Module Hooks -------------------

// AfterCDPCreated function that runs after a cdp is created
func (h Hooks) AfterCDPCreated(ctx sdk.Context, cdp cdptypes.CDP) {
  h.k.InitializeUSDFMintingClaim(ctx, cdp)
}

// BeforeCDPModified function that runs before a cdp is modified
// note that this is called immediately after interest is synchronized, and so could potentially
// be called AfterCDPInterestUpdated or something like that, if we we're to expand the scope of cdp hooks
func (h Hooks) BeforeCDPModified(ctx sdk.Context, cdp cdptypes.CDP) {
  h.k.SynchronizeUSDFMintingReward(ctx, cdp)
}
```

Jinx module hooks manage the creation and synchronization of jinx supply and borrow rewards.

```go
// ------------------- Jinx Module Hooks -------------------

// AfterDepositCreated function that runs after a deposit is created
func (h Hooks) AfterDepositCreated(ctx sdk.Context, deposit jinxtypes.Deposit) {
  h.k.InitializeJinxSupplyReward(ctx, deposit)
}

// BeforeDepositModified function that runs before a deposit is modified
func (h Hooks) BeforeDepositModified(ctx sdk.Context, deposit jinxtypes.Deposit) {
  h.k.SynchronizeJinxSupplyReward(ctx, deposit)
}

// AfterDepositModified function that runs after a deposit is modified
func (h Hooks) AfterDepositModified(ctx sdk.Context, deposit jinxtypes.Deposit) {
  h.k.UpdateJinxSupplyIndexDenoms(ctx, deposit)
}

// AfterBorrowCreated function that runs after a borrow is created
func (h Hooks) AfterBorrowCreated(ctx sdk.Context, borrow jinxtypes.Borrow) {
  h.k.InitializeJinxBorrowReward(ctx, borrow)
}

// BeforeBorrowModified function that runs before a borrow is modified
func (h Hooks) BeforeBorrowModified(ctx sdk.Context, borrow jinxtypes.Borrow) {
  h.k.SynchronizeJinxBorrowReward(ctx, borrow)
}

// AfterBorrowModified function that runs after a borrow is modified
func (h Hooks) AfterBorrowModified(ctx sdk.Context, borrow jinxtypes.Borrow) {
  h.k.UpdateJinxBorrowIndexDenoms(ctx, borrow)
}
```

Staking module hooks manage the creation and synchronization of jinx delegator rewards.

```go
// ------------------- Staking Module Hooks -------------------

// BeforeDelegationCreated runs before a delegation is created
func (h Hooks) BeforeDelegationCreated(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
  h.k.InitializeJinxDelegatorReward(ctx, delAddr)
}

// BeforeDelegationSharesModified runs before an existing delegation is modified
func (h Hooks) BeforeDelegationSharesModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
  h.k.SynchronizeJinxDelegatorRewards(ctx, delAddr)
}

// NOTE: following hooks are just implemented to ensure StakingHooks interface compliance

// BeforeValidatorSlashed is called before a validator is slashed
func (h Hooks) BeforeValidatorSlashed(ctx sdk.Context, valAddr sdk.ValAddress, fraction sdk.Dec) {}

// AfterValidatorBeginUnbonding is called after a validator begins unbonding
func (h Hooks) AfterValidatorBeginUnbonding(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
}

// AfterValidatorBonded is called after a validator is bonded
func (h Hooks) AfterValidatorBonded(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
}

// AfterDelegationModified runs after a delegation is modified
func (h Hooks) AfterDelegationModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
}

// BeforeDelegationRemoved runs directly before a delegation is deleted
func (h Hooks) BeforeDelegationRemoved(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
}

// AfterValidatorCreated runs after a validator is created
func (h Hooks) AfterValidatorCreated(ctx sdk.Context, valAddr sdk.ValAddress) {}

// BeforeValidatorModified runs before a validator is modified
func (h Hooks) BeforeValidatorModified(ctx sdk.Context, valAddr sdk.ValAddress) {}

// AfterValidatorRemoved runs after a validator is removed
func (h Hooks) AfterValidatorRemoved(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
}
```

Swap module hooks manage the creation and synchronization of Swap protocol liquidity provider rewards.

```go
// ------------------- Swap Module Hooks -------------------

func (h Hooks) AfterPoolDepositCreated(ctx sdk.Context, poolID string, depositor sdk.AccAddress, _ sdkmath.Int) {
	h.k.InitializeSwapReward(ctx, poolID, depositor)
}

func (h Hooks) BeforePoolDepositModified(ctx sdk.Context, poolID string, depositor sdk.AccAddress, sharesOwned sdkmath.Int) {
	h.k.SynchronizeSwapReward(ctx, poolID, depositor, sharesOwned)
}
```
