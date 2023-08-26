package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// MultiJINXHooks combine multiple JINX hooks, all hook functions are run in array sequence
type MultiJINXHooks []JINXHooks

// NewMultiJINXHooks returns a new MultiJINXHooks
func NewMultiJINXHooks(hooks ...JINXHooks) MultiJINXHooks {
	return hooks
}

// AfterDepositCreated runs after a deposit is created
func (h MultiJINXHooks) AfterDepositCreated(ctx sdk.Context, deposit Deposit) {
	for i := range h {
		h[i].AfterDepositCreated(ctx, deposit)
	}
}

// BeforeDepositModified runs before a deposit is modified
func (h MultiJINXHooks) BeforeDepositModified(ctx sdk.Context, deposit Deposit) {
	for i := range h {
		h[i].BeforeDepositModified(ctx, deposit)
	}
}

// AfterDepositModified runs after a deposit is modified
func (h MultiJINXHooks) AfterDepositModified(ctx sdk.Context, deposit Deposit) {
	for i := range h {
		h[i].AfterDepositModified(ctx, deposit)
	}
}

// AfterBorrowCreated runs after a borrow is created
func (h MultiJINXHooks) AfterBorrowCreated(ctx sdk.Context, borrow Borrow) {
	for i := range h {
		h[i].AfterBorrowCreated(ctx, borrow)
	}
}

// BeforeBorrowModified runs before a borrow is modified
func (h MultiJINXHooks) BeforeBorrowModified(ctx sdk.Context, borrow Borrow) {
	for i := range h {
		h[i].BeforeBorrowModified(ctx, borrow)
	}
}

// AfterBorrowModified runs after a borrow is modified
func (h MultiJINXHooks) AfterBorrowModified(ctx sdk.Context, borrow Borrow) {
	for i := range h {
		h[i].AfterBorrowModified(ctx, borrow)
	}
}
