<!--
order: 3
-->

# Messages

Users claim rewards using messages that correspond to each claim type.

```go
// MsgClaimUSDFMintingReward message type used to claim USDF minting rewards
type MsgClaimUSDFMintingReward struct {
	Sender         sdk.AccAddress `json:"sender" yaml:"sender"`
	MultiplierName string         `json:"multiplier_name" yaml:"multiplier_name"`
}

// MsgClaimJinxReward message type used to claim Jinx liquidity provider rewards
type MsgClaimJinxReward struct {
	Sender         sdk.AccAddress `json:"sender" yaml:"sender"`
	MultiplierName string         `json:"multiplier_name" yaml:"multiplier_name"`
	DenomsToClaim  []string       `json:"denoms_to_claim" yaml:"denoms_to_claim"`
}

// MsgClaimDelegatorReward message type used to claim delegator rewards
type MsgClaimDelegatorReward struct {
	Sender         sdk.AccAddress `json:"sender" yaml:"sender"`
	MultiplierName string         `json:"multiplier_name" yaml:"multiplier_name"`
	DenomsToClaim  []string       `json:"denoms_to_claim" yaml:"denoms_to_claim"`
}

// MsgClaimSwapReward message type used to claim delegator rewards
type MsgClaimSwapReward struct {
	Sender         sdk.AccAddress `json:"sender" yaml:"sender"`
	MultiplierName string         `json:"multiplier_name" yaml:"multiplier_name"`
	DenomsToClaim  []string       `json:"denoms_to_claim" yaml:"denoms_to_claim"`
}
```

## State Modifications

- Accumulated rewards for active claims are transferred from the `furydist` module account to the users account as vesting coins
- The number of coins transferred is determined by the multiplier in the message. For example, the multiplier equals 1.0, 100% of the claim's reward value is transferred. If the multiplier equals 0.5, 50% of the claim's reward value is transferred.
- The corresponding claim object is reset to zero in the store
