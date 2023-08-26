package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/percosis-labs/fury/x/incentive/types"
)

type msgServer struct {
	keeper Keeper
}

// NewMsgServerImpl returns an implementation of the incentive MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) ClaimUSDFMintingReward(goCtx context.Context, msg *types.MsgClaimUSDFMintingReward) (*types.MsgClaimUSDFMintingRewardResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	err = k.keeper.ClaimUSDFMintingReward(ctx, sender, sender, msg.MultiplierName)
	if err != nil {
		return nil, err
	}

	return &types.MsgClaimUSDFMintingRewardResponse{}, nil
}

func (k msgServer) ClaimJinxReward(goCtx context.Context, msg *types.MsgClaimJinxReward) (*types.MsgClaimJinxRewardResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	for _, selection := range msg.DenomsToClaim {
		err := k.keeper.ClaimJinxReward(ctx, sender, sender, selection.Denom, selection.MultiplierName)
		if err != nil {
			return nil, err
		}

	}

	return &types.MsgClaimJinxRewardResponse{}, nil
}

func (k msgServer) ClaimDelegatorReward(goCtx context.Context, msg *types.MsgClaimDelegatorReward) (*types.MsgClaimDelegatorRewardResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	for _, selection := range msg.DenomsToClaim {
		err := k.keeper.ClaimDelegatorReward(ctx, sender, sender, selection.Denom, selection.MultiplierName)
		if err != nil {
			return nil, err
		}
	}

	return &types.MsgClaimDelegatorRewardResponse{}, nil
}

func (k msgServer) ClaimSwapReward(goCtx context.Context, msg *types.MsgClaimSwapReward) (*types.MsgClaimSwapRewardResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	for _, selection := range msg.DenomsToClaim {
		err := k.keeper.ClaimSwapReward(ctx, sender, sender, selection.Denom, selection.MultiplierName)
		if err != nil {
			return nil, err
		}
	}

	return &types.MsgClaimSwapRewardResponse{}, nil
}

func (k msgServer) ClaimSavingsReward(goCtx context.Context, msg *types.MsgClaimSavingsReward) (*types.MsgClaimSavingsRewardResponse, error) {
	err := errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "savings claims disabled")
	return nil, err
}

func (k msgServer) ClaimEarnReward(goCtx context.Context, msg *types.MsgClaimEarnReward) (*types.MsgClaimEarnRewardResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	for _, selection := range msg.DenomsToClaim {
		err := k.keeper.ClaimEarnReward(ctx, sender, sender, selection.Denom, selection.MultiplierName)
		if err != nil {
			return nil, err
		}
	}

	return &types.MsgClaimEarnRewardResponse{}, nil
}
