package keeper

import (
	"context"

	"spectrumSharingBlockchain/x/spectrumrequest/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) CreateRequest(goCtx context.Context, msg *types.MsgCreateRequest) (*types.MsgCreateRequestResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	// _ = ctx

	requestID := k.GetNextRequestID(ctx)

	request := types.SpectrumRequest{
		Id:           requestID,
		Creator:      msg.Creator,
		Organization: msg.Organization,
		Frequency:    msg.Frequency,
		Bandwidth:    msg.Bandwidth,
		Duration:     msg.Duration,
		BidAmount:    msg.BidAmount,
		Status:       types.StatusPending,
		RequestTime:  msg.RequestTime,
	}

	k.SetSpectrumRequest(ctx, request)

	return &types.MsgCreateRequestResponse{Id: requestID}, nil
}
