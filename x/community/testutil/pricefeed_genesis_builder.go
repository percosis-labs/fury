package testutil

import (
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	jinxtypes "github.com/percosis-labs/fury/x/jinx/types"
	pricefeedtypes "github.com/percosis-labs/fury/x/pricefeed/types"
)

// lendGenesisBuilder builds the Jinx and Pricefeed genesis states for setting up Fury Lend
type lendGenesisBuilder struct {
	jinxMarkets []jinxtypes.MoneyMarket
	pfMarkets   []pricefeedtypes.Market
	prices      []pricefeedtypes.PostedPrice
}

func NewLendGenesisBuilder() lendGenesisBuilder {
	return lendGenesisBuilder{}
}

func (b lendGenesisBuilder) Build() (jinxtypes.GenesisState, pricefeedtypes.GenesisState) {
	jinxGS := jinxtypes.DefaultGenesisState()
	jinxGS.Params.MoneyMarkets = b.jinxMarkets

	pricefeedGS := pricefeedtypes.DefaultGenesisState()
	pricefeedGS.Params.Markets = b.pfMarkets
	pricefeedGS.PostedPrices = b.prices
	return jinxGS, pricefeedGS
}

func (b lendGenesisBuilder) WithMarket(denom, spotMarketId string, price sdk.Dec) lendGenesisBuilder {
	// add jinx money market
	b.jinxMarkets = append(b.jinxMarkets,
		jinxtypes.NewMoneyMarket(
			denom,
			jinxtypes.NewBorrowLimit(false, sdk.NewDec(1e15), sdk.MustNewDecFromStr("0.6")),
			spotMarketId,
			sdkmath.NewInt(1e6),
			jinxtypes.NewInterestRateModel(sdk.MustNewDecFromStr("0.05"), sdk.MustNewDecFromStr("2"), sdk.MustNewDecFromStr("0.8"), sdk.MustNewDecFromStr("10")),
			sdk.MustNewDecFromStr("0.05"),
			sdk.ZeroDec(),
		),
	)

	// add pricefeed
	b.pfMarkets = append(b.pfMarkets,
		pricefeedtypes.Market{MarketID: spotMarketId, BaseAsset: denom, QuoteAsset: "usd", Oracles: []sdk.AccAddress{}, Active: true},
	)
	b.prices = append(b.prices,
		pricefeedtypes.PostedPrice{
			MarketID:      spotMarketId,
			OracleAddress: sdk.AccAddress{},
			Price:         price,
			Expiry:        time.Now().Add(100 * time.Hour),
		},
	)

	return b
}
