package types

const (
	// ModuleName The name that will be used throughout the module
	ModuleName = "incentive"

	// StoreKey Top level store key where all module items will be stored
	StoreKey = ModuleName

	// RouterKey Top level router key
	RouterKey = ModuleName

	// DefaultParamspace default name for parameter store
	DefaultParamspace = ModuleName

	// QuerierRoute route used for abci queries
	QuerierRoute = ModuleName
)

// Key Prefixes
var (
	USDFMintingClaimKeyPrefix                     = []byte{0x01} // prefix for keys that store USDF minting claims
	USDFMintingRewardFactorKeyPrefix              = []byte{0x02} // prefix for key that stores USDF minting reward factors
	PreviousUSDFMintingRewardAccrualTimeKeyPrefix = []byte{0x03} // prefix for key that stores the blocktime
	JinxLiquidityClaimKeyPrefix                   = []byte{0x04} // prefix for keys that store Jinx liquidity claims
	JinxSupplyRewardIndexesKeyPrefix              = []byte{0x05} // prefix for key that stores Jinx supply reward indexes
	PreviousJinxSupplyRewardAccrualTimeKeyPrefix  = []byte{0x06} // prefix for key that stores the previous time Jinx supply rewards accrued
	JinxBorrowRewardIndexesKeyPrefix              = []byte{0x07} // prefix for key that stores Jinx borrow reward indexes
	PreviousJinxBorrowRewardAccrualTimeKeyPrefix  = []byte{0x08} // prefix for key that stores the previous time Jinx borrow rewards accrued
	DelegatorClaimKeyPrefix                       = []byte{0x09} // prefix for keys that store delegator claims
	DelegatorRewardIndexesKeyPrefix               = []byte{0x10} // prefix for key that stores delegator reward indexes
	PreviousDelegatorRewardAccrualTimeKeyPrefix   = []byte{0x11} // prefix for key that stores the previous time delegator rewards accrued
	SwapClaimKeyPrefix                            = []byte{0x12} // prefix for keys that store swap claims
	SwapRewardIndexesKeyPrefix                    = []byte{0x13} // prefix for key that stores swap reward indexes
	PreviousSwapRewardAccrualTimeKeyPrefix        = []byte{0x14} // prefix for key that stores the previous time swap rewards accrued
	SavingsClaimKeyPrefix                         = []byte{0x15} // prefix for keys that store savings claims
	SavingsRewardIndexesKeyPrefix                 = []byte{0x16} // prefix for key that stores savings reward indexes
	PreviousSavingsRewardAccrualTimeKeyPrefix     = []byte{0x17} // prefix for key that stores the previous time savings rewards accrued
	EarnClaimKeyPrefix                            = []byte{0x18} // prefix for keys that store earn claims
	EarnRewardIndexesKeyPrefix                    = []byte{0x19} // prefix for key that stores earn reward indexes
	PreviousEarnRewardAccrualTimeKeyPrefix        = []byte{0x20} // prefix for key that stores the previous time earn rewards accrued
)
