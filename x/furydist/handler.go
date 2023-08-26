package furydist

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
errorsmod "cosmossdk.io/errors"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	"github.com/percosis-labs/fury/x/furydist/keeper"
	"github.com/percosis-labs/fury/x/furydist/types"
)

// NewCommunityPoolMultiSpendProposalHandler
func NewCommunityPoolMultiSpendProposalHandler(k keeper.Keeper) govv1beta1.Handler {
	return func(ctx sdk.Context, content govv1beta1.Content) error {
		switch c := content.(type) {
		case *types.CommunityPoolMultiSpendProposal:
			return keeper.HandleCommunityPoolMultiSpendProposal(ctx, k, c)
		default:
			return errorsmod.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized furydist proposal content type: %T", c)
		}
	}
}
