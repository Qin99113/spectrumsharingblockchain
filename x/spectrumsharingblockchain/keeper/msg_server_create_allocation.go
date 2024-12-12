package keeper

import (
	"context"

	"spectrumSharingBlockchain/x/spectrumsharingblockchain/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) CreateAllocation(goCtx context.Context, msg *types.MsgCreateAllocation) (*types.MsgCreateAllocationResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	return &types.MsgCreateAllocationResponse{}, nil
}
