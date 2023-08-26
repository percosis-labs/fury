package keeper_test

import (
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/percosis-labs/fury/app"
	"github.com/percosis-labs/fury/x/auction/keeper"
	"github.com/percosis-labs/fury/x/auction/types"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

func TestGrpcAuctionsFilter(t *testing.T) {
	// setup
	tApp := app.NewTestApp()
	tApp.InitializeFromGenesisStates()
	auctionsKeeper := tApp.GetAuctionKeeper()
	ctx := tApp.NewContext(true, tmproto.Header{Height: 1})
	_, addrs := app.GeneratePrivKeyAddressPairs(2)

	auctions := []types.Auction{
		types.NewSurplusAuction(
			"sellerMod",
			c("mer", 12345678),
			"usdf",
			time.Date(1998, time.January, 1, 0, 0, 0, 0, time.UTC),
		).WithID(0),
		types.NewDebtAuction(
			"buyerMod",
			c("jinx", 12345678),
			c("usdf", 12345678),
			time.Date(1998, time.January, 1, 0, 0, 0, 0, time.UTC),
			c("debt", 12345678),
		).WithID(1),
		types.NewCollateralAuction(
			"sellerMod",
			c("ufury", 12345678),
			time.Date(1998, time.January, 1, 0, 0, 0, 0, time.UTC),
			c("usdf", 12345678),
			types.WeightedAddresses{
				Addresses: addrs,
				Weights:   []sdkmath.Int{sdkmath.NewInt(100)},
			},
			c("debt", 12345678),
		).WithID(2),
		types.NewCollateralAuction(
			"sellerMod",
			c("jinx", 12345678),
			time.Date(1998, time.January, 1, 0, 0, 0, 0, time.UTC),
			c("usdf", 12345678),
			types.WeightedAddresses{
				Addresses: addrs,
				Weights:   []sdkmath.Int{sdkmath.NewInt(100)},
			},
			c("debt", 12345678),
		).WithID(3),
	}
	for _, a := range auctions {
		auctionsKeeper.SetAuction(ctx, a)
	}

	qs := keeper.NewQueryServerImpl(auctionsKeeper)

	tests := []struct {
		giveName     string
		giveRequest  types.QueryAuctionsRequest
		wantResponse []types.Auction
	}{
		{
			"empty request",
			types.QueryAuctionsRequest{},
			auctions,
		},
		{
			"denom query mer",
			types.QueryAuctionsRequest{
				Denom: "mer",
			},
			auctions[0:1],
		},
		{
			"denom query usdf all",
			types.QueryAuctionsRequest{
				Denom: "usdf",
			},
			auctions,
		},
		{
			"owner",
			types.QueryAuctionsRequest{
				Owner: addrs[0].String(),
			},
			auctions[2:4],
		},
		{
			"owner and denom",
			types.QueryAuctionsRequest{
				Owner: addrs[0].String(),
				Denom: "jinx",
			},
			auctions[3:4],
		},
		{
			"owner, denom, type, phase",
			types.QueryAuctionsRequest{
				Owner: addrs[0].String(),
				Denom: "jinx",
				Type:  types.CollateralAuctionType,
				Phase: types.ForwardAuctionPhase,
			},
			auctions[3:4],
		},
	}

	for _, tc := range tests {
		t.Run(tc.giveName, func(t *testing.T) {
			res, err := qs.Auctions(sdk.WrapSDKContext(ctx), &tc.giveRequest)
			require.NoError(t, err)

			var unpackedAuctions []types.Auction

			for _, anyAuction := range res.Auctions {
				var auction types.Auction
				err := tApp.AppCodec().UnpackAny(anyAuction, &auction)
				require.NoError(t, err)

				unpackedAuctions = append(unpackedAuctions, auction)
			}

			require.Equal(t, tc.wantResponse, unpackedAuctions)
		})
	}
}
