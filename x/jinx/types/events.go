package types

// Event types for jinx module
const (
	EventTypeJinxDeposit          = "jinx_deposit"
	EventTypeJinxWithdrawal       = "jinx_withdrawal"
	EventTypeJinxBorrow           = "jinx_borrow"
	EventTypeJinxLiquidation      = "jinx_liquidation"
	EventTypeJinxRepay            = "jinx_repay"
	AttributeValueCategory        = ModuleName
	AttributeKeyDeposit           = "deposit"
	AttributeKeyDepositDenom      = "deposit_denom"
	AttributeKeyDepositCoins      = "deposit_coins"
	AttributeKeyDepositor         = "depositor"
	AttributeKeyBorrow            = "borrow"
	AttributeKeyBorrower          = "borrower"
	AttributeKeyBorrowCoins       = "borrow_coins"
	AttributeKeySender            = "sender"
	AttributeKeyRepayCoins        = "repay_coins"
	AttributeKeyLiquidatedOwner   = "liquidated_owner"
	AttributeKeyLiquidatedCoins   = "liquidated_coins"
	AttributeKeyKeeper            = "keeper"
	AttributeKeyKeeperRewardCoins = "keeper_reward_coins"
	AttributeKeyOwner             = "owner"
)
